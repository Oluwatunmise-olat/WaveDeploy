package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateEnvsCmd = &cobra.Command{
	Use:   "update-envs",
	Short: "Update project envs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updateEnvs called")
	},
}

func init() {
	rootCmd.AddCommand(updateEnvsCmd)
}
