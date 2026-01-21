package entities

import (
	"encoding/json"
	"os"

	"github.com/KambojRajan/ship/Core/common"
	"github.com/KambojRajan/ship/Core/utils"
)

type Index struct {
	Entries map[string]common.IndexEntry
}

func LoadIndex() (*Index, error) {
	ok, err := utils.ShipHasBeenInit()
	if !ok {
		return nil, err
	}

	bytes, err := os.ReadFile(utils.BaseIndexPath)
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

func SaveIndex(index *Index) error {
	b, _ := json.Marshal(index)
	return os.WriteFile(utils.BaseIndexPath, b, 0644)
}
