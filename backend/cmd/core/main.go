package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"sync"
    "reflect"
    "encoding/json"
	"runtime"
	"strings"
	"crypto/sha256"
    "os"
    "path/filepath"

	"github.com/HelixY2J/firefly/backend/pkg/discovery"
	"github.com/HelixY2J/firefly/backend/pkg/discovery/consul"
	grpcserver "github.com/HelixY2J/firefly/backend/pkg/grpc_server"
	"github.com/HelixY2J/firefly/backend/pkg/player"
	"github.com/HelixY2J/firefly/backend/pkg/registry"
	"github.com/HelixY2J/firefly/backend/pkg/websocket"
)

type NodeState struct {
    Nodes []string
    mutex sync.RWMutex
}

// Add this struct to store memory stats
type MemoryStats struct {
    Alloc      uint64  `json:"alloc"`
    TotalAlloc uint64  `json:"totalAlloc"`
    HeapAlloc  uint64  `json:"heapAlloc"`
    HeapSys    uint64  `json:"heapSys"`
}

func findProjectRoot() (string, error) {
    // Start from the current working directory
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }

    // Walk up the directory tree until we find the "firefly" root
    for {
        // Check if this directory contains the res/music folder
        if _, err := os.Stat(filepath.Join(dir, "res", "music")); err == nil {
            return dir, nil
        }

        // Move up one directory
        parent := filepath.Dir(dir)
        if parent == dir {
            return "", fmt.Errorf("could not find project root with res/music directory")
        }
        dir = parent
    }
}

func getMemoryStats() MemoryStats {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    return MemoryStats{
        Alloc:      m.Alloc,
        TotalAlloc: m.TotalAlloc,
        HeapAlloc:  m.HeapAlloc,
        HeapSys:    m.HeapSys,
    }
}

// Add this function before main()
func getMP3Files(registryService *registry.RegistryService) []string {
    files := registryService.LibraryStore.GetAllFiles()
    mp3Files := make([]string, 0)
    log.Printf("Debug: Found %d total files in LibraryStore", len(files))
    
    for _, file := range files {
        if strings.HasSuffix(file.Filename, ".wav") {
            mp3Files = append(mp3Files, file.Filename)
        }
    }
    log.Printf("Debug: Filtered to %d MP3 files", len(mp3Files))
    return mp3Files
}

// Add this function to scan directory and create FileMetadata
func loadMusicFiles(dirPath string) ([]registry.FileMetadata, error) {
    var files []registry.FileMetadata
    
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory: %v", err)
    }

    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".wav") {
            continue
        }

        filePath := filepath.Join(dirPath, entry.Name())
        
        // Read file for checksum
        data, err := os.ReadFile(filePath)
        if err != nil {
            log.Printf("Warning: Could not read file %s: %v", entry.Name(), err)
            continue
        }

        // Calculate checksum
        checksum := fmt.Sprintf("%x", sha256.Sum256(data))
        
        // Create chunk metadata (you might want to actually split the file)
        fileInfo, err := entry.Info()
        if err != nil {
            log.Printf("Warning: Could not get file info for %s: %v", entry.Name(), err)
            continue
        }

        // For now, treating entire file as one chunk
        chunk := registry.ChunkMetadata{
            Fingerprint: checksum[:16], // Using first 16 chars of checksum as fingerprint
            Size:        fileInfo.Size(),
        }

        files = append(files, registry.FileMetadata{
            Filename: entry.Name(),
            Checksum: checksum,
            Chunks:   []registry.ChunkMetadata{chunk},
        })
    }

    return files, nil
}

func main() {
	masterPort := 50051
    masterAddress := fmt.Sprintf("localhost:%d", masterPort)
    reg := registry.NewInMemoryRegistry(masterAddress)
    registryService := registry.NewRegistryService(reg)

    // Load music files from res/music directory
    projectRoot, err := findProjectRoot()
    if err != nil {
        log.Fatalf("Failed to find project root: %v", err)
    }
    musicDir := filepath.Join(projectRoot, "res", "music")
    log.Printf("Loading music from: %s", musicDir)
    files, err := loadMusicFiles(musicDir)
    if err != nil {
        log.Printf("Warning: Failed to load music files: %v", err)
    } else {
        registryService.LibraryStore.SyncFiles("master-node", files)
        log.Printf("Loaded %d music files from %s", len(files), musicDir)
    }

	consulClient, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatalf(" Failed to connect to Consul: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	instanceID := discovery.GenerateInstanceID("firefly-master", masterPort)

	err = consulClient.Register(context.Background(), instanceID, "firefly-master", masterAddress)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	log.Printf("registered Master in Consul with ID: %s at %s", instanceID, masterAddress)

	relay := websocket.NewWebSocketRelay()
    
    // Start WebSocket server for GUI connections
    go func() {
        if err := relay.StartServer(":8081"); err != nil {
            log.Fatalf("Failed to start WebSocket server: %v", err)
        } else{
			log.Println("WebSocket server started")
		}
    }()

	// Start ealth checks
	go func() {
		for {
			err := consulClient.HealthCheck(instanceID, "firefly-master")
			if err != nil {
				log.Printf("Failed to send health check: %v", err)
			} else {
				log.Println("Health check sent")
			}
			time.Sleep(5 * time.Second)
		}
	}()
	
	nodeState := &NodeState{
        Nodes: make([]string, 0),
    }

	go func() {
        for {
            clients, err := consulClient.Discover(context.Background(), "firefly-client")
            if err != nil {
                log.Printf("Failed to discover clients: %v", err)
                time.Sleep(10 * time.Second)
                continue
            }

            nodeState.mutex.Lock()
            // Check if there's a difference in nodes
            if !reflect.DeepEqual(nodeState.Nodes, clients) {
                // Update state
                nodeState.Nodes = clients

                // Prepare message for websocket
                message := map[string]interface{}{
                    "type": "nodes_update",
                    "nodes": clients,
                }
                
                // Send to websocket clients
                jsonMsg, _ := json.Marshal(message)
                relay.Broadcast(jsonMsg)
                
                log.Printf("Active Clients updated: %v", clients)
            }
            nodeState.mutex.Unlock()
            
            time.Sleep(10 * time.Second)
        }
    }()

	go func() {
		for {
			memStats := getMemoryStats()
			mp3Files := getMP3Files(registryService)
	
			// Prepare message for websocket
			message := map[string]interface{}{
				"type": "master_stats",
				"memory": memStats,
				"mp3_files": mp3Files,
			}
	
			jsonMsg, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling stats: %v", err)
				continue
			}
			log.Printf("Debug: Broadcasting message with %d MP3 files", len(mp3Files))
			relay.Broadcast(jsonMsg)
	
			time.Sleep(5 * time.Second)
		}
	}()

	server := grpcserver.NewGRPCServer(registryService)
	if err := server.Start(50051); err != nil {
		log.Fatalf(" Failed to start gRPC server: %v", err)
	}
}
