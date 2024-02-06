package structs

type GithubOauthWebhook struct {
	Code           string `json:"code"`
	InstallationId string `json:"installation_id"`
	SetupAction    string `json:"setup_action"`
	State          string `json:"state"`
}

type BaseResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
