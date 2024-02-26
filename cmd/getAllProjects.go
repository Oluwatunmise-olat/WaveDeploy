package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var getAllProjectsCmd = &cobra.Command{
	Use:   "get-projects",
	Short: "Get all projects",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		if err := getAccountProjects(accountId); err != nil {
			return fmt.Errorf("error occurred fetching all projects: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAllProjectsCmd)
}

func getAccountProjects(accountId string) error {
	accountProjects, err := projects.GetAllAccountProjects(accountId)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Project Name", "Connected Repository", "Is Live"})

	for _, project := range accountProjects {
		table.Append([]string{project.Name, project.GithubRepoUrl, strconv.FormatBool(project.IsLive)})
	}

	table.Render()
	return nil
}
