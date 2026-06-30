include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d go-rest-prod-postgres

env-down:
	@docker compose down go-rest-prod-postgres

env-cleanup:
	@read -p "Очистить все volumes? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down go-rest-prod-postgres port-forwarder && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "Файлы окружения очищены"; \
  	else \
  	  echo "Отмена очистки"; \
  	fi

env-db-port-forward:
	@docker compose up -d port-forwarder

env-db-port-close:
	@docker compose down port-forwarder

env-http-addr-clean:
	@PID=$$(lsof -t -i ${HTTP_ADDR}); \
    	if [ -n "$$PID" ]; then \
    	  kill -9 $$PID && echo "Освобождён порт ${HTTP_ADDR} (PID $$PID)"; \
    	else \
    	  echo "Порт ${HTTP_ADDR} уже свободен"; \
    	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi
	docker compose run --rm go-rest-prod-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-create action=up"; \
		exit 1; \
	fi
	docker compose run --rm go-rest-prod-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@go-rest-prod-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

logs-cleanup:
	@read -p "Очистить все log файлы? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  rm -rf ${PROJECT_ROOT}/out/logs && \
	  echo "Файлы логов очищены"; \
  	else \
  	  echo "Отмена очистки логов"; \
  	fi

run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/go-rest-prod/main.go

deploy:
	@docker compose up -d --build go-rest-prod

undeploy:
	@docker compose down go-rest-prod

ps:
	@docker compose ps

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/go-rest-prod/main.go
		-o docs \
		--parseInternal \
		--parseDependency