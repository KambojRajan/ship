package common

type IndexEntry struct {
	Path string
	Hash string
	Mode uint32
}

func (i IndexEntry) Equal(expected IndexEntry) bool {
	return i.Path == expected.Path && i.Hash == expected.Hash && i.Mode == expected.Mode
}

func (i IndexEntry) EqualWithoutMode(expected IndexEntry) bool {
	return i.Path == expected.Path && i.Hash == expected.Hash
}

type TempDirNode struct {
	Files map[string]IndexEntry
	Dirs  map[string]*TempDirNode
	Hash  string
}

func NewTempDirTree() *TempDirNode {
	return &TempDirNode{
		Files: make(map[string]IndexEntry),
		Dirs:  make(map[string]*TempDirNode),
	}
}
