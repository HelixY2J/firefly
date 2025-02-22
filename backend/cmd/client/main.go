package main

import (
	"context"
	"log"
	"time"

	"github.com/HelixY2J/firefly/backend/pkg/discovery/consul"
	grpcclient "github.com/HelixY2J/firefly/backend/pkg/grpc_client"
)

func main() {
	consulClient, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatalf(" Failed to connect to Consul: %v", err)
	}

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

	client.RegisterNode()
}
