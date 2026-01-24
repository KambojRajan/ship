package common

type IndexEntry struct {
	Path string
	Hash [20]byte
	Mode uint32
}

func (i IndexEntry) Equal(expected IndexEntry) bool {
	return i.Path == expected.Path && i.Hash == expected.Hash && i.Mode == expected.Mode
}

type TempDirNode struct {
	Files map[string]IndexEntry
	Dirs  map[string]*TempDirNode
	Hash  [20]byte
}

func NewTempDirTree() *TempDirNode {
	return &TempDirNode{
		Files: make(map[string]IndexEntry),
		Dirs:  make(map[string]*TempDirNode),
	}
}
