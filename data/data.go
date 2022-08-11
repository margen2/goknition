package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/margen2/goknition/models"
)

// LoadFaces receives a path and returns a slice of type Face that contains all of the
// IDs and files in the specified path
func LoadFaces(path string) []models.Face {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Errorf("loadfaces: %w", err))
	}

	var faces []models.Face
	for _, dir := range dirs {
		if dir.IsDir() {
			ID := dir.Name()

			filePath := filepath.Join(path, ID)
			fls, err := os.ReadDir(filePath)
			if err != nil {
				log.Fatal(fmt.Errorf("loadfaces: %w", err))
			}

			fl := fls[0]
			filePath = filepath.ToSlash(filePath)
			imageID := strings.Split(fl.Name(), ".")[0]

			faces = append(faces, models.Face{
				ID, models.Image{imageID, filePath}})

			continue
		}
		log.Fatal("loadfaces: unexpected file in path")
	}
	return faces
}

// LoadImages receives a path and returns a slice of type Image that contains
// all of the files in the specified path
func LoadImages(path string, images []models.Image) []models.Image {
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

		path = filepath.ToSlash(path)
		imageID := strings.Split(dir.Name(), ".")[0]

		images = append(images, models.Image{imageID, path})
	}
	return images
}
