package github

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"log"
	"os"
	"strings"
	"time"
)

func IsAccountConnectedAlreadyToGithub(accountId string) bool {
	githubAppRepository := respository.GithubAppsRepository{}
	githubApp, _ := githubAppRepository.GetGithubAppByAccountId(accountId)
	return githubApp != nil
}

func GetConnectToGithubUrl(accountId string) string {
	base64String := hashers.EncodeIt(accountId)
	return github.ConnectAppToGithub(base64String)
}

func GetConnectGithubRepositoryUrl() string {
	return os.Getenv("GITHUB_APP_PUBLIC_LINK") + "/installations/select_target"
}

// TODO: Validate webhook source
// Handle case of app disconnection
// https://docs.github.com/en/webhooks/using-webhooks/validating-webhook-deliveries
func CreateGithubAppIfNotExists(payload structs.GithubOauthWebhook) {
	accountRepository := respository.AccountsRepository{}
	account, err := accountRepository.GetAccountById(hashers.DecodeIt(payload.State))

	if account == nil && err != nil {
		return
	}

	githubAppRepository := respository.GithubAppsRepository{}
	githubApp, err := githubAppRepository.GetGithubAppByInstallationId(payload.InstallationId)

	if githubApp != nil && err == nil {
		return
	}

	newGithubApp := models.GithubApps{
		InstallationId: payload.InstallationId,
		AccountId:      hashers.DecodeIt(payload.State),
		Code:           payload.Code,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Id:             random.GetUUID(),
	}

	if err = githubAppRepository.CreateGithubApp(newGithubApp); err != nil {
		log.Fatalln(err)
	}
}

func GetAccountConnectedRepositories(accountId string) ([]structs.GithubAInstallationRepositories, error) {
	githubAppRepository := respository.GithubAppsRepository{}
	githubApp, _ := githubAppRepository.GetGithubAppByAccountId(accountId)

	if githubApp == nil {
		return nil, errors.New("Account not connected to github. Please connect with `wave-deploy connect-github`")
	}

	installationAuthToken, err := github.AuthenticateAsGithubAppInstallation(githubApp.InstallationId)
	if err != nil {
		//TODO: Handle better
		return nil, err
	}

	ghRepositories, err := github.GetInstallationRepositories(installationAuthToken)
	if err != nil {
		//TODO: Handle better
		return nil, err
	}

	return ghRepositories, nil
}

func GetRepositoryCloneUrl(installationId, githubUrl string) (string, error) {
	installationAuthToken, err := github.AuthenticateAsGithubAppInstallation(installationId)
	if err != nil {
		return "", err
	}

	repoUrl := strings.Split(githubUrl, "git://")[1]
	userName := strings.Split(repoUrl, "/")[1]

	cloneUrl := fmt.Sprintf("https://%s:%s@%s", userName, installationAuthToken, repoUrl)
	return cloneUrl, nil
}
