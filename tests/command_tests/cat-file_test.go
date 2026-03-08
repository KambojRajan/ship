package command_tests

import (
	"os"
	"strings"
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
	helpers.AssertNil(err)
	fileIndex := index.Entries["test.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "blob 5 hello", string(body))
}

func TestCatFile_WithSizeFlag_ShouldReturnSize(t *testing.T) {
	dir := t.TempDir()
	helpers.WriteFile(t, dir, "test.txt", []byte("hello world"))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)
	fileIndex := index.Entries["test.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-s")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "11", string(body))
}

func TestCatFile_MultipleBlobsWithDifferentContent_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	helpers.WriteFile(t, dir, "file1.txt", []byte("content1"))
	helpers.WriteFile(t, dir, "file2.txt", []byte("different content"))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)

	helpers.AssertNil(err)
	fileIndex1 := index.Entries["file1.txt"]

	helpers.AssertNil(err)
	fileIndex2 := index.Entries["file2.txt"]

	body1, err := commands.CatFile(fileIndex1.Hash, "-p")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "blob 8 content1", string(body1))

	body2, err := commands.CatFile(fileIndex2.Hash, "-p")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "blob 17 different content", string(body2))
}

func TestCatFile_EmptyFile_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	helpers.WriteFile(t, dir, "empty.txt", []byte(""))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)
	fileIndex := index.Entries["empty.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "blob 0 ", string(body))
}

func TestCatFile_LargeBlobContent_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	largeContent := strings.Repeat("x", 1000)
	helpers.WriteFile(t, dir, "large.txt", []byte(largeContent))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)

	helpers.AssertNil(err)
	fileIndex := index.Entries["large.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	expected := "blob 1000 " + largeContent
	helpers.AssertEqual(t, expected, string(body))
}

func TestCatFile_NonExistentHash_ShouldFail(t *testing.T) {
	dir := t.TempDir()
	err := commands.Init(dir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(dir)

	invalidHash := "1234567890abcdef1234567890abcdef12345678"
	_, err = commands.CatFile(invalidHash, "-p")
	helpers.AssertNotNil(err)
}

func TestCatFile_TreeObject_ShouldPass(t *testing.T) {
	// to be implemented
}

func TestCatFile_CommitObject_ShouldPass(t *testing.T) {
	// to be implemented
}

func TestCatFile_BlobWithNewlines_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	content := "line1\nline2\nline3"
	helpers.WriteFile(t, dir, "multiline.txt", []byte(content))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)
	fileIndex := index.Entries["multiline.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	expected := "blob 17 " + content
	helpers.AssertEqual(t, expected, string(body))
}

func TestCatFile_BlobWithSpecialCharacters_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	content := "Special chars: !@#$%^&*()_+-={}[]|\\:\";<>?,./"
	helpers.WriteFile(t, dir, "special.txt", []byte(content))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)
	helpers.AssertNil(err)
	fileIndex := index.Entries["special.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	if !strings.Contains(body, content) {
		t.Fatalf("expected body to contain special characters, got: %s", body)
	}
}

func TestCatFile_NestedDirectoryBlob_ShouldPass(t *testing.T) {
	dir := t.TempDir()
	helpers.WriteDir(t, dir, "nested/deep/path")
	helpers.WriteFile(t, dir, "nested/deep/path/file.txt", []byte("nested content"))
	err := commands.Init(dir)
	helpers.AssertNil(err)
	err = commands.Add(dir)
	helpers.AssertNil(err)
	index, err := entities.LoadIndex(dir)

	helpers.AssertNil(err)
	fileIndex := index.Entries["nested/deep/path/file.txt"]

	body, err := commands.CatFile(fileIndex.Hash, "-p")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, "blob 14 nested content", string(body))
}
