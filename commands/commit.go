package commands

import (
	"fmt"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/utils"
)

func Commit(message string, path string) error {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
	if err != nil {
		return err
	}
	index, err := entities.LoadIndex(repoBasePath)
	if err != nil {
		return err
	}
	if len(index.Entries) == 0 {
		return fmt.Errorf("nothing to commit")
	}
	tree := entities.Tree{}
	treeHash, err := tree.WriteTree(index)
	if err != nil {
		return err
	}

	var parents []string
	head, err := entities.ResolveHead(repoBasePath)
	if err == nil && head.Hash != "" {
		parents = append(parents, head.Hash)
	}

	author := entities.NewUserFromEnv(false)
	committer := entities.NewUserFromEnv(true)

	commit := entities.NewCommit(treeHash, parents, author, committer, message)
	commitHash, err := commit.Commit()
	if err != nil {
		return err
	}

	if err := head.UpdateRef(repoBasePath, commitHash); err != nil {
		return fmt.Errorf("failed to update ref: %w", err)
	}

	return nil
}
