# Go React Boilerplate — agent guide

This is a Go API and a React SPA in one repository, built and served as one
binary. Treat it as the product, not a pile of generated files.

Read `README.md`, `DESIGN.md`, `docs/architecture.md` and this file before
changing the shape of the app.

## Boilerplate, not the product

This repository is a starting point. If the next product does not need a part
of it, rework or remove it. Do not preserve a feature only because it is
already here.

Before deleting anything, check:

- routes and layouts
- migrations and seed data
- use cases and handlers
- navigation and permission checks
- shared components and translations
- tests

After a structural change, run `make lint` and `make test`.

## Map

| Need                              | Read                                                        |
| --------------------------------- | ----------------------------------------------------------- |
| Product scope and setup           | `README.md`                                                 |
| Request lifecycle and layers      | `docs/architecture.md`                                      |
| Auth, approval and RBAC           | `docs/auth-and-rbac.md`                                     |
| User and role management          | `docs/users-and-roles.md`                                   |
| Settings, 2FA and deletion        | `docs/settings.md`                                          |
| Locale system                     | `docs/internationalization.md`                              |
| Every route and permission        | `docs/routes-reference.md`                                  |
| Demo script                       | `docs/demo-guide.md`                                        |
| Architecture decisions            | `docs/adr/`                                                 |
| Visual rules                      | `DESIGN.md`                                                 |
| Server entrypoint                 | `cmd/server/main.go`                                        |
| Composition root                  | `internal/xmodule/`                                         |
| Domain and RBAC policies          | `internal/entity/`                                          |
| Ports                             | `internal/app/port/`                                        |
| Use cases                         | `internal/app/usecase/`                                     |
| gorm stores and migrations        | `internal/adapter/postgres/`                                |
| Argon2id, tokens, TOTP, AES-GCM   | `internal/adapter/security/`                                |
| JSON API                          | `internal/handler/httpapi/`                                 |
| Embedded React build              | `internal/handler/webui/`                                   |
| React app                         | `web/src/`                                                  |
| UI primitives                     | `web/src/components/ui/`                                    |
| Messages                          | `web/src/messages/{en,id}.json`                             |
| Tests                             | `internal/**/*_test.go`, `web/src/**/*.test.*`, `web/tests/`|
| Agent skills                      | `.agents/skills/`, `docs/agent-skills.md`                   |

## Commands

```sh
make setup          # install, .env, Postgres/Mailpit, migrate, seed, build UI
make dev            # Go API on :8080 + Vite HMR on :5173
make build          # build the UI and binaries into bin/
make test           # Go unit tests + Vitest
make test-e2e       # Playwright against a running app
make lint           # go vet, gofmt, TypeScript, i18n parity, golangci-lint
make format         # gofmt
make db-migrate     # apply embedded migrations
make db-seed        # idempotent seed
make db-reset       # destructive database rebuild
docker compose up --build
```

## Architecture rules

- `cmd/` is thin: flags, environment, signal handling, then `xmodule`.
- `internal/entity` is pure Go. No gorm, HTTP, database, or adapter imports.
  RBAC policies live here as functions.
- `internal/app/port` owns the interfaces. Adapters implement them.
- `internal/app/usecase` contains business rules and orchestration. It may use
  ports and entities only. No gorm, HTTP, or adapters.
- `internal/handler/httpapi` maps JSON ↔ use cases. No business rules, no
  storage imports.
- `internal/adapter/postgres` is the only package allowed to import gorm.
  Storage rows never escape it; everything crosses as entities.
- `internal/xmodule` is the only package that knows concrete adapters exist.
  Wiring happens there and nowhere else.
- The layer law is checked by depguard in `.golangci.yml`.

## HTTP rules

- One listener, one origin. `/api/*` is JSON; everything else is the SPA.
- `webui` serves hashed assets as immutable and the document shell as
  `no-store`.
- Mutations require `X-UI-Request: 1` and a same-origin `Origin`.
- Session cookies are `HttpOnly`, `SameSite=Lax`; `Secure` follows `APP_URL`.
- Errors are translation keys: `{ "error": "errors.forbidden" }`, plus
  `fields` for validation. Do not invent English messages in the API.
- Never trust a client-supplied user ID or role. Load the session from the
  cookie and re-check the permission close to the data.

## Data rules

- PostgreSQL is the only database.
- Schema changes need a new sequential migration in
  `internal/adapter/postgres/migrations/`. Never edit an applied migration.
- Table names are snake_case. Primary keys are UUIDs generated in Go.
- The seed is idempotent. Demo users use the password `password`.
- `internal/handler/webui/dist` is build output; do not hand-edit it.

## Auth rules

- Emails are normalized to lower case before every write.
- Passwords are 8–128 characters and Argon2id hashed. Never log a password or
  a raw token.
- Sessions store only the SHA-256 hash of the cookie token.
- Email verification and password reset tokens are one-shot and hashed.
- Registration approval is read from the `settings` table in the register
  use case, never from a client value.
- Pending and suspended users cannot create a session.
- Only `manage_roles` may assign `superadmin`. The last superadmin cannot be
  removed or demoted. HR cannot manage a superadmin.
- Sensitive actions use the password-confirmation cookie (15 minutes).
- TOTP secrets are encrypted with a key derived from `AUTH_SECRET`.

## Frontend rules

- Build on `web/src/components/ui`; do not add a second component library.
- Icons come from `@x1nx3r/altics-react-icons`. Do not import `lucide-react`.
- Use semantic tokens (`bg-background`, `text-muted-foreground`,
  `border-border`), not raw palette classes, for app surfaces.
- Use `FieldGroup` + `Field` for forms, `gap-*` instead of `space-y-*`,
  `size-*` for equal dimensions, and `cn()` for conditional classes.
- Every page has one `h1`, rendered by `PageHeader` or the auth card.
- Pages use the full content width; constrain the grid, not the page.
- The primary table sits in a bare `overflow-x-auto rounded-lg border`
  wrapper.
- Use the full Card composition (`CardHeader`, `CardContent`, `CardFooter`).
- Overlays need a title. Use `AlertDialog` for destructive confirmation.
- Send feedback through Sonner toasts. Use `FormAlert` for form-level errors.
- Server data goes through TanStack Query; mutations invalidate the relevant
  keys.
- Product screens use Inter. Marketing and auth editorial headings use
  `font-display` (Instrument Serif).
- Motion is CSS-only and compositor-friendly: animate `transform` and
  `opacity`, gate it with `motion-safe`, and keep the reduced-motion path
  fully usable.
- `web/src/lib/theme-options.ts` and `themes.css` are generated. Edit
  `opencode-themes.json` and regenerate; do not hand-edit the generated files.

## i18n rules

- Every user-facing string goes through `useTranslations`. Add keys to both
  `web/src/messages/en.json` and `id.json`.
- Do not add a `[locale]` segment. The `locale` cookie is the source of truth.
- Zod messages are translation keys (`validation.*`). Translate them in the
  form with `translateFieldError`.
- API error codes are translation keys; render them with the global `t`.
- Run `npm run check:i18n` after changing messages.

## Testing rules

- Unit-test RBAC policies, validation helpers, security primitives and pure
  helpers.
- E2E-test auth, approval, user CRUD, roles and settings with Playwright.
- Do not disable security checks to pass a test. Fix the flow.
- The E2E suite expects a running app; `make test-e2e` defaults to
  `http://127.0.0.1:3000`.

## Demo accounts

All demo accounts use the password `password`.

| Role       | Email                    |
| ---------- | ------------------------ |
| superadmin | `superadmin@example.com` |
| hr         | `hr@example.com`         |
| user       | `user@example.com`       |

The seed creates additional active, pending and suspended users.

## Agent skills

The repository ships vendor and project skills under `.agents/skills/`.
`.claude/skills/` and `.windsurf/skills/` contain relative symlinks.
`skills-lock.json` pins the vendor sources and hashes.

- Vendor: 11 Go skills from `samber/cc-skills-golang`, `shadcn`, four Docker
  skills from `docker/skills`, and `playwright-cli` from
  `microsoft/playwright-cli`.
- Project: `grb-architecture`, `grb-auth-rbac`, `grb-api-contract`,
  `grb-postgres-migrations`, `grb-release-pipeline`, `grb-i18n-themes`,
  `grb-testing`.

Load a skill when the task matches its description. Do not add the Next.js,
Prisma or Better Auth vendor skills; they do not apply. Review a vendor skill
before you trust it. See `docs/agent-skills.md` for the inventory, update
commands and licenses.
