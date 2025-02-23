package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/HelixY2J/firefly/backend/common/api"

	"github.com/HelixY2J/firefly/backend/pkg/registry"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedFireflyServiceServer
	server   *grpc.Server
	registry registry.Registry
}

func NewGRPCServer(reg registry.Registry) *GRPCServer {
	grpcServer := grpc.NewServer()
	s := &GRPCServer{
		server:   grpcServer,
		registry: reg,
	}
	pb.RegisterFireflyServiceServer(grpcServer, s)
	return s
}

func (s *GRPCServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("gRPC server started on port %d", port)
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

	missingFiles := s.registry.SyncLibrary(req.NodeId, files)

	var responseFiles []*pb.FileMetadata
	for _, file := range missingFiles {
		var chunks []*pb.ChunkMetadata
		for _, c := range file.Chunks {
			chunks = append(chunks, &pb.ChunkMetadata{
				Fingerprint: c.Fingerprint,
				Size:        c.Size,
			})
		}
		responseFiles = append(responseFiles, &pb.FileMetadata{
			Filename: file.Filename,
			Checksum: file.Checksum,
			Chunks:   chunks,
		})
	}

	return &pb.SyncLibraryResponse{
		MissingFiles: responseFiles,
	}, nil
}

func (s *GRPCServer) RequestPlayback(ctx context.Context, req *pb.PlaybackRequest) (*pb.PlaybackResponse, error) {
	return &pb.PlaybackResponse{}, nil
}

func (s *GRPCServer) SyncPlayback(ctx context.Context, req *pb.SyncPlaybackCommand) (*pb.SyncPlaybackResponse, error) {
	return &pb.SyncPlaybackResponse{}, nil
}
