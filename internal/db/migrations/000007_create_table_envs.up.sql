CREATE TABLE IF NOT EXISTS envs (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    `key` TEXT NOT NULL,
    value TEXT NOT NULL, -- note: this are encrypted with varying keys per account
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX envs_account_id_index ON envs(account_id);
CREATE INDEX envs_project_id_index ON envs(project_id);