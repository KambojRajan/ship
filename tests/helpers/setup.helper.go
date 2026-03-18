package helpers

import (
	"os"
	"testing"
)

type SetupInfo struct {
	RepoDir        string
	IncorrectDir   string
	RandomFilePath string
}

func Setup(t *testing.T) *SetupInfo {
	repoDir := t.TempDir()

	randomFilePath := t.TempDir() + "/random.txt"

	return &SetupInfo{
		RepoDir:        repoDir,
		IncorrectDir:   "Dummy",
		RandomFilePath: randomFilePath,
	}
}

func BurnDown(t *testing.T) {
	// Ensure we're in a safe directory before cleaning up
	homeDir, err := os.UserHomeDir()
	if err == nil {
		os.Chdir(homeDir)
	}

	err = os.RemoveAll(t.TempDir())
	if err != nil {
		return
	}
	t.Cleanup(func() {
		err := os.RemoveAll(t.TempDir())
		if err != nil {
			return
		}
	})
}
