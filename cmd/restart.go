package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a deployed(live) project",
	PreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		IsAuthenticated(ctx, "", cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not Implemented")
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
