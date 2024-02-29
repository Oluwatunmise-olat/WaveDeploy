package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var updateEnvsCmd = &cobra.Command{
	Use:   "update-envs",
	Short: "Update project envs",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkGitHubConnection(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		accountId := getAccountID(cmd)
		projectName := getProjectName(cmd)

		s := initializeSpinner("Updating Application Envs ", "\n")
		envs, err := updateProjectEnvs(accountId, projectName)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error occurred updating project envs: %w", err)
		}

		if err := reDeployProject(DeploymentOptions{
			Envs: envs,
		}, projectName, s); err != nil {
			s.Stop()
			return fmt.Errorf("error occurred redeploying project after envs update: %w", err)
		}
		s.FinalMSG = "Application Envs Updated 🪐\n"
		s.Stop()
		return nil
	},
	SilenceUsage: true,
	Example:      "wave-deploy update-envs -n <PROJECT NAME>",
}

func init() {
	rootCmd.AddCommand(updateEnvsCmd)
	updateEnvsCmd.Flags().StringP("name", "n", "", "Project Name")
	updateEnvsCmd.MarkFlagRequired("name")
}

func updateProjectEnvs(accountId, projectName string) (ProjectEnvs, error) {
	project, err := projects.GetProjectByName(accountId, projectName)

	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	accountUUID, _ := uuid.Parse(project.AccountId)
	if !project.IsLive {
		return nil, errors.New("can only update deployed project envs")
		//return nil, nil
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

func reDeployProject(opts DeploymentOptions, projectName string, s *spinner.Spinner) error {
	vmUser, ipv4Addr, privateKeyPath := promptDeploymentCredentialsDetails()

	opts.PublicIPV4Addr = ipv4Addr
	opts.VmUser = vmUser
	opts.PrivateKeyPath = privateKeyPath

	s.Start()
	client, err := establishSSHConnection(opts)
	if err != nil {
		return err
	}
	defer client.Close()

	// Update One replica at a time
	redeployCommand := fmt.Sprintf("sudo docker service update %s --force --update-parallelism 1 --update-delay %s", projectName, "5s")
	for key, value := range opts.Envs {
		redeployCommand += fmt.Sprintf(" --env-add %s=%s", key, value)
	}
	_, err = client.Run(redeployCommand)
	if err != nil {
		return err
	}

	return nil
}
