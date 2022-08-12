package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/margen2/goknition/models"
)

func loadFaces(path string) ([]models.Face, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("loadfaces/os.readdir: %w", err)
	}

	var faces []models.Face
	for _, dir := range dirs {
		if dir.IsDir() {
			ID := dir.Name()

			filePath := filepath.Join(path, ID)
			fls, err := os.ReadDir(filePath)
			if err != nil {
				return nil, fmt.Errorf("loadfaces/range/os.readdir: %w", err)
			}

			fmt.Println(len(fls))
			if len(fls) != 1 {
				return nil, fmt.Errorf("loadfaces/range: ID folder has an invalid amount of files: %s", filePath)
			}

			fl := fls[0]
			imageID := strings.Split(fl.Name(), ".")[0]
			faces = append(faces, models.Face{
				ID, models.Image{imageID, filePath}})

			continue
		}
		return nil, fmt.Errorf("loadfaces: unexpected file in path")
	}
	return faces, nil
}

var images []models.Image

func loadImages(path string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("loadimages/os.readdir: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			path := filepath.Join(path, dir.Name())
			err = loadImages(path)
			if err != nil {
				return fmt.Errorf("loadimages/range dirs: %w", err)
			}
			continue
		}
		imageID := strings.Split(dir.Name(), ".")[0]
		images = append(images, models.Image{imageID, path})
	}

	return nil
}

// Load receives two different paths, one to an images folder, and the other to a faces folder. It returns
//a slice of type Image and a slice of type Face with all of the images and faces in the corresponding paths.
func Load(imagesPath, facesPath string) ([]models.Image, []models.Face, error) {
	err := loadImages(imagesPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load: %w", err)
	}
	defer func() {
		images = images[cap(images):]
	}()

	faces, err := loadFaces(facesPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load: %w", err)
	}

	return images, faces, nil
}
