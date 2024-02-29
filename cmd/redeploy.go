package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var redeployCmd = &cobra.Command{
	Use:   "redeploy",
	Short: "Redeploy a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("redeploy coming soon")
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy redeploy -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(redeployCmd)
	redeployCmd.Flags().StringP("name", "n", "", "Project Name")
	redeployCmd.MarkFlagRequired("name")
}
