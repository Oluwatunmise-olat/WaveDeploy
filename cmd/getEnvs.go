package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd/flags"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var getEnvsCmd = &cobra.Command{
	Use:   "get-envs",
	Short: "Get project envs",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		accountId := getAccountID(cmd)

		if err := getProjectEnvs(accountId); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(getEnvsCmd)
	flags.InitializeProjectNameFlag(getEnvsCmd)
}

func getProjectEnvs(accountId string) error {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(flags.GetProjectName()))

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
