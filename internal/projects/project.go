package projects

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"github.com/google/uuid"
	"time"
)

type UpdateProjectAndCreateEnvsPayload struct {
	Envs                 []models.Envs
	AccountId            uuid.UUID
	ProjectId            uuid.UUID
	UpdateProjectPayload models.Projects
}

func IsProjectNameTaken(accountId string, projectName string) (bool, error) {
	projectRepository := respository.ProjectsRepository{}
	return projectRepository.ProjectExistsWithName(projectName, accountId)
}

func CreateProject(accountId string, projectName string, ghRepo *structs.GithubAInstallationRepositories) error {
	projectRepository := respository.ProjectsRepository{}

	project := models.Projects{
		Id:            random.GetUUID(),
		IsLive:        false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		AccountId:     accountId,
		Name:          projectName,
		GithubBranch:  ghRepo.DefaultBranch,
		GithubRepoUrl: ghRepo.GitUrl,
		GithubCommit:  "",
	}

	if err := projectRepository.CreateProject(project); err != nil {
		return errors.New("An error occurred creating project")
	}

	return nil
}

func UpdateProjectAndCreateEnvs(upace UpdateProjectAndCreateEnvsPayload) error {
	dbTransaction := respository.DBTransaction()
	projectRepository := respository.ProjectsRepository{}
	envRepository := respository.EnvsRepository{}

	if err := projectRepository.UpdateProject(respository.UpdateProjectPayload{
		ProjectId: upace.ProjectId,
		AccountId: upace.AccountId,
		Project:   upace.UpdateProjectPayload,
		Trx:       dbTransaction,
	}); err != nil {
		fmt.Println(err)
		return errors.New("An error occurred updating project")
	}

	if len(upace.Envs) != 0 {
		if err := envRepository.CreateMultipleProjectEnvs(upace.Envs, dbTransaction); err != nil {
			return errors.New("An error occurred setting project envs")
		}
	}

	return nil
}

func GetProjectByName(accountId, projectName string) (*models.Projects, error) {
	projectRepository := respository.ProjectsRepository{}
	return projectRepository.GetProjectByNameAndAccount(projectName, accountId)
}

func UpdateProject(updateProjectPayload models.Projects, projectId, accountId uuid.UUID) error {
	projectRepository := respository.ProjectsRepository{}
	return projectRepository.UpdateProject(respository.UpdateProjectPayload{
		ProjectId: projectId,
		AccountId: accountId,
		Project:   updateProjectPayload,
	})
}

func GetProjectById(projectId, accountId string) (*models.Projects, error) {
	projectRepository := respository.ProjectsRepository{}
	return projectRepository.GetProjectById(projectId, accountId)
}

func GetProjectEnvs(projectId, accountId uuid.UUID) ([]models.Envs, error) {
	envRepository := respository.EnvsRepository{}
	return envRepository.GetEnvs(projectId, accountId)
}

func DeleteProjectEnvs(projectId, accountId uuid.UUID) error {
	envRepository := respository.EnvsRepository{}
	return envRepository.DeleteEnvs(projectId, accountId)
}

func CreateBatchProjectEnvs(envs []models.Envs) error {
	envRepository := respository.EnvsRepository{}
	if err := envRepository.CreateMultipleProjectEnvs(envs, nil); err != nil {
		return errors.New("An error occurred setting project envs")
	}

	return nil
}
