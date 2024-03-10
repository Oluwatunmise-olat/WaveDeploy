package http_server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"log"
	"net/http"
	"os"
)

func HandleOauthCallbackWebhook(writer http.ResponseWriter, request *http.Request) {
	baseResponse := structs.BaseResponse{Status: true, Message: "Event processed successfully"}
	jsonPayload, _ := json.Marshal(&baseResponse)
	writer.Header().Set("Content-Type", "application/json")

	//if !validateGithubSignature(request) {
	//	if _, err := writer.Write(jsonPayload); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}

	query := request.URL.Query()
	payload := structs.GithubOauthWebhook{
		State:          query.Get("state"),
		Code:           query.Get("code"),
		InstallationId: query.Get("installation_id"),
		SetupAction:    query.Get("setup_action"),
	}

	if payload.InstallationId != "" && payload.SetupAction == "install" && payload.State != "" {
		github.CreateGithubAppIfNotExists(payload)
	}

	if _, err := writer.Write(jsonPayload); err != nil {
		log.Fatalln(err)
	}
}

func validateGithubSignature(request *http.Request) bool {
	signature := request.Header.Get("x-hub-signature-256")
	if signature == "" {
		return false
	}

	hasher := hmac.New(sha256.New, []byte(os.Getenv("GITHUB_APP_WEBHOOK_SECRET")))
	body, _ := json.Marshal(request.Body)
	hasher.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(hasher.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
