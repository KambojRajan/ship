package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashBytes(data []byte) string {
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash[:])
}

func HashString(data string) string {
	return HashBytes([]byte(data))
}
