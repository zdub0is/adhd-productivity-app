# ADHD Productivity App

A self-hosted, ADHD-friendly productivity app combining a Sunsama-style daily
planning ritual, Todoist-style simplicity with an open API-first design,
Lunatask-style journaling and mood/energy check-ins, and ML-driven
recommendations trained on the user's own logged patterns.

Built as a portfolio project to demonstrate stack growth (Go, Python,
RabbitMQ, Redis), clean service architecture, and visible reasoning via
[Architecture Decision Records](docs/adr) and a [devlog](../../tree/devlog).

## Status

Early scaffolding — see the [Projects board](../../projects) and
[open issues](../../issues) for current progress.

## Architecture

_The Go API service and Postgres now exist (phase 2). This section will grow
to describe the Python ML microservice, RabbitMQ event plumbing, and Redis
caching layer as later phases land. See [docs/adr](docs/adr) for decisions
made so far._

The v1 API is a single Go service backed by Postgres, covering tasks, the
daily planning ritual, mood/energy check-ins, and journaling. Every `/api/v1`
route requires an `X-API-Key` header; keys are minted via
`POST /api/v1/auth/keys`, itself gated by an `X-Admin-Token` bootstrap secret
(see [`.env.example`](.env.example)).

### Running locally

```sh
cp .env.example .env   # fill in POSTGRES_PASSWORD and ADMIN_BOOTSTRAP_TOKEN
docker compose up --build
```

The API applies its own Postgres migrations on startup and listens on
`:8090`.

### Deploying to the Raspberry Pi

The stack targets a Raspberry Pi (`bethel@aogami.local`) that already runs a
Second Life viewer container, within a ~128GB budget on a 256GB SD card.
Both `postgres:16-alpine` and the Go build's `golang:1.27-alpine` /
`alpine:3.20` base images are multi-arch and run on arm64 as-is — no
Pi-specific image changes are needed.

From a machine with SSH access to the Pi (the script runs over SSH; it
doesn't need to run on the Pi itself):

```sh
./scripts/deploy-pi.sh [branch]   # defaults to main
```

This clones the repo on first run (or fast-forwards an existing checkout to
`origin/<branch>`), then runs `docker compose up -d --build` and polls
`/healthz`. It refuses to proceed if `~/adhd-productivity-app/.env` is
missing on the Pi — copy `.env.example` there and fill in
`POSTGRES_PASSWORD` / `ADMIN_BOOTSTRAP_TOKEN` first. `api` binds `8090` and
Postgres is only reachable inside the compose network, so neither should
collide with the existing SL viewer container's ports.

## Repository layout

```
api/                     Go API service
  cmd/server/             entrypoint
  internal/config/        env-based configuration
  internal/db/            Postgres connection + embedded migrations
  internal/apikey/        API key issuance and auth middleware
  internal/tasks/         /api/v1/tasks
  internal/planning/      /api/v1/planning/*
  internal/checkins/      /api/v1/checkins
  internal/journal/       /api/v1/journal
docs/adr/                 Architecture Decision Records
docker-compose.yml        Go API + Postgres for local/Pi deployment
scripts/deploy-pi.sh      Deploy the stack to the Raspberry Pi over SSH
```

## Devlog

Development notes and reasoning live on the [`devlog`](../../tree/devlog)
branch, published via GitHub Pages.
