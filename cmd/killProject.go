package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/files"
	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"github.com/melbahja/goph"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var killProjectCmd = &cobra.Command{
	Use:   "kill-project",
	Short: "Stop a deployed project",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		projectName := getProjectName(cmd)

		s := initializeSpinner("killing Project ", "")
		if err := killProject(accountId, projectName, s); err != nil {
			s.Stop()
			return fmt.Errorf("error occurred halting project: %w", err)
		}
		s.Stop()
		fmt.Println("Project Killed 🫸🏾🫷🏾")
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy kill-project -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(killProjectCmd)

	killProjectCmd.Flags().StringP("name", "n", "", "Project Name")
	killProjectCmd.MarkFlagRequired("name")
}

func killProject(accountId, projectName string, s *spinner.Spinner) error {
	project, err := projects.GetProjectByName(accountId, projectName)
	accountUUID, _ := uuid.Parse(accountId)

	if err != nil {
		return fmt.Errorf("project not found")
	}

	if !project.IsLive {
		return errors.New("Project is not deployed")
	}

	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()

	if s != nil {
		s.Start()
	}

	client, err := establishSSHConnection(DeploymentOptions{
		VmUser:         vmUser,
		PublicIPV4Addr: ipv4Addr,
		PrivateKeyPath: privateKeyPath,
	})

	if err != nil {
		return err
	}
	defer client.Close()

	killCommand := fmt.Sprintf("sudo docker service rm %s", project.Name)
	_, _ = client.Run(killCommand)

	deleteStaleContainersCommand := fmt.Sprintf("sudo docker rm $(sudo docker ps -a --filter ancestor=%s -q)", "projectName")
	_, _ = client.Run(deleteStaleContainersCommand)

	deleteImageCommand := fmt.Sprintf("sudo docker rmi %s:latest", project.Name)
	_, _ = client.Run(deleteImageCommand)

	// Cleanup app directory
	_, _ = client.Run(fmt.Sprintf("sudo rm -rf /home/%s/app", vmUser))

	updateData := map[string]interface{}{
		"is_live": false,
	}

	_ = projects.UpdateProject(updateData, project.Id, accountUUID)

	// revert default caddy config
	if err = revertToDefaultCaddyConfig(client); err != nil {
		return err
	}

	return nil
}

func revertToDefaultCaddyConfig(client *goph.Client) error {
	rootPath := files.GetCurrentPathRootDirectory()
	configPath := path.Join(rootPath, "/webserver/Caddyfile")

	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return errors.New("Error occurred reading default caddy config file")
	}

	command := fmt.Sprintf("echo \"%s\" | sudo tee /etc/caddy/Caddyfile", configContent)

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
