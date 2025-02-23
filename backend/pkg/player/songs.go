package player

import (
	"log"
	"os"

	pb "github.com/HelixY2J/firefly/backend/common/api"
)

func GetAvailableSongs() []*pb.FileMetadata {
	files := []*pb.FileMetadata{}

	songDir := "songs/"
	songFiles, err := os.ReadDir(songDir)
	if err != nil {
		log.Printf("Could not read songs directory: %v", err)
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

// func GetAvailableSongsMap() map[string]bool {
// 	songMap := make(map[string]bool)

// 	songDir := "cmd/client/songs/"
// 	songFiles, err := os.ReadDir(songDir)
// 	if err != nil {
// 		log.Printf("cant read songs directory: %v", err)
// 		return songMap
// 	}

// 	for _, file := range songFiles {
// 		if !file.IsDir() {
// 			songMap[file.Name()] = true
// 		}
// 	}

// 	return songMap
// }
