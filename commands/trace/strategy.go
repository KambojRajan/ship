package trace

import (
	"fmt"
	"sort"
)

type Strategy interface {
	Name() string
	Description() string
	ValidateArgs(args []string) error
	Execute(ctx *ExecContext) error
}

type ExecContext struct {
	RepoBasePath string
	Args         []string
}

var Registry = map[string]Strategy{}

func init() {
	for _, s := range []Strategy{
		CommitStrategy{},
		AddStrategy{},
		StatusStrategy{},
	} {
		Registry[s.Name()] = s
	}
}

func SupportedOps() []string {
	ops := make([]string, 0, len(Registry))
	for k := range Registry {
		ops = append(ops, k)
	}
	sort.Strings(ops)
	return ops
}

func Dispatch(operation string, ctx *ExecContext) error {
	s, ok := Registry[operation]
	if !ok {
		return fmt.Errorf("unknown operation %q — supported: %v", operation, SupportedOps())
	}
	if err := s.ValidateArgs(ctx.Args); err != nil {
		return err
	}
	return s.Execute(ctx)
}
