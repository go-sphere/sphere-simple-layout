# Sphere Simple Layout

`sphere-simple-layout` is the smallest official Sphere project template. It
demonstrates the Proto-first HTTP pipeline with Gin and Wire without a database,
authentication, Swagger, dashboard, provider SDK, or deployment integration.

## Capabilities

- Protobuf and Buf API contracts.
- Generated Gin-compatible HTTP handlers.
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
