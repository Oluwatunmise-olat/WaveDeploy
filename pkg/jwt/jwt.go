package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"time"
)

func GenerateJWT(accountId uuid.UUID) (string, error) {
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

func GenerateGithubAppJWT(privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + (10 * 60),
		"iss": os.Getenv("GITHUB_APP_ID"),
		"alg": "RS256",
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unauthorized. Please login with `wave-deploy login`\n")
		}
		return []byte(os.Getenv("APP_KEY")), nil
	})

	if err != nil {
		return "", errors.New("Unauthorized. Please login with `wave-deploy login`\n")
	}

	if !token.Valid {
		return "", errors.New("Unauthorized. Please login with `wave-deploy login`\n")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Unauthorized. Please login with `wave-deploy login`\n")
	}
	// Type Assertion
	return claims["account_id"].(string), nil
}
