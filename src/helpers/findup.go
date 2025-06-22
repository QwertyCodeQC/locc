package helpers

import (
	"os"
	"path/filepath"
)

func FindUp(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil // File found
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached the root directory
		}
		dir = parent
	}
	return "", nil
}