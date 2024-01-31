package cmd

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of your wave-deploy account",
	Run: func(cmd *cobra.Command, args []string) {
		s := initializeSpinner(" Logging Off ...", "Successfully Logged Out\n")
		s.Start()
		auth.LogoutAccount()
		s.Stop()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
