package commands

import (
	"os"
	"path/filepath"

	entities "github.com/KambojRajan/ship/Core/Entities"
	"github.com/KambojRajan/ship/Core/common"
	"github.com/KambojRajan/ship/Core/utils"
)

func Add(path string) error {
	_, err := utils.ShipHasBeenInit(path)
	if err != nil {
		return err
	}
	index, err := entities.LoadIndex(path)
	if err != nil {
		return err
	}

	currentFiles := make(map[string]bool)

	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == utils.RootShipDir {
			return filepath.SkipDir
		}

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		hash := utils.HashBytes(data)
		if !utils.ObjectExists(hash) {
			if err := utils.StoreObject(hash, data); err != nil {
				return err
			}
		}

		rel, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}

		mode := utils.GetMode(info)
		index.Entries[rel] = common.IndexEntry{
			Path: rel,
			Hash: hash,
			Mode: mode,
		}

		currentFiles[rel] = true

		return nil
	})

	if err != nil {
		return err
	}

	for path := range index.Entries {
		if !currentFiles[path] {
			delete(index.Entries, path)
		}
	}

	if err = entities.SaveIndex(path, index); err != nil {
		return err
	}

	return nil
}
