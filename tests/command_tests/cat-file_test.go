package command_tests

import (
	"path/filepath"
	"testing"

	"github.com/KambojRajan/ship/commands"
	entities "github.com/KambojRajan/ship/core/Entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func TestCatFile_blobHash_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	helpers.WriteFile(t, dir, "test.txt", []byte("hello"))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	filePath, err := filepath.EvalSymlinks(filepath.Join(dir, "test.txt"))
	helpers.AssertNil(err)
	fileIndex := index.Entries[filePath]

	body, err := commands.CatFile(fileIndex.Hash)
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "hello", string(body))
}
