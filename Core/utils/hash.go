package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashBytes(data []byte) string {
	h := sha1.Sum(data)

	return fmt.Sprintf("%x", h[:])
}
