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

	commits, err := entities.LoadCommits(repoBasePath)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Println("No commits found")
		return nil
	}

	for _, commit := range commits {
		// Print commit hash (first 7 chars like git)
		// For now we'll just print the full commit details since we only load HEAD
		fmt.Printf("commit %s\n", getCommitHash(repoBasePath, commit))
		fmt.Printf("Author: %s <%s> %d %s\n", commit.Author.Name, commit.Author.Email, commit.Author.Timestamp.Unix(), commit.Author.Timestamp.Format("-0700"))
		fmt.Printf("Committer: %s <%s> %d %s\n", commit.Committer.Name, commit.Committer.Email, commit.Committer.Timestamp.Unix(), commit.Committer.Timestamp.Format("-0700"))

		if len(commit.ParentHashes) > 0 {
			fmt.Printf("Parent: %s\n", commit.ParentHashes[0])
		}

		fmt.Printf("\n    %s\n\n", commit.Message)
	}

	return nil
}

// getCommitHash retrieves the current HEAD commit hash
func getCommitHash(repoPath string, commit *entities.Commit) string {
	head, err := entities.ResolveHead(repoPath)
	if err != nil {
		return "unknown"
	}
	return head.Hash
}

