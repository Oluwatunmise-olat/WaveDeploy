package cmd

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/projects"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/hashers"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type ProjectEnvs map[string]string

var projectName string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a project",
	PreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		IsAuthenticated(ctx, "Checking github connection", cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		accountId := cmd.Context().Value("accountId").(string)

		project, err := preDeploymentChecks(accountId)
		if err != nil {
			fmt.Errorf(err.Error())
			return
		}

		initializeDeployment(cmd, project.Id)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	deployCmd.MarkFlagRequired("name")
}

func preDeploymentChecks(accountId string) (*models.Projects, error) {
	projectRecord, err := projects.GetProjectByName(accountId, projectName)

	if err != nil && projectRecord == nil {
		return nil, errors.New(fmt.Sprintf("Project with name(%s) not found", projectName))
	}

	if projectRecord != nil && projectRecord.IsLive {
		return nil, errors.New("Project already deployed. If you want to restart the project run `wave-deploy restart -n <project name>`")
	}

	return projectRecord, nil
}

// TODO:: Major refactors once bearest minimum impl is completed
func initializeDeployment(cmd *cobra.Command, projectId uuid.UUID) {
	accountId := cmd.Context().Value("accountId").(string)
	var updatedPayload models.Projects

	hasCustomBuildPromptCmd := Prompt{label: "> Custom build command (y/n)?: "}
	hasCustomRunPromptCmd := Prompt{label: "> Custom run command (y/n)?: "}
	setEnvsPromptCmd := Prompt{label: "> Set Envs (y/n)?: "}

	buildPromptCmd := Prompt{errorMessage: "Please provide a build command", label: "Build command: "}
	runPromptCmd := Prompt{errorMessage: "Please provide a run command", label: "Run command: "}

	var buildCommand string
	var runCommand string
	var Envs ProjectEnvs = make(ProjectEnvs)

	buildValue := GetPromptInput(hasCustomBuildPromptCmd, nil)
	if buildValue == "y" {
		buildCommand = GetPromptInput(buildPromptCmd, nil)
	}

	runValue := GetPromptInput(hasCustomRunPromptCmd, nil)
	if runValue == "y" {
		runCommand = GetPromptInput(runPromptCmd, nil)
	}

	setEnvValue := GetPromptInput(setEnvsPromptCmd, nil)
	if setEnvValue == "y" {
		var rePromptForEnv bool = true

		for ok := true; ok; ok = rePromptForEnv != false {
			envKeyPromptCmd := Prompt{label: "Enter the environment variable key: "}
			envKeyCommand := GetPromptInput(envKeyPromptCmd, nil)

			envValuePromptCmd := Prompt{
				label: "Enter the environment variable value: ",
				mask:  '*',
			}
			envValueCommand := GetPromptInput(envValuePromptCmd, nil)

			Envs[envKeyCommand] = envValueCommand

			rePromptCmd := Prompt{label: "> Set More Envs (y/n)?: "}
			rePromptEnvValue := GetPromptInput(rePromptCmd, nil)

			if rePromptEnvValue == "n" {
				rePromptForEnv = false
			}
		}
	}

	if buildCommand != "" {
		updatedPayload.BuildCommand = buildCommand
	}

	if runCommand != "" {
		updatedPayload.RunCommand = runCommand
	}

	var envsPayload []models.Envs
	if Envs != nil {
		accountUUID, _ := uuid.Parse(accountId)

		for key, value := range Envs {
			encryptedValue, err := hashers.EncryptIt(value, os.Getenv("APP_KEY"))
			if err != nil {
				fmt.Errorf(err.Error())
				os.Exit(1)
			}

			envRecord := models.Envs{
				Key:       key,
				Value:     encryptedValue,
				Id:        random.GetUUID(),
				AccountId: accountUUID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ProjectId: projectId,
			}

			envsPayload = append(envsPayload, envRecord)
		}
	}

	accountIdToUUId, _ := uuid.Parse(accountId)
	if err := projects.UpdateProjectAndCreateEnvs(projects.UpdateProjectAndCreateEnvsPayload{
		Envs:                 envsPayload,
		UpdateProjectPayload: updatedPayload,
		AccountId:            accountIdToUUId,
		ProjectId:            projectId,
	}); err != nil {
		fmt.Println(err)
		return
	}

	// Containerize application
	// -> create a dockerfile ✅
	// -> add public key to .authorized-keys. How about *fingerprint?
	// -> scp docker file and setup vm (install docker and all needed deps)
	// -> Orchestrate application with or without replicas (docker swarm)

}
