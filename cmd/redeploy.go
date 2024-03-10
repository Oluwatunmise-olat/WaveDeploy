package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/account"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/github"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var redeployCmd = &cobra.Command{
	Use:   "redeploy",
	Short: "Redeploy a project",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		projectName := getProjectName(cmd)

		project, err := preReDeploymentChecks(accountId, projectName)
		if err != nil {
			return fmt.Errorf("error occurred redeploying project: %w", err)
		}

		s := initializeSpinner("Redeploying Application ", "\n")
		if err := redeployProject(project, s); err != nil {
			s.Stop()
			return fmt.Errorf("error occurred redeploying project: %w", err)
		}
		s.FinalMSG = "Application Redeployed Successfully 🦍\n"
		s.Stop()
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy redeploy -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(redeployCmd)
	redeployCmd.Flags().StringP("name", "n", "", "Project Name")
	redeployCmd.MarkFlagRequired("name")
}

func preReDeploymentChecks(accountId, projectName string) (*models.Projects, error) {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(projectName))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project: %v", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project with name '%s' not found", projectName)
	}
	if !project.IsLive {
		return nil, errors.New("Cannot redeploy project that is not already deployed")
	}
	return project, nil
}

func redeployProject(project *models.Projects, s *spinner.Spinner) error {
	accountUUID, _ := uuid.Parse(project.AccountId)

	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()
	client, err := establishSSHConnection(DeploymentOptions{
		VmUser:         vmUser,
		PublicIPV4Addr: ipv4Addr,
		PrivateKeyPath: privateKeyPath,
	})
	remoteHomeDir := fmt.Sprintf("/home/%s", vmUser)
	remoteAppDirectory := fmt.Sprintf("/home/%s/app/%s", vmUser, project.Name)
	dockerFileScriptPath := fmt.Sprintf("%s/app/.builder/generate-dockerfile.sh", remoteHomeDir)

	if err != nil {
		return err
	}
	defer client.Close()
	s.Start()

	// Update Tag
	updateCurrentImageTagCommand := fmt.Sprintf("sudo docker tag %s:%s %s:backup", project.Name, "latest", project.Name)
	_, _ = client.Run(updateCurrentImageTagCommand)

	installationId, err := account.GetAccountInstallationId(project.AccountId)
	if err != nil {
		return fmt.Errorf("Failed to get GitHub installation ID: %v", err)
	}

	githubCloneUrl, err := github.GetRepositoryCloneUrl(installationId, project.GithubRepoUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to github repository.\n Please connect account to github if disconnected with `wave-deploy connect-github`")
	}
	// Pull App Latest Changes
	_, _ = client.Run(fmt.Sprintf("sudo git -C /home/%s/app/%s remote set-url origin %s", vmUser, project.Name, githubCloneUrl))
	_, _ = client.Run(fmt.Sprintf("sudo git -C /home/%s/app/%s pull origin %s", vmUser, project.Name, project.GithubBranch))

	// ReBuild Application
	args := []string{
		"-n", project.Name,
		"-w", remoteHomeDir + "/app",
		"-p", remoteAppDirectory,
		"-o", fmt.Sprintf("%s/app/.builder", remoteHomeDir),
	}

	if project.BuildCommand != "" {
		args = append(args, "-b", fmt.Sprintf(`"%s"`, project.BuildCommand))
	}

	if project.RunCommand != "" {
		args = append(args, "-s", fmt.Sprintf(`"%s"`, project.RunCommand))
	}

	envs, _ := projects.GetProjectEnvs(project.Id, accountUUID)

	if len(envs) > 0 {
		for _, row := range envs {
			value, _ := hashers.DecryptIt(row.Value, os.Getenv("APP_KEY"))

			args = append(args, "-e", fmt.Sprintf("%s=%s", row.Key, value))
		}
	}

	args = append([]string{dockerFileScriptPath}, args...)
	dockerFileGenerationCommand := "sudo " + strings.Join(args, " ")

	_, err = client.Run(dockerFileGenerationCommand)
	if err != nil {
		return err
	}

	// Build Docker Image
	_, err = client.Run(
		fmt.Sprintf(
			"sudo docker build -t %s:latest -f %s %s",
			project.Name,
			remoteAppDirectory+"/Dockerfile.wavedeploy",
			remoteAppDirectory),
	)
	if err != nil {
		return err
	}

	// Redeploy Command
	updateCmd := fmt.Sprintf("sudo docker service update --replicas %d --image %s:%s --update-parallelism 1 --update-delay 5s --force %s", project.Replicas, project.Name, "latest", project.Name)
	projectEnvs, err := projects.GetProjectEnvs(project.Id, accountUUID)

	if len(projectEnvs) > 0 {
		for key, value := range projectEnvs {
			updateCmd += fmt.Sprintf(" --env %s=%s", key, value)
		}
	}

	// If failure, rollback
	_, err = client.Run(updateCmd)
	if err != nil {
		_, err := client.Run(fmt.Sprintf("sudo docker service rollback %s", project.Name))
		if err != nil {
			return err
		}
	}

	_, _ = client.Run("sudo docker container prune -f")

	deleteStaleContainersCommand := fmt.Sprintf("sudo docker rm $(sudo docker ps -a --filter ancestor=%s -q)", project.Name)
	_, _ = client.Run(deleteStaleContainersCommand)

	deletePrevImageCommand := fmt.Sprintf("sudo docker rmi %s:backup", project.Name)
	_, _ = client.Run(deletePrevImageCommand)

	return nil
}
