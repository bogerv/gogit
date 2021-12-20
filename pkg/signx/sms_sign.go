package signx

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func MD5(msg string) string {
	sum := md5.Sum([]byte(msg))
	return hex.EncodeToString(sum[:])
}

func Sha1(msg string) string {
	sum := sha1.Sum([]byte(msg))
	return hex.EncodeToString(sum[:])
}
