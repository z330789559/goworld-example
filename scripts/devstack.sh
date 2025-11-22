#!/usr/bin/env bash
set -euo pipefail

CMD=${1:-help}

compose_file=$(cat <<'YAML'
version: '3.9'
services:
  redis:
    image: redis:7
    ports:
      - "6379:6379"
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: mahjong
      POSTGRES_PASSWORD: mahjong
      POSTGRES_DB: mahjong
    ports:
      - "5432:5432"
YAML
)

run_compose() {
  echo "$compose_file" | docker compose -f - "$@"
}

case "$CMD" in
  up)
    run_compose up -d
    ;;
  down)
    run_compose down
    ;;
  help|*)
    echo "Usage: $0 [up|down]"
    ;;
esac
