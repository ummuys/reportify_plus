
####################################################

GREEN  := \033[1;32m
RESET  := \033[0m


SCRIPTS_PATH=scripts
UPDATE_SECRETS=$(SCRIPTS_PATH)/update_secrets.sh
COMPOSE=docker-compose.yaml
ENV_PATH=../.env
BASE_COMMAND=-docker compose --env-file $(ENV_PATH)

####################################################


up: 
	@./$(UPDATE_SECRETS)
	@-docker network create reportify_network
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

down: 
	@$(BASE_COMMAND) -f $(COMPOSE) down

update:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

status:
	@$(BASE_COMMAND) -f $(COMPOSE) ps -a

drop-auth-volume:
	@$(BASE_COMMAND) -f $(COMPOSE) down -v
	

update-gateway-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm gateway-service
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

update-auth-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm auth-service
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

update-report-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm docker-app
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

update-report-worker-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm docker-app
	@$(BASE_COMMAND) -f $(COMPOSE) up -d
	



# 			TESTING			# 
srv-test:
	@echo "$(GREEN)[TEST]$(RESET) Running tests..."
	@gotestsum --format=github-actions ./../internal/service

test-all: tools fmt vet lint race vuln
	@echo "$(GREEN)[OK]$(RESET) No problema :)"

tools:
	@echo "$(GREEN)[SETUP]$(RESET) Installing Go tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(GREEN)[OK]$(RESET) All tools installed into $(TOOLS_DIR)"

fmt:
	@echo "$(GREEN)[TEST]$(RESET) Formatting code..."
	@go fmt ./../...
	@goimports -w .

vet:
	@echo "$(GREEN)[TEST]$(RESET) Running go vet..."
	@go vet ./../...

lint:
	@echo "$(GREEN)[TEST]$(RESET) Running lint..."
	@golangci-lint run ./../...

race:
	@echo "$(GREEN)[TEST]$(RESET) Running race cheking..."
	@go test ./../... -race -count=1

vuln:
	@echo "$(GREEN)[TEST]$(RESET) Running vuln..."
	@govulncheck ./../...



# CI

code-review:
	@cd .. && ./deployment/scripts/pre-commit.sh

code-auto-format:
	@cd .. && gofumpt -w $$(find . -name '*.go' -type f)