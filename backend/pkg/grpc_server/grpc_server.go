package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/HelixY2J/firefly/backend/common/api"

	"github.com/HelixY2J/firefly/backend/pkg/registry"
	"github.com/HelixY2J/firefly/backend/pkg/websocket"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedFireflyServiceServer
	server   *grpc.Server
	registry registry.Registry
	relay    *websocket.WebSocketRelay
}

func NewGRPCServer(reg registry.Registry, relay *websocket.WebSocketRelay) *GRPCServer {
	grpcServer := grpc.NewServer()
	s := &GRPCServer{
		server:   grpcServer,
		registry: reg,
		relay:    relay,
	}
	pb.RegisterFireflyServiceServer(grpcServer, s)
	return s
}

func (s *GRPCServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("gRPC server started on %s", address)
	return s.server.Serve(listener)
}

func (s *GRPCServer) Stop() {
	log.Println("Shutting down gRPC server...")
	s.server.GracefulStop()
}

func (s *GRPCServer) RegisterNode(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.registry.RegisterNode(ctx, req)
}

func (s *GRPCServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	err := s.registry.Heartbeat(ctx, req)
	if err != nil {
		return &pb.HeartbeatResponse{Success: false}, err
	}
	return &pb.HeartbeatResponse{Success: true}, nil
}

func (s *GRPCServer) SyncLibrary(ctx context.Context, req *pb.SyncLibraryRequest) (*pb.SyncLibraryResponse, error) {
	files := make([]registry.FileMetadata, len(req.Files))
	for i, f := range req.Files {
		chunks := make([]registry.ChunkMetadata, len(f.Chunks))
		for j, c := range f.Chunks {
			chunks[j] = registry.ChunkMetadata{
				Fingerprint: c.Fingerprint,
				Size:        c.Size,
			}
		}
		files[i] = registry.FileMetadata{
			Filename: f.Filename,
			Checksum: f.Checksum,
			Chunks:   chunks,
		}
	}
	s.registry.SyncLibrary(req.NodeId, files)

	return &pb.SyncLibraryResponse{}, nil

}

func (s *GRPCServer) RequestPlayback(ctx context.Context, req *pb.PlaybackRequest) (*pb.PlaybackResponse, error) {
	return &pb.PlaybackResponse{}, nil
}

func (s *GRPCServer) SyncPlayback(req *pb.SyncPlaybackCommand, stream pb.FireflyService_SyncPlaybackServer) error {
	log.Printf("Client %s started listening for playback commands", req.NodeId)

	for {
		command := s.relay.GetLastCommand()
		if command == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		if command != "PLAY" && command != "PAUSE" {
			log.Printf("Invalid playback command received: %s", command)
			continue
		}

		if s.registry.SyncPlayback(req.NodeId, command) {
			cmd := &pb.SyncPlaybackResponse{
				Filename: "",
				Status:   command,
			}
			log.Printf("Sending playback command: %s", command)

			if err := stream.Send(cmd); err != nil {
				log.Printf("Failed to send playback command: %v", err)
				return err
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// func (s *GRPCServer) SyncPlayback(req *pb.SyncPlaybackCommand, stream pb.FireflyService_SyncPlaybackServer) error {
// 	log.Printf("Client %s started listening for playback commands", req.NodeId)

// 	availableSongs := s.registry.GetAvailableSongs()
// 	if len(availableSongs) == 0 {
// 		log.Println(" No songs available for playback.")
// 		return nil
// 	}

// 	// Pick a song to play
// 	for _, song := range availableSongs {
// 		cmd := &pb.SyncPlaybackResponse{
// 			Filename: song,
// 			Status:   "PLAY",
// 		}
// 		log.Printf(" Sending playback command: %s", song)
// 		if err := stream.Send(cmd); err != nil {
// 			log.Printf("Failed to send playback command: %v", err)
// 			return err
// 		}
// 		log.Printf("Sent playback command: %s", song)
// 		time.Sleep(5 * time.Second)
// 	}

// 	return nil
// }
