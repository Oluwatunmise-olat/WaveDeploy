package cmd

import (
	"context"
	"errors"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/account"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	PromptTemplate = promptui.PromptTemplates{
		Prompt:  "{{ . | green }}",
		Valid:   "{{ . | blue }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}
)

func GetPromptInput(promptCommand Prompt, validator func(string) error) string {
	if validator == nil {
		validator = BasePromptValidator(promptCommand.errorMessage)
	}

	prompt := &promptui.Prompt{
		Label:     promptCommand.label,
		Templates: &PromptTemplate,
		Validate:  validator,
		Mask:      promptCommand.mask,
		IsConfirm: promptCommand.confirm,
	}
	value, err := prompt.Run()

	if err != nil {
		os.Exit(1)

	}

	return value
}

func GetPromptSelector(promptCommand Prompt, validator func(string) error) string {
	if validator == nil {
		validator = BasePromptValidator(promptCommand.errorMessage)
	}

	prompt := &promptui.Select{
		Label: promptCommand.label,
		Items: promptCommand.items,
	}

	_, value, err := prompt.Run()

	if err != nil {
		os.Exit(1)

	}

	return value
}

func BasePromptValidator(errorMessage string) func(string) error {
	return func(value string) error {
		if len(value) <= 0 {
			return errors.New(errorMessage)
		}
		return nil
	}
}

func initializeSpinner(prefix, finalMessage string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	if len(prefix) > 0 {
		s.Prefix = prefix
	}

	if len(finalMessage) > 0 {
		s.FinalMSG = finalMessage
	}

	return s
}

func IsAuthenticated(ctx context.Context, msg string, cobraCmd *cobra.Command) {
	s := initializeSpinner(msg, "")
	s.Start()
	accountId, err := account.GetAuthTokenDetails()
	if err != nil {
		s.FinalMSG = err.Error()
		s.Stop()
		os.Exit(1)
	}
	ctx = context.WithValue(ctx, "accountId", accountId)
	cobraCmd.SetContext(ctx)
	s.Stop()
}
