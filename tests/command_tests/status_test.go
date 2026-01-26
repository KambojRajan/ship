package command_tests

import (
	"testing"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/tests/helpers"
)

func Test_BasicFunctionalityShouldWork(t *testing.T) {
	dir := t.TempDir()
	err := commands.Init(dir)
	helpers.AssertNil(err)

	files := []string{"test.txt", "test2.txt"}
	for _, file := range files {
		helpers.WriteFile(t, dir, file, []byte("content"))
	}
	err = commands.Add(dir)
	helpers.AssertNil(err)
	commands.Status(dir)
}
