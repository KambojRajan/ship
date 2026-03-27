package entities

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Index struct {
	Entries map[string]common.IndexEntry
}

func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]common.IndexEntry),
	}
}

func (i *Index) AddIndex(entry common.IndexEntry) {
	if i.Entries == nil {
		i.Entries = make(map[string]common.IndexEntry)
	}
	i.Entries[entry.Path] = entry
}

func (i *Index) Equal(expected *Index) bool {
	if i == nil || expected == nil {
		return i == expected
	}
	if len(i.Entries) != len(expected.Entries) {
		return false
	}
	for k, v := range i.Entries {
		expectedEntry, exists := expected.Entries[k]
		if !exists || !v.Equal(expectedEntry) {
			return false
		}
	}
	return true
}

func LoadIndex(reportPath string) (*Index, error) {
	indexPath := filepath.Join(reportPath, utils.RootIndexPath)

	bytes, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	index := NewIndex()
	err = json.Unmarshal(bytes, index)
	if err != nil {
		return nil, err
	}

	if index.Entries == nil {
		index.Entries = make(map[string]common.IndexEntry)
	}

	return index, nil
}

func (i *Index) Save(path string) error {
	if i == nil {
		return os.ErrInvalid
	}

	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	indexPath := filepath.Join(path, utils.RootIndexPath)
	return os.WriteFile(indexPath, b, 0644)
}
