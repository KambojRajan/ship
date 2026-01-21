package entities

type IndexEntry struct {
	Path string
	Hash [20]byte
	Mode uint32
}

type Index struct {
	Entries map[string]IndexEntry
}
