package data

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/margen2/goknition/backend/models"
)

// LoadFaces receives a path and returns all of the faces within the given path
func LoadFaces(path string) ([]models.Face, error) {
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

			if len(fls) != 1 {
				return nil, fmt.Errorf("loadfaces/range: ID folder has an invalid amount of files: %s", filePath)
			}

			fl := fls[0]
			faces = append(faces, models.Face{
				0, ID, models.Image{0, fl.Name(), filePath}})

			continue
		}
		return nil, fmt.Errorf("loadfaces: unexpected file in path")
	}
	return faces, nil
}

var images []models.Image

func getImages(path string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("os.readdir: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			path := filepath.Join(path, dir.Name())
			err = getImages(path)
			if err != nil {
				return fmt.Errorf("getimages: %w", err)
			}
			continue
		}
		images = append(images, models.Image{0, dir.Name(), path})
	}

	return nil
}

// loadimages receives a path to an images folder. It returns a slice of type Image
// with all of the images in the corresponding path.
func Loadimages(imagesPath string) ([]models.Image, error) {
	err := getImages(imagesPath)
	if err != nil {
		return nil, fmt.Errorf("getimages: %w", err)
	}
	defer func() {
		images = nil
	}()

	return images, nil
}

// CopyImages copies all the given images to the specified path.
func CopyImages(faceID, copy string, images []models.Image) error {

	err := os.Mkdir(copy+`\`+faceID, fs.ModeDir)
	if err != nil {
		return fmt.Errorf("os.writefile: %w", err)
	}

	for _, image := range images {
		fl, err := os.Open(image.Path + `\` + image.Filename)
		if err != nil {
			return fmt.Errorf("os.open: %w", err)
		}
		defer fl.Close()

		fs, err := fl.Stat()
		if err != nil {
			return fmt.Errorf("fl.stat: %w", err)
		}

		b := make([]byte, fs.Size())
		n, err := fl.Read(b)
		if n < 1 || err != nil {
			return fmt.Errorf("no bytes read or err fl.read: %w", err)
		}

		err = os.WriteFile(copy+`\`+faceID+`\`+image.Filename, b, 0666)
		if err != nil {
			return fmt.Errorf("os.writefile: %w", err)
		}
	}

	return nil
}
