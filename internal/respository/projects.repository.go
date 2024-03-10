package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		DB.Model(&models.Projects{}).
		Where("account_id = ? AND name = ? and deleted_at is null", accountId, strings.ToLower(projectName)).
		First(&project).
		Select("id").
		Error

	return project.Name != "", err
}

func (pr *ProjectsRepository) CreateProject(project models.Projects) error {
	err := pr.initializeProjectsRepository().DB.Model(&models.Projects{}).Create(project).Error
	return err
}

func (pr *ProjectsRepository) UpdateProject(up UpdateProjectPayload) error {
	dbExecutor := pr.initializeProjectsRepository().DB
	if up.Trx != nil {
		dbExecutor = up.Trx
	}

	err := dbExecutor.Model(&models.Projects{}).
		Where("account_id = ? and id = ?", up.AccountId, up.ProjectId).
		Updates(up.Project).
		Error

	return err
}

func (pr *ProjectsRepository) DeleteProject(projectId, accountId uuid.UUID, trx *gorm.DB) error {
	dbExecutor := pr.initializeProjectsRepository().DB
	if trx != nil {
		dbExecutor = trx
	}

	err := dbExecutor.Model(&models.Projects{}).Unscoped().Delete("account_id = ? and project_id = ?", accountId, projectId).Error
	return err
}

func (pr *ProjectsRepository) GetProjectByNameAndAccount(name, accountId string) (*models.Projects, error) {
	var project models.Projects

	err := pr.initializeProjectsRepository().
		DB.Model(&models.Projects{}).
		Where("account_id = ? and name = ? and deleted_at is null", accountId, name).
		First(&project).Error
	return &project, err
}

func (pr *ProjectsRepository) GetProjectById(projectId, accountId string) (*models.Projects, error) {
	var project models.Projects
	err := pr.initializeProjectsRepository().
		DB.Model(&models.Projects{}).
		Where("account_id = ? and id = ? and deleted_at is null", accountId, projectId).
		First(&project).Error

	return &project, err
}

func (pr *ProjectsRepository) GetAllProjectsByAccountId(accountId string) ([]models.Projects, error) {
	var projects []models.Projects

	err := pr.initializeProjectsRepository().
		DB.Model(&models.Projects{}).
		Select("name", "github_repo_url", "is_live", "id").
		Where("account_id = ?", accountId).
		Find(&projects).
		Error

	return projects, err
}
