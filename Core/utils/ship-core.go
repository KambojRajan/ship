package utils

import (
	"os"
	"path/filepath"
)

func ShipHasBeenInit(currentPath string) (bool, error) {
	path := filepath.Join(currentPath, RootShipDir)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return false, err
	}

	return info.IsDir(), nil
}
