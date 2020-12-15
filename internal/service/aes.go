package service

import (
	"crypto/aes"
	"encoding/hex"
	"github.com/systemz/wp-atrd-task/internal/config"
)

// https://golangdocs.com/aes-encryption-decryption-in-golang
func EncryptWithAes128(plaintext string) (err error, result string) {
	c, err := aes.NewCipher([]byte(config.AES_KEY))
	if err != nil {
		return
	}
	out := make([]byte, len(plaintext))
	// TODO check encrypted string length
	c.Encrypt(out, []byte(plaintext))
	// return as string
	result = hex.EncodeToString(out)
	return
}

func DecryptWithAes128(str string) (err error, plaintext string) {
	encryptedText, _ := hex.DecodeString(str)
	c, err := aes.NewCipher([]byte(config.AES_KEY))
	if err != nil {
		return
	}
	pt := make([]byte, len(str))
	c.Decrypt(pt, encryptedText)
	plaintext = string(pt[:])
	return
}
