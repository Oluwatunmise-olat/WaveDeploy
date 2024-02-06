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
		"iat": time.Now().Unix() - 60,        // issued at time, 60 seconds in the past
		"exp": time.Now().Unix() + (10 * 60), // expiration time (10 minutes)
		"iss": os.Getenv("GITHUB_APP_ID"),    // GitHub App's identifier
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("APP_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth token tampered\n")
		}

		return []byte(os.Getenv("APP_KEY")), nil
	})

	if err != nil {
		return "", errors.New("auth token expired. please login again\n")
	}

	if !token.Valid {
		return "", errors.New("auth token expired. please login again\n")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("unable to extract claims\n")
	}
	// Type Assertion
	return claims["account_id"].(string), nil
}
