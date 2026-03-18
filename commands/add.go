package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

func Add(paths ...string) error {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(paths...)
	if err != nil {
		return err
	}

	if repoBasePath == "" {
		return fmt.Errorf("not a ship repository (or any of the parent directories)")
	}

	repoBasePath, err = filepath.EvalSymlinks(repoBasePath)
	if err != nil {
		return fmt.Errorf("error resolving repository path: %w", err)
	}

	// Change to repo directory to ensure HashObject writes to correct .ship location
	oldCwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	if err := os.Chdir(repoBasePath); err != nil {
		return fmt.Errorf("error changing to repository directory: %w", err)
	}
	defer func() {
		os.Chdir(oldCwd)
	}()

	index, err := entities.LoadIndex(repoBasePath)
	if err != nil {
		return err
	}

	relPathMap, err := getRepoRelativePath(repoBasePath, paths...)
	if err != nil {
		return err
	}

	existingFiles := make(map[string]bool)

	for path := range relPathMap {
		if err := processPath(repoBasePath, path, index, existingFiles); err != nil {
			return err
		}
	}

	for entry := range index.Entries {
		if !existingFiles[entry] {
			delete(index.Entries, entry)
		}
	}

	return index.Save(repoBasePath)
}

func processPath(baseRepoPath, givenPath string, index *entities.Index, existingFiles map[string]bool) error {
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

		relPath, err := filepath.Rel(baseRepoPath, path)
		if err != nil {
			return err
		}

		existingFiles[relPath] = true
		index.AddIndex(common.IndexEntry{
			Path: relPath,
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
