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

The baseline's first line, `# module <path>`, names the Go module it was
recorded for: generated code embeds the module path, so a baseline only
matches its own module. `make init` records the baseline last, which gives a
project scaffolded under another module its own; `make codegen-check` asks for
`make codegen-baseline` when go.mod and the baseline disagree. `make init` uses
the committed `buf.lock` as is; refresh it deliberately with `make gen/deps`.

## HTTP Server Limits

`httpsrv.NewServer` takes `httpsrv.Options`, which is embedded in each
server's HTTP config, so these keys sit next to `address` and `cors` in
`api.http` of `config.json`. `0` selects the default and a negative value
disables the limit.

- `read_timeout_seconds` (default 30) bounds reading a whole request, body
  included; `idle_timeout_seconds` (default 120) bounds an idle keep-alive
  connection.
- `max_body_bytes` (default 4 MiB) caps a request body; reading past it
  fails, and a JSON request that does is answered with 413.
- `trusted_proxies` lists the reverse proxies (IPs or CIDRs) whose
  `X-Forwarded-For` header is honoured. Set it when the service runs behind a
  proxy: otherwise every client, and every per-IP rate limit, shares the
  proxy's address. An invalid entry fails configuration loading.

## Upgrade Notes

Template changes after the stdx migration:

- `internal/pkg/httpsrv` sets read and idle timeouts, a request body cap and
  trusted proxies from `httpsrv.Options` (see HTTP Server Limits), and
  `NewServer(name, addr, opts)` takes it as a third argument.
- The unused `environments` config field is removed; it had no effect, and an
  existing `config.json` that still carries it keeps loading.
- `.dockerignore` excludes `go.work` and `go.work.sum`, so a local workspace no
  longer breaks `docker build`.

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

The process timezone is now set explicitly: `cmd/app/main.go` calls
`boot.InitTimezone(boot.DefaultTimezone)` (`Asia/Shanghai`) before startup and
exits if the zone cannot be loaded. Newer sphere releases no longer set it from
a package `init()`, so a project that keeps an older `main.go` runs in the host
timezone after upgrading sphere. Merge `cmd/app/main.go` (it is `mixed`), then
pass another IANA zone or delete the call to keep the host default.
