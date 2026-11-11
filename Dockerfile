# syntax=docker/dockerfile:1

# Multi-stage build for the Go binaries. One image per binary, selected with
# --build-arg CMD=api|worker. Static, stripped, non-root, no shell.
# Base images are pinned by digest (Dependabot keeps them current), the same
# reason every GitHub Action in .github/workflows is pinned by SHA.

FROM golang:1.24-alpine@sha256:8bee1901f1e530bfb4a7850aa7a479d17ae3a18beb6e09064ed54cfd245b7191 AS build
ARG CMD=api
WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/${CMD}

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
ARG CMD=api
LABEL org.opencontainers.image.source="https://github.com/rj-karan/AegisFlow" \
      org.opencontainers.image.title="aegisflow-${CMD}" \
      org.opencontainers.image.licenses="Apache-2.0"
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
