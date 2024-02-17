package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Oluwatunmise-olat/WaveDeploy/internal/account"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/files"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/google/uuid"
	"github.com/melbahja/goph"
)

type ProjectEnvs map[string]string

var projectName string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a project",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		accountId := getAccountID(cmd)
		project, err := preDeploymentChecks(accountId)
		if err != nil {
			log.Fatal(err)
			return
		}
		deployProject(cmd, project)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	deployCmd.MarkFlagRequired("name")
}

func checkGitHubConnection(cmd *cobra.Command) {
	IsAuthenticated(cmd.Context(), "Checking GitHub connection", cmd)
}

func getAccountID(cmd *cobra.Command) string {
	return cmd.Context().Value("accountId").(string)
}

func preDeploymentChecks(accountID string) (*models.Projects, error) {
	project, err := projects.GetProjectByName(accountID, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project: %v", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project with name '%s' not found", projectName)
	}
	if project.IsLive {
		return nil, errors.New("project is already deployed")
	}
	return project, nil
}

func deployProject(cmd *cobra.Command, project *models.Projects) {
	accountId := getAccountID(cmd)
	updatedPayload := models.Projects{}
	buildCommand, runCommand, Envs, err := promptDeploymentOptions()
	if err != nil {
		log.Fatal(err)
	}

	updatedPayload.BuildCommand = buildCommand
	updatedPayload.RunCommand = runCommand

	envsPayload, err := createEnvRecords(accountId, project, Envs)
	if err != nil {
		log.Fatal(err)
	}

	if err = updateProjectAndCreateEnvs(accountId, project, updatedPayload, envsPayload); err != nil {
		log.Fatal(err)
	}

	buildPath, err := buildApplicationDockerfile(accountId, project.Id, Envs)
	if err != nil {
		log.Fatal("Failed to build application Dockerfile: ", err)
	}

	// TODO:: Collect vm credentials (username, privateKeyPath, publicIp)
	deploymentOptions := DeploymentOptions{
		Envs:           Envs,
		BuildPath:      buildPath,
		VmUser:         "******",
		PublicIPV4Addr: "******",
		PrivateKeyPath: "******",
	}
	if err = deployAndStartApplication(deploymentOptions); err != nil {
		log.Fatal("Failed to deploy and start application:", err)
	}

	fmt.Println("Application deployed successfully!")
}

func promptDeploymentOptions() (string, string, ProjectEnvs, error) {
	buildCommand, err := promptForCommand("Build", "build")
	runCommand, err := promptForCommand("Run", "run")
	envs, err := promptForEnvVariables()

	return buildCommand, runCommand, envs, err
}

func promptForCommand(action, commandType string) (string, error) {
	cmd := Prompt{
		label: fmt.Sprintf("> Custom %s command (y/n)?: ", action),
	}
	value := GetPromptInput(cmd, nil)
	if value != "y" {
		return "", nil
	}

	prompt := Prompt{
		errorMessage: fmt.Sprintf("Please provide a %s command", action),
		label:        fmt.Sprintf("%s command: ", action),
	}

	return GetPromptInput(prompt, nil), nil
}

func promptForEnvVariables() (ProjectEnvs, error) {
	envs := make(ProjectEnvs)
	cmd := Prompt{label: "> Set Envs (y/n)?: "}
	setEnvValue := GetPromptInput(cmd, nil)
	if setEnvValue != "y" {
		return envs, nil
	}
	for {
		envKey := GetPromptInput(Prompt{label: "Enter the environment variable key: "}, nil)
		envValue := GetPromptInput(Prompt{label: "Enter the environment variable value: ", mask: '*'}, nil)
		envs[envKey] = envValue
		rePrompt := GetPromptInput(Prompt{label: "> Set More Envs (y/n)?: "}, nil)
		if rePrompt != "y" {
			break
		}
	}
	return envs, nil
}

func createEnvRecords(accountID string, project *models.Projects, Envs ProjectEnvs) ([]models.Envs, error) {
	envsPayload := make([]models.Envs, 0, len(Envs))
	accountUUID, _ := uuid.Parse(accountID)
	for key, value := range Envs {
		encryptedValue, err := hashers.EncryptIt(value, os.Getenv("APP_KEY"))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt environment variable value: %v", err)
		}
		envRecord := models.Envs{
			Key:       key,
			Value:     encryptedValue,
			Id:        random.GetUUID(),
			AccountId: accountUUID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ProjectId: project.Id,
		}
		envsPayload = append(envsPayload, envRecord)
	}
	return envsPayload, nil
}

func updateProjectAndCreateEnvs(accountID string, project *models.Projects, updatedPayload models.Projects, envsPayload []models.Envs) error {
	accountIdToUUId, _ := uuid.Parse(accountID)
	payload := projects.UpdateProjectAndCreateEnvsPayload{
		Envs:                 envsPayload,
		UpdateProjectPayload: updatedPayload,
		AccountId:            accountIdToUUId,
		ProjectId:            project.Id,
	}
	return projects.UpdateProjectAndCreateEnvs(payload)
}

func buildApplicationDockerfile(accountID string, projectId uuid.UUID, projectEnvs ProjectEnvs) (string, error) {
	project, _ := projects.GetProjectById(projectId.String(), accountID)
	rootDirectory := files.GetCurrentPathRootDirectory()
	scriptPath := filepath.Join(rootDirectory, "/../../scripts")
	executionPathBytes, _ := exec.Command("mktemp", "-d").Output()
	githubCloneUrl := ""
	appPath, _ := os.Getwd()

	executionPath := strings.TrimSpace(string(executionPathBytes))

	if project.GithubRepoUrl != "" {
		installationId, err := account.GetAccountInstallationId(accountID)
		if err != nil {
			return "", fmt.Errorf("failed to get GitHub installation ID: %v", err)
		}
		githubCloneUrl, err = github.GetRepositoryCloneUrl(installationId, project.GithubRepoUrl)
		if err != nil {
			return "", fmt.Errorf("failed to get GitHub repository clone URL: %v", err)
		}

		appPath = "."
	}
	args := []string{"-n", project.Name, "-w", string(executionPath), "-p", appPath}

	if githubCloneUrl != "" {
		args = append(args, "-l", githubCloneUrl)
	}

	if project.BuildCommand != "" {
		args = append(args, "-b", fmt.Sprintf(`"%s"`, project.BuildCommand))
	}

	if project.RunCommand != "" {
		args = append(args, "-s", fmt.Sprintf(`"%s"`, project.RunCommand))
	}

	for key, value := range projectEnvs {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	args = append([]string{filepath.Join(scriptPath, "generate-dockerfile.sh")}, args...)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command("/bin/bash", args...)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr)

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return "", err
	}

	err := cmd.Wait()

	return executionPath, err
}

// TODO:: Add nginx or caddy as extra load balancing layer (ssl termination)
// TODO:: Add letsEncrypt
func deployAndStartApplication(opts DeploymentOptions) error {
	remoteWorkDir := fmt.Sprintf("/home/%s/app/.builder", opts.VmUser)
	fullBuildPath := filepath.Join(opts.BuildPath, ".nixpacks")

	auth, err := goph.Key(opts.PrivateKeyPath, "")
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %v", err)
	}

	client, err := goph.New(opts.VmUser, opts.PublicIPV4Addr, auth)
	if err != nil {
		return fmt.Errorf("failed to create SSH client: %v", err)
	}
	defer client.Close()

	client.Run(fmt.Sprintf("mkdir -p %s", remoteWorkDir))
	err = UploadFileToRemoteRecursively(client, fullBuildPath, remoteWorkDir)
	if err != nil {
		return err
	}
	exec.Command("rm -rf", opts.BuildPath)

	// Change to Work Directory
	client.Run(fmt.Sprintf("cd %s", fullBuildPath))
	// Build Docker Image
	client.Run(fmt.Sprintf("sudo docker build -t %s:latest .", projectName))
	// Initialize Docker Swarm
	client.Run("sudo docker swarm init")

	// Deploy Command
	// TODO:: Make replicas setting flexible
	createCmd := fmt.Sprintf("docker service create --name %s --replicas 1 --publish published=8080,target=80", projectName)
	// Environment variables map
	for key, value := range opts.Envs {
		createCmd += fmt.Sprintf(" --env %s=%s", key, value)
	}
	createCmd += fmt.Sprintf(" %s:latest", projectName)

	// Deploy service
	//out, err := client.Run(createCmd)
	return nil
}

func UploadFileToRemoteRecursively(client *goph.Client, localFolder, remoteFolder string) error {
	return filepath.Walk(localFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(localFolder, path)

		if len(relPath) > 2 {
			if err != nil {
				return err
			}

			remoteFilePath := filepath.Join(remoteFolder, relPath)
			err = client.Upload(path, remoteFilePath)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
