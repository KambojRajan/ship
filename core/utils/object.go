package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KambojRajan/ship/core/common"
)

func HashObject(data []byte, objectType common.ObjectType, write bool) ([20]byte, error) {
	header := fmt.Sprintf("%s %d", objectType.String(), len(data))

	var store bytes.Buffer
	store.WriteString(header)
	store.WriteByte(0)
	store.Write(data)

	h := sha1.Sum(store.Bytes())
	hash := fmt.Sprintf("%x", h)

	if !write {
		return h, nil
	}

	folder := hash[0:2]
	file := hash[2:]

	objectDir := filepath.Join(RootObjectDir, folder)
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return h, err
	}

	objectPath := filepath.Join(objectDir, file)

	if _, err := os.Stat(objectPath); err == nil {
		return h, nil
	}
	out, err := os.Create(objectPath)
	if err != nil {
		return h, err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	zw := zlib.NewWriter(out)
	defer func(zw *zlib.Writer) {
		err := zw.Close()
		if err != nil {

		}
	}(zw)

	if _, err := zw.Write(store.Bytes()); err != nil {
		return h, err
	}

	return h, nil
}

func StoreObject(hash [20]byte, data []byte) error {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(".ship", "objects", hashStr)
	return os.WriteFile(file, data, 0644)
}

func ObjectExists(hash [20]byte, baseRepoPath string) bool {
	hashStr := fmt.Sprintf("%x", hash[:])
	file := filepath.Join(baseRepoPath, RootObjectDir, hashStr[:2], hashStr[2:])

	_, err := os.Stat(file)
	return err != nil
}

func GetMode(info os.FileInfo) uint32 {
	var mode uint32 = 100644
	if info.Mode()&0111 != 0 {
		mode = 100755
	}
	return mode
}
