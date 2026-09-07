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
```

## Devlog

Development notes and reasoning live on the [`devlog`](../../tree/devlog)
branch, published via GitHub Pages.
