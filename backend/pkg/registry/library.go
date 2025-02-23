package registry

import (
	"sync"

	"github.com/HelixY2J/firefly/backend/pkg/player"
)

type LibraryStore struct {
	mu    sync.RWMutex
	files map[string][]string
}

type FileMetadata struct {
	Filename string
	Checksum string
	Chunks   []ChunkMetadata
}

type ChunkMetadata struct {
	Fingerprint string
	Size        int64
}

func NewLibraryStore() *LibraryStore {
	return &LibraryStore{
		files: make(map[string][]string),
	}
}

func (l *LibraryStore) SyncFiles(nodeID string, files []FileMetadata) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// []FileMetadata - return
	// var missingFiles []FileMetadata
	// for _, file := range files {
	// 	if _, exists := l.files[file.Checksum]; !exists {
	// 		missingFiles = append(missingFiles, file)
	// 		l.files[file.Checksum] = file
	// 	}
	// }
	// return missingFiles

	for _, file := range files {
		l.files[file.Filename] = append(l.files[file.Filename], nodeID)
	}
}

func (s *LibraryStore) GetAvailableSongs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	localFiles := player.GetAvailableSongs() // Fetch actual files from disk

	var availableSongs []string
	for _, file := range localFiles {
		availableSongs = append(availableSongs, file.Filename)
	}
	return availableSongs
}
