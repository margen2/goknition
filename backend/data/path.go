package data

import (
	"os"
	"path/filepath"
)

func GetCwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func ListFolders(path string) ([]string, error) {
	path = filepath.Clean(path)
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var folders []string
	for _, v := range dirs {
		if v.IsDir() {
			folders = append(folders, v.Name())
		}
	}

	return folders, nil
}

func GoBack(path string) string {
	return filepath.Dir(path)
}
