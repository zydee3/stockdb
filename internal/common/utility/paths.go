package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateParentDir(path string) error {
	directory := filepath.Dir(path)

	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {

		// The directory does not exist, create it
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create socket directory: %s", err.Error())
		}
	}

	return nil
}
