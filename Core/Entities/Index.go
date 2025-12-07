package entities

type IndexEntry struct {
	Path     string
	BlobHash string
}

type Index struct {
	Entries map[string]IndexEntry
}
