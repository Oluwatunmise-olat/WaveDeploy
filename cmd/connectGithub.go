package cmd

import (
	"context"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/auth"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/spf13/cobra"
)

var connectGithubCmd = &cobra.Command{
	Use:   "connect-github",
	Short: "A brief description of connect github",
	PreRun: func(cmd *cobra.Command, args []string) {
		s := initializeSpinner("Checking github connection", "")
		s.Start()
		accountId, err := auth.GetAuthTokenDetails()
		if err != nil {
			s.FinalMSG = err.Error()
			s.Stop()
		}
		ctx := cmd.Context()
		ctx = context.WithValue(ctx, "accountId", accountId)
		cmd.SetContext(ctx)
		s.Stop()
	},
	Run: func(cmd *cobra.Command, args []string) {
		connectToGithub(cmd)
	},
}

func init() {
	rootCmd.AddCommand(connectGithubCmd)
}

func connectToGithub(cmd *cobra.Command) {
	s := initializeSpinner("", "")
	s.Start()
	accountId := cmd.Context().Value("accountId")

	accountConnected := github.IsAccountConnectedAlreadyToGithub(accountId.(string))
	if accountConnected {
		s.FinalMSG = "Github account connected"
		s.Stop()
		return
	}

	connectToGithubLink := github.GetConnectToGithubUrl(accountId.(string))
	s.FinalMSG = fmt.Sprintf("Click on this link to authenticate with GitHub: %s\n", connectToGithubLink)
	s.Stop()
}
