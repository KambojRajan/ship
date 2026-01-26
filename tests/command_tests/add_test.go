package command_tests

import (
	"os"
	"testing"

	"github.com/KambojRajan/ship/commands"
	entities "github.com/KambojRajan/ship/core/Entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func TestAdd_ToEmptyInitDir_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	helpers.BurnDown(t)
}

func TestAdd_ToCurrentInitDir_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)

	originalDir, err := os.Getwd()
	helpers.AssertNil(err)

	err = os.Chdir(info.RepoDir)
	helpers.AssertNil(err)
	defer func() {
		err := os.Chdir(originalDir)
		if err != nil {
			return
		}
	}()

	err = commands.Init(".")
	helpers.AssertNil(err)

	err = commands.Add(".")
	helpers.AssertNil(err)

	helpers.BurnDown(t)
}

func TestAddTwice_ToInitDir_ShouldPass(t *testing.T) {

	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)
	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	helpers.BurnDown(t)
}

func TestAdd_ToUninitializedDir_ShouldFail(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Add(info.RepoDir)
	helpers.AssertNotNil(err)

	helpers.BurnDown(t)
}

func TestAdd_WithNoFiles_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	index, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)
	helpers.AssertEqual(t, len(index.Entries), 0)

	helpers.BurnDown(t)
}

func TestAdd_WithNestedDirectories_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)

	dirs := []string{"dir1", "dir2/subdir1", "dir2/subdir2", "a/b/c/d/e/f/g/h"}
	for _, name := range dirs {
		helpers.WriteDir(t, info.RepoDir, name)
	}

	files := []string{
		"file1.txt",
		"dir2/file2.txt",
		"dir2/subdir1/file3.txt",
		"dir1/file4.txt",
		"a/b/c/d/e/f/g/h/deep.txt",
	}
	for _, name := range files {
		helpers.WriteFile(t, info.RepoDir, name, []byte("content for "+name))
	}

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	for _, name := range files {
		helpers.AssertFileInIndex(t, info.RepoDir, name)
	}

	helpers.BurnDown(t)
}

func TestAdd_WithMixedContent_ShouldAddAll(t *testing.T) {
	info := helpers.Setup(t)

	helpers.WriteFile(t, info.RepoDir, "text.txt", []byte("text content"))
	helpers.WriteFile(t, info.RepoDir, "binary.bin", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) // PNG header
	helpers.WriteFile(t, info.RepoDir, "empty.txt", []byte(""))
	helpers.WriteFile(t, info.RepoDir, "large.dat", make([]byte, 1024*1024)) // 1MB file
	helpers.WriteFile(t, info.RepoDir, "script.sh", []byte("#!/bin/bash\necho test"))
	helpers.WriteFile(t, info.RepoDir, "file1.txt", []byte("file 1"))
	helpers.WriteFile(t, info.RepoDir, "file2.txt", []byte("file 2"))

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	files := []string{"text.txt", "binary.bin", "empty.txt", "large.dat", "script.sh", "file1.txt", "file2.txt"}
	for _, name := range files {
		helpers.AssertFileInIndex(t, info.RepoDir, name)
	}

	helpers.BurnDown(t)
}

func TestAdd_WithSpecialFilenames_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	filenames := []string{
		"file-with-dash.txt",
		"file_with_underscore.txt",
		"file.multiple.dots.txt",
		"file with spaces.txt",
		"another file.txt",
		".gitignore",
		".env",
		".editorconfig",
		"😀emoji.txt",
		"this_is_a_very_long_filename_that_tests_the_limits_of_what_can_be_handled_by_the_system_abcdefghijklmnopqrstuvwxyz_0123456789_repeated_many_times.txt",
	}

	for _, name := range filenames {
		helpers.WriteFile(t, info.RepoDir, name, []byte("content"))
	}

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	for _, name := range filenames {
		helpers.AssertFileInIndex(t, info.RepoDir, name)
	}

	helpers.BurnDown(t)
}

func TestAdd_WithModifiedFile_ShouldUpdateIndex(t *testing.T) {
	info := helpers.Setup(t)

	dirs := []string{"dir1", "dir2/subdir1"}
	for _, name := range dirs {
		helpers.WriteDir(t, info.RepoDir, name)
	}

	files := []string{"file1.txt", "dir2/file2.txt", "dir1/file3.txt"}
	for _, name := range files {
		helpers.WriteFile(t, info.RepoDir, name, []byte("content for "+name))
	}

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)
	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	initialIndex, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)
	// Use relative path instead of absolute path
	relPath := "file1.txt"
	initialHash := initialIndex.Entries[relPath].Hash

	helpers.WriteFile(t, info.RepoDir, "file1.txt", []byte("modified content"))

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	finalIndex, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)
	finalHash := finalIndex.Entries[relPath].Hash

	helpers.AssertNotEqualIndex(t, initialIndex, finalIndex)
	helpers.AssertNotEqual(t, initialHash, finalHash)

	helpers.BurnDown(t)
}

func TestAdd_WithDeletedFile_ShouldHandle(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, info.RepoDir, "file1.txt", []byte("content for file1.txt"))
	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)
	helpers.AssertFileInIndex(t, info.RepoDir, "file1.txt")

	helpers.DeleteFile(t, info.RepoDir, "file1.txt")

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	index, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)
	_, exists := index.Entries["file1.txt"]
	helpers.AssertEqual(t, false, exists)

	helpers.BurnDown(t)
}

func TestAdd_WithIdenticalContent_ShouldReuseObject(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	content := []byte("identical content")
	helpers.WriteFile(t, info.RepoDir, "file1.txt", content)
	helpers.WriteFile(t, info.RepoDir, "file2.txt", content)

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	index, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)

	helpers.AssertFileInIndex(t, info.RepoDir, "file1.txt")
	helpers.AssertFileInIndex(t, info.RepoDir, "file2.txt")

	hash1 := index.Entries["file1.txt"].Hash
	hash2 := index.Entries["file2.txt"].Hash

	if hash1 != hash2 {
		t.Fatalf("expected identical content to have same hash, got %v and %v", hash1, hash2)
	}

	helpers.BurnDown(t)
}

func TestAdd_SkipsShipDirectory_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, info.RepoDir, "file1.txt", []byte("content"))

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	index, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)

	for path := range index.Entries {
		if len(path) >= 5 && path[:5] == ".ship" {
			t.Fatalf("expected .ship directory to be skipped, but found entry: %s", path)
		}
	}

	helpers.BurnDown(t)
}

func TestAdd_AfterIndexCorruption_ShouldHandle(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, info.RepoDir, "file1.txt", []byte("content"))
	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, info.RepoDir, ".ship/index", []byte("corrupted data"))

	err = commands.Add(info.RepoDir)
	helpers.AssertNotNil(err)

	helpers.BurnDown(t)
}

func TestAdd_WithRegularFileMode_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	helpers.WriteFile(t, info.RepoDir, "regular.txt", []byte("regular file"))

	err = commands.Add(info.RepoDir)
	helpers.AssertNil(err)

	index, err := entities.LoadIndex(info.RepoDir)
	helpers.AssertNil(err)

	helpers.AssertFileInIndex(t, info.RepoDir, "regular.txt")
	helpers.AssertNil(err)
	helpers.AssertEqual(t, uint32(100644), index.Entries["regular.txt"].Mode)

	helpers.BurnDown(t)
}
