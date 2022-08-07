package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Face struct {
	ID    string
	Image Image
}

type Image struct {
	ID    string
	Path  string
	Bytes []byte
}

// LoadFaces receives a path and returns a slice of type Face that contains all of the
// IDs and files in the specified path
func LoadFaces(path string) []Face {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Errorf("loadfaces: %w", err))
	}

	var faces []Face
	for _, dir := range dirs {
		if dir.IsDir() {
			ID := dir.Name()

			filePath := filepath.Join(path, ID)
			fls, err := os.ReadDir(filePath)
			if err != nil {
				log.Fatal(fmt.Errorf("loadfaces: %w", err))
			}

			fl := fls[0]

			filePath = filepath.Join(filePath, fl.Name())
			b := getImageBytes(filePath)
			imageID := strings.Split(fl.Name(), ".")[0]

			faces = append(faces, Face{
				ID, Image{imageID, filePath, b}})

			continue
		}
		log.Fatal("loadfaces: unexpected file in path")
	}
	return faces
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
		imageID := strings.Split(dir.Name(), ".")[0]

		images = append(images, Image{imageID, filePath, b})
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
