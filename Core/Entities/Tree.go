package entities

import (
	"bytes"
	"sort"
	"strings"

	"github.com/KambojRajan/ship/Core/common"
	"github.com/KambojRajan/ship/Core/utils"
)

type Tree struct {
	Nodes []Node
}

func (t *Tree) WriteTree(index Index) ([20]byte, error) {
	root := buildTempTree(index)

	return writeTreeRecursive(root)
}

func writeTreeRecursive(root common.TempDirNode) ([20]byte, error) {
	for _, dir := range root.Dirs {
		childHash, err := writeTreeRecursive(*dir)
		if err != nil {
			return [20]byte{}, err
		}
		dir.Hash = childHash
	}

	treeBytes := serializeTree(root)

	return utils.HashObject(treeBytes, common.TREE, true)
}

func buildTempTree(index Index) common.TempDirNode {
	root := newTempDirNode()

	for _, entry := range index.Entries {
		insertEntryIntoTree(root, entry)
	}
	return *root
}

func insertEntryIntoTree(root *common.TempDirNode, entry common.IndexEntry) {
	pathChunks := strings.Split(entry.Path, "/")
	curr := root
	for i := 0; i < len(pathChunks)-1; i++ {
		dirName := pathChunks[i]
		if curr.Dirs[dirName] == nil {
			curr.Dirs[dirName] = newTempDirNode()
		}
		curr = curr.Dirs[dirName]
	}

	filename := pathChunks[len(pathChunks)-1]
	curr.Files[filename] = entry
}

func newTempDirNode() *common.TempDirNode {
	return &common.TempDirNode{
		Files: make(map[string]common.IndexEntry),
		Dirs:  make(map[string]*common.TempDirNode),
		Hash:  [20]byte{},
	}
}

func serializeTree(root common.TempDirNode) []byte {
	var entries []Node

	for name, file := range root.Files {
		entries = append(entries, Node{
			Mode: file.Mode,
			Name: name,
			Hash: file.Hash,
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
		buffer.Write(entry.Hash[:])
	}

	return buffer.Bytes()
}
