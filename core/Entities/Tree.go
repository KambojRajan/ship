package entities

import (
	"bytes"
	"log"
	"sort"
	"strings"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Tree struct {
	Nodes []Node
}

func (t *Tree) WriteTree(index *Index) (string, error) {
	log.Println("[WriteTree] Starting tree write operation")
	log.Printf("[WriteTree] Index has %d entries", len(index.Entries))

	root := buildTempDirTree(index)
	log.Println("[WriteTree] Temporary directory tree built successfully")

	hash, err := writeTreeRecursive(root)
	if err != nil {
		log.Printf("[WriteTree] Error writing tree: %v", err)
		return "", err
	}

	log.Printf("[WriteTree] Tree written successfully with hash: %s", hash)
	return hash, nil
}

func writeTreeRecursive(root *common.TempDirNode) (string, error) {
	log.Printf("[writeTreeRecursive] Processing node with %d directories and %d files", len(root.Dirs), len(root.Files))

	for dirName, dir := range root.Dirs {
		log.Printf("[writeTreeRecursive] Processing child directory: %s", dirName)
		childHash, err := writeTreeRecursive(dir)
		if err != nil {
			log.Printf("[writeTreeRecursive] Error processing directory %s: %v", dirName, err)
			return "", err
		}
		dir.Hash = childHash
		log.Printf("[writeTreeRecursive] Directory %s hashed as: %s", dirName, childHash)
	}

	treeBytes := serializeTree(root)
	log.Printf("[writeTreeRecursive] Serialized tree to %d bytes", len(treeBytes))

	hash, err := utils.HashObject(treeBytes, common.TREE, true)
	if err != nil {
		log.Printf("[writeTreeRecursive] Error hashing tree object: %v", err)
		return "", err
	}

	log.Printf("[writeTreeRecursive] Node hashed successfully: %s", hash)
	return hash, nil
}

func serializeTree(root *common.TempDirNode) []byte {
	log.Println("[serializeTree] Starting tree serialization")
	var entries []Node

	for name, files := range root.Files {
		log.Printf("[serializeTree] Adding file entry: %s (mode: %o, hash: %s)", name, files.Mode, files.Hash)
		entries = append(entries, Node{
			Mode: files.Mode,
			Hash: files.Hash,
			Name: name,
		})
	}

	for name, dir := range root.Dirs {
		log.Printf("[serializeTree] Adding directory entry: %s (mode: 040000, hash: %s)", name, dir.Hash)
		entries = append(entries, Node{
			Mode: 040000,
			Name: name,
			Hash: dir.Hash,
		})
	}

	log.Printf("[serializeTree] Total entries before sorting: %d", len(entries))
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	log.Println("[serializeTree] Entries sorted alphabetically")

	var buffer bytes.Buffer

	for _, entry := range entries {
		buffer.WriteString(entry.modeString() + " ")
		buffer.WriteString(entry.Name)
		buffer.WriteByte(0)
		buffer.Write([]byte(entry.Hash))
	}

	log.Printf("[serializeTree] Serialization complete, buffer size: %d bytes", buffer.Len())
	return buffer.Bytes()
}

func buildTempDirTree(index *Index) *common.TempDirNode {
	log.Println("[buildTempDirTree] Building temporary directory tree from index")
	root := &common.TempDirNode{
		Dirs:  make(map[string]*common.TempDirNode),
		Files: make(map[string]common.IndexEntry),
	}

	if index.Entries == nil || len(index.Entries) == 0 {
		log.Println("[buildTempDirTree] Index is empty, returning empty tree")
		return root
	}

	log.Printf("[buildTempDirTree] Processing %d index entries", len(index.Entries))
	for _, entry := range index.Entries {
		pathParts := strings.Split(entry.Path, utils.SEPARATOR)
		log.Printf("[buildTempDirTree] Processing entry: %s (depth: %d)", entry.Path, len(pathParts))
		curr := root

		for i := 0; i < len(pathParts)-1; i++ {
			dirName := pathParts[i]
			if curr.Dirs[dirName] == nil {
				log.Printf("[buildTempDirTree] Creating new directory node: %s", dirName)
				curr.Dirs[dirName] = &common.TempDirNode{
					Dirs:  make(map[string]*common.TempDirNode),
					Files: make(map[string]common.IndexEntry),
				}
			}
			curr = curr.Dirs[dirName]
		}
		filename := pathParts[len(pathParts)-1]
		log.Printf("[buildTempDirTree] Adding file: %s (hash: %s)", filename, entry.Hash)
		curr.Files[filename] = entry
	}

	log.Println("[buildTempDirTree] Temporary directory tree build complete")
	return root
}
