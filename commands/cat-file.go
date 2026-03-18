package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/KambojRajan/ship/core/utils"
)

func CatFile(hash, flag string) (string, error) {
	folder := hash[0:2]
	file := hash[2:]

	path := filepath.Join(utils.RootObjectDir, folder, file)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	defer func(zr io.ReadCloser) {
		err := zr.Close()
		if err != nil {
			_ = err
		}
	}(zr)

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return "", err
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) != 2 {
		return "", fmt.Errorf(utils.ErrInvalidObjectFormat)
	}

	header := string(parts[0])

	headerFields := strings.Split(header, " ")
	objectType := headerFields[0]
	objectSize := headerFields[1]
	content := string(parts[1])
	if len(headerFields) != 2 {
		return "", fmt.Errorf(utils.ErrInvalidObjectHeader)
	}
	switch flag {
	case utils.CatFileFormatPretty:
		return fmt.Sprintf("%s %s %s", objectType, objectSize, content), nil
	case utils.CatFileContentSize:
		return objectSize, nil
	case utils.CatFileFormatTree:
		return content, nil
	default:
		return string(decompressed), nil
	}
}
