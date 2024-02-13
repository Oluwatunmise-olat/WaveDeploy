package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a deployed(live) project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restart called")
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
