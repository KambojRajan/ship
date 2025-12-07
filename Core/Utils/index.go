package utils

import (
	"encoding/json"
	"os"

	entities "github.com/KambojRajan/ship/Core/Entities"
)

func LoadIndex() (*entities.Index, error) {
	path := ".ship/index"

	// This check if this path exists or not.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
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
	return os.WriteFile(".ship/index", b, 0644)
}
