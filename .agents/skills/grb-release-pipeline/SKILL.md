---
name: grb-release-pipeline
description: Use when changing the Vite build, prerendering, go:embed, the Dockerfile, Compose, runtime app name or locale injection, or the release and health-check path. Covers the single-binary pipeline from React source to a running container.
---

# Build and release pipeline

The production artifact is one Go binary. It embeds the React build from
`internal/handler/webui/dist`.

## Pipeline

1. `web/index.html` is the shell. It contains the `__APP_NAME__` and
   `__APP_LOCALE__` placeholders.
2. Vite builds the client into `internal/handler/webui/dist`.
3. Vite builds `web/src/entry-server.tsx` into `web/dist-ssr`.
4. `web/scripts/prerender.mjs` renders `/`, `/login`, `/register`,
   `/forgot-password`, `/reset-password`, `/verify-email`, `/two-factor` and
   `/403`.
5. Go embeds the whole `dist` directory in `internal/handler/webui/server.go`.
6. At serve time, Go replaces the app-name and locale placeholders.
7. `webui` serves a prerendered page when it exists, then falls back to the
   SPA shell.

## Commands

```bash
make webui        # build client, SSR bundle and prerendered pages
make build        # webui plus all Go binaries
make docker-up    # full Compose stack
```

Do not hand-edit `internal/handler/webui/dist`.

## Runtime configuration

- `APP_NAME` is injected into every HTML document.
- `APP_LOCALE` is injected into every HTML document.
- `APP_URL` is used for email links, origin checks and secure cookies.
- `LISTEN` is the bind address.
- Compose uses `COMPOSE_APP_URL` to avoid inheriting the dev `APP_URL`.

## Docker

- Stage `web` installs with `npm ci` and runs the build.
- Stage `gobuild` downloads modules, copies the source and the built UI, and
  compiles `server`, `seed` and `migrate`.
- Stage `runner` is Alpine, runs as uid 10001 and exposes 3000.
- The `seed` service must override the entrypoint with `["/app/seed"]`.
- The app health check calls `/api/health`.
- Migrations run at boot and in the seed service.

## Release checklist

- [ ] `make lint` passes.
- [ ] `make test` passes.
- [ ] Playwright passes against the built app.
- [ ] `internal/handler/webui/dist` contains the new hashed assets.
- [ ] `docker compose build` succeeds.
- [ ] `docker compose up -d --wait` reports healthy.
- [ ] `GET /api/health` returns `database: up`.
- [ ] Email links point at `APP_URL` or `COMPOSE_APP_URL`.
- [ ] `AUTH_SECRET` is at least 32 characters in production.
- [ ] `APP_URL` uses HTTPS so cookies are Secure.
