package command_tests

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer
	callErr := fn()
	_ = writer.Close()
	os.Stdout = originalStdout

	output, readErr := io.ReadAll(reader)
	_ = reader.Close()
	if readErr != nil {
		t.Fatalf("failed to read stdout: %v", readErr)
	}

	return string(output), callErr
}

func captureStderr(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	originalStderr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	os.Stderr = writer
	callErr := fn()
	_ = writer.Close()
	os.Stderr = originalStderr

	output, readErr := io.ReadAll(reader)
	_ = reader.Close()
	if readErr != nil {
		t.Fatalf("failed to read stderr: %v", readErr)
	}

	return string(output), callErr
}

func TestPour_OnEmptyRepository_ShouldPrintNoCommitsFound(t *testing.T) {
	info := helpers.Setup(t)

	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	output, err := captureStdout(t, func() error {
		return commands.Pour(info.RepoDir)
	})
	helpers.AssertNil(err)

	if !strings.Contains(output, "No commits found") {
		t.Fatalf("expected empty repository message, got: %q", output)
	}

	helpers.BurnDown(t)
}

func TestPour_ShouldLoadAndPrintCommitHistory(t *testing.T) {
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()

	info := helpers.Setup(t)
	dir := info.RepoDir
	t.Setenv("SHIP_AUTHOR_NAME", "Test Author")
	t.Setenv("SHIP_AUTHOR_EMAIL", "author@example.com")
	t.Setenv("SHIP_AUTHOR_DATE", "2026-03-18T09:00:00Z")
	t.Setenv("SHIP_COMMITTER_NAME", "Test Committer")
	t.Setenv("SHIP_COMMITTER_EMAIL", "committer@example.com")
	t.Setenv("SHIP_COMMITTER_DATE", "2026-03-18T09:05:00Z")

	err := commands.Init(dir)
	helpers.AssertNil(err)
	helpers.WriteFile(t, dir, "file.txt", []byte("content v1"))
	err = commands.Add(dir)
	helpers.AssertNil(err)
	err = commands.Commit("first commit", dir)
	helpers.AssertNil(err)

	head, err := entities.ResolveHead(dir)
	helpers.AssertNil(err)
	firstHash := head.Hash

	helpers.WriteFile(t, dir, "file.txt", []byte("content v2"))
	helpers.WriteFile(t, dir, "second.txt", []byte("new file"))
	err = commands.Add(dir)
	helpers.AssertNil(err)
	err = commands.Commit("second commit", dir)
	helpers.AssertNil(err)

	head, err = entities.ResolveHead(dir)
	helpers.AssertNil(err)
	secondHash := head.Hash

	commits, err := entities.LoadCommits(dir)
	helpers.AssertNil(err)
	helpers.AssertEqual(t, 2, len(commits))
	helpers.AssertEqual(t, secondHash, commits[0].Hash)
	helpers.AssertEqual(t, firstHash, commits[1].Hash)
	helpers.AssertEqual(t, 1, len(commits[0].ParentHashes))
	helpers.AssertEqual(t, firstHash, commits[0].ParentHashes[0])

	output, err := captureStdout(t, func() error {
		return commands.Pour(dir)
	})
	helpers.AssertNil(err)

	if !strings.Contains(output, "commit "+secondHash) {
		t.Fatalf("expected pour output to contain head commit hash %s, got: %q", secondHash, output)
	}
	if !strings.Contains(output, "commit "+firstHash) {
		t.Fatalf("expected pour output to contain parent commit hash %s, got: %q", firstHash, output)
	}
	if !strings.Contains(output, "Parent: "+firstHash) {
		t.Fatalf("expected pour output to contain parent hash %s, got: %q", firstHash, output)
	}
	if !strings.Contains(output, "    second commit") || !strings.Contains(output, "    first commit") {
		t.Fatalf("expected pour output to contain both commit messages, got: %q", output)
	}
	if strings.Index(output, "commit "+secondHash) > strings.Index(output, "commit "+firstHash) {
		t.Fatalf("expected newer commit %s to appear before parent %s, got: %q", secondHash, firstHash, output)
	}

	stderrOutput, err := captureStderr(t, func() error {
		return commands.Pour(dir)
	})
	helpers.AssertNil(err)
	if strings.TrimSpace(stderrOutput) != "" {
		t.Fatalf("expected pour to produce no stderr output, got: %q", stderrOutput)
	}

	helpers.BurnDown(t)

}
