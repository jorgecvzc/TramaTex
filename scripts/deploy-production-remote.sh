#!/usr/bin/env bash
# TramaTex - Deploy to production on remote Linux host (DigitalOcean)
# Runs on the production server. Pulls pre-built images from GHCR — never builds locally.

set -euo pipefail

PROJECT_DIR="${PROJECT_DIR:-/opt/tramatex}"
COMPOSE_FILE="${COMPOSE_FILE:-docker/docker-compose.remote.yml}"
ENV_FILE="${ENV_FILE:-docker/.env}"
CHECKOUT_REF="${CHECKOUT_REF:-origin/master}"
PRESERVE_DATABASE="${PRESERVE_DATABASE:-true}"
GHCR_USER="${GHCR_USER:-}"
GHCR_TOKEN="${GHCR_TOKEN:-}"

usage() {
  cat <<'EOF'
Usage: scripts/deploy-production-remote.sh [options]

Options:
  --project-dir <path>      Repo path on remote host (default: /opt/tramatex)
  --compose-file <path>     Compose file path (default: docker/docker-compose.remote.yml)
  --env-file <path>         Env file path (default: docker/.env)
  --checkout-ref <ref>      Ref to align before deploy (default: origin/master)
  --no-checkout             Skip git fetch/checkout/reset step
  --wipe-database           Remove DB volume (DESTRUCTIVE — resets all demo data)
  --ghcr-user <user>        GHCR username for docker login
  --ghcr-token <token>      GHCR token/PAT for docker login
  -h, --help                Show this help

Examples:
  scripts/deploy-production-remote.sh
  scripts/deploy-production-remote.sh --no-checkout
  scripts/deploy-production-remote.sh --wipe-database
EOF
}

WIPE_DATABASE="false"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project-dir)   PROJECT_DIR="$2";  shift 2 ;;
    --compose-file)  COMPOSE_FILE="$2"; shift 2 ;;
    --env-file)      ENV_FILE="$2";     shift 2 ;;
    --checkout-ref)  CHECKOUT_REF="$2"; shift 2 ;;
    --no-checkout)   CHECKOUT_REF="";   shift   ;;
    --wipe-database) WIPE_DATABASE="true"; PRESERVE_DATABASE="false"; shift ;;
    --ghcr-user)     GHCR_USER="$2";    shift 2 ;;
    --ghcr-token)    GHCR_TOKEN="$2";   shift 2 ;;
    -h|--help)       usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 1 ;;
  esac
done

cd "$PROJECT_DIR"

echo "[1/5] Preparing repository in $PROJECT_DIR"
if [[ -n "$CHECKOUT_REF" ]]; then
  git fetch origin
  git checkout -B master "$CHECKOUT_REF"
  echo "Using commit: $(git rev-parse HEAD)"
else
  echo "Skipping git checkout step"
fi

echo "[2/5] Stopping containers"
if [[ "$PRESERVE_DATABASE" == "true" ]]; then
  echo "Preserving database (no volume removal)"
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down --remove-orphans
else
  echo "WARNING: Removing database volume — all data will be lost"
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down -v --remove-orphans
fi

echo "[3/5] Logging in to GHCR and pulling images"
if [[ -n "$GHCR_USER" && -n "$GHCR_TOKEN" ]]; then
  echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
fi

docker image rm -f \
  ghcr.io/jorgecvzc/tramatex-api:latest \
  ghcr.io/jorgecvzc/tramatex-frontend:latest || true

docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" pull

echo "[4/5] Starting services"
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --force-recreate

echo "[5/5] Waiting for API health check..."
for i in $(seq 1 36); do
  if wget -qO- http://localhost/api/health > /dev/null 2>&1; then
    echo "API is healthy"
    break
  fi
  if [[ $i -eq 36 ]]; then
    echo "ERROR: API did not become healthy after 3 minutes" >&2
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
    exit 1
  fi
  sleep 5
done

docker image prune -f
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

echo "Deploy to production finished successfully."
