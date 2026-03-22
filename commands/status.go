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

func Status(currentDir string) {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(currentDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	repoBasePath, err = filepath.EvalSymlinks(repoBasePath)
	if err != nil {
		fmt.Printf("error resolving repository path: %v\n", err)
		return
	}

	index, err := entities.LoadIndex(repoBasePath)
	if err != nil {
		return
	}

	newIndex, err := recalculateIndex(repoBasePath, index)
	if err != nil {
		return
	}

	for path, indexEntry := range index.Entries {
		newIndexEntry, exists := newIndex.Entries[path]
		if !exists {
			fmt.Println(colorPath(path, utils.Unstaged, false))
			continue
		}
		if indexEntry.Equal(newIndexEntry) {
			fmt.Println(colorPath(path, utils.Staged, false))
		} else {
			fmt.Println(colorPath(path, utils.Unstaged, false))
		}
	}

	for path := range newIndex.Entries {
		if _, exists := index.Entries[path]; !exists {
			fmt.Println(colorPath(path, utils.Unstaged, false))
		}
	}
}

func recalculateIndex(repoRoot string, previous *entities.Index) (*entities.Index, error) {
	index := entities.NewIndex()

	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == utils.RootShipDir {
			return fs.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash, err := utils.HashObject(data, common.BLOB, false)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		relPath = filepath.Clean(relPath)

		var previousMode *uint32
		if previous != nil {
			if existing, ok := previous.Entries[relPath]; ok {
				mode := existing.Mode
				previousMode = &mode
			}
		}

		index.AddIndex(common.IndexEntry{
			Path: relPath,
			Hash: hash,
			Mode: utils.GetMode(info, previousMode),
		})

		return nil
	})

	return index, err
}

func colorPath(path string, state string, use256 bool) string {
	switch state {
	case utils.Staged:
		return utils.ShipBlue + path + utils.Reset
	case utils.Unstaged:
		if use256 {
			return utils.Sand256 + path + utils.Reset
		}
		return utils.Sand + path + utils.Reset
	default:
		return path
	}
}
