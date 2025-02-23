package grpcclient

import (
	"context"
	"log"

	pb "github.com/HelixY2J/firefly/backend/common/api"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.FireflyServiceClient
}

func NewClient(masterAddr string) *GRPCClient {
	conn, err := grpc.Dial(masterAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}

	client := pb.NewFireflyServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}
}

func (c *GRPCClient) RegisterNode() string {
	req := &pb.RegisterRequest{
		NodeId:   "",
		NodeType: "client",
	}

	resp, err := c.client.RegisterNode(context.Background(), req)
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}

	log.Printf("Registered with Master, Assigned ID: %s, Master URL: %s", resp.AssignedId, resp.MasterUrl)
	return resp.AssignedId
}

func (c *GRPCClient) Close() {
	c.conn.Close()
}

func (c *GRPCClient) SyncLibrary(nodeID string, files []*pb.FileMetadata) (*pb.SyncLibraryResponse, error) {
	req := &pb.SyncLibraryRequest{
		NodeId: nodeID,
		Files:  files,
	}
	return c.client.SyncLibrary(context.Background(), req)
}

func (c *GRPCClient) ListenForPlayback() {
	stream, err := c.client.SyncPlayback(context.Background(), &pb.SyncPlaybackCommand{})
	if err != nil {
		log.Fatalf(" Failed to start playback listener: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf(" Error receiving playback command: %v", err)
		}

		log.Printf(" Received playback command: Start playing %s", resp.Filename)

		// addd logic for music player
	}
}
