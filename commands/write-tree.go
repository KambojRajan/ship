package commands

import (
	"bytes"
	"sort"
	"strings"

	entities "github.com/KambojRajan/ship/Core/Entities"
	"github.com/KambojRajan/ship/Core/temp"
)

func WriteTree(index entities.Index) ([20]byte, error) {
	root := BuildTempTree(index)

	return writeTreeRecursive(root)
}

func writeTreeRecursive(root temp.TempDirNode) ([20]byte, error) {
	for _, dir := range root.Dirs {
		childHash, err := writeTreeRecursive(*dir)
		if err != nil {
			return [20]byte{}, err
		}
		dir.Hash = childHash
	}

	treeBytes := serializeTree(root)

	return HashObject(treeBytes, temp.TREE, true)
}

func BuildTempTree(index entities.Index) temp.TempDirNode {
	root := newTempDirNode()

	for _, entry := range index.Entries {
		insertEntryIntoTree(root, entry)
	}
	return *root
}

func insertEntryIntoTree(root *temp.TempDirNode, entry entities.IndexEntry) {
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

func newTempDirNode() *temp.TempDirNode {
	return &temp.TempDirNode{
		Files: make(map[string]entities.IndexEntry),
		Dirs:  make(map[string]*temp.TempDirNode),
		Hash:  [20]byte{},
	}
}

func serializeTree(root temp.TempDirNode) []byte {
	var entries []entities.Node

	for name, file := range root.Files {
		entries = append(entries, entities.Node{
			Mode: file.Mode,
			Name: name,
			Hash: file.Hash,
		})
	}

	for name, dir := range root.Dirs {
		entries = append(entries, entities.Node{
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
		buffer.WriteString(entry.ModeString() + " ")
		buffer.WriteString(entry.Name)
		buffer.WriteByte(0)
		buffer.Write(entry.Hash[:])
	}

	return buffer.Bytes()
}
