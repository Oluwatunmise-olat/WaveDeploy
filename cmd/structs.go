package cmd

type Prompt struct {
	label        string
	errorMessage string
	mask         rune
	confirm      bool
	items        []string
	allowEdit    bool
}

type DeploymentOptions struct {
	BuildPath      string
	PrivateKeyPath string
	PublicIPV4Addr string
	VmUser         string
	Envs           ProjectEnvs
}
