---
name: grb-auth-rbac
description: Use when changing authentication, sessions, email verification, password reset, TOTP two-factor authentication, registration approval, roles, permissions, user access checks or account deletion. Covers the app-owned auth model and the RBAC rules that must stay true.
---

# Auth and RBAC

Auth is app-owned. There is no Better Auth, NextAuth or external identity
service. The main files are:

- `internal/app/usecase/auth/` — register, login, verification, reset, 2FA
- `internal/app/usecase/users/` — account management
- `internal/app/usecase/roles/` — role matrix
- `internal/entity/authorize.go` — RBAC policies
- `internal/handler/httpapi/auth.go` — cookies and auth routes
- `docs/auth-and-rbac.md` — the full flow

## Session model

- The cookie is `gr_session`. It contains a random 32-byte token.
- The `sessions` table stores only the SHA-256 hash of the token.
- The cookie is `HttpOnly`, `SameSite=Lax`. `Secure` follows `APP_URL`.
- A session is valid for five days and slides forward at most once a day.
- Logout deletes the row. Password reset and status changes revoke sessions.

## Login and 2FA

1. Rate-limit the login attempt by IP.
2. Verify the password with Argon2id.
3. Block `pending` and `suspended` accounts.
4. If TOTP is verified, return a short-lived signed challenge and set
   `gr_2fa`. Do not create a session yet.
5. On `/api/auth/two-factor/verify`, validate the TOTP or a backup code, then
   create the session.
6. Backup codes are single-use. Five failures lock the factor for five
   minutes.

## Email tokens

- Verification and reset tokens are random, one-shot and stored as hashes.
- Verification tokens expire in one hour.
- Reset tokens expire in one hour.
- A reset updates the password and revokes every session.
- The mailer logs the link when SMTP is not configured.

## Registration approval

- The switch writes `settings.registration_approval`.
- The register use case reads the setting on the server.
- `true` creates the account as `pending`.
- `pending` accounts cannot sign in.
- A user manager calls `POST /api/users/{id}/verify` to approve.
- Never accept an approval flag from a request body.

## RBAC rules

| Permission     | Superadmin | HR  | User |
| -------------- | ---------- | --- | ---- |
| `view_admin`   | yes        | yes | no   |
| `manage_users` | yes        | yes | no   |
| `manage_roles` | yes        | no  | no   |

Policies in `internal/entity/authorize.go`:

- HR may manage regular users. Only a superadmin may manage a superadmin.
- Only `manage_roles` may assign the `superadmin` role.
- The last superadmin cannot be deleted or demoted.
- An actor cannot delete itself from the admin screen.
- The `superadmin` role is immutable.
- The `user` role can change label and permissions, but not its name.
- A role with users cannot be deleted.

## Adding a permission

1. Add the constant and label to `internal/entity/role.go`.
2. Add it to `AllPermissions`.
3. Add it to the seeded role definitions where necessary.
4. Re-run `make db-seed`.
5. Gate the API route in `internal/handler/httpapi/server.go`.
6. Gate the navigation item in `web/src/config/navigation.ts`.
7. Gate the page with `RequirePermission` in `web/src/app.tsx`.
8. Add the label to both message catalogues under `permissions`.
9. Add a policy test in `internal/entity/authorize_test.go`.

## Sensitive actions

- Account deletion requires the `gr_password_confirmed` cookie.
- The cookie is signed, `HttpOnly` and lasts 15 minutes.
- Verify the password on the server again. Never trust a client claim.

## Tests

- Go unit tests cover the policies and security primitives.
- Playwright covers the superadmin, HR and user access differences.
- The approval test proves that a pending account is blocked and then
  admitted.
