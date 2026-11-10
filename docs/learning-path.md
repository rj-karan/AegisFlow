# Learning path

AegisFlow is built in the open as a way to learn DevOps by doing. Each milestone introduces
a set of concepts; this page lists them in the order the code meets them, with one good
source each. Tick them off as you go.

## Before M1: how a real repo is run

- [ ] Branch protection, squash merge, linear history. [GitHub docs](https://docs.github.com/repositories/configuring-branches-and-merges-in-your-repository)
- [ ] Conventional Commits and why squash-merge makes them matter. [conventionalcommits.org](https://www.conventionalcommits.org/)
- [ ] Issue forms, PR templates, CODEOWNERS, milestones, Projects. [GitHub community health files](https://docs.github.com/communities)
- [ ] Architecture Decision Records. [adr.github.io](https://adr.github.io/)
- [ ] Choosing a licence: MIT vs Apache-2.0 (patent grant). [choosealicense.com](https://choosealicense.com/)

## M1 Foundation: Go services done properly

- [ ] Go modules, `cmd/` vs `internal/`, `go vet`, `gofmt`, `golangci-lint`. [Go project layout](https://go.dev/doc/modules/layout)
- [ ] `net/http` server timeouts and why the defaults are unsafe. [Cloudflare: the complete guide to Go net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- [ ] Graceful shutdown with `context` and `signal.NotifyContext`; what SIGTERM means in Kubernetes. [Kubernetes pod lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination)
- [ ] Liveness vs readiness probes. [Kubernetes probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [ ] Structured logging with `log/slog`; request IDs. [Go blog: structured logging](https://go.dev/blog/slog)
- [ ] Table-driven tests, `httptest`, the race detector. [Go testing](https://go.dev/doc/tutorial/add-a-test)
- [ ] Trusting proxy headers: why `X-Forwarded-For` is spoofable. [MDN](https://developer.mozilla.org/docs/Web/HTTP/Headers/X-Forwarded-For)
- [ ] PostgreSQL migrations (goose) and compile-time-checked SQL (sqlc). [sqlc docs](https://docs.sqlc.dev/)
- [ ] Background jobs, retries, dead-letter queues with Redis (asynq). [asynq](https://github.com/hibiken/asynq)
- [ ] GitHub Apps vs OAuth Apps, installation tokens, webhook HMAC verification. [GitHub Apps docs](https://docs.github.com/apps)
- [ ] Multi-tenancy: scoping every query by workspace. [OWASP: broken access control](https://owasp.org/Top10/A01_2021-Broken_Access_Control/)
- [ ] Docker multi-stage builds, distroless images, non-root containers. [Docker best practices](https://docs.docker.com/build/building/best-practices/)

## M7 Operate: the DevOps core (runs alongside M1 to M6)

- [ ] Hardened GitHub Actions: least-privilege `permissions`, SHA pinning, `pull_request_target` risks. [GitHub security hardening](https://docs.github.com/actions/security-for-github-actions/security-guides/security-hardening-for-github-actions), [zizmor](https://docs.zizmor.sh/)
- [ ] Terraform: providers, state, remote backends, `plan` vs `apply`, modules. [Terraform tutorials](https://developer.hashicorp.com/terraform/tutorials)
- [ ] Kubernetes basics: Pod, Deployment, Service, Ingress, ConfigMap, Secret, namespaces. [Kubernetes concepts](https://kubernetes.io/docs/concepts/)
- [ ] kind and k3s: what a "distribution" is. [kind](https://kind.sigs.k8s.io/), [k3s](https://docs.k3s.io/)
- [ ] Helm charts, values, templating, `helm lint`. [Helm docs](https://helm.sh/docs/)
- [ ] Pod security: `securityContext`, read-only root FS, dropping capabilities, resource requests and limits. [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [ ] GitOps and Argo CD: app-of-apps, sync waves, self-heal, prune. [Argo CD docs](https://argo-cd.readthedocs.io/)
- [ ] Observability: RED metrics, Prometheus data model, PromQL, Grafana dashboards as code, Loki. [Prometheus docs](https://prometheus.io/docs/introduction/overview/), [Grafana Loki](https://grafana.com/docs/loki/latest/)
- [ ] SLOs and error budgets. [Google SRE workbook: implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ ] Supply-chain security: SBOM, Sigstore keyless signing, SLSA provenance, attestations. [SLSA](https://slsa.dev/), [Sigstore](https://www.sigstore.dev/)
- [ ] Secrets in git: SOPS and age. [SOPS](https://github.com/getsops/sops)
- [ ] Multi-arch images and QEMU; why arm64 matters on free cloud tiers. [Docker multi-platform](https://docs.docker.com/build/building/multi-platform/)
- [ ] cert-manager and Let's Encrypt. [cert-manager](https://cert-manager.io/docs/)
- [ ] Runbooks and incident drills. [Google SRE book: managing incidents](https://sre.google/sre-book/managing-incidents/)

## M2 to M6: the product

- [ ] Parsing manifests (`package.json`, `pyproject.toml`, `go.mod`) and the YAML of GitHub Actions. [Workflow syntax](https://docs.github.com/actions/writing-workflows/workflow-syntax-for-github-actions)
- [ ] Policy as code with OPA and Rego, unit-testing policies. [OPA docs](https://www.openpolicyagent.org/docs/latest/)
- [ ] Sandboxing untrusted code: containers, seccomp, network egress control. [gVisor](https://gvisor.dev/)
- [ ] Deterministic compilers and golden-file tests.
- [ ] Evaluating a system: baselines, metrics, seeded failures.
