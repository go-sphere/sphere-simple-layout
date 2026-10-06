# Sphere Simple Layout

`sphere-simple-layout` is the smallest official Sphere project template. It
demonstrates the Proto-first HTTP pipeline with the stdx (net/http) engine and
Wire, without a database, authentication, Swagger, dashboard, provider SDK, or
deployment integration.

## Capabilities

- Protobuf and Buf API contracts.
- Generated HTTP handlers on the stdx (net/http) engine.
- A minimal greet service.
- Wire dependency injection.
- Docker and multi-architecture build targets.

## Workflow

```shell
make init
make run
```

During development use `make gen/all`, `make check`, and `make build`. Run
`make help` for the exact supported targets. This layout intentionally has no
`gen/docs`, `run/swag`, `gen/db`, or `deploy` target.

## Structure and Ownership

- `proto/**` contains handwritten API contracts.
- `api/**` is generated.
- `internal/service/**` contains service implementations.
- `internal/server/**` contains HTTP construction and route registration.
- `cmd/app` composes the application through Wire.

Read `.sphere/layout.json` and `AGENTS.md` before extending or synchronizing the
layout. Unclassified paths are project-owned by default.

## Generator Versions and Codegen Baseline

`codegen.versions` pins every tool `make install` installs: the go-sphere
protoc plugins, Buf, protoc-gen-go, Swag, Wire, golangci-lint, and sphere-cli.
Change versions only there; the Makefile and the codegen scripts both read it.

Generated `api/**` is not committed, so `codegen.sha256` records the SHA-256
of every generated `api/**` file as the tracked regression baseline:

- `make codegen-check` compares the current `api/**` with the baseline. The CI
  workflow runs it right after `make gen/all`.
- `make codegen-verify` regenerates `api/**` with the pinned plugins in a
  temporary tool directory and checks that the generated packages build, that
  generation is idempotent, and that the output matches the baseline. The
  Codegen workflow runs it on every push.
- `make codegen-baseline` rewrites the baseline from the current `api/**`.

After changing Proto files or bumping a generator, run
`make install && make gen/all && make codegen-baseline` and commit
`codegen.sha256` with the change, so reviewers see which generated files
moved. A failing check without such a change means the installed tools do not
match `codegen.versions`.

## Upgrade Notes

This revision is a breaking template change: generated handlers are served by
the `stdx` (net/http) engine from `httpx` / `httpx/stdx` v0.0.5, pinned together
with `github.com/go-sphere/sphere` v0.0.6. The Gin adapter is gone.

- Generated `api/**` code and its error envelopes use the `httpz.*` types
  instead of the former `ginx.*` names.
- `internal/pkg/httpsrv` exports `NewServer(name, addr) httpx.Engine`, backed by
  `stdx` over `net/http`, and `UseCORS` registers CORS on that engine.
- Middleware is registered on the engine rather than a Gin router, so CORS and
  panic recovery also cover paths no route matched.

A project generated from an earlier revision should merge
`internal/pkg/httpsrv/**`, `internal/server/*/web.go`, and the regenerated
`api/**` outputs at its next layout sync, then run `make gen/all` and
`make check`.
