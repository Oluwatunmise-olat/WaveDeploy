package http_server

import (
	"encoding/json"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"log"
	"net/http"
)

func HandleOauthCallbackWebhook(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	payload := structs.GithubOauthWebhook{
		State:          query.Get("state"),
		Code:           query.Get("code"),
		InstallationId: query.Get("installation_id"),
		SetupAction:    query.Get("setup_action"),
	}

	writer.Header().Set("Content-Type", "application/json")
	baseResponse := structs.BaseResponse{Status: true, Message: "Event processed successfully"}

	if payload.InstallationId != "" && payload.SetupAction == "install" && payload.State != "" {
		github.CreateGithubAppIfNotExists(payload)
	}

	jsonPayload, _ := json.Marshal(&baseResponse)
	if _, err := writer.Write(jsonPayload); err != nil {
		log.Fatalln(err)
	}
}
