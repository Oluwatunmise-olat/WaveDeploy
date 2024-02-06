package cmd

import (
	"errors"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"os"
	"time"
)

var (
	PromptTemplate = promptui.PromptTemplates{
		Prompt:  "{{ . }}",
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
	}
	value, err := prompt.Run()
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
