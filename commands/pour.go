package commands

import (
	"fmt"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/utils"
)

func Pour(path string) error {
	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
	if err != nil {
		return err
	}
	if repoBasePath == "" {
		return fmt.Errorf("not a ship repository (or any of the parent directories)")
	}

	commits, err := entities.LoadCommits(repoBasePath)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Println("No commits found")
		return nil
	}

	for _, commit := range commits {
		fmt.Printf("commit %s\n", commit.Hash)
		fmt.Printf("Author: %s <%s> %d %s\n", commit.Author.Name, commit.Author.Email, commit.Author.Timestamp.Unix(), commit.Author.Timestamp.Format("-0700"))
		fmt.Printf("Committer: %s <%s> %d %s\n", commit.Committer.Name, commit.Committer.Email, commit.Committer.Timestamp.Unix(), commit.Committer.Timestamp.Format("-0700"))

		if len(commit.ParentHashes) > 0 {
			fmt.Printf("Parent: %s\n", commit.ParentHashes[0])
		}

		fmt.Printf("\n    %s\n\n", commit.Message)
	}

	return nil
}
