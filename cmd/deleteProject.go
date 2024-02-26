package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteProjectCmd = &cobra.Command{
	Use:   "delete-project",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteProject called")
	},
}

func init() {
	rootCmd.AddCommand(deleteProjectCmd)
}

// kill project if live
// delete project and related resources
