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
		IsAuthenticated(cmd.Context(), "", cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		createProject(cmd)
	},
}

func init() {
	rootCmd.AddCommand(createProjectCmd)
}

func createProject(cmd *cobra.Command) {
	_projectName := PromptForProjectName()

	accountId := cmd.Context().Value("accountId").(string)
	nameTaken, _ := projects.IsProjectNameTaken(accountId, _projectName)

	if nameTaken {
		fmt.Println("Project with name already exists")
		return
	}

	selectedRepository := PromptForGithubRepository(accountId)

	err := projects.CreateProject(accountId, _projectName, &selectedRepository)
	if err != nil {
		fmt.Println("Error creating project:", err)
		return
	}

	fmt.Println("Project created ✨")
}

// PromptForProjectName Prompts the user for the project name
func PromptForProjectName() string {
	projectNamePromptCmd := Prompt{errorMessage: "Please provide a project name", label: "> Project name: "}
	return GetPromptInput(projectNamePromptCmd, nil)
}

// PromptForGithubRepository Prompts user to select a GitHub repository
func PromptForGithubRepository(accountId string) structs.GithubAInstallationRepositories {
	codeSourceGHPromptCmd := Prompt{label: "> Select project from github (y/n)?: "}
	codeSourceIsGH := GetPromptInput(codeSourceGHPromptCmd, nil)
	var selectedRepository structs.GithubAInstallationRepositories

	if codeSourceIsGH == "y" {
		ghRepositories, err := github.GetAccountConnectedRepositories(accountId)
		if err != nil {
			fmt.Println("Error fetching GitHub repositories:", err)
			return structs.GithubAInstallationRepositories{}
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

	return selectedRepository
}
