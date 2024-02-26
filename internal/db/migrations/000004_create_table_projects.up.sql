CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    `name` VARCHAR(500) NOT NULL,
    github_repo_url VARCHAR(5000) NULL,
    github_branch VARCHAR(5000) NULL,
    github_commit VARCHAR(5000) NULL,
    is_live BOOL DEFAULT FALSE,
    replicas INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX projects_account_id_index ON projects(account_id);