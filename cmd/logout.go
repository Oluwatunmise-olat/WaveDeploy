package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/auth"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/github"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of your wave-deploy account",
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := github.GetAllAccountInstallations()
		fmt.Println(data)
		s := initializeSpinner(" Logging Off ... ", "Successfully Logged Out\n")
		s.Start()
		auth.LogoutAccount()
		s.Stop()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
