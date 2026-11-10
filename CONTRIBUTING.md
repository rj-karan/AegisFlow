# Contributing

Thanks for helping build AegisFlow. This project runs like a product team: every change
starts as an issue, lands as a pull request, and is tracked on the project board.

## Workflow

1. **Find or open an issue.** Use the templates. Link it to a milestone (M1 to M6).
2. **Branch from `main`.** Name it `<type>/<issue-number>-<short-slug>`, for example
   `feat/12-github-app-webhook` or `fix/40-webhook-dedup`.
3. **Commit small and often** using [Conventional Commits](https://www.conventionalcommits.org/):
   `feat(api): verify webhook HMAC signatures`. Types: `feat`, `fix`, `docs`, `test`,
   `refactor`, `chore`, `ci`, `build`.
4. **Open a pull request** against `main`. Fill in the template. Reference the issue with
   `Closes #12` so it closes on merge.
5. **CI must be green.** Lint, tests, actionlint and zizmor run on every PR.
6. **Squash-merge.** The PR title becomes the commit message, so make it a good one.

`main` is protected: no direct pushes, PRs only.

## Code standards

- Go: `gofmt`, `go vet`, `golangci-lint` clean. Errors wrapped with `fmt.Errorf("...: %w", err)`.
  Table-driven tests. No global state.
- TypeScript: strict mode, ESLint clean, components small and typed.
- SQL: every table has `workspace_id`; every query filters by it.
- Generated workflows and this repo's own workflows: read-only `permissions`, every `uses:`
  pinned to a full commit SHA with a version comment.

## Definition of done

- Tests added or updated.
- Docs updated if behaviour changed (README, `docs/`, or a new ADR in `docs/adr/`).
- Issue closed by the PR.

## Architecture decisions

Anything that changes the stack or a cross-cutting rule gets an ADR in `docs/adr/`.
Copy the format of [ADR-0001](docs/adr/0001-go-for-server-side.md).
