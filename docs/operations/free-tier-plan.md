# Running AegisFlow for free

Researched 2026-11-10. Prices and limits change; each line links to its source.

## The principle

Two environments share the same Terraform modules, Helm chart and Argo CD apps:

| Environment | Where | Cost | Needs a card? |
|---|---|---|---|
| `local` | `kind` cluster in Docker on the laptop | 0 | No |
| `oci` | `k3s` on one Oracle Cloud Always Free arm64 VM | 0 | Yes, $0 verification hold |

Everything else is a free tier or free for public repositories. If Oracle signup is
blocked, the fallback is Azure for Students or staying on `local`; the GitOps, monitoring
and release work is identical either way.

## What is free and why it is enough

### Compute

- **Oracle Cloud Always Free**: 2 OCPU and 12 GB RAM of Ampere A1 (arm64) compute, 200 GB
  block storage, forever. Oracle halved this from 4 OCPU / 24 GB in August 2026. A
  single-node k3s with api, worker, web, Postgres, Redis, Argo CD and a trimmed
  monitoring stack fits in 12 GB. Signup requires a real credit or debit card (no
  prepaid or virtual cards) for a temporary verification hold; Always Free usage is never
  charged. Sources: [terminalbytes](https://terminalbytes.com/oracle-cloud-free-tier-changes-2026/),
  [Hacker News thread](https://news.ycombinator.com/item?id=49183750),
  [signup requirements](https://grandavehousing.calpoly.edu/news/oracle-cloud-free-tier-can).
- **Azure for Students**: $100 credit for 12 months with a .edu email, no card. AKS's
  free tier control plane costs $0; you pay only worker nodes (about $30/month for the
  smallest), so the credit lasts roughly three months. Good for a time-boxed cloud demo,
  not for always-on. Sources: [GitHub Student Pack 2026](https://capitalandcompute.net/blog/github-student-developer-pack-2026/),
  [free control planes](https://k8s-pricing.pages.dev/rankings/free-control-plane/).
- **Not available any more**: DigitalOcean left the Student Pack (credits ended July 31
  2026) and AWS has no Student Pack offer. Source: [capitalandcompute](https://capitalandcompute.net/blog/github-student-developer-pack-2026/).
- **Laptop**: `kind` runs a multi-node cluster in Docker in under a minute and is what
  Kubernetes itself uses for CI. `k3d` is faster to start; `kind` has better upstream
  parity. We use `kind`. Source: [k3d vs kind 2026](https://sanj.dev/post/k3d-vs-kind-local-kubernetes/).

### CI, registry, releases

- **GitHub Actions** is unmetered on public repositories: no minute cap on standard
  runners. Source: [GitHub billing docs](https://docs.github.com/billing/managing-billing-for-your-products/managing-billing-for-github-actions).
- **GitHub Container Registry** storage and bandwidth are free for public packages.
  Source: same page.
- **Sigstore / cosign** keyless signing and **GitHub artifact attestations** are free for
  public repos. Source: [GoReleaser supply chain guide](https://goreleaser.com/blog/supply-chain-security/).

### Observability

- Self-hosted in the cluster: `kube-prometheus-stack`, Loki in single-binary mode and
  Alloy for log shipping, with CPU requests tuned down for one node. Free and the most
  learning. Source: [DevOps Odyssey on OCI free tier](https://januschung.github.io/devops-odyssey/oci-k3s-monitoring),
  [k3s monitoring](https://oneuptime.com/blog/post/2026-02-02-k3s-monitoring/markdown).
- Managed alternative: **Grafana Cloud Free**: 10k metric series, 50 GB logs, 50 GB
  traces, 14-day retention, no card, does not expire. Useful as a remote-write target
  when the laptop cluster is off. Source: [Grafana Cloud free plan](https://costbench.com/software/observability/grafana-cloud/free-plan/).
- **New Relic for students**: 500 GB/month via the Student Pack. Not planned, listed for
  completeness. Source: [New Relic students](https://newrelic.com/students).

### Domain, TLS, tunnels

- Free `.me` domain for a year from Namecheap in the Student Pack; TLS from Let's Encrypt
  via cert-manager.
- Local GitHub App development: **smee.io** forwards webhooks to localhost with no
  account (GitHub's own docs use it). For exposing the laptop cluster on a stable
  hostname, **Cloudflare Tunnel** is free. Source: [webhook tools comparison](https://hookdeck.com/webhooks/platforms/best-webhook-testing-tools-local-development).

## What hiring managers look for

The 2026 guides agree on one thing: a deployed, end-to-end repo beats certificates.
The repeated checklist is PR to CI to image to GitOps sync with no manual `kubectl`;
Terraform with a remote backend you can destroy and recreate; Prometheus plus Grafana
showing real data; Argo CD demonstrating self-healing. M7 is built to tick every item.
Sources: [devopsboys](https://devopsboys.com/blog/devops-portfolio-projects-to-get-hired-2026),
[artech](https://www.artech.com/blog/devops-roles-2026-skills-resume-portfolio/),
[engidock](https://www.engidock.com/blog/devops-portfolio-projects).

## Tool choices and why

| Need | Choice | Why not the alternative |
|---|---|---|
| GitOps | Argo CD | Flux is excellent but Argo's UI makes demos and learning clearer; both are CNCF graduated. [Comparison](https://oneuptime.com/blog/post/2026-02-26-argocd-vs-fluxcd-2026/markdown) |
| IaC | Terraform | Industry default on job posts; the `kind` and `oci` providers cover both environments |
| Local cluster | kind | Upstream fidelity; k3d is a fine swap |
| Cloud cluster | k3s | Single binary, light enough for 12 GB, same API as any Kubernetes |
| Logs | Loki + Alloy | Lightest stack that still gives label-based log search |
| Secrets | SOPS + age | Works with plain git and Argo CD; no cloud KMS needed |
| Releases | GoReleaser + cosign + Syft | Standard Go supply-chain trio; keyless signing needs no key management |
