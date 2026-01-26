package entities

import (
	"bytes"
	"sort"
	"strings"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Tree struct {
	Nodes []Node
}

func (t *Tree) WriteTree(index *Index) (string, error) {
	root := buildTempDirTree(index)

	return writeTreeRecursive(root)
}

func writeTreeRecursive(root *common.TempDirNode) (string, error) {
	for _, dir := range root.Dirs {
		childHash, err := writeTreeRecursive(dir)
		if err != nil {
			return "", err
		}
		dir.Hash = childHash
	}
	treeBytes := serializeTree(root)

	return utils.HashObject(treeBytes, common.TREE, true)
}

func serializeTree(root *common.TempDirNode) []byte {
	var entries []Node

	for name, files := range root.Files {
		entries = append(entries, Node{
			Mode: files.Mode,
			Hash: files.Hash,
			Name: name,
		})
	}

	for name, dir := range root.Dirs {
		entries = append(entries, Node{
			Mode: 040000,
			Name: name,
			Hash: dir.Hash,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	var buffer bytes.Buffer

	for _, entry := range entries {
		buffer.WriteString(entry.modeString() + " ")
		buffer.WriteString(entry.Name)
		buffer.WriteByte(0)
		buffer.Write([]byte(entry.Hash))
	}
	return buffer.Bytes()
}

func buildTempDirTree(index *Index) *common.TempDirNode {
	root := &common.TempDirNode{
		Dirs:  make(map[string]*common.TempDirNode),
		Files: make(map[string]common.IndexEntry),
	}

	if index.Entries == nil || len(index.Entries) == 0 {
		return root
	}

	for _, entry := range index.Entries {
		pathParts := strings.Split(entry.Path, utils.SEPARATOR)
		curr := root

		for i := 0; i < len(pathParts)-1; i++ {
			dirName := pathParts[i]
			if curr.Dirs[dirName] == nil {
				curr.Dirs[dirName] = &common.TempDirNode{}
			}
			curr = curr.Dirs[dirName]
		}
		filename := pathParts[len(pathParts)-1]
		curr.Files[filename] = entry
	}
	return root
}
