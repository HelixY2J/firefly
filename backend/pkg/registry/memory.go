package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/HelixY2J/firefly/backend/common/api"
	"github.com/google/uuid"
)

type InMemoryRegistry struct {
	mu            sync.RWMutex
	nodes         map[string]*NodeInfo
	masterURL     string
	libraryStore  *LibraryStore
	playingStatus map[string]bool
	pauseStatus   map[string]bool
}

func NewInMemoryRegistry(masterURL string) *InMemoryRegistry {
	return &InMemoryRegistry{
		nodes:         make(map[string]*NodeInfo),
		masterURL:     masterURL,
		libraryStore:  NewLibraryStore(),
		playingStatus: make(map[string]bool),
		pauseStatus:   make(map[string]bool),
	}
}

func (r *InMemoryRegistry) RegisterNode(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeID := uuid.New().String()
	r.nodes[nodeID] = &NodeInfo{
		NodeID:   nodeID,
		Address:  "unknown",
		LastSeen: time.Now(),
		IsMaster: req.NodeType == "master",
	}

	return &pb.RegisterResponse{
		AssignedId: nodeID,
		MasterUrl:  r.masterURL,
	}, nil
}

func (r *InMemoryRegistry) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[req.NodeId]
	if !exists {
		return fmt.Errorf("node not found")
	}
	node.LastSeen = time.Now()
	return nil
}

func (r *InMemoryRegistry) GetActiveNodes() []*NodeInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeNodes []*NodeInfo
	for _, node := range r.nodes {
		activeNodes = append(activeNodes, node)
	}
	return activeNodes
}

func (r *InMemoryRegistry) CleanupInactiveNodes() {
	r.mu.Lock()
	defer r.mu.Unlock()
	timeout := 30 * time.Second
	now := time.Now()

	for id, node := range r.nodes {
		if now.Sub(node.LastSeen) > timeout {
			delete(r.nodes, id)
		}
	}
}

func (r *InMemoryRegistry) SyncLibrary(nodeID string, files []FileMetadata) {
	r.libraryStore.SyncFiles(nodeID, files)
}

func (r *InMemoryRegistry) GetAvailableSongs() []string {
	return r.libraryStore.GetAvailableSongs()
}

func (r *InMemoryRegistry) SetPlayingStatus(nodeID string, isPlaying bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.playingStatus[nodeID] = isPlaying
	if !isPlaying {
		r.pauseStatus[nodeID] = false
	}
}

func (r *InMemoryRegistry) IsPlaying(nodeID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.playingStatus[nodeID]
}

func (r *InMemoryRegistry) SetPauseStatus(nodeID string, isPaused bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pauseStatus[nodeID] = isPaused
}

func (r *InMemoryRegistry) IsPaused(nodeID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.pauseStatus[nodeID]
}

func (r *InMemoryRegistry) CanSendPlaybackCommand(nodeID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return !r.playingStatus[nodeID]
}

func (r *InMemoryRegistry) SyncPlayback(nodeID string, action string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if action == "PLAY" && !r.playingStatus[nodeID] {
		r.playingStatus[nodeID] = true
		r.pauseStatus[nodeID] = false
		return true
	} else if action == "PAUSE" && r.playingStatus[nodeID] {
		r.pauseStatus[nodeID] = true
		r.playingStatus[nodeID] = false
		return true
	}
	return false
}

func (r *InMemoryRegistry) FinishPlayback(nodeID string) {
	r.SetPlayingStatus(nodeID, false)
}
func (r *InMemoryRegistry) HandleWebSocketCommand(nodeID string, command string) bool {
	return r.SyncPlayback(nodeID, command)
}
