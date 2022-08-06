package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ID struct {
	Name string
	Image
}
type Image struct {
	ImageID string
	Path    string
	Bytes   []byte
}

// LoadIds receives a path and returns a slice of type ID that contains all of the
// ids and files in the specified path
func LoadIDs(path string) []ID {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Errorf("loadids: %w", err))
	}

	var IDs []ID
	for _, dir := range dirs {
		if dir.IsDir() {
			log.Fatal("loadids: file is dir")
		}

		filePath := filepath.Join(path, dir.Name())
		b := getImageBytes(filePath)
		IDs = append(IDs, ID{
			dir.Name(), Image{dir.Name(), filePath, b}})
	}
	return IDs
}

// LoadImages receives a path and returns a slice of type Image that contains
// all of the files in the specified path
func LoadImages(path string, images []Image) []Image {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Errorf("copyimagefiles: %w", err))
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			path = filepath.Join(path, dir.Name())
			LoadImages(path, images)
			break
		}

		filePath := filepath.Join(path, dir.Name())
		b := getImageBytes(filePath)
		images = append(images, Image{dir.Name(), filePath, b})
	}
	return images
}

func getImageBytes(filePath string) []byte {
	fl, err := os.Open(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("getimagebytes/os.open: %w", err))
	}
	defer fl.Close()

	fileInfo, err := fl.Stat()
	if err != nil {
		log.Fatal(fmt.Errorf("getimagebytes/fl.open: %w", err))
	}
	var size int64 = fileInfo.Size()
	b := make([]byte, size)

	n, err := fl.Read(b)
	if err != nil || n == 0 {
		log.Fatal(fmt.Errorf("getimagebytes/fl.read: %w", err))
	}

	return b
}
