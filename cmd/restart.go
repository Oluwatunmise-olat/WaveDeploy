package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
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
		// TODO:: Clear up after testing
		projectRecord, err := github.GetRepositoryCloneUrl("47014654", "git://github.com/Oluwatunmise-olat/go_math.git")
		_ = err
		fmt.Println(projectRecord)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
