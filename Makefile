.PHONY: all
all: help

default: help

COMPOSE_PROJECT_NAME ?= mes

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## format, test, and package the Spring Boot monolith
	@mvn -Dmaven.repo.local=.m2/repository -f app/backend/pom.xml clean spotless:check test package

.PHONY: test
test: ## run monolith tests
	@mvn -Dmaven.repo.local=.m2/repository -f app/backend/pom.xml test

.PHONY: format
format: ## format monolith Java sources
	@mvn -Dmaven.repo.local=.m2/repository -f app/backend/pom.xml spotless:apply

##@ Build

.PHONY: watch-frontend
watch-frontend:
	@cd app/frontend && pnpm dev --host 0.0.0.0

##@ Development Env

.PHONY: env-start
env-start:  ## launch all middleware software as the docker
	@docker-compose -p $(COMPOSE_PROJECT_NAME) up -d

.PHONY: env-stop
env-stop: ## stop all docker
	@docker-compose -p $(COMPOSE_PROJECT_NAME) down

.PHONY: clean
clean: ## clean monolith/frontend build output
	@mvn -Dmaven.repo.local=.m2/repository -f app/backend/pom.xml clean
	@cd app/frontend && rm -rf dist .vite

##@ Open Browser

.PHONY: open
open: ## open the frontend in the default browser
	@open "http://localhost:5173/"

.PHONY: monolith
monolith: ## build and run the Spring Boot monolith
	@./start.sh
