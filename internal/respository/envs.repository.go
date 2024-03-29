package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
)

var (
	envsRepository        *EnvsRepository
	envRepositoryInitOnce sync.Once
)

func (ar *EnvsRepository) initializeEnvsRepository() *EnvsRepository {
	envRepositoryInitOnce.Do(func() {
		envsRepository = &EnvsRepository{
			DB: db.DB,
		}
	})

	return envsRepository
}

func (er *EnvsRepository) CreateProjectEnvs(env models.Envs, trx *gorm.DB) error {
	dbExecutor := er.initializeEnvsRepository().DB
	if trx != nil {
		dbExecutor = trx
	}

	err := dbExecutor.Model(&models.Envs{}).Create(env).Error
	return err
}

func (er *EnvsRepository) CreateMultipleProjectEnvs(envs []models.Envs, trx *gorm.DB) error {
	dbExecutor := er.initializeEnvsRepository().DB
	if trx != nil {
		dbExecutor = trx
	}

	err := dbExecutor.Model(&models.Envs{}).Create(envs).Error
	return err
}

func (er *EnvsRepository) GetEnvs(projectId, accountId uuid.UUID) ([]models.Envs, error) {
	var envs []models.Envs

	err := er.initializeEnvsRepository().
		DB.Model(&models.Envs{}).
		Where("account_id = ? and project_id = ? and deleted_at is null", accountId, projectId).
		Find(&envs).
		Error

	return envs, err
}

func (er *EnvsRepository) DeleteEnvs(projectId, accountId uuid.UUID, trx *gorm.DB) error {
	dbExecutor := er.initializeEnvsRepository().DB
	if trx != nil {
		dbExecutor = trx
	}

	var envs models.Envs
	err := dbExecutor.
		Model(&models.Envs{}).
		Where("account_id = ? and project_id = ?", accountId, projectId).
		Unscoped().
		Delete(&envs).
		Error

	return err
}
