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

	"github.com/HelixY2J/firefly/backend/pkg/discovery"
	"github.com/HelixY2J/firefly/backend/pkg/discovery/consul"
	grpcserver "github.com/HelixY2J/firefly/backend/pkg/grpc_server"
	"github.com/HelixY2J/firefly/backend/pkg/registry"
	"github.com/HelixY2J/firefly/backend/pkg/websocket"
)

type NodeState struct {
    Nodes []string
    mutex sync.RWMutex
}

func main() {
	masterPort := 50051
	masterAddress := fmt.Sprintf("localhost:%d", masterPort)
	// Initialize Consul client
	reg := registry.NewInMemoryRegistry(masterAddress)
	registryService := registry.NewRegistryService(reg)

	registryService.LibraryStore.SyncFiles("master-node", []registry.FileMetadata{
		{
			Filename: "server_song.mp3",
			Checksum: "server123",
			Chunks: []registry.ChunkMetadata{
				{Fingerprint: "server_chunk1", Size: 2048},
			},
		},
		{
			Filename: "test_song.mp3",
			Checksum: "abc123",
			Chunks: []registry.ChunkMetadata{
				{Fingerprint: "chunk1_hash", Size: 1024},
			},
		},
		{
			Filename: "background_music.mp3",
			Checksum: "server456",
			Chunks: []registry.ChunkMetadata{
				{Fingerprint: "server_chunk2", Size: 3072},
			},
		},
	})
	log.Println("loaded server files into the library")

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

	server := grpcserver.NewGRPCServer(registryService)
	if err := server.Start(50051); err != nil {
		log.Fatalf(" Failed to start gRPC server: %v", err)
	}
}
