.PHONY: db-migrate-make db-migrate-up db-migrate-down help

define HELP_MESSAGE
Usage:
	make db-migrate-make migration_filename=<filename>
Example:
    make migrate migration_name=create_table_accounts
endef
export HELP_MESSAGE

help:
	@echo "$$HELP_MESSAGE"
	@echo ""
	@echo "Available targets:"
	@egrep -E '^[a-zA-Z0-9_.-]+:' $(MAKEFILE_LIST) | grep -vE 'help:|Example:|Usage:|.PHONY:' | awk 'BEGIN {FS = ":*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

migration_filename :=

db-migrate-make:
	@if [ -z "$(migration_filename)" ]; then \
  		echo "Error: migration_filename is required."; \
        echo "$$HELP_MESSAGE"; \
        exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(migration_filename)

db-migrate-up:

db-migrate-down: