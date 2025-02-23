package registry

import (
	"context"

	pb "github.com/HelixY2J/firefly/backend/common/api"
)

type RegistryService struct {
	Registry     Registry
	LibraryStore *LibraryStore
}

func NewRegistryService(reg Registry) *RegistryService {
	return &RegistryService{
		Registry:     reg,
		LibraryStore: NewLibraryStore(),
	}
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

func (s *RegistryService) SyncLibrary(nodeID string, files []FileMetadata) {
	s.LibraryStore.SyncFiles(nodeID, files)
}

func (s *RegistryService) GetAvailableSongs() []string {
	return s.LibraryStore.GetAvailableSongs()
}

func (s *RegistryService) SyncPlayback(nodeID string, action string) bool {
	return s.Registry.SyncPlayback(nodeID, action)
}
