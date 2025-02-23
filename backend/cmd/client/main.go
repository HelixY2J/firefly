package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/HelixY2J/firefly/backend/pkg/discovery"
	"github.com/HelixY2J/firefly/backend/pkg/discovery/consul"
	grpcclient "github.com/HelixY2J/firefly/backend/pkg/grpc_client"
	"github.com/HelixY2J/firefly/backend/pkg/player"
)

var (
	service = "firefly-client"
)

func main() {
	consulClient, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatalf(" Failed to connect to Consul: %v", err)
	}
	rand.Seed(time.Now().UnixNano())
	clientPort := 50052 + rand.Intn(8)
	clientAddress := fmt.Sprintf("localhost:%d", clientPort)
	instanceID := discovery.GenerateInstanceID(service, clientPort)

	err = consulClient.Register(context.Background(), instanceID, service, clientAddress)
	if err != nil {
		log.Fatalf("dFailed to register client in Consul: %v", err)
	}
	log.Printf("reg client in Consul with ID: %s", instanceID)

	go func() {
		for {
			err := consulClient.HealthCheck(instanceID, service)
			if err != nil {
				log.Printf(" Failed to send health check: %v", err)
			} else {
				log.Println("client health check sent")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	var masterAddr []string
	for i := 0; i < 10; i++ {
		masterAddr, err = consulClient.Discover(context.Background(), "firefly-master")
		if err != nil {
			log.Printf("Error querying Consul: %v", err)
		} else if len(masterAddr) > 0 {
			break
		}
		log.Println("Waiting for master to register in Consul...")
		time.Sleep(1 * time.Second)
	}

	if len(masterAddr) == 0 {
		log.Fatalf(" Could not find master node after retries. Last error: %v", err)
	}

	log.Printf(" Found master at: %s", masterAddr[0])

	client := grpcclient.NewClient(masterAddr[0])
	defer client.Close()

	nodeID := client.RegisterNode()

	// files := []*pb.FileMetadata{
	// 	{
	// 		Filename: "test_song.mp3",
	// 		Checksum: "abc123",
	// 		Chunks: []*pb.ChunkMetadata{
	// 			{Fingerprint: "chunk1_hash", Size: 1024},
	// 		},
	// 	},
	// }

	files := player.GetAvailableSongs()

	resp, err := client.SyncLibrary(nodeID, files)
	if err != nil {
		log.Fatalf(" SyncLibrary failed: %v", err)
	}

	log.Printf("SyncLibrary successful, missing files: %v", resp.MissingFiles)
	log.Println("Client is now listening for playback commands...")
	client.ListenForPlayback(nodeID)
}
