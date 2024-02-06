CREATE TABLE IF NOT EXISTS github_apps (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    installation_id VARCHAR(250) NOT NULL,
    code VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX github_apps_account_id_index ON github_apps(account_id);
CREATE INDEX github_apps_installation_id_index ON github_apps(installation_id);