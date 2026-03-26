package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReadObjectContent reads a stored object by hash (relative to repoBasePath)
// and returns its raw content (header stripped).
func ReadObjectContent(repoBasePath, hash string) ([]byte, error) {
	hash = strings.TrimSpace(hash)
	if len(hash) < 3 {
		return nil, fmt.Errorf("invalid object hash: %q", hash)
	}

	objectPath := filepath.Join(repoBasePath, RootObjectDir, hash[:2], hash[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, err
	}

	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("zlib open %s: %w", hash, err)
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("zlib read %s: %w", hash, err)
	}

	nullIdx := bytes.IndexByte(decompressed, 0)
	if nullIdx == -1 {
		return nil, fmt.Errorf("invalid object %s: missing null separator", hash)
	}

	return decompressed[nullIdx+1:], nil
}
