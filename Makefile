.PHONY: db-migrate-make db-migrate-up db-migrate-down help

define HELP_MESSAGE
Usage:
	make db-migrate-make migration_filename=<filename>
	make db-migrate-up mysql_username=<mysql_username> mysql_password=<mysql_password>
	make db-migrate-down mysql_username=<mysql_username> mysql_password=<mysql_password>
Example:
    make migrate migration_name=create_table_accounts
    make db-migrate-up mysql_username=test mysql_password=test_password
    make db-migrate-down mysql_username=test mysql_password=test_password
endef
export HELP_MESSAGE

help:
	@echo "$$HELP_MESSAGE"
	@echo ""
	@echo "Available targets:"
	@egrep -E '^[a-zA-Z0-9_.-]+:' $(MAKEFILE_LIST) | grep -vE 'help:|Example:|Usage:|.PHONY:' | awk 'BEGIN {FS = ":*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

migration_filename :=
mysql_username :=
mysql_password :=

db-migrate-make:
	@if [ -z "$(migration_filename)" ]; then \
  		echo "Error: migration_filename is required."; \
        echo "$$HELP_MESSAGE"; \
        exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(migration_filename)

db-migrate-up:
	@if [ -z "$(mysql_username)" ] || [ -z "$(mysql_password)" ]; then \
  		echo "Error: mysql_username and mysql_password are required."; \
        echo "$$HELP_MESSAGE"; \
        exit 1; \
	fi
	migrate -path db/migrations -database "mysql://$(mysql_username):$(mysql_password)@tcp(localhost:3324)/wave_deploy_db" -verbose up

db-migrate-down:
	@if [ -z "$(mysql_username)" ] || [ -z "$(mysql_password)" ]; then \
  		echo "Error: mysql_username and mysql_password are required."; \
        echo "$$HELP_MESSAGE"; \
        exit 1; \
	fi
	migrate -path db/migrations -database "mysql://$(mysql_username):$(mysql_password)@tcp(localhost:3324)/wave_deploy_db" -verbose down