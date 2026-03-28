package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/utils"
)

func Diff(path string, cached bool) error {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
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

	if cached {
		return diffCached(repoBasePath, index)
	}
	return diffWorkingTree(repoBasePath, index)
}

func diffWorkingTree(repoBasePath string, index *entities.Index) error {
	for filePath, entry := range index.Entries {
		stagedBytes, err := utils.ReadObjectContent(repoBasePath, entry.Hash)
		if err != nil {
			return fmt.Errorf("read staged object for %s: %w", filePath, err)
		}

		absPath := filepath.Join(repoBasePath, filepath.FromSlash(filePath))
		workingBytes, err := os.ReadFile(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				workingBytes = []byte{}
			} else {
				return fmt.Errorf("read working tree file %s: %w", filePath, err)
			}
		}

		diff := utils.DiffLines(
			utils.SplitLines(string(stagedBytes)),
			utils.SplitLines(string(workingBytes)),
		)
		if out := utils.FormatUnifiedDiff(filePath, filePath, diff, 3); out != "" {
			fmt.Print(out)
		}
	}
	return nil
}

func diffCached(repoBasePath string, index *entities.Index) error {
	head, err := entities.ResolveHead(repoBasePath)
	if err != nil {
		return err
	}

	headFiles := map[string]string{}
	if strings.TrimSpace(head.Hash) != "" {
		headFiles, err = readCommitFiles(repoBasePath, strings.TrimSpace(head.Hash))
		if err != nil {
			return err
		}
	}

	seen := map[string]bool{}

	for filePath, entry := range index.Entries {
		seen[filePath] = true

		if headHash, ok := headFiles[filePath]; ok && headHash == entry.Hash {
			continue
		}

		stagedBytes, err := utils.ReadObjectContent(repoBasePath, entry.Hash)
		if err != nil {
			return fmt.Errorf("read staged object for %s: %w", filePath, err)
		}

		var oldLines []string
		oldPath := "/dev/null"
		if headHash, ok := headFiles[filePath]; ok {
			headBytes, err := utils.ReadObjectContent(repoBasePath, headHash)
			if err != nil {
				return err
			}
			oldLines = utils.SplitLines(string(headBytes))
			oldPath = filePath
		}

		diff := utils.DiffLines(oldLines, utils.SplitLines(string(stagedBytes)))
		if out := utils.FormatUnifiedDiff(oldPath, filePath, diff, 3); out != "" {
			fmt.Print(out)
		}
	}

	for filePath, headHash := range headFiles {
		if seen[filePath] {
			continue
		}
		headBytes, err := utils.ReadObjectContent(repoBasePath, headHash)
		if err != nil {
			return err
		}
		diff := utils.DiffLines(utils.SplitLines(string(headBytes)), nil)
		if out := utils.FormatUnifiedDiff(filePath, "/dev/null", diff, 3); out != "" {
			fmt.Print(out)
		}
	}

	return nil
}

func readCommitFiles(repoBasePath, commitHash string) (map[string]string, error) {
	commitContent, err := utils.ReadObjectContent(repoBasePath, commitHash)
	if err != nil {
		return nil, fmt.Errorf("read commit %s: %w", commitHash, err)
	}

	treeHash := ""
	for _, line := range strings.SplitN(string(commitContent), "\n\n", 2)[0:1] {
		for _, l := range strings.Split(line, "\n") {
			if strings.HasPrefix(l, "tree ") {
				treeHash = strings.TrimSpace(strings.TrimPrefix(l, "tree "))
				break
			}
		}
	}
	if treeHash == "" {
		return nil, fmt.Errorf("no tree hash in commit %s", commitHash)
	}

	files := make(map[string]string)
	if err := readTreeRecursive(repoBasePath, treeHash, "", files); err != nil {
		return nil, err
	}
	return files, nil
}

func readTreeRecursive(repoBasePath, treeHash, prefix string, files map[string]string) error {
	treeContent, err := utils.ReadObjectContent(repoBasePath, strings.TrimSpace(treeHash))
	if err != nil {
		return fmt.Errorf("read tree %s: %w", treeHash, err)
	}

	data := treeContent
	for len(data) > 0 {
		spaceIdx := bytes.IndexByte(data, ' ')
		if spaceIdx < 0 {
			break
		}
		mode := string(data[:spaceIdx])
		data = data[spaceIdx+1:]

		nullIdx := bytes.IndexByte(data, 0)
		if nullIdx < 0 {
			break
		}
		name := string(data[:nullIdx])
		data = data[nullIdx+1:]

		if len(data) < 40 {
			break
		}
		hash := string(data[:40])
		data = data[40:]

		entryPath := name
		if prefix != "" {
			entryPath = prefix + "/" + name
		}

		if mode == "40000" {
			if err := readTreeRecursive(repoBasePath, hash, entryPath, files); err != nil {
				return err
			}
		} else {
			files[entryPath] = hash
		}
	}
	return nil
}
