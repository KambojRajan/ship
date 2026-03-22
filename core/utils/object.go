package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/KambojRajan/ship/core/common"
)

func HashObject(data []byte, objectType common.ObjectType, write bool) (string, error) {
	header := fmt.Sprintf("%s %d", objectType.String(), len(data))

	var store bytes.Buffer
	store.WriteString(header)
	store.WriteByte(0)
	store.Write(data)

	h := sha1.Sum(store.Bytes())
	hash := fmt.Sprintf("%x", h)

	if !write {
		return hash, nil
	}

	folder := hash[0:2]
	file := hash[2:]

	objectDir := filepath.Join(RootObjectDir, folder)
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return hash, err
	}

	objectPath := filepath.Join(objectDir, file)

	if _, err := os.Stat(objectPath); err == nil {
		return hash, nil
	}

	tempFile, err := os.CreateTemp(objectDir, file+".tmp-*")
	if err != nil {
		return hash, err
	}
	tempPath := tempFile.Name()
	cleanupTemp := true
	defer func() {
		if cleanupTemp {
			_ = os.Remove(tempPath)
		}
	}()

	zw := zlib.NewWriter(tempFile)

	if _, err := zw.Write(store.Bytes()); err != nil {
		_ = zw.Close()
		_ = tempFile.Close()
		return hash, err
	}

	if err := zw.Close(); err != nil {
		_ = tempFile.Close()
		return hash, err
	}

	if err := tempFile.Close(); err != nil {
		return hash, err
	}

	if err := os.Rename(tempPath, objectPath); err != nil {
		if _, statErr := os.Stat(objectPath); statErr == nil {
			return hash, nil
		}
		return hash, err
	}

	cleanupTemp = false

	return hash, nil
}

func StoreObject(hash string, data []byte) error {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(".ship", "objects", hashStr)
	return os.WriteFile(file, data, 0644)
}

func ObjectExists(hash string, baseRepoPath string) bool {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(baseRepoPath, RootObjectDir, hashStr[:2], hashStr[2:])

	_, err := os.Stat(file)
	return err == nil
}

func IsExecutableMode(mode uint32) bool {
	return mode == GitFileModeExecutable
}

func GetMode(info os.FileInfo, previousMode *uint32) uint32 {
	if info.Mode()&0o111 != 0 {
		return GitFileModeExecutable
	}

	if runtime.GOOS == "windows" && previousMode != nil && IsExecutableMode(*previousMode) {
		return GitFileModeExecutable
	}

	return GitFileModeRegular
}
