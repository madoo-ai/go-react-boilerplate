# Go React Boilerplate

Opinionated Go + React boilerplate with cookie auth, RBAC, i18n, 38 color
themes and a shadcn/ui shell. The Next.js sibling provides the same product
surface; this one is a Go API with a static React app embedded in the binary.

This is a core-only scaffold: auth, roles and permissions, user management,
settings, locale, theme and the app shell. Domain modules are deliberately
left out so the next project can add its own.

## What is inside

- **One Go binary.** The Vite build is embedded with `go:embed`; the Go server
  owns `/api/*` and the SPA on one origin. Node is a development/build
  dependency only.
- **Auth.** Email/password, email verification, password reset, TOTP 2FA with
  backup codes, server-side sessions, rate limiting and a confirm-password
  gate.
- **Registration approval.** Optional approval mode keeps new accounts
  `pending` until a user manager verifies them.
- **RBAC.** Seeded `superadmin`, `hr` and `user` roles, a permission registry,
  user and role management, protected superadmins, navigation gating and
  privilege-escalation checks inside the Go use cases.
- **Users and roles.** Searchable, filterable and paginated users; create,
  edit, delete, verify; role permission matrix; protected system roles.
- **Settings.** Profile, password, appearance, security (2FA) and account
  deletion with a recent-password check.
- **i18n.** Indonesian default and English, cookie-based with no URL prefix.
  `npm run check:i18n` fails when the catalogues drift.
- **Appearance.** Light/dark/system plus the Stripe-inspired default and 37
  OpenCode palettes. Both choices are cookie-backed with no flash on load.
- **Landing page.** Public entry point with a live theme lab, a dashboard
  preview, the feature tour and demo accounts. Marketing and auth routes are
  prerendered at build time for better SEO.
- **Frontend.** React 19, Vite, React Router, Tailwind CSS 4, shadcn-style
  primitives, Altics icons, React Query, Sonner and Recharts.
- **Data.** PostgreSQL 18, gorm and embedded golang-migrate migrations.
- **Quality.** Go test, Vitest, Playwright, TypeScript strict, go vet,
  gofmt, golangci-lint layer rules and i18n parity.
- **Ops.** Multi-stage Dockerfile, Compose stack with Postgres and Mailpit,
  Makefile and ADRs.
- **Agent skills.** 24 versioned skills for Go, shadcn, Docker, Playwright and
  this project's architecture, auth/RBAC, API contract, migrations, release
  pipeline, i18n/themes and testing. They live under `.agents/skills/` and are
  linked into `.claude/skills/` and `.windsurf/skills/`.

## Stack

| Layer         | Choice                                                        |
| ------------- | ------------------------------------------------------------- |
| API           | Go 1.24+, stdlib `net/http`                                   |
| UI            | React 19, Vite 8, Tailwind CSS 4, shadcn-style primitives     |
| Data          | PostgreSQL 18, gorm, golang-migrate                           |
| Auth          | App-owned cookie sessions, Argon2id, TOTP (pquerna/otp)       |
| Authorization | App-owned RBAC tables and guards                              |
| Forms         | Zod 4, React Hook Form                                        |
| Server state  | TanStack Query                                                |
| i18n          | Cookie locale, typed message catalogues                       |
| Tests         | Go test, Vitest, Testing Library, Playwright                  |
| Runtime       | One Go binary, Alpine container                               |

## Demo accounts

All demo accounts use the password `password`.

| Account      | Email                    | Role         | Can do                            |
| ------------ | ------------------------ | ------------ | --------------------------------- |
| Super Admin  | `superadmin@example.com` | `superadmin` | Everything                        |
| HR Staff     | `hr@example.com`         | `hr`         | Manage users, cannot manage roles |
| Regular User | `user@example.com`       | `user`       | Self-service only                 |

The seed also creates eleven more users, including pending and suspended
accounts, so the dashboard and user list have realistic data.

## Requirements

- Go 1.24 or newer
- Node.js 22 or newer (build/dev only)
- Docker with Compose (for Postgres and Mailpit)

## Quick start

```bash
make setup   # install, create .env, start Postgres/Mailpit, migrate, seed, build UI
make dev     # Go API on :8080, Vite with HMR on :5173
```

Open http://localhost:5173. Vite proxies `/api` to the Go server.

Or run the production-shaped stack:

```bash
docker compose up --build
```

- App: http://localhost:3000
- Mailpit: http://localhost:8025
- Postgres: `127.0.0.1:54330` (Compose only)

The Compose `seed` service runs the migrations and the seed once before the
app starts.

## Commands

| Command              | Purpose                                                        |
| -------------------- | -------------------------------------------------------------- |
| `make setup`         | Install, write `.env`, start Postgres/Mailpit, migrate, seed    |
| `make dev`           | Go API + Vite HMR                                               |
| `make build`         | Build the React app and all Go binaries into `bin/`             |
| `make test`          | Go unit tests and Vitest                                        |
| `make test-e2e`      | Playwright against a running app (`E2E_BASE_URL`, default :3000)|
| `make lint`          | `go vet`, `gofmt` check, TypeScript and i18n parity             |
| `make format`        | `gofmt` the Go trees                                            |
| `make db-up`         | Start Postgres and Mailpit                                      |
| `make db-migrate`    | Apply embedded migrations                                       |
| `make db-seed`       | Seed permissions, roles, settings and demo users                |
| `make db-reset`      | Delete the volume and rebuild the database (destructive)        |
| `make docker-up`     | Build and start the full Compose stack                          |
| `make clean`         | Remove `bin/` and built React assets                            |

## Routes

| Route                    | Access                            |
| ------------------------ | --------------------------------- |
| `/`                      | Public landing page (prerendered) |
| `/login`                 | Public                            |
| `/register`              | Public                            |
| `/forgot-password`       | Public                            |
| `/reset-password`        | Public (token in query)           |
| `/verify-email`          | Public (verifies a token)         |
| `/two-factor`            | Public (challenge cookie)         |
| `/confirm-password`      | Authenticated                     |
| `/dashboard`             | Authenticated                     |
| `/admin/users`           | `manage_users`                    |
| `/admin/users/create`    | `manage_users`                    |
| `/admin/users/:id/edit`  | `manage_users`                    |
| `/admin/roles`           | `manage_roles`                    |
| `/settings/profile`      | Authenticated                     |
| `/settings/password`     | Authenticated                     |
| `/settings/appearance`   | Authenticated                     |
| `/settings/security`     | Authenticated                     |
| `/403`                   | Public error page                 |

The API lives under `/api`; see `docs/routes-reference.md` for the full list.

## Where things live

```
cmd/server/              daemon entrypoint
cmd/seed/                idempotent seed
cmd/migrate/             migration runner
internal/
    entity/              domain types, RBAC policies, pure helpers
    app/port/            driven port interfaces
    app/usecase/         auth, users, roles, settings, dashboard, health
    adapter/postgres/    gorm stores, embedded SQL migrations, seed
    adapter/security/    Argon2id, tokens, AES-GCM, TOTP
    adapter/mailer/      SMTP delivery with a logged fallback link
    adapter/ratelimit/   in-memory sliding window
    handler/httpapi/     cookie-authenticated JSON API
    handler/webui/       embedded Vite build + SPA fallback
    xmodule/             config, logger, composition root, HTTP server
web/                     React + Vite source
    src/components/ui/   shadcn-style primitives
    src/components/      app shell, auth, users, roles, settings, marketing
    src/routes/          pages
    src/lib/             API client, auth, i18n, theme, helpers
docs/                    architecture, flows, ADRs
DESIGN.md                visual contract
```

Layer flow:

```
handler/httpapi → app/usecase → app/port ← adapter/*
```

`internal/xmodule` is the only package that knows concrete adapters exist.

## Auth and RBAC

- Passwords use Argon2id. Only the PHC string is stored.
- Sessions are server-side rows; the cookie carries a random token and the
  database stores only its SHA-256 hash. The cookie is `HttpOnly`,
  `SameSite=Lax` and `Secure` when `APP_URL` is HTTPS.
- Mutations require `X-UI-Request: 1` and a same-origin `Origin`; this is the
  CSRF boundary for the SPA.
- Email verification and reset links are one-shot tokens; only their hashes
  are stored.
- TOTP secrets are encrypted at rest with AES-256-GCM derived from
  `AUTH_SECRET`. Backup codes are stored as SHA-256 hashes.
- Pending and suspended users cannot create a session.
- Only actors with `manage_roles` may assign `superadmin`; the last
  superadmin cannot be deleted or demoted.
- HR can manage regular users but cannot touch a superadmin.
- The `superadmin` role is immutable. The default `user` role can change its
  label and permissions, but not its name.
- Sensitive account deletion needs a password confirmation cookie that lasts
  15 minutes.

See `docs/auth-and-rbac.md` for the flows.

## i18n

- Messages: `web/src/messages/en.json` and `web/src/messages/id.json`.
- The `locale` cookie is the source of truth; there is no `/[locale]` path.
- `web/src/lib/i18n.tsx` provides `useTranslations` and interpolation.
- `npm run check:i18n` compares the two catalogues.

## Themes

- `web/src/lib/theme-options.ts` is generated from
  `web/src/lib/opencode-themes.json`.
- `web/src/styles/themes.css` holds the 38 palettes.
- The `theme` cookie stores the mode; `color-theme` stores the palette.
- An inline script in `web/index.html` applies both before first paint.

## Configuration

Copy `.env.example` to `.env` (the Makefile does this). The important keys:

| Key                 | Purpose                                              |
| ------------------- | ---------------------------------------------------- |
| `APP_NAME`          | Brand name in the UI and emails                      |
| `APP_URL`           | Public origin used in email links and origin checks  |
| `APP_LOCALE`        | Default locale (`id` or `en`)                        |
| `LISTEN`            | API listen address (`127.0.0.1:8080` in dev)         |
| `DATABASE_URL`      | PostgreSQL DSN                                       |
| `AUTH_SECRET`       | HMAC/cipher key, at least 32 characters              |
| `TRUSTED_ORIGINS`   | Comma-separated extra origins allowed to mutate      |
| `SMTP_*`, `MAIL_FROM` | Mail delivery; with no SMTP host the links are logged |
| `RATE_LIMIT_PER_MINUTE` | Auth endpoint budget per IP                      |
| `COOKIE_SECURE`     | Override the derived Secure cookie flag              |

In production set `APP_URL` to the HTTPS origin, `AUTH_SECRET` to a random
32+ character value and `APP_ENV=production`.

## Testing

```bash
make test        # Go unit tests + Vitest
make test-e2e    # Playwright; start the app first
make test-e2e-install  # once, to download Chromium
```

The Playwright suite covers the landing page, auth gates, login, the
superadmin/HR/user access differences, dashboard, users, roles, settings,
2FA and the registration approval flow. Set `E2E_BASE_URL` to target another
origin.

## Docker and deployment

- The Dockerfile builds the React app, embeds it in the Go binaries, and runs
  on Alpine as a non-root user.
- `docker compose up --build` starts Postgres, Mailpit, a one-shot seed and
  the app.
- The container listens on `0.0.0.0:3000`. The Compose file publishes it on
  loopback by default; put a reverse proxy in front for a public deployment.
- Set `COMPOSE_APP_URL=https://app.example.com` for the Compose stack. Email
  links, origin checks and secure cookies use that public URL.
- Migrations run automatically at boot. `cmd/migrate` is also available for
  release pipelines.

## Not included

Employee documents, attendance and leave requests are intentionally absent.
The Next.js, Prisma and Better Auth vendor skills from the sibling project do
not apply to this stack and are not shipped. Deploy pipelines, backups and
monitoring vary by host and are left out. The included GitHub Actions
workflow runs the quality gates only.

## Agent skills

The repository ships vendor and project skills under `.agents/skills/`.
`.claude/skills/` and `.windsurf/skills/` are relative symlinks to the same
directories. `skills-lock.json` pins the vendor sources and content hashes.

```bash
npx skills list
npx skills update
make skills
```

See `docs/agent-skills.md` for the inventory, sources, licenses and the
install commands. The project skills (`grb-*`) cover the Go layout, auth/RBAC,
API contract, migrations, release pipeline, i18n/themes and testing.

## Docs

- `docs/architecture.md`
- `docs/auth-and-rbac.md`
- `docs/users-and-roles.md`
- `docs/settings.md`
- `docs/internationalization.md`
- `docs/routes-reference.md`
- `docs/demo-guide.md`
- `docs/agent-skills.md`
- `docs/adr/README.md`
- `AGENTS.md`
- `DESIGN.md`

## License

MIT.
