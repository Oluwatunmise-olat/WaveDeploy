package hashers

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncodeStringToBase64(payload string) string {
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

func DecodeBase64ToString(b64 string) string {
	decodedBytes, _ := base64.StdEncoding.DecodeString(b64)
	return string(decodedBytes)
}
