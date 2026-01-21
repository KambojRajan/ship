package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func StoreObject(hash [20]byte, data []byte) error {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(".ship", "objects", hashStr)
	return os.WriteFile(file, data, 0644)
}

func ObjectExists(hash [20]byte) bool {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(".ship", "objects", hashStr)

	_, err := os.Stat(file)
	return err != nil
}

func GetMode(info os.FileInfo) uint32 {
	var mode uint32 = 100644
	if info.Mode()&0111 != 0 {
		mode = 100755
	}
	return mode
}
