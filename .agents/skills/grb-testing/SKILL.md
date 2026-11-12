---
name: grb-testing
description: Use when adding or running tests for the Go React boilerplate. Covers Go unit tests, Vitest, Playwright, the demo accounts, rate limits, cleanup rules and the Docker test stack.
---

# Testing

Three layers, one command each.

## Go unit tests

- Same package as the code.
- Cover policies, validation, config and security primitives.
- Use table-driven tests.
- Do not mock gorm or HTTP internals. Put fakes at port boundaries.
- Run `go test ./cmd/... ./internal/...`.

## Vitest

- `web/src/**/*.test.{ts,tsx}`.
- Use Testing Library with `jsdom`.
- Cover i18n, helpers and pure UI logic.
- Run `cd web && npm test`.

## Playwright

- Specs: `web/tests/e2e/*.spec.ts`.
- Default base URL: `http://127.0.0.1:3000`.
- Override with `E2E_BASE_URL`.
- Demo accounts use password `password`.
- Set the `locale=en` cookie in `beforeEach` for English assertions.
- Use unique emails and role names with a timestamp.
- Clean up users, roles and settings in `finally`.
- Re-login as the superadmin before admin cleanup after an API test signs in
  as a lower-privileged account.
- The auth rate limit is 20 requests per minute per IP by default. Keep the
  suite under it or raise `RATE_LIMIT_PER_MINUTE` for a test run.

## Run

```bash
make test              # Go + Vitest
make test-e2e-install  # once
make test-e2e          # against a running app
```

For a full local check:

```bash
docker compose up --build -d
make test-e2e
```

## Before you finish

- [ ] `make lint` passes.
- [ ] `make test` passes.
- [ ] The affected Playwright spec passes.
- [ ] The seed is back to 14 users, 3 roles, approval off and 2FA off.
- [ ] No test artifacts remain in `web/test-results` or `web/playwright-report`.
