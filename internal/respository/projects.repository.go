package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"strings"
	"sync"
)

var (
	projectsRepository        *ProjectsRepository
	projectRepositoryInitOnce sync.Once
)

func (pr *ProjectsRepository) initializeProjectsRepository() *ProjectsRepository {
	projectRepositoryInitOnce.Do(func() {
		projectsRepository = &ProjectsRepository{
			DB: db.DB,
		}
	})

	return projectsRepository
}

func (pr *ProjectsRepository) ProjectExistsWithName(projectName string, accountId string) (bool, error) {
	var project models.Projects
	err := pr.
		initializeProjectsRepository().
		DB.
		Where("account_id = ? AND name = ?", accountId, strings.ToLower(projectName)).
		First(&project).
		Select("id").
		Error

	return project.Name != "", err
}

func (pr *ProjectsRepository) CreateProject(project models.Projects) error {
	err := pr.initializeProjectsRepository().DB.Create(project).Error
	return err
}

func (pr *ProjectsRepository) UpdateProject(up UpdateProjectPayload) error {
	dbExecutor := pr.initializeProjectsRepository().DB
	if up.Trx != nil {
		dbExecutor = up.Trx
	}

	err := dbExecutor.
		Where("account_id = ? and id = ?", up.AccountId, up.ProjectId).
		Updates(up.Project).Error

	return err
}

func (pr *ProjectsRepository) GetProjectByNameAndAccount(name, accountId string) (*models.Projects, error) {
	var project models.Projects
	err := pr.initializeProjectsRepository().DB.Where("account_id = ? and name = ?", accountId, name).First(&project).Error
	return &project, err
}
