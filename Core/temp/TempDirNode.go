package temp

import entities "github.com/KambojRajan/ship/Core/Entities"

type TempDirNode struct {
	Files map[string]entities.IndexEntry
	Dirs  map[string]*TempDirNode
	Hash  [20]byte
}
