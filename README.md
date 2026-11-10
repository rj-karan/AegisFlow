<div align="center">

# AegisFlow

**Secure CI/CD for GitHub repositories, generated, validated and kept safe by policy.**

Connect a repo. AegisFlow reads the code, recommends the right scanners with reasons,
compiles a hardened GitHub Actions pipeline, and opens a pull request.
When CI breaks later, it diagnoses the failure and proposes a repair that
**cannot weaken security**, because every change passes a deterministic policy gate
before a human merges it.

[![CI](https://github.com/rj-karan/AegisFlow/actions/workflows/ci.yml/badge.svg)](https://github.com/rj-karan/AegisFlow/actions/workflows/ci.yml)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![Next.js](https://img.shields.io/badge/dashboard-Next.js-000?logo=nextdotjs)](web/)
[![Project board](https://img.shields.io/badge/kanban-project%20board-6f42c1?logo=github)](https://github.com/users/rj-karan/projects/2)

</div>

---

## Why

Most small teams either have no CI security at all, or have workflows that quietly drift
into an unsafe state: unpinned third-party Actions, `write-all` tokens, unsafe
`pull_request_target` triggers, jobs running in the wrong folder of a monorepo.

AI "fix my CI" bots can make it worse. Told only to "make CI green", an AI will happily
delete the failing scanner or add `continue-on-error: true`.

In March 2026 the popular `aquasecurity/trivy-action` was compromised and exfiltrated
runner secrets. One organisation spent 48 hours SHA-pinning Actions across 68 repos.
AegisFlow exists to make that kind of hardening the default, and to keep it that way.

## What it does

```
 connect ─▶ analyze ─▶ recommend ─▶ plan ─▶ compile ─▶ gate ─▶ validate ─▶ PR ─▶ monitor ─▶ repair ─┐
                                      ▲                                                             │
                                      └─────────────── a repair edits the plan, never the YAML ◀────┘
```

| Stage | What happens |
|---|---|
| **Analyze** | Parsers read `package.json`, `pyproject.toml`, `go.mod`, Dockerfiles, Terraform and existing workflows. Every detected fact points at the file that proves it. The repo's code is never executed here. |
| **Recommend** | A rule-based ranker over a tool catalog (license, stacks, runtime cost, credentials) picks tests, lint, SAST, dependency, secret and container scanners. AI only writes the plain-English *why*. |
| **Plan and compile** | Recommendations become a provider-neutral pipeline plan. A deterministic compiler turns it into GitHub Actions YAML: read-only `permissions`, every `uses:` pinned to a full commit SHA, each job in the right folder. |
| **Gate** | The gate-weakening detector diffs old plan vs new plan and blocks anything that makes CI less safe: stage removed, `continue-on-error`, tests narrowed, permissions widened, Action unpinned, unsafe trigger, threshold lowered, secret exposed. |
| **Validate** | Lint (`actionlint`, `zizmor`) and a sandboxed run with no secrets and restricted network. Jobs that need real credentials are marked `not validated`, never "passing". |
| **Install** | Everything arrives as a pull request. AegisFlow never pushes to the default branch and never merges its own PRs. |
| **Monitor and repair** | `workflow_run` webhooks feed a classifier (`missing_dependency`, `runtime_version`, `deprecated_action`, `secret_missing`, ...). A fix is proposed against the plan, re-gated, re-validated, and opened as a small repair PR with logs and rollback notes. Real `test_failure`s go back to the developer. |

## Architecture

```
┌──────────────────────── control plane ────────────────────────┐
│  Next.js dashboard  ──▶  Go API (chi)  ──▶  PostgreSQL         │
│                               │                                │
│                               ▼                                │
│                        Redis + asynq queue                     │
└───────────────────────────────┬────────────────────────────────┘
                                │ jobs
┌───────────────────────────────▼─── worker plane (Docker) ─────┐
│   analyzer     ·     validator (sandbox)     ·     repair      │
│   Go binaries; untrusted repo code is only ever touched here   │
└────────────────────────────────────────────────────────────────┘
```

| Layer | Choice |
|---|---|
| Dashboard | Next.js, TypeScript, Tailwind |
| API | Go, `chi` router, `sqlc` typed queries, `goose` migrations |
| Database | PostgreSQL, workspace-scoped on every table |
| Queue | Redis + `asynq` (retries, dead-letter, scheduling) |
| Workers | Go binaries in disposable, resource-limited containers |
| Policy | OPA embedded as a Go library, `actionlint`, `zizmor` |
| LLM | Pluggable (Claude, OpenAI, local). Bring your own key. Works without AI. |
| Local dev | Docker Compose, one command |

## Quickstart

```bash
git clone https://github.com/rj-karan/AegisFlow.git
cd AegisFlow
cp .env.example .env
docker compose up --build
```

Dashboard: http://localhost:3000 · API docs: http://localhost:8000/docs

See [docs/development.md](docs/development.md) for running each service on its own.
Why Go instead of the Python in the original design: [ADR-0001](docs/adr/0001-go-for-server-side.md).

## Roadmap

Built in vertical slices. Each milestone ships backend, UI, tests and docs together and is
done only when its exit test passes on real repositories.

| Milestone | Scope | Exit test |
|---|---|---|
| [M1 Foundation](docs/milestones/M1.md) | Monorepo, auth, workspaces, GitHub App read access | Lists a test org's repos |
| [M2 Analyzer](docs/milestones/M2.md) | Detectors for Python, JS/TS, Go, Java; profiles with evidence | 8 fixture repos detected correctly |
| [M3 Recommend + plan](docs/milestones/M3.md) | Tool catalog, ranking, plan IR, Actions compiler | Generated YAML passes actionlint |
| [M4 Gate + validate](docs/milestones/M4.md) | OPA rules, gate-weakening detector, sandbox runner | Every bad fixture blocked |
| [M5 PR + monitor](docs/milestones/M5.md) | Branch, commit, open PR; webhooks and run history | A real PR merged |
| [M6 Repair + eval](docs/milestones/M6.md) | Failure classifier, repair PRs, 30-repo benchmark | Evaluation vs template-only and LLM-only |
| [M7 Operate](docs/milestones/M7.md) (alongside M1–M6) | Terraform, Kubernetes, Helm, Argo CD GitOps, Prometheus/Grafana/Loki, signed releases | Merge to `main` is live on the cluster with no manual `kubectl`, verified by `gh attestation verify` |

Progress is tracked on the [project board](https://github.com/users/rj-karan/projects/2)
and in [issues](https://github.com/rj-karan/AegisFlow/issues). Everything runs for free on a
student budget: see [docs/operations/free-tier-plan.md](docs/operations/free-tier-plan.md).
Learning DevOps along the way? Follow [docs/learning-path.md](docs/learning-path.md).

## Security model

AegisFlow holds write access to other people's repositories, so it is itself a target.
All repo content, logs and AI output are treated as untrusted data, never as instructions.
Secrets never enter prompts. Webhooks are HMAC-verified and deduplicated. Every API call
and every query is workspace-scoped. High-risk actions require explicit human approval and
land in an append-only audit log. Full threat model: [docs/design.md](docs/design.md#security-model).

Found a vulnerability? See [SECURITY.md](SECURITY.md).

## Contributing

Issues and pull requests are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) for the
branch, commit and review conventions, and pick an issue labelled
[`good first issue`](https://github.com/rj-karan/AegisFlow/labels/good%20first%20issue).

## License

[Apache-2.0](LICENSE)
