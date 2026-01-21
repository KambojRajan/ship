package common

type IndexEntry struct {
	Path string
	Hash [20]byte
	Mode uint32
}

type TempDirNode struct {
	Files map[string]IndexEntry
	Dirs  map[string]*TempDirNode
	Hash  [20]byte
}
