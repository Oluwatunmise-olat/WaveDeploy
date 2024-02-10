package structs

type GithubOauthWebhook struct {
	Code           string `json:"code"`
	InstallationId string `json:"installation_id"`
	SetupAction    string `json:"setup_action"`
	State          string `json:"state"`
}

type GithubAuthenticateAsInstallationSuccessResponse struct {
	Token       string `json:"token"`
	ExpiresAt   string `json:"expires_at"`
	Permissions struct {
	} `json:"permissions"`
}

type GithubAInstallationRepositoriesSuccessResponse struct {
	TotalCount   int                               `json:"total_count"`
	Repositories []GithubAInstallationRepositories `json:"repositories"`
}

type GithubAInstallationRepositories struct {
	FullName      string `json:"full_name"`
	Name          string `json:"name"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	GitUrl        string `json:"git_url"`
	SSHUrl        string `json:"ssh_url"`
	CloneUrl      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type GithubAuthenticateAsInstallationFailureResponse struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
