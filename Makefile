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
	  docker compose down go-rest-prod-postgres && \
	  rm -rf out/pgdata && \
	  echo "Файлы окружения очищены"; \
  	else \
  	  echo "Отмена очистки"; \
  	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

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