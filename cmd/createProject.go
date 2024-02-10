package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new web project",
	PreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		IsAuthenticated(ctx, "", cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		createProject(cmd)
	},
}

func init() {
	rootCmd.AddCommand(createProjectCmd)
}

func createProject(cmd *cobra.Command) {
	projectNamePromptCmd := Prompt{errorMessage: "Please provide a project name", label: "> Project name: "}
	projectName := GetPromptInput(projectNamePromptCmd, nil)

	accountId := cmd.Context().Value("accountId").(string)
	nameTaken, _ := projects.IsProjectNameTaken(accountId, projectName)

	if nameTaken {
		fmt.Println("Project with name already exists")
		return
	}

	codeSourceGHPromptCmd := Prompt{label: "> Select project from github (y/n)?: "}
	codeSourceIsGH := GetPromptInput(codeSourceGHPromptCmd, nil)

	var selectedRepository structs.GithubAInstallationRepositories

	if codeSourceIsGH == "y" {
		ghRepositories, err := github.GetAccountConnectedRepositories(accountId)

		if err != nil {
			fmt.Errorf(err.Error())
			return
		}

		var ghRepositoryNames []string
		for _, repo := range ghRepositories {
			ghRepositoryNames = append(ghRepositoryNames, repo.Name)
		}

		selectRepoPromptCmd := Prompt{label: "> Select a repository ", items: ghRepositoryNames}
		selectedRepositoryName := GetPromptSelector(selectRepoPromptCmd, nil)

		for _, repo := range ghRepositories {
			if repo.Name == selectedRepositoryName {
				selectedRepository = repo
				break
			}
		}
	}

	err := projects.CreateProject(accountId, projectName, &selectedRepository)
	if err != nil {
		fmt.Sprintf(err.Error())
	}
}

// No linked project found. Run railway link to connect to a project, and a service.
// project name
//> Select a project royal-back
//> Select an environment production
//> Select a service Abstract-Syntax-Tree
