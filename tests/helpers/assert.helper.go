package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	entities "github.com/KambojRajan/ship/core/Entities"
)

func AssertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("expected %s to exist", path)
		}
		t.Fatalf("error checking %s: %v", path, err)
	}
}

func AssertNil(err error) {
	if err != nil {
		panic("expected error")
	}
}

func AssertNotNil(err error) {
	if err == nil {
		panic(err)
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if expected == actual {
		t.Fatalf("expected %v to not be equal to %v", expected, actual)
	}
}

func WriteFile(t *testing.T, dir, filename string, content []byte) {
	t.Helper()
	fullPath := filepath.Join(dir, filename)

	// Create parent directories if needed
	parentDir := filepath.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatalf("failed to create parent directory %s: %v", parentDir, err)
	}

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", fullPath, err)
	}
}

func WriteDir(t *testing.T, dir, filename string) {
	t.Helper()
	fullPath := filepath.Join(dir, filename)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", fullPath, err)
	}
}

func DeleteFile(t *testing.T, dir, filename string) {
	t.Helper()
	fullPath := filepath.Join(dir, filename)
	if err := os.Remove(fullPath); err != nil {
		t.Fatalf("failed to delete file %s: %v", fullPath, err)
	}
}

func AssertFileInIndex(t *testing.T, repoDir, filename string) {
	t.Helper()

	index, err := entities.LoadIndex(repoDir)
	if err != nil {
		t.Fatalf("failed to load index: %v", err)
	}

	relPath := filepath.Join(repoDir, filename)

	absRelPath, _ := filepath.EvalSymlinks(relPath)
	if _, ok := index.Entries[absRelPath]; !ok {
		t.Fatalf("file %s not found in index", filename)
	}
}

func AppendToFile(path string, content string) error {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("opening file for append failed: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic(err.Error())
		}
	}(file)

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("writing to file failed: %w", err)
	}
	return nil
}

func AssertEqualIndex(t *testing.T, expected, actual *entities.Index) {
	t.Helper()
	flag := actual.Equal(expected)
	AssertEqual(t, flag, true)
}

func AssertNotEqualIndex(t *testing.T, expected, actual *entities.Index) {
	t.Helper()
	flag := actual.Equal(expected)
	AssertEqual(t, !flag, true)
}
