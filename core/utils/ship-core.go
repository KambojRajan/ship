package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ShipHasBeenInitRecursive(paths []string) (repoPath string, err error) {
	for _, path := range paths {
		pathToCheck, err := searchBaseRepo(path)
		if err != nil {
			return "", err
		}
		if pathToCheck != "" {
			return pathToCheck, nil
		}
	}
	return "", nil
}

func searchBaseRepo(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	splits := strings.Split(absPath, string(os.PathSeparator))

	for i := len(splits); i > 0; i-- {
		currentPath := strings.Join(splits[:i], string(os.PathSeparator))
		if currentPath == "" {
			currentPath = string(os.PathSeparator)
		}

		if shipHasBeenInit(currentPath) {
			return currentPath, nil
		}
	}
	return "", nil
}

func shipHasBeenInit(path string) bool {
	shipDir := filepath.Join(path, RootShipDir)
	info, err := os.Stat(shipDir)
	return err == nil && info.IsDir()
}
