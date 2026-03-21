package commands

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Purge(repoBasePath string) error {
	objectDir := filepath.Join(repoBasePath, ".ship", "objects")

	return filepath.WalkDir(objectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.Contains(d.Name(), ".tmp-") {
			info, _ := d.Info()
			if time.Since(info.ModTime()) > 30*time.Minute {
				return os.Remove(path)
			}
		}
		return nil
	})
}
