.PHONY: build up down restart logs rebuild-app infra-up infra-down backup

# ── Infra (DB + pgAdmin) ──────────────────────────────────────────────────────
infra-up:
	docker network create trends-net 2>/dev/null || true
	docker compose -f docker-compose.infra.yml up -d

infra-down:
	docker compose -f docker-compose.infra.yml down

# ── App ───────────────────────────────────────────────────────────────────────
build:
	docker compose build

up:
	docker network create trends-net 2>/dev/null || true
	docker compose up -d

down:
	docker compose down

restart: down build up

logs:
	docker compose logs -f app

rebuild-app:
	docker compose build app && docker compose up -d --no-deps app

# ── Backup ────────────────────────────────────────────────────────────────────
backup:
	./scripts/backup-db.sh
