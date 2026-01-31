package helpers

import (
	"os"
	"testing"

	"github.com/KambojRajan/ship/core/utils"
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
	err := os.RemoveAll(t.TempDir())
	if err != nil {
		return
	}
	t.Cleanup(func() {
		err := os.RemoveAll(t.TempDir())
		if err != nil {
			return
		}
	})
	err = os.RemoveAll(utils.RootShipDir)
	if err != nil {
		return
	}
}
