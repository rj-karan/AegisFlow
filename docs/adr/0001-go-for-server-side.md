# ADR-0001: Go for the API and workers

- Status: accepted
- Date: 2026-11-10

## Context

The original project overview (`docs/AegisFlow-Project-Overview.pdf`) chose FastAPI for the
API and Python for the workers, mainly so both planes share one language and the LLM SDKs
are familiar.

## Decision

Use Go for everything server-side: the API (`cmd/api`), the queue workers (`cmd/worker`),
and the policy engine (OPA embedded as a library). Keep Next.js for the dashboard.

## Why

- The tooling AegisFlow orchestrates is Go: OPA, actionlint, Trivy, Gitleaks, OSV-Scanner,
  the GitHub CLI, Docker, Kubernetes. Embedding OPA and actionlint as libraries removes
  shell-outs and version drift.
- Static binaries make tiny, disposable worker images with no interpreter or dependency
  tree inside the sandbox.
- Goroutines fit webhook fan-out and parallel repository analysis.
- `sqlc` gives compile-time checked SQL, which matters when every query must be
  workspace-scoped.
- The maintainer wants to learn the Go DevOps ecosystem.

## Consequences

- LLM calls use the official Anthropic and OpenAI Go SDKs. Fewer examples online than
  Python, but the surface we need (one chat completion per explanation or repair) is small.
- Queue is Redis + `asynq` instead of RQ. Same semantics: retries, dead-letter, scheduling.
- Migrations with `goose`, HTTP with `chi`, DB access with `pgx`.
- Any teammate joining the server side needs Go. The design doc's table of tech choices
  is superseded by `README.md#architecture`.
