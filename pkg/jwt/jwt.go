package jwt

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func GenerateJWT(accountId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(2 * time.Hour).Unix()
	claims["account_id"] = accountId

	tokenString, err := token.SignedString([]byte(os.Getenv("APP_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT() {}
