# Routes reference

## UI routes

| Path                       | Layout        | Access                        |
| -------------------------- | ------------- | ----------------------------- |
| `/`                        | Marketing     | Public, prerendered           |
| `/login`                   | Auth          | Public, prerendered           |
| `/register`                | Auth          | Public, prerendered           |
| `/forgot-password`         | Auth          | Public, prerendered           |
| `/reset-password`          | Auth          | Public, prerendered           |
| `/verify-email`            | Auth          | Public, prerendered           |
| `/two-factor`              | Auth          | Public, prerendered           |
| `/confirm-password`        | Auth          | Session required              |
| `/dashboard`               | App shell     | Session required              |
| `/admin/users`             | App shell     | `manage_users`                |
| `/admin/users/create`      | App shell     | `manage_users`                |
| `/admin/users/:id/edit`    | App shell     | `manage_users`                |
| `/admin/roles`             | App shell     | `manage_roles`                |
| `/settings`                | App shell     | Redirects to profile          |
| `/settings/profile`        | App shell     | Session required              |
| `/settings/password`       | App shell     | Session required              |
| `/settings/appearance`     | App shell     | Session required              |
| `/settings/security`       | App shell     | Session required              |
| `/403`                     | Standalone    | Public                        |
| `*`                        | Standalone    | Public, not found             |

## API routes

All API responses are JSON. Errors use `{ "error": "translation.key" }`.

### Public

| Method | Path                        | Purpose                              |
| ------ | --------------------------- | ------------------------------------ |
| GET    | `/api/health`               | Liveness and database readiness      |
| POST   | `/api/auth/register`        | Create an account                    |
| POST   | `/api/auth/login`           | Start a session or a 2FA challenge   |
| POST   | `/api/auth/logout`          | Delete the session and clear cookies |
| POST   | `/api/auth/verify-email`    | Consume an email verification token  |
| POST   | `/api/auth/forgot-password` | Request a reset link                 |
| POST   | `/api/auth/reset-password`  | Consume a reset token                |
| POST   | `/api/auth/two-factor/verify` | Finish a 2FA login                 |

### Authenticated

| Method | Path                                 | Purpose                          |
| ------ | ------------------------------------ | -------------------------------- |
| GET    | `/api/auth/session`                  | Current session and user         |
| POST   | `/api/auth/verify-email/request`     | Resend the verification email    |
| POST   | `/api/auth/confirm-password`         | Set the confirm-password cookie  |
| GET    | `/api/settings/account`              | Deletion state for the account   |
| PATCH  | `/api/settings/profile`              | Update the profile               |
| POST   | `/api/settings/password`             | Change the password              |
| POST   | `/api/settings/two-factor/enable`    | Start TOTP enrollment            |
| POST   | `/api/settings/two-factor/confirm`   | Verify the first TOTP code       |
| POST   | `/api/settings/two-factor/disable`   | Disable TOTP                     |
| DELETE | `/api/settings/account`              | Delete the account               |

### Admin

| Method | Path                                   | Permission     |
| ------ | -------------------------------------- | -------------- |
| GET    | `/api/dashboard`                       | `view_admin`   |
| GET    | `/api/users`                           | `manage_users` |
| POST   | `/api/users`                           | `manage_users` |
| GET    | `/api/users/{id}`                      | `manage_users` |
| PATCH  | `/api/users/{id}`                      | `manage_users` |
| DELETE | `/api/users/{id}`                      | `manage_users` |
| POST   | `/api/users/{id}/verify`               | `manage_users` |
| PATCH  | `/api/admin/registration-approval`     | `manage_users` |
| GET    | `/api/roles`                           | `manage_roles` |
| POST   | `/api/roles`                           | `manage_roles` |
| PATCH  | `/api/roles/{id}`                      | `manage_roles` |
| DELETE | `/api/roles/{id}`                      | `manage_roles` |

## Mutation contract

Every `POST`, `PATCH` and `DELETE` request must send:

```
Content-Type: application/json
X-UI-Request: 1
```

If the request has an `Origin` header, the origin must match the host,
`APP_URL`, or `TRUSTED_ORIGINS`.

## Cookies

| Cookie            | Purpose                               | Lifetime |
| ----------------- | ------------------------------------- | -------- |
| `gr_session`      | Session token                         | 7 days   |
| `gr_2fa`          | Signed 2FA challenge                  | 5 min    |
| `gr_password_confirmed`| Signed password confirmation          | 15 min   |
| `theme`           | Light, dark or system                 | 1 year   |
| `color-theme`     | Palette ID                            | 1 year   |
| `locale`          | `id` or `en`                          | 1 year   |
| `sidebar_state`   | Sidebar collapsed state               | 7 days   |
