.PHONY: restart, migrate, update, clean, test, update_test, migrate_test

restart:
	docker compose restart

migrate:
	docker compose down
	docker volume rm postgres_chat_app_dev
	docker compose --env-file .env.development up -d --build

update:
	docker compose --env-file .env.development up -d --build

clean:
	docker compose down
	docker volume rm postgres_chat_app
	docker volume rm postgres_chat_app_dev
	docker volume rm postgres_chat_app_test
	docker volume rm postgres_chat_app_test_dev

test:
	docker container exec -it api_chat_app_test go test ./test

update_test:
	docker compose --env-file .env.test up -d --build

migrate_test:
	docker compose down
	docker volume rm postgres_chat_app_test_dev
	docker compose --env-file .env.test up -d --build
