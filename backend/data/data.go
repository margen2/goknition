package data

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/margen2/goknition/backend/models"
)

// LoadFaces receives a path and returns all of the faces within the given path
func LoadFaces(path string) ([]models.Face, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("os.readdir: %w", err)
	}

	var faces []models.Face
	for _, dir := range dirs {
		if !dir.IsDir() {
			return nil, fmt.Errorf("unexpected file in path")
		}

		var face models.Face

		face.FaceID = dir.Name()
		dirPath := filepath.Join(path, dir.Name())
		images, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("os.readdir: %w", err)
		}

		for _, image := range images {
			if image.IsDir() {
				return nil, fmt.Errorf("unexpected dir in path")
			}

			if strings.ToLower(filepath.Ext(image.Name())) != ".jpg" {
				continue
			}
			face.Images = append(face.Images, models.Image{ID: 0, Filename: image.Name(), Path: dirPath})

		}
		faces = append(faces, face)
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

		if strings.ToLower(filepath.Ext(dir.Name())) != ".jpg" {
			continue
		}
		images = append(images, models.Image{ID: 0, Filename: dir.Name(), Path: path})
	}

	return nil
}

// Loadimages receives a path to an images folder. It returns a slice of type Image
// with all the images found in the corresponding path.
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

	err := os.Mkdir(filepath.Join(copy, faceID), fs.ModeDir)
	if err != nil {
		return fmt.Errorf("os.writefile: %w", err)
	}

	for _, image := range images {
		fl, err := os.Open(filepath.Join(image.Path, image.Filename))
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

		err = os.WriteFile(filepath.Join(copy, faceID, image.Filename), b, 0666)
		if err != nil {
			return fmt.Errorf("os.writefile: %w", err)
		}
	}

	return nil
}
