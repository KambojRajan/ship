package utils

import (
	"encoding/json"
	"os"

	entities "github.com/KambojRajan/ship/Core/Entities"
)

func LoadIndex() (*entities.Index, error) {
	ok, err := ShipHasBeenInit()
	if !ok {
		return nil, err
	}

	bytes, err := os.ReadFile(BASE_INDEX_PATH)
	if err != nil {
		return nil, err
	}

	index := entities.Index{}
	err = json.Unmarshal(bytes, &index)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

func SaveIndex(index *entities.Index) error {
	b, _ := json.Marshal(index)
	return os.WriteFile(BASE_INDEX_PATH, b, 0644)
}
