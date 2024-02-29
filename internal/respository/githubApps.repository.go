package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"sync"
)

var (
	githubAppsRepository         *GithubAppsRepository
	githubAppsRepositoryInitOnce sync.Once
)

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
	err := gar.initializeGithubAppsRepository().DB.Model(&models.GithubApps{}).First(&githubApp, "installation_id = ? and deleted_at is null", installationId).Error
	return &githubApp, err
}

func (gar *GithubAppsRepository) GetGithubAppInstallationIdByAccountId(accountId string) (string, error) {
	var githubApp models.GithubApps
	err := gar.initializeGithubAppsRepository().DB.Model(&models.GithubApps{}).First(&githubApp, "account_id = ? and deleted_at is null", accountId).Select("installation_id").Error
	return githubApp.InstallationId, err
}

func (gar *GithubAppsRepository) GetGithubAppByAccountId(accountId string) (*models.GithubApps, error) {
	var githubApp models.GithubApps
	err := gar.initializeGithubAppsRepository().DB.Model(&models.GithubApps{}).First(&githubApp, "account_id = ? and deleted_at is null", accountId).Error
	return &githubApp, err
}

func (gar *GithubAppsRepository) CreateGithubApp(githubApps models.GithubApps) error {
	err := gar.initializeGithubAppsRepository().DB.Model(&models.GithubApps{}).Create(githubApps).Error
	return err
}
