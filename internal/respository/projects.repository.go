package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var (
	projectsRepository        *ProjectsRepository
	projectRepositoryInitOnce sync.Once
)

type ProjectsRepository struct {
	DB *gorm.DB
}

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
