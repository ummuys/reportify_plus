
####################################################

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
SCRIPTS_PATH=scripts
UPDATE_SECRETS=$(SCRIPTS_PATH)/update_secrets.sh
COMPOSE=docker-compose.yaml
ENV_PATH=../.env
BASE_COMMAND=-docker compose --env-file $(ENV_PATH)
GOBIN := $(shell go env GOPATH)/bin
LINT_BIN := $(GOBIN)/golangci-lint

####################################################

.PHONY: up update-up down update status drop-all-volume code-review ci-lint ci-test ci-race ci-docker-build test

up:
	@./$(UPDATE_SECRETS)
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

update-up:
	@./$(UPDATE_SECRETS)
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@$(BASE_COMMAND) -f $(COMPOSE) up -d --build


down:
	@$(BASE_COMMAND) -f $(COMPOSE) down

update:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

status:
	@$(BASE_COMMAND) -f $(COMPOSE) ps -a

drop-all-volume:
	@$(BASE_COMMAND) -f $(COMPOSE) down -v


# CI

code-review:
	@cd $(ROOT_DIR) && ./deployment/scripts/pre_commit.sh

ci-lint:
	@echo "Running golangci-lint..."
	@if [ ! -x "$(LINT_BIN)" ]; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2; \
	fi
	@cd $(ROOT_DIR) && "$(LINT_BIN)" run ./...

ci-test:
	@echo "Running unit tests for auth/report/report-worker services..."
	@cd $(ROOT_DIR) && go test -count=1 ./services/auth/internal/service/... ./services/report/internal/service/... ./services/report-worker/internal/service/...

ci-race:
	@echo "Running race tests for auth/report/report-worker services..."
	@cd $(ROOT_DIR) && go test -race -count=1 ./services/auth/internal/service/... ./services/report/internal/service/... ./services/report-worker/internal/service/...

ci-docker-build:
	@echo "Building report image..."
	@cd $(ROOT_DIR) && docker build -f deployment/Dockerfile --target report -t reportify-report:ci .

test:
	@echo "Testing AuthService . . ."
	@richgo test ../services/auth/internal/service/... -cover

	@echo "Testing ReportService . . ."
	@richgo test ../services/report/internal/service/... -cover
	@echo "Testing ReportWorkerService . . ."
	@richgo test ../services/report-worker/internal/service/... -coverprofile=worker_cover.out

