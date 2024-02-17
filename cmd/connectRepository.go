package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/spf13/cobra"
)

var connectRepositoryCmd = &cobra.Command{
	Use:   "connect-repository",
	Short: "A brief description of your command",
	PreRun: func(cmd *cobra.Command, args []string) {
		IsAuthenticated(cmd.Context(), "", cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		connectGithubRepositoryToAppInstallation(cmd)
	},
}

func init() {
	rootCmd.AddCommand(connectRepositoryCmd)
}

func connectGithubRepositoryToAppInstallation(cmd *cobra.Command) {
	accountId := cmd.Context().Value("accountId")
	s := initializeSpinner("", "")
	s.Start()
	accountConnected := github.IsAccountConnectedAlreadyToGithub(accountId.(string))

	if !accountConnected {
		s.FinalMSG = "Account not connected to github. Please connect with `wave-deploy connect-github`"
		s.Stop()
		return
	}

	connectToGithubLink := github.GetConnectGithubRepositoryUrl()
	s.FinalMSG = fmt.Sprintf("Click on this link to connect a repository to account: %s\n", connectToGithubLink)
	s.Stop()
}
