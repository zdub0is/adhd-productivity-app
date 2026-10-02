#!/usr/bin/env bash
# Deploy the phase 2 stack (Go API + Postgres) to the Raspberry Pi.
#
# Run this from any machine with SSH access to the Pi (it does the work
# over SSH, so it does not need to run on the Pi itself).
#
# Usage:
#   ./scripts/deploy-pi.sh [branch]
#
# Env overrides:
#   PI_HOST   ssh target (default: bethel@aogami.local)
#   PI_DIR    checkout path on the Pi (default: ~/adhd-productivity-app)
set -euo pipefail

PI_HOST="${PI_HOST:-bethel@aogami.local}"
PI_DIR="${PI_DIR:-adhd-productivity-app}"
BRANCH="${1:-main}"
REPO_URL="https://github.com/zdub0is/adhd-productivity-app.git"

echo "==> Deploying branch '$BRANCH' to $PI_HOST:$PI_DIR"

ssh "$PI_HOST" bash -s -- "$PI_DIR" "$BRANCH" "$REPO_URL" <<'REMOTE'
set -euo pipefail
DIR="$1"
BRANCH="$2"
REPO_URL="$3"

if [ ! -d "$DIR/.git" ]; then
  echo "==> No existing checkout, cloning..."
  git clone "$REPO_URL" "$DIR"
fi

cd "$DIR"
git fetch origin "$BRANCH"
git checkout "$BRANCH"
git reset --hard "origin/$BRANCH"

if [ ! -f .env ]; then
  cat >&2 <<EOF
ERROR: $DIR/.env is missing.
Copy .env.example to .env on the Pi and fill in POSTGRES_PASSWORD and
ADMIN_BOOTSTRAP_TOKEN before deploying:
  cp .env.example .env && \$EDITOR .env
EOF
  exit 1
fi

echo "==> Building and starting containers"
docker compose up -d --build

echo "==> Container status"
docker compose ps
REMOTE

echo "==> Waiting for health check"
for i in $(seq 1 15); do
  if ssh "$PI_HOST" 'curl -sf localhost:8090/healthz' >/dev/null 2>&1; then
    echo "==> API is healthy"
    exit 0
  fi
  sleep 2
done

echo "==> API did not become healthy in time; check logs with:" >&2
echo "    ssh $PI_HOST 'cd $PI_DIR && docker compose logs api'" >&2
exit 1
