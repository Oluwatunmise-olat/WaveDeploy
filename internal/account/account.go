package account

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/files"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/jwt"
	"github.com/google/uuid"
	"os"
)

func AuthenticateAccount(email, password string) error {
	accountRepository := respository.AccountsRepository{}
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
	authTokenPath := getCredentialsPath()

	if err := files.WriteToFileWithOverride("", authTokenPath); err != nil {
		panic(err)
	}
}

func GetAuthTokenDetails() (string, error) {
	authTokenPath := getCredentialsPath()
	content, err := files.GetFileContent(authTokenPath)
	if err != nil {
		return "", err
	}

	accountId, err := jwt.VerifyJWT(string(content))
	if err != nil {
		return "", err
	}

	accountRepository := respository.AccountsRepository{}
	account, err := accountRepository.GetAccountById(accountId)

	if account == nil && err != nil {
		return "", errors.New("Unauthorized. Please create an account with `wave-deploy create-account`")
	}
	return accountId, err
}

func getCredentialsPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/.wave-deploy/credentials", homeDir)
}

func saveAuthenticationToken(accountId uuid.UUID) {
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

func GetAccountInstallationId(accountId string) (string, error) {
	githubAppRepository := respository.GithubAppsRepository{}

	installationId, err := githubAppRepository.GetGithubAppInstallationIdByAccountId(accountId)
	if err != nil {
		return "", errors.New("Account not connected to github. Please connect with `wave-deploy connect-github`")
	}

	return installationId, nil

}
