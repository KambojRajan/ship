package utils

import (
	"os"
	"path/filepath"
)

func StoreObject(hash string, data []byte) error {
	file := filepath.Join(".ship", "objects", hash)
	return os.WriteFile(file, data, 0644)
}

func ObjectExists(hash string) bool {
	file := filepath.Join(".ship", "objects", hash)

	_, err := os.Stat(file)

	return err != nil
}
