package utility

import (
	"os"
	"path/filepath"
)

func CreateParentDir(path string, perm os.FileMode) error {
	directory := filepath.Dir(path)

	// The directory does not exist, create it
	if err := os.MkdirAll(directory, perm); err != nil {
		return err
	}

	return nil
}
