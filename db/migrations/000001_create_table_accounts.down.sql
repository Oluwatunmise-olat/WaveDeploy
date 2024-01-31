CREATE TABLE IF NOT EXISTS "accounts" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    "username" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "last_auth_at" TIMESTAMP,
    "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT (now()),
    "deleted_at" TIMESTAMP NULL
);