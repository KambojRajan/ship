package commands

import (
	"fmt"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/core/trace"
	"github.com/KambojRajan/ship/core/utils"
)

func Commit(message string, path string) error {
	end := trace.Step("ShipHasBeenInitRecursive")
	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
	end(err)
	if err != nil {
		return err
	}
	end = trace.Step("LoadIndex")
	index, err := entities.LoadIndex(repoBasePath)
	end(err)
	if err != nil {
		return err
	}
	if len(index.Entries) == 0 {
		return fmt.Errorf("nothing to commit")
	}

	// Workload context — visible in trace summary
	trace.Meta("staged_entries", fmt.Sprintf("%d", len(index.Entries)))
	tree := entities.Tree{}
	end = trace.Step("WriteTree")
	treeHash, err := tree.WriteTree(index)
	end(err)
	if err != nil {
		return err
	}

	var parents []string
	end = trace.Step("ResolveHead")
	head, err := entities.ResolveHead(repoBasePath)
	end(err)
	if err == nil && head.Hash != "" {
		parents = append(parents, head.Hash)
	}
	trace.Meta("parent_commits", fmt.Sprintf("%d", len(parents)))

	end = trace.Step("NewUserFromEnv(author)")
	author := entities.NewUserFromEnv(false)
	end(nil)
	end = trace.Step("NewUserFromEnv(committer)")
	committer := entities.NewUserFromEnv(true)
	end(nil)

	end = trace.Step("NewCommit")
	commit := entities.NewCommit(treeHash, parents, author, committer, message, repoBasePath)
	end(nil)
	end = trace.Step("CommitObject")
	commitHash, err := commit.Commit()
	end(err)
	if err != nil {
		return err
	}

	end = trace.Step("UpdateRef")
	if err := head.UpdateRef(repoBasePath, commitHash); err != nil {
		end(err)
		return fmt.Errorf("failed to update ref: %w", err)
	}
	end(nil)

	return nil
}
