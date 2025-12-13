package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/Core/utils"
)

func HashObject(path string, write bool) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	header := fmt.Sprintf("blob %d", len(data))

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

	objectDir := filepath.Join(utils.BASE_OBJECT_DIR, folder)
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return "", err
	}

	objectPath := filepath.Join(objectDir, file)

	if _, err := os.Stat(objectPath); err == nil {
		return hash, nil
	}
	out, err := os.Create(objectPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	zw := zlib.NewWriter(out)
	if _, err := zw.Write(store.Bytes()); err != nil {
		return "", err
	}
	zw.Close()

	return hash, nil
}
