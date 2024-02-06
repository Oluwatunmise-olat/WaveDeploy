package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"gorm.io/gorm"
	"sync"
)

var (
	githubAppsRepository         *GithubAppsRepository
	githubAppsRepositoryInitOnce sync.Once
)

type GithubAppsRepository struct {
	DB *gorm.DB
}

func (gar *GithubAppsRepository) initializeGithubAppsRepository() *GithubAppsRepository {
	githubAppsRepositoryInitOnce.Do(func() {
		githubAppsRepository = &GithubAppsRepository{
			DB: db.DB,
		}
	})

	return githubAppsRepository
}

func (gar *GithubAppsRepository) GetGithubAppByInstallationId(installationId string) (*models.GithubApps, error) {
	var githubApp models.GithubApps
	err := gar.initializeGithubAppsRepository().DB.First(&githubApp, "installation_id = ?", installationId).Error
	return &githubApp, err
}

func (gar *GithubAppsRepository) GetGithubAppByAccountId(accountId string) (*models.GithubApps, error) {
	var githubApp models.GithubApps
	err := gar.initializeGithubAppsRepository().DB.First(&githubApp, "account_id = ?", accountId).Error
	return &githubApp, err
}

func (gar *GithubAppsRepository) CreateGithubApp(githubApps models.GithubApps) error {
	err := gar.initializeGithubAppsRepository().DB.Create(githubApps).Error
	return err
}
