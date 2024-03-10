package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	"html/template"
)

// todo:
// Dynamic ports on vm

type ProjectEnvs map[string]string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a project",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		projectName := getProjectName(cmd)

		project, err := preDeploymentChecks(accountId, projectName)
		if err != nil {
			return fmt.Errorf("error occurred deploying project: %w", err)
		}
		deployProject(cmd, project)
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy deploy -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("name", "n", "", "Project Name")
	deployCmd.MarkFlagRequired("name")
}

func checkGitHubConnection(cmd *cobra.Command) {
	IsAuthenticated(cmd.Context(), "Checking GitHub connection", cmd)
}

func preDeploymentChecks(accountId, projectName string) (*models.Projects, error) {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(projectName))
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
	updatedPayload := make(map[string]interface{})

	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()

	replicas := promptForReplicaCommand()
	buildCommand, runCommand, Envs, err := promptDeploymentOptions()

	if err != nil {
		log.Fatal(err)
	}

	updatedPayload["build_command"] = buildCommand
	updatedPayload["run_command"] = runCommand
	updatedPayload["replicas"] = replicas

	envsPayload, err := createEnvRecords(accountId, project, Envs)
	if err != nil {
		log.Fatal(err)
	}

	if err = updateProjectAndCreateEnvs(accountId, project, updatedPayload, envsPayload); err != nil {
		log.Fatal(err)
	}

	deploymentOptions := DeploymentOptions{
		Envs:           Envs,
		VmUser:         vmUser,
		PublicIPV4Addr: ipv4Addr,
		PrivateKeyPath: privateKeyPath,
		RemoteHomeDir:  fmt.Sprintf("/home/%s", vmUser),
		RemoteAppDir:   fmt.Sprintf("/home/%s/app/.builder", vmUser),
		Replicas:       replicas,
	}

	s := initializeSpinner("Building Application ", "Application Build Successful\n")
	s.Start()
	remoteAppDir, err := buildApplicationDockerfile(BuildApplicationOptions{
		AccountId:            accountId,
		ProjectId:            project.Id,
		ProjectUpdatePayload: updatedPayload,
		Envs:                 Envs,
		DeploymentOptions:    deploymentOptions,
	})
	if err != nil {
		s.FinalMSG = "Application Build Failed\n"
		s.Stop()
		fmt.Sprintf("Failed to build application: %v", err)
		return
	}
	s.Stop()

	s.Prefix = "Deploying Application "
	s.Start()
	if err = deployAndStartApplication(deploymentOptions, remoteAppDir, project.Name); err != nil {
		s.FinalMSG = "Application Deployment Failed\n"
		s.Stop()
		log.Fatal("Failed to deploy and start application:", err)
	}
	accountUUId, _ := uuid.Parse(accountId)
	s.Stop()

	if err = projects.UpdateProject(map[string]interface{}{"is_live": true}, project.Id, accountUUId); err != nil {
		log.Fatal("Application deployment failed 😪")
	}

	fmt.Println("Application deployed successfully! ⚡️")
}

func promptDeploymentOptions() (string, string, ProjectEnvs, error) {
	buildCommand, err := promptForCommand("Build", "build")
	runCommand, err := promptForCommand("Run", "run")
	envs, err := promptForEnvVariables()

	return buildCommand, runCommand, envs, err
}

func promptDeploymentCredentialsDetails() (string, string, string) {
	vmUserCommand := promptForVmCommands("User")
	vmIpCommand := promptForVmCommands("Public Ipv4 Address")
	vmPrivateKeyPathCommand := promptForVmCommands("Private Key Path")

	return vmUserCommand, vmIpCommand, vmPrivateKeyPathCommand
}

func promptForVmCommands(title string) string {
	cmd := Prompt{
		label:        fmt.Sprintf("Vm %s: ", title),
		errorMessage: fmt.Sprintf("Please provide vm %s", strings.ToLower(title)),
	}
	value := GetPromptInput(cmd, nil)

	return value
}

func promptForReplicaCommand() int {
	defaultReplicas := 1

	cmd := Prompt{
		label: fmt.Sprintf("> Replicas Count (default: %d): ", defaultReplicas),
	}

	for {
		value := GetPromptInput(cmd, nil)
		if value == "" {
			return defaultReplicas
		}

		count, err := strconv.Atoi(value)
		if err == nil {
			return count
		}

		fmt.Println("Invalid input. Please enter a numeric value for the replica count.")
	}
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
	uniqueKeys := make(map[string]bool)

	for key, value := range Envs {
		if _, ok := uniqueKeys[key]; ok {
			// Duplicated key, skip
			continue
		}

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

		uniqueKeys[key] = true
	}
	return envsPayload, nil
}

func updateProjectAndCreateEnvs(accountID string, project *models.Projects, updatedPayload map[string]interface{}, envsPayload []models.Envs) error {
	accountIdToUUId, _ := uuid.Parse(accountID)
	payload := projects.UpdateProjectAndCreateEnvsPayload{
		Envs:                 envsPayload,
		UpdateProjectPayload: updatedPayload,
		AccountId:            accountIdToUUId,
		ProjectId:            project.Id,
	}

	return projects.UpdateProjectAndCreateEnvs(payload)
}

func buildApplicationDockerfile(opts BuildApplicationOptions) (string, error) {
	project, _ := projects.GetProjectById(opts.ProjectId.String(), opts.AccountId)
	//ghRepoUrl := strings.Split(project.GithubRepoUrl, "/")
	//ghRepoName := strings.ReplaceAll(ghRepoUrl[len(ghRepoUrl)-1], ".git", "")

	appRootDirectory := files.GetCurrentPathRootDirectory()
	scriptPath := filepath.Join(appRootDirectory, "/../scripts")
	dockerFileScriptPath := fmt.Sprintf("%s/generate-dockerfile.sh", opts.DeploymentOptions.RemoteAppDir)
	vmSetupScriptPath := fmt.Sprintf("%s/setup-ubuntu-vm.sh", opts.DeploymentOptions.RemoteAppDir)
	appRemoteDirectory := opts.DeploymentOptions.RemoteHomeDir + fmt.Sprintf("/app/%s", project.Name)

	client, err := establishSSHConnection(opts.DeploymentOptions)
	if err != nil {
		return "", err
	}
	defer client.Close()

	_, err = client.Run(fmt.Sprintf("mkdir -p %s", opts.DeploymentOptions.RemoteAppDir))
	if err != nil {
		return "", err
	}

	err = uploadFileToRemoteRecursively(client, scriptPath, opts.DeploymentOptions.RemoteAppDir)
	if err != nil {
		return "", err
	}

	err = makeRemoteFileExecutableMode(client, dockerFileScriptPath)
	if err != nil {
		return "", err
	}

	err = makeRemoteFileExecutableMode(client, vmSetupScriptPath)
	if err != nil {
		return "", err
	}

	installationId, err := account.GetAccountInstallationId(opts.AccountId)
	if err != nil {
		return "", fmt.Errorf("Failed to get GitHub installation ID: %v", err)
	}

	githubCloneUrl, err := github.GetRepositoryCloneUrl(installationId, project.GithubRepoUrl)
	if err != nil {
		errCtx := err
		_ = errCtx
		return "", fmt.Errorf("Failed to get GitHub repository clone URL.\n Please connect account to github if disconnected with `wave-deploy connect-github`")
	}

	args := []string{
		"-n", project.Name,
		"-w", opts.DeploymentOptions.RemoteHomeDir + "/app",
		"-p", appRemoteDirectory,
		"-o", opts.DeploymentOptions.RemoteAppDir,
	}

	if githubCloneUrl != "" {
		args = append(args, "-l", githubCloneUrl)
	}

	if project.BuildCommand != "" {
		args = append(args, "-b", fmt.Sprintf(`"%s"`, project.BuildCommand))
	}

	if project.RunCommand != "" {
		args = append(args, "-s", fmt.Sprintf(`"%s"`, project.RunCommand))
	}

	for key, value := range opts.Envs {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	args = append([]string{dockerFileScriptPath}, args...)
	dockerFileGenerationCommand := "sudo " + strings.Join(args, " ")
	vmSetupCommand := fmt.Sprintf("sudo %s", vmSetupScriptPath)

	_, err = client.Run(vmSetupCommand)
	if err != nil {
		return "", err
	}

	out, err := client.Run(dockerFileGenerationCommand)
	fmt.Println(string(out))
	if err != nil {
		return "", err
	}

	return appRemoteDirectory, err
}

func deployAndStartApplication(opts DeploymentOptions, remoteAppDirectory, projectName string) error {
	client, err := establishSSHConnection(opts)
	if err != nil {
		return err
	}
	defer client.Close()

	deleteStaleContainersCommand := fmt.Sprintf("sudo docker rm $(sudo docker ps -a --filter ancestor=%s -q)", projectName)
	_, _ = client.Run(deleteStaleContainersCommand)

	deleteImageCommand := fmt.Sprintf("sudo docker rmi %s:latest", projectName)
	_, _ = client.Run(deleteImageCommand)

	// Build Docker Image
	_, err = client.Run(
		fmt.Sprintf(
			"sudo docker build -t %s:latest -f %s %s",
			projectName,
			remoteAppDirectory+"/Dockerfile.wavedeploy",
			remoteAppDirectory),
	)
	if err != nil {
		return err
	}

	// Initialize Docker Swarm
	msg, err := client.Run("sudo docker swarm init")
	if err != nil {
		if !(strings.Contains(string(msg), "This node is already part of a swarm")) {
			return err
		} else {
			err = nil
		}
	}

	// Stop any running service
	client.Run(fmt.Sprintf("sudo docker service rm %s", projectName))

	// Deploy Command
	createCmd := fmt.Sprintf("sudo docker service create --name %s --replicas %d --publish 8080:8080 --env PORT=8080", projectName, opts.Replicas)
	// Environment variables map
	for key, value := range opts.Envs {
		createCmd += fmt.Sprintf(" --env %s=%s", key, value)
	}
	createCmd += fmt.Sprintf(" %s:latest", projectName)

	// Deploy service
	_, err = client.Run(createCmd)
	if err != nil {
		return err
	}

	if err = setupAndReloadApiWebServer(client, 8080); err != nil {
		return err
	}

	return err
}

func establishSSHConnection(opts DeploymentOptions) (client *goph.Client, err error) {
	privateKeyBytes, err := os.ReadFile(opts.PrivateKeyPath)
	if err != nil {
		return nil, errors.New("An error occurred reading private key file")
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.New("An error occurred parsing private key")
	}

	client, err = goph.NewConn(&goph.Config{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		User:     opts.VmUser,
		Callback: ssh.InsecureIgnoreHostKey(),
		Timeout:  60 * time.Second,
		Port:     22,
		Addr:     opts.PublicIPV4Addr,
	})
	return
}

func uploadFileToRemoteRecursively(client *goph.Client, localFolder, remoteFolder string) error {
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
			client.Run(fmt.Sprintf("sudo rm -rf %s", remoteFilePath))
			err = client.Upload(path, remoteFilePath)

			if err != nil {
				return err
			}
		}
		return nil
	})
}

func makeRemoteFileExecutableMode(client *goph.Client, path string) (err error) {
	_, err = client.Run(fmt.Sprintf("sudo chmod 0100 %s", path))
	return
}

func getDynamicPort(publicIp string) {

}

func setupAndReloadApiWebServer(client *goph.Client, port int) error {
	rootPath := files.GetCurrentPathRootDirectory()
	templatePath := path.Join(rootPath, "/webserver/tpl/web-server.caddy.tmpl")

	payload := WebServerTmpl{
		EXTERNAL_PORT:           80,
		INTERNAL_LISTENING_PORT: port,
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, payload); err != nil {
		return err
	}

	command := "echo \"" + buffer.String() + "\" | sudo tee /etc/caddy/Caddyfile"

	_, err = client.Run(command)
	if err != nil {
		return err
	}

	_, err = client.Run("sudo systemctl reload caddy")
	if err != nil {
		return err
	}

	return nil
}
