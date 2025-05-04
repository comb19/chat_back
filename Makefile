.PHONY: migrate, update

migrate:
	docker compose down
	docker volume rm chat_back_db_store_dev
	docker compose up -d --build

update:
	docker compose up -d --build