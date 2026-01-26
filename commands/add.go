package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	entities "github.com/KambojRajan/ship/core/Entities"
	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

func Add(paths ...string) error {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(paths)
	if err != nil {
		return err
	}

	if repoBasePath == "" {
		return fmt.Errorf("not a ship repository (or any of the parent directories)")
	}

	index, err := entities.LoadIndex(repoBasePath)
	if err != nil {
		return err
	}

	relPathMap, err := getRepoRelativePath(repoBasePath, paths...)
	if err != nil {
		return err
	}

	for path, relPath := range relPathMap {
		if err := processPath(repoBasePath, path, relPath, index); err != nil {
			return err
		}
	}

	return index.Save(repoBasePath)
}

func processPath(baseRepoPath, givenPath, relPath string, index *entities.Index) error {
	err := filepath.Walk(givenPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == utils.RootShipDir {
			return fs.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash, err := utils.HashObject(data, common.BLOB, false)
		if err != nil {
			return err
		}

		if exists := utils.ObjectExists(hash, baseRepoPath); exists {
			return nil
		}

		hash, err = utils.HashObject(data, common.BLOB, true)

		if err != nil {
			return err
		}

		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}

		index.AddIndex(common.IndexEntry{
			Path: path,
			Hash: hash,
			Mode: utils.GetMode(info),
		})

		return nil
	})
	return err
}

func getRepoRelativePath(repoBasePath string, paths ...string) (map[string]string, error) {
	relPaths := make(map[string]string)
	for _, path := range paths {
		if _, err := os.Lstat(path); err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("path does not exist: %s", path)
			}
			return nil, fmt.Errorf("error accessing path %s: %w", path, err)
		}

		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, fmt.Errorf("error resolving symlinks for %s: %w", path, err)
		}

		absPath, err := filepath.Abs(realPath)
		if err != nil {
			return nil, fmt.Errorf("error converting to absolute path %s: %w", realPath, err)
		}

		relPath, err := filepath.Rel(repoBasePath, absPath)
		if err != nil {
			return nil, fmt.Errorf("error computing relative path for %s: %w", absPath, err)
		}

		relPaths[absPath] = relPath
	}
	return relPaths, nil
}
