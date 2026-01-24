package entities

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/Core/common"
	"github.com/KambojRajan/ship/Core/utils"
)

type Index struct {
	Entries map[string]common.IndexEntry
}

func (i Index) Equal(expected *Index) bool {
	flag := len(i.Entries) == len(expected.Entries)
	for k, v := range i.Entries {
		flag = flag && expected.Entries[k].Equal(v)
	}
	return flag
}

func LoadIndex(currentPath string) (*Index, error) {
	ok, err := utils.ShipHasBeenInit(currentPath)
	if !ok {
		return nil, err
	}

	indexPath := filepath.Join(currentPath, utils.RootIndexPath)

	bytes, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	index := Index{}
	err = json.Unmarshal(bytes, &index)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

func SaveIndex(path string, index *Index) error {
	b, _ := json.Marshal(index)
	path = filepath.Join(path, utils.RootIndexPath)
	return os.WriteFile(path, b, 0644)
}
