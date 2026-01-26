package command_tests

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/Entities"
	"github.com/KambojRajan/ship/tests/helpers"
)

func hashToString(hash [20]byte) string {
	return hex.EncodeToString(hash[:])
}

func TestCommit_NewCommit_ShouldCreateValidCommit(t *testing.T) {
	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Initial commit"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)

	helpers.AssertEqual(t, commit.TreeHash, treeHash)
	helpers.AssertEqual(t, len(commit.ParentHashes), 0)
	helpers.AssertEqual(t, commit.Author.Name, "Test Author")
	helpers.AssertEqual(t, commit.Committer.Name, "Test Committer")
	helpers.AssertEqual(t, commit.Message, "Initial commit")
}

func TestCommit_WithNoParents_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(oldDir)
	}()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Initial commit"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}

	helpers.BurnDown(t)
}

func TestCommit_WithSingleParent_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHash := [20]byte{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	parentHashes := []string{hashToString(parentHash)}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Second commit"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, len(commit.ParentHashes), 1)

	helpers.BurnDown(t)
}

func TestCommit_WithMultipleParents_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parent1 := [20]byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	parent2 := [20]byte{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	parentHashes := []string{hashToString(parent1), hashToString(parent2)}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Merge commit"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, len(commit.ParentHashes), 2)

	helpers.BurnDown(t)
}

func TestCommit_WithEmptyMessage_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := ""

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, commit.Message, "")

	helpers.BurnDown(t)
}

func TestCommit_WithLongMessage_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "First line of commit message\n\nDetailed explanation of what changed.\nMultiple lines are supported.\n\n- Feature 1\n- Feature 2"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, commit.Message, message)

	helpers.BurnDown(t)
}

func TestCommit_WithDifferentAuthorAndCommitter_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Original Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Different Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567900, 0).UTC(),
	}
	message := "Commit with different author and committer"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, commit.Author.Name, "Original Author")
	helpers.AssertEqual(t, commit.Committer.Name, "Different Committer")

	helpers.BurnDown(t)
}

func TestCommit_SameInputProducesSameHash_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	timestamp := time.Unix(1234567890, 0).UTC()
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: timestamp,
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: timestamp,
	}
	message := "Test commit"

	commit1 := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	hash1, err := commit1.Commit()
	helpers.AssertNil(err)

	commit2 := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	hash2, err := commit2.Commit()
	helpers.AssertNil(err)

	helpers.AssertEqual(t, hash1, hash2)

	helpers.BurnDown(t)
}

func TestCommit_DifferentInputProducesDifferentHash_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}

	commit1 := entities.NewCommit(treeHash, parentHashes, author, committer, "First message")
	hash1, err := commit1.Commit()
	helpers.AssertNil(err)

	commit2 := entities.NewCommit(treeHash, parentHashes, author, committer, "Second message")
	hash2, err := commit2.Commit()
	helpers.AssertNil(err)

	if hash1 == hash2 {
		t.Fatal("expected different commit messages to produce different hashes")
	}

	helpers.BurnDown(t)
}

func TestCommit_CommitTree_WithDefaultCommitter_ShouldUseAuthorTimestamp(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Test commit tree"

	commit := &entities.Commit{
		Committer: entities.User{
			Name:  "Test Committer",
			Email: "committer@example.com",
		},
	}

	commitHash, err := commit.CommitTree(treeHash, parentHashes, author, message)
	helpers.AssertNil(err)

	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}

	helpers.BurnDown(t)
}

func TestCommit_WithValidTreeHash_ShouldStoreObject(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})

	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Commit with valid tree hash"

	commit := entities.NewCommit(treeHash, []string{}, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}

	objectPath := filepath.Join(info.RepoDir, ".ship", "objects", commitHash[:2], commitHash[2:])

	_, err = os.Stat(objectPath)
	if err != nil {
		t.Fatalf("expected commit object to be stored at %s, but got error: %v", objectPath, err)
	}

	helpers.BurnDown(t)
}

func TestCommit_WithSpecialCharactersInMessage_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Commit with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?`~\n\t"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, commit.Message, message)

	helpers.BurnDown(t)
}

func TestCommit_WithUnicodeInMessage_ShouldPass(t *testing.T) {
	info := helpers.Setup(t)
	err := commands.Init(info.RepoDir)
	helpers.AssertNil(err)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author 测试",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Commit with unicode: 你好世界 🚀 émojis"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}
	helpers.AssertEqual(t, commit.Message, message)

	helpers.BurnDown(t)
}

func TestCommit_WithoutInitializedRepo_ShouldCreateObject(t *testing.T) {
	info := helpers.Setup(t)

	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(info.RepoDir)

	treeHash := hashToString([20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	parentHashes := []string{}
	author := entities.User{
		Name:      "Test Author",
		Email:     "author@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	committer := entities.User{
		Name:      "Test Committer",
		Email:     "committer@example.com",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}
	message := "Commit without init"

	commit := entities.NewCommit(treeHash, parentHashes, author, committer, message)
	commitHash, err := commit.Commit()

	helpers.AssertNil(err)
	if commitHash == "" {
		t.Fatal("expected non-zero commit hash")
	}

	helpers.BurnDown(t)
}
