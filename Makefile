.PHONY: migrate, update, cleanup, test

migrate:
	docker compose down
	docker volume rm chat_back_db_store_dev
	docker compose up -d --build

update:
	docker compose up -d --build

cleanup:
	docker compose down
	docker volume rm chat_back_db_store
	docker volume rm chat_back_db_store_dev
	docker compose up -d --build

test:
	docker container exec -it chat_back-app-1 go test ./test