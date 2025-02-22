package registry

import (
	"context"

	pb "github.com/HelixY2J/firefly/backend/common/api"
)

type RegistryService struct {
	Registry Registry
}

func NewRegistryService(reg Registry) *RegistryService {
	return &RegistryService{Registry: reg}
}

func (s *RegistryService) RegisterNode(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.Registry.RegisterNode(ctx, req)
}

func (s *RegistryService) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) error {
	return s.Registry.Heartbeat(ctx, req)
}

func (s *RegistryService) GetActiveNodes() []*NodeInfo {
	return s.Registry.GetActiveNodes()
}

func (s *RegistryService) CleanupInactiveNodes() {
	s.Registry.CleanupInactiveNodes()
}
