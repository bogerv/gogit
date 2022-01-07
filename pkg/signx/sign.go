package signx

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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

func Sha224(msg string) string {
	sum := sha256.Sum224([]byte(msg))
	return hex.EncodeToString(sum[:])
}

func Sha256(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return hex.EncodeToString(sum[:])
}

func Sha384(msg string) string {
	sum := sha512.Sum384([]byte(msg))
	return hex.EncodeToString(sum[:])
}

func Sha512(msg string) string {
	sum := sha512.Sum512([]byte(msg))
	return hex.EncodeToString(sum[:])
}
