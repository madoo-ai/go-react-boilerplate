# 0003 — Cookie sessions, Argon2id and app-owned auth

## Status

Accepted.

## Context

The Next.js app uses Better Auth. There is no direct Go equivalent. The auth
needs are known: email/password, verification, reset, 2FA and approval.

## Decision

Implement auth in the application:

- Argon2id password hashes in PHC format.
- Server-side session rows. The cookie holds a random token. The database
  stores its SHA-256 hash.
- One-shot hashed tokens for email verification and password reset.
- TOTP with encrypted secrets and hashed backup codes.
- A signed password-confirmation cookie for sensitive actions.
- An in-memory rate limiter for auth endpoints.

## Consequences

- Revocation, approval and 2FA state are simple and auditable.
- The team owns the security code and must review it.
- Session checks cost one database read per API request. The join loads the
  role and permissions at the same time.
- A future OAuth provider can add an identity table without a session model
  change.
