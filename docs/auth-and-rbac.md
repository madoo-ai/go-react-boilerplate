# Auth and RBAC

## Registration

1. The client sends `POST /api/auth/register`.
2. The API normalizes the email and validates the name and password.
3. The API checks that the email is available.
4. The API reads `settings.registration_approval`.
5. The user status becomes `pending` or `active`.
6. The API hashes the password with Argon2id.
7. The API creates the user and the credential in one transaction.
8. The API creates a one-shot verification token and sends the email.
9. The client signs in when the status is `active`.
10. The client redirects to login with an approval message when the status is
    `pending`.

`autoSignIn` does not exist. The client must call `/api/auth/login` after a
successful registration.

## Login

1. The API applies the login rate limit.
2. The API loads the user and the password hash.
3. If the user does not exist, the API still verifies a dummy hash. This keeps
   the response time similar.
4. The API rejects `pending` and `suspended` accounts.
5. If TOTP is enabled, the API returns `twoFactorRequired` and sets a signed
   challenge cookie for 5 minutes.
6. Otherwise the API creates a session row and sets the session cookie.

Session lifetime is 7 days. A session updates its expiry at most once a day.
When `rememberMe` is false, the cookie has no `Max-Age` and ends with the
browser session.

## Two-factor authentication

Enable flow:

1. The user sends the current password to
   `POST /api/settings/two-factor/enable`.
2. The API creates a TOTP secret and 10 backup codes.
3. The API encrypts the secret with AES-256-GCM.
4. The API stores SHA-256 hashes of the backup codes.
5. The client shows the QR code, the secret and the backup codes once.
6. The user sends a 6-digit code to
   `POST /api/settings/two-factor/confirm`.
7. The API marks the factor as verified and enables the user flag.

Login flow:

1. The password check succeeds.
2. The API returns a challenge instead of a session.
3. The client sends the TOTP code or a backup code to
   `POST /api/auth/two-factor/verify`.
4. The API validates the challenge, the code and the account status.
5. The API creates the session and clears the challenge.

After 5 failed codes, the factor locks for 5 minutes. A backup code is removed
after it is used.

## Email verification

- The register use case creates a token for `email_verification`.
- The email link points to `APP_URL/verify-email?token=...`.
- The SPA sends the token to `POST /api/auth/verify-email`.
- The API consumes the token and sets `email_verified`.
- Tokens expire after one hour.
- The mailer logs the link when SMTP is not configured.

Login does not require a verified email. The user can resend the link from the
verify page.

## Password reset

1. The client sends the email to `POST /api/auth/forgot-password`.
2. The API always reports success. This prevents account enumeration.
3. If the user exists, the API creates a one-shot reset token.
4. The email link points to `APP_URL/reset-password?token=...`.
5. The client sends the token and the new password to
   `POST /api/auth/reset-password`.
6. The API updates the hash and revokes every session for that user.

## Confirm-password gate

Some actions need a recent password check. Account deletion uses it.

1. The user sends the password to `POST /api/auth/confirm-password`.
2. The API verifies the password.
3. The API sets `gr_password_confirmed` with an HMAC value for the user and the
   issue time.
4. The cookie lasts 15 minutes.
5. The delete endpoint checks the cookie again on the server.

The cookie is `HttpOnly`. JavaScript cannot read or change the value.

## RBAC model

A user has one role. Roles have many permissions.

| Permission     | Label                         | Superadmin | HR  | User |
| -------------- | ----------------------------- | ---------- | --- | ---- |
| `view_admin`   | Access the admin area         | yes        | yes | no   |
| `manage_users` | Create, edit and delete users | yes        | yes | no   |
| `manage_roles` | Create, edit and delete roles | yes        | no  | no   |

Policies live in `internal/entity/authorize.go`:

- `CanManageUser`: an actor needs `manage_users`. Only a superadmin can manage
  a superadmin.
- `CanDeleteUser`: adds self-delete and last-superadmin protection.
- `CanAssignRole`: only `manage_roles` can assign `superadmin`.
- `CanManageRole`: `manage_roles` is necessary. The superadmin role is
  immutable.

The use cases call these policies close to the data. The sidebar is a
presentation filter, not a security boundary.

## Navigation gating

`web/src/config/navigation.ts` maps items to permissions. The sidebar filters
items with the permissions from the session response.

## Session invalidation

- Logout deletes the session row.
- Password reset deletes every session for the user.
- A password change deletes every other session.
- A status change to `pending` or `suspended` deletes every session.
- User deletion cascades to sessions and credentials.
