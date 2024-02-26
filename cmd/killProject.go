package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd/flags"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/google/uuid"
	"strings"

	"github.com/spf13/cobra"
)

var killProjectCmd = &cobra.Command{
	Use:   "kill-project",
	Short: "Stop a deployed project",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		accountId := getAccountID(cmd)
		if err := killProject(accountId); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Project Killed 🫸🏾🫷🏾")
	},
}

func init() {
	rootCmd.AddCommand(killProjectCmd)
	flags.InitializeProjectNameFlag(updateEnvsCmd)
}

func killProject(accountId string) error {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(flags.GetProjectName()))
	accountUUID, _ := uuid.Parse(accountId)

	if err != nil {
		return fmt.Errorf("project not found")
	}

	if !project.IsLive {
		return errors.New("Project is not deployed")
	}

	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()
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

	deleteImageCommand := fmt.Sprintf("sudo docker rmi %s:latest", "")
	_, _ = client.Run(deleteImageCommand)

	projects.UpdateProject(models.Projects{
		IsLive: false,
	}, project.Id, accountUUID)

	return nil
}
