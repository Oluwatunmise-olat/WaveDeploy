package github

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
)

func IsAccountConnectedAlreadyToGithub(accountId string) bool {
	githubAppRepository := respository.InitializeGithubAppsRepository()
	githubApp, _ := githubAppRepository.GetGithubAppByAccountId(accountId)
	return githubApp == nil
}

func GetConnectToGithubUrl(accountId string) string {
	base64String := hashers.EncodeStringToBase64(accountId)
	return github.ConnectAppToGithub(base64String)
}
