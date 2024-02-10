package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/spf13/cobra"
)

var connectGithubCmd = &cobra.Command{
	Use:   "connect-github",
	Short: "Connect your github account to your wave-deploy account for seamless deployment",
	PreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		IsAuthenticated(ctx, "Checking github connection", cmd)
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
		s.FinalMSG = "Github account connected\n"
		s.Stop()
		return
	}

	connectToGithubLink := github.GetConnectToGithubUrl(accountId.(string))
	s.FinalMSG = fmt.Sprintf("Click on this link to authenticate with GitHub: %s\n", connectToGithubLink)
	s.Stop()
}
