# Settings

Every settings endpoint requires an active session.

## Profile

`PATCH /api/settings/profile`

Fields: name, email, phone, job title and department.

Rules:

- The name is 2 to 120 characters.
- The email must be valid and available.
- An email change resets `email_verified` and sends a new verification link.
- Empty optional fields are stored as null.
- The response is the fresh session payload, thus the UI can update its cache.

## Password

`POST /api/settings/password`

Fields: `currentPassword` and `newPassword`.

Rules:

- The new password is 8 to 128 characters.
- The current password must match.
- The API stores an Argon2id hash.
- Every other session is revoked. The current cookie stays valid.

## Appearance

Appearance is client-side only.

- `theme` cookie: `light`, `dark` or `system`.
- `color-theme` cookie: the palette ID.
- The locale lives in the `locale` cookie.

An inline script in `web/index.html` applies the mode and palette before first
paint. The server does not store these values.

## Security

### Enable 2FA

1. `POST /api/settings/two-factor/enable` with the current password.
2. The API returns the secret, the `otpauth://` URI, a QR data URL and 10
   backup codes.
3. Scan the QR code in an authenticator app.
4. `POST /api/settings/two-factor/confirm` with a 6-digit code.
5. The API marks the factor as verified.

The backup codes are shown once. The API stores only their hashes.

### Disable 2FA

`POST /api/settings/two-factor/disable` with the current password.

The API deletes the TOTP row and clears the user flag.

### Backup codes

Each code works once and is removed after use. Codes are normalized for case,
spaces and dashes before hashing.

## Account deletion

1. The user opens `/confirm-password?next=/settings/profile`.
2. The password confirmation endpoint sets a 15-minute cookie.
3. The profile page shows the delete action when the cookie is valid.
4. `DELETE /api/settings/account` checks the cookie again on the server.
5. The last superadmin cannot delete the account.
6. The API deletes the user. Sessions, credentials and TOTP rows cascade.
7. The client clears the session and goes to `/login?deleted=1`.

The delete endpoint never accepts a client claim that the password was
confirmed. It reads the signed cookie.

## Account overview

`GET /api/settings/account` returns:

```json
{ "canDelete": true, "confirmed": false }
```

`canDelete` is false for the last superadmin. `confirmed` reflects the
password-confirmation cookie.
