package trace

import "github.com/KambojRajan/ship/commands"

type StatusStrategy struct{}

func (StatusStrategy) Name() string        { return "status" }
func (StatusStrategy) Description() string { return "trace the status pipeline (walk, hash comparison, diff)" }

func (StatusStrategy) ValidateArgs(_ []string) error { return nil }

func (StatusStrategy) Execute(ctx *ExecContext) error {
	path := ctx.RepoBasePath
	if len(ctx.Args) > 0 {
		path = ctx.Args[0]
	}
	commands.Status(path)
	return nil
}
