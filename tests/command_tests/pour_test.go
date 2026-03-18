package command_tests

import (
	"os"
	"testing"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func Test_pour(t *testing.T) {
	// Save and restore working directory to handle any side effects from previous tests
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()

	info := helpers.Setup(t)
	dir := info.RepoDir

	err := commands.Init(dir)
	helpers.AssertNil(err)
	file := "file.txt"
	helpers.WriteFile(t, dir, file, []byte("content"))
	err = commands.Add(dir)
	helpers.AssertNil(err)
	err = commands.Commit("dummy", dir)
	helpers.AssertNil(err)

	_, err = entities.LoadCommits(dir)
	helpers.AssertNil(err)

}
