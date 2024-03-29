package cmd

import (
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
	Replicas       int
}

type BuildApplicationOptions struct {
	AccountId            string
	ProjectId            uuid.UUID
	ProjectUpdatePayload map[string]interface{}
	Envs                 ProjectEnvs
	DeploymentOptions    DeploymentOptions
}

type WebServerTmpl struct {
	EXTERNAL_PORT           int
	INTERNAL_LISTENING_PORT int
}

type SpaTmpl struct {
	EXTERNAL_PORT    int
	APPLICATION_PATH string
}
