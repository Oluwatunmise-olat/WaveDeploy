package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"strings"
)

var envUpdateProjectName string

var updateEnvsCmd = &cobra.Command{
	Use:   "update-envs",
	Short: "Update project envs",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		accountId := getAccountID(cmd)

		envs, err := updateProjectEnvs(accountId)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := reDeployProject(DeploymentOptions{
			Envs: envs,
		}); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateEnvsCmd)
	updateEnvsCmd.Flags().StringVarP(&envUpdateProjectName, "name", "n", "", "Project name")
	updateEnvsCmd.MarkFlagRequired("name")
}

func updateProjectEnvs(accountId string) (ProjectEnvs, error) {
	project, err := projects.GetProjectByName(accountId, strings.TrimSpace(envUpdateProjectName))

	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	accountUUID, _ := uuid.Parse(project.AccountId)
	if !project.IsLive {
		return nil, errors.New("can only update deployed project envs")
	}

	err = projects.DeleteProjectEnvs(project.Id, accountUUID)
	if err != nil {
		return nil, err
	}

	envs, err := promptForEnvVariables()
	if err != nil {
		return nil, err
	}
	envsPayload, err := createEnvRecords(accountId, project, envs)
	if err != nil {
		return nil, err
	}

	err = projects.CreateBatchProjectEnvs(envsPayload)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func reDeployProject(opts DeploymentOptions) error {
	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()

	opts.PublicIPV4Addr = ipv4Addr
	opts.VmUser = vmUser
	opts.PrivateKeyPath = privateKeyPath

	client, err := establishSSHConnection(opts)
	if err != nil {
		return err
	}
	defer client.Close()

	redeployCommand := fmt.Sprintf("sudo docker service update --name %s", envUpdateProjectName)
	for key, value := range opts.Envs {
		redeployCommand += fmt.Sprintf(" --env-add %s=%s", key, value)
	}
	_, err = client.Run(redeployCommand)
	if err != nil {
		return err
	}

	// Update One replica at a time
	rollingUpdatesCommand := fmt.Sprintf("sudo docker service --force --update-parallelism 1 --update-delay %s %s", "5s", projectName)
	_, err = client.Run(rollingUpdatesCommand)
	if err != nil {
		return err
	}

	return nil
}
