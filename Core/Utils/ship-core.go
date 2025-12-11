package utils

import (
	"fmt"
	"os"
)

func ShipHasBeenInit() (bool, error) {
	currentRoot, err := os.Getwd()
	if err != nil {
		return false, err
	}

	path := fmt.Sprintf("%s/%s", currentRoot, BASE_SHIP_DIR)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
