package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your wave-deploy account",
	Run: func(cmd *cobra.Command, args []string) {
		initializeAuthentication()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func initializeAuthentication() {
	emailPromptCmd := Prompt{errorMessage: "Please provide your email address", label: "Email Address: "}
	passwordPromptCmd := Prompt{errorMessage: "Please provide your password", label: "Password: ", mask: '*'}

	email := GetPromptInput(emailPromptCmd, nil)
	password := GetPromptInput(passwordPromptCmd, nil)

	s := initializeSpinner(" Authenticating ...", fmt.Sprintf("Logged in as %s\n", email))

	s.Start()
	if err := auth.AuthenticateAccount(email, password); err != nil {
		s.FinalMSG = err.Error()
	}
	s.Stop()
}
