#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd "${SCRIPT_DIR}"

if [ -f ./.env ]; then
	echo "Loading environment from .env"
	set -a
	source ./.env
	set +a
fi

APP_ROOT="${SCRIPT_DIR}/quickstart/eino_assistant"
cd "${APP_ROOT}"

if [ -f ./.env ]; then
	echo "Loading module environment from quickstart/eino_assistant/.env"
	set -a
	source ./.env
	set +a
fi

exec go run ./cmd/einoagentcli/main.go "$@"
