package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/KambojRajan/ship/Core/utils"
)

func CateFile(hash string) error {
	folder := hash[0:2]
	file := hash[2:]

	if shipInitDone, err := utils.ShipHasBeenInit(); !shipInitDone {
		return err
	}

	path := fmt.Sprintf(utils.RootObjectDir+"/%v/%v", folder, file)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer func(zr io.ReadCloser) {
		err := zr.Close()
		if err != nil {

		}
	}(zr)

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return err
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) != 2 {
		return fmt.Errorf(utils.ErrInvalidObjectFormat)
	}

	header := string(parts[0])
	body := parts[1]

	headerFields := strings.Split(header, " ")
	if len(headerFields) != 2 {
		return fmt.Errorf(utils.ErrInvalidObjectHeader)
	}

	objectType := headerFields[0]

	switch objectType {
	case utils.BLOB:
		fmt.Print(string(body))
	case utils.COMMIT:
		fmt.Print(string(body))
	case utils.TREE:
		return fmt.Errorf(utils.ErrTreeNotImplemented)
	default:
		return fmt.Errorf(utils.ErrUnknownType, objectType)
	}

	return nil
}
