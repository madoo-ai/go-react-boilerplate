#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -f .env ]; then
  echo "scripts/dev.sh: .env is missing; run make env first" >&2
  exit 1
fi

set -a
# shellcheck disable=SC1091
. ./.env
set +a

go run ./cmd/server &
api_pid=$!

cleanup() {
  kill "$api_pid" 2>/dev/null || true
  wait "$api_pid" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

cd web
npm run dev
