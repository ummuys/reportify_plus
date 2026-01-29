
####################################################

SCRIPTS_PATH=scripts
UPDATE_SECRETS=$(SCRIPTS_PATH)/update_secrets.sh
COMPOSE=docker-compose.yaml
ENV_PATH=../.env
BASE_COMMAND=-docker compose --env-file $(ENV_PATH)

####################################################


up: 
	@./$(UPDATE_SECRETS)
	@-docker network create reportify-network
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
	@cd .. && ./deployment/scripts/pre_commit.sh

code-auto-format:
	@cd .. && gofumpt -w $$(find . -name '*.go' -type f)


test:
	@echo "Testing AuthService . . ."
	@richgo test ../services/auth/internal/service/... -cover

	@echo "Testing ReportService . . ."
	@richgo test ../services/report/internal/service/... -cover
	@echo "Testing ReportWorkerService . . ."
	@richgo test ../services/report-worker/internal/service/... -coverprofile=worker_cover.out