package projects

import (
	"errors"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/respository"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/random"
	"github.com/Oluwatunmise-olat/WaveDeploy/pkg/structs"
	"time"
)

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
