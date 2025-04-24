package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateParentDir(path string) error {
	directory := filepath.Dir(path)

	// The directory does not exist, create it
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create socket directory: %s", err.Error())
	}

	return nil
}
