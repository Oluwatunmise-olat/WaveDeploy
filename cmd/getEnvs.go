package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var getEnvsCmd = &cobra.Command{
	Use:   "get-envs",
	Short: "Get project envs",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		projectName := getProjectName(cmd)

		if err := getProjectEnvs(accountId, projectName); err != nil {
			return fmt.Errorf("error occurred fetching project envs: %w", err)
		}

		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy get-envs -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(getEnvsCmd)

	getEnvsCmd.Flags().StringP("name", "n", "", "Project Name")
	getEnvsCmd.MarkFlagRequired("name")
}

func getProjectEnvs(accountId, projectName string) error {
	project, err := projects.GetProjectByName(accountId, projectName)

	if err != nil {
		return fmt.Errorf("project not found")
	}

	accountUUID, _ := uuid.Parse(project.AccountId)
	envs, err := projects.GetProjectEnvs(project.Id, accountUUID)

	if err != nil {
		return fmt.Errorf("project has no envs set")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Env Key", "Env Value"})

	for _, env := range envs {
		value, _err := hashers.DecryptIt(env.Value, os.Getenv("APP_KEY"))
		if _err != nil {
			return _err
		}

		table.Append([]string{env.Key, value})
	}

	table.Render()
	return nil
}
