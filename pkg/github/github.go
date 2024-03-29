package github

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/files"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/http"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/jwt"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"io"
	"net/url"
	"os"
)

var (
	httpclient = http.HttpClient("https://api.github.com")
)

func ConnectAppToGithub(state string) string {
	connectGithubUrl := os.Getenv("GITHUB_APP_PUBLIC_LINK")
	parsedUrl, _ := url.Parse(connectGithubUrl)
	parsedUrl = parsedUrl.JoinPath("/installations/new")

	urlParams := url.Values{}
	urlParams.Add("state", state)

	parsedUrl.RawQuery = urlParams.Encode()
	return parsedUrl.String()
}

func AuthenticateAsGithubAppInstallation(installationId string) (string, error) {
	urlPath := fmt.Sprintf("/app/installations/%s/access_tokens", installationId)
	token, err := GetGithubAppJWT()

	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        fmt.Sprintf("Bearer %s", token),
		"X-GitHub-Api-Version": "2022-11-28",
	}
	httpclient.Headers = headers
	response, err := httpclient.SendPostRequest(urlPath, "")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 400 {
		//TODO:: Handle better
		return "", errors.New("....")
	}

	var responsePayload structs.GithubAuthenticateAsInstallationSuccessResponse
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		//TODO:: Handle better
		return "", errors.New(err.Error())
	}

	return responsePayload.Token, nil
}

func GetGithubAppJWT() (string, error) {
	bytes, err := files.GetFileContent(os.Getenv("PRIVATE_KEY_PATH"))

	if err != nil {
		return "", errors.New("an error occurred reading private key file")
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("invalid private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New("error parsing private key")
	}

	token, err := jwt.GenerateGithubAppJWT(privateKey)
	if err != nil {
		return "", errors.New("An error occurred signing app jwt")
	}

	return token, nil
}

func GetInstallationRepositories(accessToken string) ([]structs.GithubAInstallationRepositories, error) {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        fmt.Sprintf("Bearer %s", accessToken),
		"X-GitHub-Api-Version": "2022-11-28",
	}
	httpclient.Headers = headers
	response, err := httpclient.SendGetRequest("/installation/repositories")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, errors.New("An error occurred fetching installed repositories")
	}

	var responsePayload structs.GithubAInstallationRepositoriesSuccessResponse
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		//	//TODO:: Handle better
		return nil, errors.New(err.Error())
	}

	return responsePayload.Repositories, nil
}
