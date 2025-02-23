package player

import (
	"log"
	"os"
	"path/filepath"

	pb "github.com/HelixY2J/firefly/backend/common/api"
)

func GetMasterSongs() []*pb.FileMetadata {
	files := []*pb.FileMetadata{}
	wd, _ := os.Getwd()
	songDir := filepath.Join(wd, "songs")
	songFiles, err := os.ReadDir(songDir)
	if err != nil {
		log.Printf(" Could not read master songs directory: %v", err)
		return files
	}

	for _, file := range songFiles {
		if !file.IsDir() {
			log.Printf("found song: %s", file.Name())
			files = append(files, &pb.FileMetadata{
				Filename: file.Name(),
				Checksum: "unknown",
			})
		}
	}

	return files
}
