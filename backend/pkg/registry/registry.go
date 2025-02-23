package registry

import (
	"context"

	pb "github.com/HelixY2J/firefly/backend/common/api"
)

type Registry interface {
	RegisterNode(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) error
	GetActiveNodes() []*NodeInfo
	CleanupInactiveNodes()
	GetAvailableSongs() []string
	SyncLibrary(nodeID string, files []FileMetadata)
}
