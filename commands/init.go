package commands

import (
	"fmt"
	"os"
	"path/filepath"

	utils "github.com/KambojRajan/ship/Core/Utils"
)

func Init(path string) error {
	var root string
	var err error

	if path == "" {
		root, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		root, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	}

	shipDir := filepath.Join(root, ".ship")

	if _, err = os.ReadFile(shipDir); !os.IsNotExist(err) {
		return fmt.Errorf(utils.REPO_ALREADY_EXISTS, shipDir)
	}

	if err := os.MkdirAll(filepath.Join(shipDir, "objects"), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(shipDir, "index"), []byte(`{"entries":{}}`), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(shipDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return err
	}

	return nil
}
