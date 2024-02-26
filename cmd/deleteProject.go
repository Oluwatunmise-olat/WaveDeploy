package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd/flags"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/google/uuid"
	"strings"

	"github.com/spf13/cobra"
)

var deleteProjectCmd = &cobra.Command{
	Use:   "delete-project",
	Short: "Delete a project (Deployed Or Not)",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		if err := deleteProject(accountId); err != nil {
			return fmt.Errorf("error occurred deleting project: %w", err)
		}

		fmt.Println("Project deleted successfully 🏮")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteProjectCmd)
}

func deleteProject(accountId string) error {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(flags.GetProjectName()))
	if err != nil {
		return fmt.Errorf("project not found")
	}

	accountUUID, _ := uuid.Parse(accountId)

	if project.IsLive {
		if err = killProject(accountId); err != nil {
			return err
		}
	}

	if err = projects.DeleteProjectAndRelatedResources(project.Id, accountUUID); err != nil {
		return err
	}

	return nil
}
