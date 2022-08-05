package data

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type ID struct {
	Name string
	Picture
}
type Picture struct {
	PictureID string
	Path      string
	File      *fs.DirEntry
}
type Matches struct {
	Picture
	Ids []ID
}

var picturesRaw []Picture
var picturesIDs []ID

// LoadIds receives a path and returns a slice of type ID that contains all of the
// ids and files in the specified path
func LoadIDs(path string) []ID {
	xdir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	copyIDFiles(xdir, path)
	return picturesIDs
}

// LoadImages receives a path and returns a slice of type Picture that contains
// all of the files in the specified path
func LoadImages(path string) []Picture {
	xdir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	copyImageFiles(xdir, path)
	return picturesRaw
}

func copyIDFiles(xdir []fs.DirEntry, path string) {
	for _, dir := range xdir {
		if dir.IsDir() {
			path = filepath.Join(path, dir.Name())
			xdir2, err := os.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			copyIDFiles(xdir2, path)
			break
		}
		name := filepath.Base(path)
		picturesIDs = append(picturesIDs, ID{name, Picture{dir.Name(), filepath.Join(path, dir.Name()), &dir}})
	}
}

func copyImageFiles(xdir []fs.DirEntry, path string) {
	for _, dir := range xdir {
		if dir.IsDir() {
			path = filepath.Join(path, dir.Name())
			xdir2, err := os.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			copyImageFiles(xdir2, path)
			break
		}
		picturesRaw = append(picturesRaw, Picture{dir.Name(), filepath.Join(path, dir.Name()), &dir})
	}
}
