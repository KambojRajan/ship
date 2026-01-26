package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/utils"
)

func Init(path string) error {
	if path == "" || path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return fmt.Errorf(utils.ErrFailedToGetWorkingDir, err)
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf(utils.ErrFailedToAccessPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf(utils.ErrPathNotDirectory, path)
	}

	shipDir := filepath.Join(path, utils.RootShipDir)
	objectDir := filepath.Join(path, utils.RootObjectDir)
	refsHeadsDir := filepath.Join(path, utils.RootShipDir, "refs", "heads")
	refsTagsDir := filepath.Join(path, utils.RootShipDir, "refs", "tags")
	indexPath := filepath.Join(path, utils.RootIndexPath)
	headPath := filepath.Join(path, utils.RootHEADPath)

	if _, err := os.Stat(shipDir); err == nil {
		fmt.Println("Reinitialized existing Ship repository")
		return nil
	}

	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return fmt.Errorf(utils.ErrFailedToCreateObjectsDir, err)
	}

	if err := os.MkdirAll(refsHeadsDir, 0755); err != nil {
		return fmt.Errorf(utils.ErrFailedToCreateRefsHeadsDir, err)
	}

	if err := os.MkdirAll(refsTagsDir, 0755); err != nil {
		return fmt.Errorf(utils.ErrFailedToCreateRefsTagsDir, err)
	}

	if err := os.WriteFile(indexPath, []byte(`{"entries":{}}`), 0644); err != nil {
		return fmt.Errorf(utils.ErrFailedToCreateIndexFile, err)
	}

	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf(utils.ErrFailedToCreateHEADFile, err)
	}

	fmt.Printf("Initialized empty Ship repository in %s\n", shipDir)
	return nil
}
