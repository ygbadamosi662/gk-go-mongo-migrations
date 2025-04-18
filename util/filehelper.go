package util

import (
	"os"
	"path/filepath"
)

// CreateDirIfNotExist checks if a directory exists, and creates it if it doesn't.
func CreateDirIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

// FileExists checks if a file exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// JoinPaths safely joins parts of a file path.
func JoinPaths(parts ...string) string {
	return filepath.Join(parts...)
}
