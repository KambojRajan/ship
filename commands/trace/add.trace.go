package trace

import (
	"fmt"

	"github.com/KambojRajan/ship/commands"
)

type AddStrategy struct{}

func (AddStrategy) Name() string        { return "add" }
func (AddStrategy) Description() string { return "trace the staging pipeline (walk, hash, index update)" }

func (AddStrategy) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ship trace add <path> [paths...]")
	}
	return nil
}

func (AddStrategy) Execute(ctx *ExecContext) error {
	return commands.Add(ctx.Args...)
}
