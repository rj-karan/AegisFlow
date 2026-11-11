# Development

## Prerequisites

- Go 1.24+
- Node 22+ with `corepack enable` (provides pnpm)
- Docker with Compose v2
- `make`

## Everything at once

```bash
cp .env.example .env
make dev          # docker compose up --build: postgres, redis, api
make dev-logs     # tail logs
make dev-down     # stop and delete volumes
```

The API is on http://127.0.0.1:8080. Ports bind to loopback only, so nothing is exposed on
your LAN. Worker and web services join `compose.yaml` with issues #5 and #6.

## Services on their own

```bash
docker compose up postgres redis   # dependencies only
make run                           # go run ./cmd/api, reads .env values from your shell
```

Load `.env` into a shell with `set -a; source .env; set +a`.

## Container image

`make image` builds `aegisflow-api:dev` from the multi-stage `Dockerfile`: Go builder on
Alpine, final image on distroless `static` (about 17 MB, no shell, non-root). Compose runs
it with a read-only root filesystem, `no-new-privileges` and all capabilities dropped.
Base images are pinned by digest; Dependabot proposes updates weekly.

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
