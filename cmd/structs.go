package cmd

type Prompt struct {
	label        string
	errorMessage string
	mask         rune
	confirm      bool
	items        []string
	allowEdit    bool
}
