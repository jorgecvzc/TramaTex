#!/usr/bin/env bash
# TramaTex - Rebuild total de staging en host remoto Linux (pcele)
# Ejecutar en el servidor remoto, dentro o fuera del repo.

set -euo pipefail

PROJECT_DIR="${PROJECT_DIR:-/opt/tramatex}"
COMPOSE_FILE="${COMPOSE_FILE:-docker/docker-compose.remote.yml}"
ENV_FILE="${ENV_FILE:-docker/.env}"
CHECKOUT_REF="${CHECKOUT_REF:-origin/staging}"
PRESERVE_DATABASE="${PRESERVE_DATABASE:-false}"
REMOVE_IMAGES="${REMOVE_IMAGES:-true}"

usage() {
  cat <<'EOF'
Usage: scripts/rebuild-staging-remote.sh [options]

Options:
  --project-dir <path>      Repo path on remote host (default: /opt/tramatex)
  --compose-file <path>     Compose file path (default: docker/docker-compose.remote.yml)
  --env-file <path>         Env file path (default: docker/.env)
  --checkout-ref <ref>      Ref to align staging branch before rebuild (default: origin/staging)
  --no-checkout             Skip git fetch/checkout/reset step
  --preserve-database       Keep DB volume/data (no -v)
  --skip-image-remove       Do not remove API/Frontend/Postgres images before pull
  -h, --help                Show this help

Examples:
  scripts/rebuild-staging-remote.sh --checkout-ref origin/staging
  scripts/rebuild-staging-remote.sh --checkout-ref origin/chore/staging-deploy-scripts
  scripts/rebuild-staging-remote.sh --no-checkout --preserve-database
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project-dir)
      PROJECT_DIR="$2"
      shift 2
      ;;
    --compose-file)
      COMPOSE_FILE="$2"
      shift 2
      ;;
    --env-file)
      ENV_FILE="$2"
      shift 2
      ;;
    --checkout-ref)
      CHECKOUT_REF="$2"
      shift 2
      ;;
    --no-checkout)
      CHECKOUT_REF=""
      shift
      ;;
    --preserve-database)
      PRESERVE_DATABASE="true"
      shift
      ;;
    --skip-image-remove)
      REMOVE_IMAGES="false"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

cd "$PROJECT_DIR"

echo "[1/5] Preparing repository in $PROJECT_DIR"
if [[ -n "$CHECKOUT_REF" ]]; then
  git fetch origin
  git checkout -B staging "$CHECKOUT_REF"
  echo "Using commit: $(git rev-parse HEAD)"
else
  echo "Skipping git checkout step"
fi

echo "[2/5] Stopping containers"
if [[ "$PRESERVE_DATABASE" == "true" ]]; then
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down --remove-orphans
else
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down -v --remove-orphans
fi

echo "[3/5] Refreshing images"
if [[ "$REMOVE_IMAGES" == "true" ]]; then
  docker image rm -f \
    ghcr.io/jorgecvzc/tramatex-api:latest \
    ghcr.io/jorgecvzc/tramatex-frontend:latest \
    postgres:15-alpine || true
fi

docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" pull

echo "[4/5] Starting services"
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --force-recreate

echo "[5/5] Final cleanup and status"
docker image prune -f
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

echo "Rebuild finished successfully."
