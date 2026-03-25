package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/trace"
	"github.com/KambojRajan/ship/core/utils"
)

func Status(currentDir string) {
	end := trace.Step("ShipHasBeenInitRecursive")
	repoBasePath, err := utils.ShipHasBeenInitRecursive(currentDir)
	end(err)
	if err != nil {
		fmt.Println(err)
		return
	}

	end = trace.Step("EvalSymlinks(repo)")
	repoBasePath, err = filepath.EvalSymlinks(repoBasePath)
	end(err)
	if err != nil {
		fmt.Printf("error resolving repository path: %v\n", err)
		return
	}

	end = trace.Step("LoadIndex")
	index, err := entities.LoadIndex(repoBasePath)
	end(err)
	if err != nil {
		return
	}

	end = trace.Step("recalculateIndex")
	newIndex, err := recalculateIndex(repoBasePath, index)
	end(err)
	if err != nil {
		return
	}

	// Workload context — visible in trace summary
	trace.Meta("tracked_files", fmt.Sprintf("%d", len(index.Entries)))
	trace.Meta("disk_files", fmt.Sprintf("%d", len(newIndex.Entries)))

	end = trace.Step("compareEntries")
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
	end(nil)
}

func recalculateIndex(repoRoot string, previous *entities.Index) (*entities.Index, error) {
	index := entities.NewIndex()

	end := trace.Step("WalkDir(repo)")
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
	end(err)

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
