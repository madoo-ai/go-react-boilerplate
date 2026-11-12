---
name: grb-architecture
description: Use when adding or changing a feature in the Go React boilerplate, deciding where code belongs, or reviewing the hexagonal layer boundaries. Covers entity, port, use case, adapter, HTTP handler, React route and test placement.
---

# Go React boilerplate architecture

The project has one Go module and one React app. The Go server owns `/api/*`
and the embedded SPA on one origin. The layer law is in `AGENTS.md` and is
checked by depguard in `.golangci.yml`.

## Where code goes

| Change                         | Location                                  |
| ------------------------------ | ----------------------------------------- |
| Domain type or RBAC policy     | `internal/entity/`                        |
| Repository or service contract | `internal/app/port/`                      |
| Business rule or orchestration | `internal/app/usecase/<domain>/`          |
| SQL, gorm rows, migrations     | `internal/adapter/postgres/`              |
| Password, token, TOTP, crypto  | `internal/adapter/security/`              |
| SMTP                           | `internal/adapter/mailer/`                |
| JSON route and DTO             | `internal/handler/httpapi/`               |
| Embedded SPA and SPA fallback  | `internal/handler/webui/`                 |
| Wiring and config              | `internal/xmodule/`                       |
| React page                     | `web/src/routes/`                         |
| React domain component         | `web/src/components/<domain>/`             |
| API call and DTO types         | `web/src/api/`                            |
| Shared hook or helper          | `web/src/lib/`                            |
| Message catalogue              | `web/src/messages/{en,id}.json`           |

Dependencies point inward:

```
handler → usecase → port ← adapter
```

## Add a feature end to end

1. Add the domain type or policy in `internal/entity/`.
2. Add the smallest port method the use case needs in `internal/app/port/`.
3. Put the business rule in `internal/app/usecase/<domain>/`.
4. Implement the port in `internal/adapter/postgres/` or another adapter.
5. Add or extend a migration when the schema changes.
6. Add a DTO and handler in `internal/handler/httpapi/`.
7. Register the route and the required permission.
8. Add the API call and types in `web/src/api/`.
9. Build the page or component with TanStack Query and the UI primitives.
10. Add the message keys to both catalogues.
11. Add Go, Vitest or Playwright coverage at the right level.
12. Run `make lint` and `make test`.

## Hard rules

- Do not import `internal/adapter` from `entity`, `port` or `usecase`.
- Do not import gorm outside `internal/adapter/postgres`, `internal/xmodule`
  or `cmd`.
- Do not put business rules in `internal/handler/httpapi`.
- Do not instantiate adapters outside `internal/xmodule`.
- Do not hand-edit `internal/handler/webui/dist`.
- React server state goes through TanStack Query. Do not add a second store.
- Use the existing `web/src/components/ui` primitives.

## Orientation

- `README.md` is the product and setup guide.
- `docs/architecture.md` is the request lifecycle.
- `docs/adr/0001-hexagonal-layers.md` explains the structure.
- `internal/xmodule/dependencies.go` is the only composition root.
