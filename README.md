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
