package github

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"log"
	"time"
)

func IsAccountConnectedAlreadyToGithub(accountId string) bool {
	githubAppRepository := respository.GithubAppsRepository{}
	githubApp, _ := githubAppRepository.GetGithubAppByAccountId(accountId)
	return githubApp != nil
}

func GetConnectToGithubUrl(accountId string) string {
	base64String := hashers.EncodeStringToBase64(accountId)
	return github.ConnectAppToGithub(base64String)
}

func CreateGithubAppIfNotExists(payload structs.GithubOauthWebhook) {
	accountRepository := respository.AccountsRepository{}
	account, err := accountRepository.GetAccountById(hashers.DecodeBase64ToString(payload.State))

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
		AccountId:      hashers.DecodeBase64ToString(payload.State),
		Code:           payload.Code,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Id:             random.GetUUID(),
	}

	if err = githubAppRepository.CreateGithubApp(newGithubApp); err != nil {
		log.Fatalln(err)
	}
}
