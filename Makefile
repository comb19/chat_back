include .env
.PHONY: init, restart, migrate, update, clean, test, dev-migrate, test-migrate

init:
	$(MAKE) clean
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=DEVELOPMENT \
	docker compose up -d --build
	docker compose down db-dev db-test-dev

restart:
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=DEVELOPMENT \
	docker compose restart app db db-test

update:
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=DEVELOPMENT \
	docker compose up app -d --build

clean:
	docker compose down -v

test:
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${TEST_POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=TEST \
	docker compose restart app	
	docker container exec -it ${APP_CONTAINER_NAME} go test ./test

migrate:
	$(MAKE) dev-migrate
	$(MAKE) test-migrate

dev-migrate:
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=DEVELOPMENT \
	docker compose up db-dev -d
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=DEVELOPMENT \
	docker compose up migrate
	docker compose down db-dev

test-migrate:
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${TEST_POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=TEST \
	docker compose up db-test-dev -d
	env DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${TEST_POSTGRES_HOSTNAME}:5432/${POSTGRES_DB}?sslmode=disable" ENV=TEST \
	docker compose up test-migrate
	docker compose down db-test-dev