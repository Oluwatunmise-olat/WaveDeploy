package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
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

	err := dbExecutor.Create(env).Error
	return err
}

func (er *EnvsRepository) CreateMultipleProjectEnvs(envs []models.Envs, trx *gorm.DB) error {
	dbExecutor := er.initializeEnvsRepository().DB
	if trx != nil {
		dbExecutor = trx
	}

	err := dbExecutor.Create(envs).Error
	return err
}
