package main

import (
	"context"
	"log"
	"time"

	"github.com/HelixY2J/firefly/backend/pkg/discovery"
	"github.com/HelixY2J/firefly/backend/pkg/discovery/consul"
	grpcserver "github.com/HelixY2J/firefly/backend/pkg/grpc_server"
	"github.com/HelixY2J/firefly/backend/pkg/registry"
	"github.com/HelixY2J/firefly/backend/pkg/websocket"
)

func main() {
	// Initialize Consul client
	reg := registry.NewInMemoryRegistry("localhost:50051")
	registryService := registry.NewRegistryService(reg)

	consulClient, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatalf(" Failed to connect to Consul: %v", err)
	}
	instanceID := discovery.GenerateInstanceID("firefly-master")

	err = consulClient.Register(context.Background(), instanceID, "firefly-master", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	relay := websocket.NewWebSocketRelay()
    
    // Start WebSocket server for GUI connections
    go func() {
        if err := relay.StartServer(":8080"); err != nil {
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

	server := grpcserver.NewGRPCServer(registryService)
	if err := server.Start(50051); err != nil {
		log.Fatalf(" Failed to start gRPC server: %v", err)
	}
}
