# 0001 — Hexagonal layers in one Go module

## Status

Accepted.

## Context

The Next.js sibling keeps server code, UI and operations in one repository.
The Go port needs the same small scope without premature service boundaries.

## Decision

Use one Go module with these layers:

- `internal/entity` for the domain
- `internal/app/port` for interfaces
- `internal/app/usecase` for business rules
- `internal/adapter/*` for mechanics
- `internal/handler/*` for delivery
- `internal/xmodule` for composition

No `go.work`. No second module. Depguard enforces the direction.

## Consequences

- The code is easy to read and test.
- Adapters can change without a use case change.
- A future service split has a clear seam at the ports.
- A small boilerplate pays a little structure cost. This is deliberate.
