package trace

import (
	"fmt"

	"github.com/KambojRajan/ship/commands"
)

type CommitStrategy struct{}

func (CommitStrategy) Name() string        { return "commit" }
func (CommitStrategy) Description() string { return "trace the commit pipeline (tree, object write, ref update)" }

func (CommitStrategy) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ship trace commit <message>")
	}
	return nil
}

func (CommitStrategy) Execute(ctx *ExecContext) error {
	return commands.Commit(ctx.Args[0], ctx.RepoBasePath)
}
