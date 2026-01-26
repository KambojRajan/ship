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

func CatFile(args ...string) (string, error) {
	hash := args[0]
	flag := utils.CatFileDefaultFormat
	if len(args) > 1 {
		flag = args[1]
	}
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
			_ = fmt.Errorf(err.Error())
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
	objectSize := headerFields[1]
	if len(headerFields) != 2 {
		return "", fmt.Errorf(utils.ErrInvalidObjectHeader)
	}
	switch flag {
	case utils.CatFileDefaultFormat:
		return string(parts[1]), nil
	case utils.CatFileFormatBlob:
		return string(decompressed), nil
	case utils.CatFileContentSize:
		return objectSize, nil
	case utils.CatFileFormatTree:
		return string(parts[1]), nil
	}
	return "", nil
}
