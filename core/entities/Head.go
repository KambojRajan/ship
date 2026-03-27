package entities

import (
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/utils"
)

type Head struct {
	Ref  string
	Hash string
}

func ResolveHead(repoDir string) (Head, error) {
	headPath := filepath.Join(repoDir, utils.RootHEADPath)
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return Head{}, err
	}

	ref := string(headContent)

	refPath := filepath.Join(repoDir, utils.RootShipDir, utils.MainHeadPath)
	hashBytes, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Head{
				Ref:  ref,
				Hash: "",
			}, nil
		}
		return Head{}, err
	}

	return Head{
		Ref:  ref,
		Hash: string(hashBytes),
	}, nil
}

func (h Head) Write(repoDir string) error {
	headPath := filepath.Join(repoDir, utils.RootHEADPath)
	return os.WriteFile(headPath, []byte(h.Ref), 0644)
}

func (h Head) UpdateRef(repoDir string, commitHash string) error {
	refPath := filepath.Join(repoDir, utils.RootShipDir, utils.MainHeadPath)
	return os.WriteFile(refPath, []byte(commitHash), 0644)
}
