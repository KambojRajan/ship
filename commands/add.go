package commands

import (
	"os"
	"path/filepath"

	entities "github.com/KambojRajan/ship/Core/Entities"
	"github.com/KambojRajan/ship/Core/utils"
)

func Add(path string) error {
	index, err := utils.LoadIndex()
	if err != nil {
		return err
	}

	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".ship" {
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

		index.Entries[rel] = entities.IndexEntry{
			Path:     rel,
			BlobHash: hash,
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err = utils.SaveIndex(index); err != nil {
		return err
	}

	return nil
}
