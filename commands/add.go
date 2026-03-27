package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/trace"
	"github.com/KambojRajan/ship/core/utils"
)

type addTarget struct {
	absPath string
	relPath string
	isDir   bool
}

func Add(paths ...string) error {
	end := trace.Step("ShipHasBeenInitRecursive")
	repoBasePath, err := utils.ShipHasBeenInitRecursive(paths...)
	end(err)
	if err != nil {
		return err
	}

	if repoBasePath == "" {
		return fmt.Errorf("not a ship repository (or any of the parent directories)")
	}

	end = trace.Step("EvalSymlinks(repo)")
	repoBasePath, err = filepath.EvalSymlinks(repoBasePath)
	end(err)
	if err != nil {
		return fmt.Errorf("error resolving repository path: %w", err)
	}

	// Change to repo directory to ensure HashObject writes to correct .ship location
	end = trace.Step("Getwd")
	oldCwd, err := os.Getwd()
	end(err)
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	end = trace.Step("Chdir(repo)")
	if err := os.Chdir(repoBasePath); err != nil {
		end(err)
		return fmt.Errorf("error changing to repository directory: %w", err)
	}
	end(nil)
	defer func() {
		os.Chdir(oldCwd)
	}()

	end = trace.Step("LoadIndex")
	index, err := entities.LoadIndex(repoBasePath)
	end(err)
	if err != nil {
		return err
	}

	end = trace.Step("getRepoRelativePath")
	targets, err := getRepoRelativePath(repoBasePath, paths...)
	end(err)
	if err != nil {
		return err
	}

	existingFiles := make(map[string]bool)

	end = trace.Step("processTargets")
	for _, target := range targets {
		if err := processPath(repoBasePath, target.absPath, index, existingFiles); err != nil {
			end(err)
			return err
		}
	}
	end(nil)

	// Workload context — visible in trace summary
	trace.Meta("files_staged", fmt.Sprintf("%d", len(existingFiles)))
	trace.Meta("targets_resolved", fmt.Sprintf("%d", len(targets)))

	cleanupIndexEntries(index, existingFiles, targets)

	end = trace.Step("SaveIndex")
	err = index.Save(repoBasePath)
	end(err)
	return err
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

		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash, err := utils.HashObject(data, common.BLOB, true)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(baseRepoPath, resolvedPath)
		if err != nil {
			return err
		}
		relPath = filepath.Clean(relPath)

		var previousMode *uint32
		if existing, ok := index.Entries[relPath]; ok {
			mode := existing.Mode
			previousMode = &mode
		}

		existingFiles[relPath] = true
		index.AddIndex(common.IndexEntry{
			Path: relPath,
			Hash: hash,
			Mode: utils.GetMode(info, previousMode),
		})

		return nil
	})
	return err
}

func cleanupIndexEntries(index *entities.Index, existingFiles map[string]bool, targets []addTarget) {
	for entry := range index.Entries {
		if existingFiles[entry] {
			continue
		}

		for _, target := range targets {
			if pathWithinScope(entry, target) {
				delete(index.Entries, entry)
				break
			}
		}
	}
}

func pathWithinScope(entry string, target addTarget) bool {
	entry = filepath.Clean(entry)
	scope := filepath.Clean(target.relPath)

	if scope == "." {
		return true
	}

	if !target.isDir {
		return entry == scope
	}

	return entry == scope || strings.HasPrefix(entry, scope+string(os.PathSeparator))
}

func getRepoRelativePath(repoBasePath string, paths ...string) ([]addTarget, error) {
	targets := make([]addTarget, 0, len(paths))
	seen := make(map[string]bool)

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
		absPath = filepath.Clean(absPath)

		relPath, err := filepath.Rel(repoBasePath, absPath)
		if err != nil {
			return nil, fmt.Errorf("error computing relative path for %s: %w", absPath, err)
		}
		relPath = filepath.Clean(relPath)

		info, err := os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("error accessing path %s: %w", absPath, err)
		}

		if seen[absPath] {
			continue
		}

		targets = append(targets, addTarget{
			absPath: absPath,
			relPath: relPath,
			isDir:   info.IsDir(),
		})
		seen[absPath] = true
	}

	return targets, nil
}
