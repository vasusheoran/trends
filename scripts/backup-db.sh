#!/usr/bin/env bash
set -euo pipefail

NAS_PATH="${NAS_BACKUP_PATH:-/mnt/tank/data/trends}"
KEEP_DAYS="${BACKUP_KEEP_DAYS:-7}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DUMP_FILE="/tmp/trends_${TIMESTAMP}.sql.gz"
DB_CONTAINER="${DB_CONTAINER:-trends-db-1}"

echo "[$(date)] Starting backup..."

# Dump from running container, stream directly to compressed file
docker exec "$DB_CONTAINER" \
    pg_dump -U trends trends | gzip > "$DUMP_FILE"

mkdir -p "$NAS_PATH"
cp "$DUMP_FILE" "$NAS_PATH/"
rm "$DUMP_FILE"

# Rotate: delete dumps older than KEEP_DAYS
find "$NAS_PATH" -name "trends_*.sql.gz" -mtime +"$KEEP_DAYS" -delete

echo "[$(date)] Backup complete: $NAS_PATH/trends_${TIMESTAMP}.sql.gz"
