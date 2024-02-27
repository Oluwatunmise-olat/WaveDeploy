package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/google/uuid"
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
		projectName := getProjectName(cmd)

		if err := deleteProject(accountId, projectName); err != nil {
			return fmt.Errorf("error occurred deleting project: %w", err)
		}

		fmt.Println("Project deleted successfully 🏮")
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy delete-project -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(deleteProjectCmd)

	deleteProjectCmd.Flags().StringP("name", "n", "", "Project Name")
	deleteProjectCmd.MarkFlagRequired("name")
}

func deleteProject(accountId, projectName string) error {
	project, err := projects.GetProjectByName(accountId, projectName)
	if err != nil {
		return fmt.Errorf("project not found")
	}

	accountUUID, _ := uuid.Parse(accountId)

	if project.IsLive {
		if err = killProject(accountId, project.Name); err != nil {
			return err
		}
	}

	if err = projects.DeleteProjectAndRelatedResources(project.Id, accountUUID); err != nil {
		return err
	}

	return nil
}
