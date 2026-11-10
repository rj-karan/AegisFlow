# Development

## Prerequisites

- Go 1.24+
- Node 22+ with `corepack enable` (provides pnpm)
- Docker with Compose v2
- `make`

## Everything at once

```bash
cp .env.example .env
make dev          # docker compose up --build
```

## Services on their own

```bash
make api          # go run ./cmd/api
make worker       # go run ./cmd/worker
make web          # cd web && pnpm dev
```

Postgres and Redis can stay in Compose: `docker compose up postgres redis`.

## Tests and lint

```bash
make test         # go test ./... + pnpm test
make lint         # golangci-lint, eslint, actionlint
```

## Database

```bash
make migrate      # goose up
make sqlc         # regenerate typed queries from internal/db/queries/*.sql
```

## Branch and PR flow

See [CONTRIBUTING.md](../CONTRIBUTING.md).
