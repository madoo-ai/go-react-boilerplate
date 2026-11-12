# Architecture

## Summary

The repository has one Go module. The Go server owns the JSON API and serves
the React build from the same origin. The build is embedded with `go:embed`,
thus the production image has no Node runtime.

```
web/ (React + Vite) ──build──> internal/handler/webui/dist ──embed──> Go binary
                                                                       │
browser ── /api/* ─────────────────────────────────────────────────────┤
        └─ everything else ────────────────────────────────────────────┘
```

## Layers

| Layer                      | Responsibility                                      |
| -------------------------- | --------------------------------------------------- |
| `internal/entity`          | Domain types, RBAC policies, validation helpers     |
| `internal/app/port`        | Interfaces owned by the application                 |
| `internal/app/usecase`     | Business rules and orchestration                    |
| `internal/adapter/*`       | Mechanics: Postgres, SMTP, Argon2id, TOTP, limiter  |
| `internal/handler/httpapi` | JSON HTTP adapter                                   |
| `internal/handler/webui`   | Embedded React build and SPA fallback               |
| `internal/xmodule`         | Config, logger, composition, HTTP server            |

Dependencies point inward:

```
handler → usecase → port ← adapter
```

The depguard rules in `.golangci.yml` enforce this direction.

## Request lifecycle

1. `cmd/server` reads flags and environment.
2. `xmodule.InitDependencies` applies migrations and opens Postgres.
3. The composition root builds adapters, use cases and handlers.
4. `xmodule.Serve` starts one listener.
5. `/api/*` goes to `httpapi`. The rest goes to `webui`.
6. `httpapi` loads the session cookie, calls a use case and writes JSON.
7. `webui` serves a prerendered page or the SPA shell.

## API error model

Every error is a translation key:

```json
{ "error": "errors.forbidden" }
```

Validation errors also carry field keys:

```json
{
  "error": "validation.failed",
  "fields": { "email": "validation.email" }
}
```

The React client renders the key with the active locale. This keeps the API
language-neutral.

## Sessions

Sessions are rows in Postgres. The cookie contains a random token. The
database stores only the SHA-256 hash of the token. The session join loads
the user, role and permissions in one request.

`httpapi` reads the cookie on every API request. It does not trust a user ID
from the client.

## Mutations and CSRF

The SPA uses cookie authentication. To block cross-site writes:

- Every mutation must send `X-UI-Request: 1`.
- If the request has an `Origin` header, the origin must match the request
  host, `APP_URL`, or an entry in `TRUSTED_ORIGINS`.

A cross-origin form cannot set the custom header. A cross-origin fetch that
sets the header needs a CORS preflight, and the server does not allow CORS.

## Database

- PostgreSQL is the only store.
- Migrations are embedded and run at boot.
- gorm rows stay inside `internal/adapter/postgres`.
- UUID primary keys are generated in Go.
- `settings` is a key/value table so small switches need no migration.

## Frontend

- React Router owns client routing.
- TanStack Query owns server state.
- React Hook Form and Zod own form state.
- The auth provider loads `/api/auth/session` once and caches it.
- Mutations update or invalidate the relevant query keys.
- The theme and locale are cookies, not database rows.

## Prerendering

The build has three steps:

1. Vite builds the client into `internal/handler/webui/dist`.
2. Vite builds `src/entry-server.tsx` into `web/dist-ssr`.
3. `web/scripts/prerender.mjs` renders the public routes and writes real HTML
   files.

`webui` serves `/<route>/index.html` when it exists, then falls back to the
SPA shell. The client uses `createRoot`, thus a prerendered page becomes
interactive without a hydration contract.

## Testing

- Go unit tests cover policies, config and security primitives.
- Vitest covers i18n and helpers.
- Playwright covers the browser flows against a running server.
