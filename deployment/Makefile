
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

down: 
	@$(BASE_COMMAND) -f $(COMPOSE) down

update:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@$(BASE_COMMAND) -f $(COMPOSE) up -d

status:
	@$(BASE_COMMAND) -f $(COMPOSE) ps -a

drop-all-volume:
	@$(BASE_COMMAND) -f $(COMPOSE) down -v
	

update-gateway-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm gateway-service

update-auth-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm auth-service

update-report-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm report-service

update-report-worker-image:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm report-worker

update-all-images:
	@$(BASE_COMMAND) -f $(COMPOSE) down
	@-docker image rm auth-service
	@-docker image rm gateway-service
	@-docker image rm report-service
	@-docker image rm report-worker



update-all-images: update-auth-image update-gateway-image

# CI

code-review:
	@cd .. && ./deployment/scripts/pre-commit.sh

code-auto-format:
	@cd .. && gofumpt -w $$(find . -name '*.go' -type f)