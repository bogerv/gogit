package aesutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/go-errors/errors"
	"io"
	"sync"
)

type AesKey struct {
	key       []byte
	syncMutex sync.Mutex
}

var commonKey = &AesKey{
	key: []byte("lOvInG_@*CoInETH20190423"),
}

func SetAesKey(key string) {
	commonKey.syncMutex.Lock()
	commonKey.key = []byte(key)
	commonKey.syncMutex.Unlock()
}

// AesCFBEncrypt encrypt plain text
func AesCFBEncrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(commonKey.key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	return hex.EncodeToString(cipherText), nil
}

// AesCFBDecrypt decrypt text
func AesCFBDecrypt(decryptText string) (string, error) {
	cipherText, err := hex.DecodeString(decryptText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(commonKey.key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
