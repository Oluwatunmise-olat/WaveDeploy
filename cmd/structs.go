package cmd

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/google/uuid"
)

type Prompt struct {
	label        string
	errorMessage string
	mask         rune
	confirm      bool
	items        []string
	allowEdit    bool
}

type DeploymentOptions struct {
	PrivateKeyPath string
	PublicIPV4Addr string
	VmUser         string
	Envs           ProjectEnvs
	RemoteAppDir   string
	RemoteHomeDir  string
}

type BuildApplicationOptions struct {
	AccountId            string
	ProjectId            uuid.UUID
	ProjectUpdatePayload models.Projects
	Envs                 ProjectEnvs
	DeploymentOptions    DeploymentOptions
}

type WebServerTmpl struct {
	EXTERNAL_PORT           int
	INTERNAL_LISTENING_PORT int
}

type SPATmpl struct {
	EXTERNAL_PORT    int
	APPLICATION_PATH string
}
