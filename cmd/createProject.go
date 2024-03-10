package cmd

import (
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"github.com/spf13/cobra"
	"strings"
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
	Example: "wave-deploy create-project",
}

func init() {
	rootCmd.AddCommand(createProjectCmd)
}

func createProject(cmd *cobra.Command) {
	_projectName := strings.TrimSpace(PromptForProjectName())

	accountId := cmd.Context().Value("accountId").(string)
	nameTaken, _ := projects.IsProjectNameTaken(accountId, _projectName)

	if nameTaken {
		fmt.Println("Project with name already exists")
		return
	}

	projectType := PromptForProjectType()

	selectedRepository := PromptForGithubRepository(accountId)
	if selectedRepository.GitUrl == "" {
		return
	}

	err := projects.CreateProject(accountId, _projectName, projectType, &selectedRepository)
	if err != nil {
		fmt.Println("Error creating project: Please connect account to github with `wave-deploy connect-github`")
		return
	}

	fmt.Println("Project created ✨")
}

// PromptForProjectName Prompts the user for the project name
func PromptForProjectName() string {
	projectNamePromptCmd := Prompt{errorMessage: "Please provide a project name", label: "> Project name: "}
	return GetPromptInput(projectNamePromptCmd, nil)
}

// PromptForProjectType prompts the user for the project type
func PromptForProjectType() string {
	selectRProjectTypePromptCmd := Prompt{label: "> Select project type ", items: []string{string(models.API), string(models.SPA)}}
	return GetPromptSelector(selectRProjectTypePromptCmd, nil)
}

// PromptForGithubRepository Prompts user to select a GitHub repository
func PromptForGithubRepository(accountId string) structs.GithubAInstallationRepositories {
	var selectedRepository structs.GithubAInstallationRepositories

	ghRepositories, err := github.GetAccountConnectedRepositories(accountId)
	if err != nil {
		fmt.Println("Error fetching GitHub repositories: Please connect account to github with `wave-deploy connect-github`")
		return structs.GithubAInstallationRepositories{GitUrl: ""}
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

	return selectedRepository
}
