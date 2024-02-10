CREATE TABLE IF NOT EXISTS project_artifacts (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    tag VARCHAR(2000) NOT NULL,
    is_live BOOL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX project_artifacts_project_id_index ON project_artifacts(project_id);
CREATE INDEX project_artifacts_account_id_index ON project_artifacts(account_id);