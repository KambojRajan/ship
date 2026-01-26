package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashBytes(data []byte) [20]byte {
	return sha1.Sum(data)
}

func HashString(data string) [20]byte {
	return HashBytes([]byte(data))
}

func GetHashString(hash [20]byte) string {
	return fmt.Sprintf("%x", hash[:])
}
