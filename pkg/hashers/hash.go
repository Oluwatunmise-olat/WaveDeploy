package hashers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"io"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncodeIt(payload string) string {
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

func DecodeIt(b64 string) string {
	decodedBytes, _ := base64.StdEncoding.DecodeString(b64)
	return string(decodedBytes)
}

// https://earthly.dev/blog/cryptography-encryption-in-go/
func EncryptIt(text, secretKey string) (string, error) {
	aesBlock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cipheredText := gcmInstance.Seal(nonce, nonce, []byte(text), nil)

	return EncodeIt(string(cipheredText)), nil
}

func DecryptIt(ciphered, secretKey string) (string, error) {
	decoded := DecodeIt(ciphered)

	aesBlock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := decoded[:nonceSize], decoded[nonceSize:]
	originalText, err := gcmInstance.Open(nil, []byte(nonce), []byte(cipheredText), nil)
	if err != nil {
		return "", err
	}

	return string(originalText), nil
}
