package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/account"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a wave-deploy account",
	Run: func(cmd *cobra.Command, args []string) {
		createAccount()
	},
	Example: "wave-deploy create-account",
}

func init() {
	rootCmd.AddCommand(createAccountCmd)
}

func createAccount() {
	userName, email, password := promptAccountCreationDetails()
	email = strings.ToLower(email)

	accountWithEmail := account.GetAccountByEmail(email)
	if accountWithEmail.Id != uuid.Nil {
		fmt.Println("Account with email already exist")
		return
	}

	hashedPswd, _ := hashers.HashPassword(password)

	if err := account.CreateAccount(models.Accounts{
		Id:         random.GetUUID(),
		UserName:   userName,
		Email:      email,
		Password:   hashedPswd,
		LastAuthAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}); err != nil {
		fmt.Println("Unable to create account. Please try again later")
		return
	}

	fmt.Println("Account Created 🏄🏽‍")
}

func promptAccountCreationDetails() (string, string, string) {
	UsernameCommand := promptForCommands("Username", false)
	EmailCommand := promptForCommands("Email", false)
	PasswordCommand := promptForCommands("Password", true)

	return UsernameCommand, EmailCommand, PasswordCommand
}

func promptForCommands(title string, mask bool) string {
	cmd := Prompt{
		label:        fmt.Sprintf("%s: ", title),
		errorMessage: fmt.Sprintf("Please provide %s", strings.ToLower(title)),
	}

	if mask == true {
		cmd.mask = '*'
	}

	value := GetPromptInput(cmd, nil)

	return value
}
