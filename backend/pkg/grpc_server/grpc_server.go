package grpcserver

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	grpcServer := grpc.NewServer()
	return &GRPCServer{
		server: grpcServer,
	}
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
