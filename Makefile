.PHONY: build up down restart logs rebuild-app

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

restart: down build up

logs:
	docker compose logs -f app

rebuild-app:
	docker compose build app && docker compose up -d --no-deps app
