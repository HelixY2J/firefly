package grpcclient

import (
	"context"
	"io"
	"log"

	pb "github.com/HelixY2J/firefly/backend/common/api"
	"github.com/HelixY2J/firefly/backend/pkg/player"
	"github.com/faiface/beep/speaker"
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

func (c *GRPCClient) ListenForPlayback(nodeID string) {
	for {
		stream, err := c.client.SyncPlayback(context.Background(), &pb.SyncPlaybackCommand{NodeId: nodeID})
		if err != nil {
			log.Printf(" Failed to start playback listener: %v", err)
			continue
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println(" Playback stream ended. Reconnecting...")
				break
			}
			if err != nil {
				log.Printf("eerror receiving playback command: %v", err)
				break
			}

			log.Printf("[GRPCClient] received playback command: %s - %s", resp.Status, resp.Filename)

			// Play or stop music based on the command
			if resp.Status == "PLAY" {
				go player.PlaySong(resp.Filename)
			} else if resp.Status == "STOP" {
				log.Printf("Stopping playback: %s", resp.Filename)
				speaker.Clear()
			}
		}
	}
}
