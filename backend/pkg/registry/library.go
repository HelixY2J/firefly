package registry

import (
	"sync"
)

type LibraryStore struct {
	mu    sync.RWMutex
	files map[string]FileMetadata
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
		files: make(map[string]FileMetadata),
	}
}

func (l *LibraryStore) SyncFiles(nodeID string, files []FileMetadata) []FileMetadata {
	l.mu.Lock()
	defer l.mu.Unlock()

	var missingFiles []FileMetadata
	for _, file := range files {
		if _, exists := l.files[file.Checksum]; !exists {
			missingFiles = append(missingFiles, file)
			l.files[file.Checksum] = file
		}
	}
	return missingFiles
}
