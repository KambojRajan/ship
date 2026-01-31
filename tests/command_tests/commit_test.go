package command_tests

import (
	"testing"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func Test_CommitCommand(t *testing.T) {
	dir := t.TempDir()
	err := commands.Init(dir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, dir, "test.txt", []byte("hello"))
	err = commands.Add(dir)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)
	helpers.AssertNotNil(index)

	err = commands.Commit("dummy", dir)
	helpers.AssertNil(err)
}
