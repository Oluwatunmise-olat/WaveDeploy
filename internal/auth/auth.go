package auth

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/files"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/jwt"
	"os"
)

func AuthenticateAccount(email, password string) error {
	accountRepository := respository.InitializeAccountsRepository()
	account, err := accountRepository.GetAccountByEmail(email)

	if err != nil {
		return errors.New("account not found\n")
	}

	validPassword := hashers.ComparePasswordHash(password, account.Password)
	if !validPassword {
		return errors.New("incorrect credentials\n")
	}
	saveAuthenticationToken(account.Id)

	return nil
}

func LogoutAccount() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	authTokenPath := fmt.Sprintf("%s/.wave-deploy/credentials", homeDir)

	err = files.WriteToFileWithOverride("", authTokenPath)
	if err != nil {
		panic(err)
	}
}

func saveAuthenticationToken(accountId string) {
	token, err := jwt.GenerateJWT(accountId)
	if err != nil {
		panic(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	authTokenPath := fmt.Sprintf("%s/.wave-deploy/credentials", homeDir)

	err = files.WriteToFileWithOverride(token, authTokenPath)
	if err != nil {
		panic(err)
	}
}

func GetAuthenticationToken() (string, error) {
	return "", nil
}
