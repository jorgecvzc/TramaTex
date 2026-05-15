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
SKIP_GIT="${SKIP_GIT:-false}"
BUILD_SOURCE="${BUILD_SOURCE:-false}"

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
  --build-source            Build images from local source instead of pulling from GHCR
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
    --build-source)
      BUILD_SOURCE="true"
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

# Redirigir el almacenamiento temporal de Docker a la partición con espacio (/)
# Esto evita el error "no space left on device" cuando /home está lleno.
export DOCKER_CONFIG="$PROJECT_DIR/.docker_config"
mkdir -p "$DOCKER_CONFIG"
# Asegurar permisos: si la carpeta ya existía (ej. creada por root), recuperamos la propiedad para el usuario actual
if [ -d "$DOCKER_CONFIG" ]; then
    sudo chown -R $(id -u):$(id -g) "$DOCKER_CONFIG" 2>/dev/null || true
fi

echo "[1/5] Preparing repository in $PROJECT_DIR"
if [[ "$SKIP_GIT" == "true" ]]; then
  echo "Skipping git step (already done by caller)"
elif [[ -n "$CHECKOUT_REF" ]]; then
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

# Login to GHCR if credentials are present in the env file
if [[ -f "$ENV_FILE" ]]; then
  GHCR_USER_VAL=$(grep -E '^GHCR_USER=' "$ENV_FILE" | cut -d= -f2- | tr -d '\r' || true)
  GHCR_TOKEN_VAL=$(grep -E '^GHCR_TOKEN=' "$ENV_FILE" | cut -d= -f2- | tr -d '\r' || true)
  if [[ -n "$GHCR_USER_VAL" && -n "$GHCR_TOKEN_VAL" ]]; then
    echo "Logging in to GHCR as $GHCR_USER_VAL..."
    echo "$GHCR_TOKEN_VAL" | docker login ghcr.io -u "$GHCR_USER_VAL" --password-stdin
  else
    echo "INFO: GHCR_USER/GHCR_TOKEN not set in $ENV_FILE — skipping login (assuming public images)"
  fi
fi

if [[ "$BUILD_SOURCE" == "true" ]]; then
  echo "INFO: Building images from source..."
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build
else
  echo "INFO: Pulling pre-built images from GHCR..."
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" pull
fi

echo "[4/5] Starting services"
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --force-recreate

echo "[5/5] Final cleanup and status"
docker image prune -f

echo "Waiting for API health check..."
for i in $(seq 1 30); do
  if wget -qO- http://localhost/api/health > /dev/null 2>&1; then
    echo "API is healthy"
    break
  fi
  if [[ $i -eq 30 ]]; then
    echo "WARN: API health check timed out after 60s"
  fi
  sleep 2
done

docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

echo "Rebuild finished successfully."
