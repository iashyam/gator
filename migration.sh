 #!/usr/bin/env zsh

 # migration.sh - wrapper around goose to run migrations up or down
 # Usage:
 #   ./migration.sh up   # run all pending up migrations
 #   ./migration.sh down # run all pending down migrations
 #
 # The script uses $DATABASE_URL if set, otherwise you can provide
 # DB components via env vars: DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME.

# set -euo pipefail

 MODE="${1:-}"
 if [[ -z "$MODE" ]]; then
	 echo "Usage: $0 up|down" >&2
	 exit 2
 fi

 if [[ "$MODE" != "up" && "$MODE" != "down" ]]; then
	 echo "Invalid mode: $MODE. Expected 'up' or 'down'." >&2
	 exit 2
 fi

 # Prefer DATABASE_URL if provided
 if [[ -n "${DATABASE_URL:-}" ]]; then
	 DSN="$DATABASE_URL"
 else
	 : ${DB_USER:="postgres"}
	 : ${DB_PASS:=""}
	 : ${DB_HOST:="localhost"}
	 : ${DB_PORT:="5432"}
	 : ${DB_NAME:="gator"}

	 if [[ -n "$DB_PASS" ]]; then
		 DSN="postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
	 else
		 DSN="postgres://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
	 fi
 fi

 MIGRATIONS_DIR="./sql/schema"

 echo "Running goose ${MODE} against: ${DSN} (migrations dir: ${MIGRATIONS_DIR})"

 # Ensure goose is installed (the CLI is required)
 if ! command -v goose >/dev/null 2>&1; then
	 echo "goose CLI not found. Install with: go install github.com/pressly/goose/v3/cmd/goose@latest" >&2
	 exit 3
 fi

 case "$MODE" in
	 up)
		 goose -dir "$MIGRATIONS_DIR" postgres "$DSN" up
		 ;;
	 down)
		 goose -dir "$MIGRATIONS_DIR" postgres "$DSN" down
		 ;;
 esac

 echo "goose ${MODE} finished"