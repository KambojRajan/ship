package utils

import (
	"crypto/sha1"
)

func HashBytes(data []byte) [20]byte {
	return sha1.Sum(data)
}
