# AI Agent Guide

## Layout Profile

This is the minimal Sphere layout: Protobuf/Buf, generated HTTP handlers, Gin,
Wire, Docker, and one greet example. It has no database, authentication,
Swagger, dashboard, Telegram, WeChat, or deployment scripts.

## Ownership and Extension

Read `.sphere/layout.json` before modifying files. Never edit generated paths
by hand. Treat mixed paths as three-way merge seams and assume every
unclassified path is project-owned.

Add contracts under `proto/<domain>/v1`, business logic under
`internal/biz/<domain>` when needed, and service implementations under
`internal/service/<domain>`. Product logic must not be placed in layout-owned
CI, generation, app-bootstrap, or HTTP-adapter files.

The canonical family rules and AI update algorithm live in
`sphere-layout/docs/LAYOUT_CONTRACT.md`.

## Workflow

- `make gen/all` regenerates Proto and Wire outputs.
- `make test` runs Go tests.
- `make lint` checks Go and Buf without rewriting files.
- `make check` verifies dependencies, formatting, lint, and tests.
- `make build` builds the application.

After changing constructors, provider sets, or Proto, regenerate before
testing. Delivery requires `make check`, `make build`, and review of tracked
generated changes.
