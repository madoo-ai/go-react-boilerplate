# Demo guide

Start the app with the production-shaped stack:

```bash
docker compose up --build
```

Open http://localhost:3000.

All demo accounts use the password `password`.

| Account      | Email                    |
| ------------ | ------------------------ |
| Super Admin  | `superadmin@example.com` |
| HR Staff     | `hr@example.com`         |
| Regular User | `user@example.com`       |

## 1. Landing page

- Change the mode with the header toggle.
- Pick a palette in the live theme rail.
- Reload the page. The choice persists and there is no flash.
- Change the language. The page switches between Indonesian and English.

## 2. Superadmin

1. Sign in as `superadmin@example.com`.
2. The dashboard shows four stats, two charts and five recent users.
3. Open Users. Search for `Siti`, filter by role or status, and page through
   the list.
4. Create a user. Assign the `user` role and the `pending` status.
5. Use the row menu to verify the pending user. The status becomes `active`.
6. Turn on registration approval. Sign out and register a new account. The
   new account stays pending and cannot sign in.
7. Open Roles. Create a `support` role with `view_admin`. The permission
   matrix updates.
8. Open Settings. Change the profile, the palette and the password.
9. Open Security. Enable 2FA. Scan the QR code. Save the backup codes.

## 3. HR

1. Sign in as `hr@example.com`.
2. Users is available. Roles is not in the sidebar.
3. Open `/admin/roles` directly. The app shows the access-denied page.
4. Open a superadmin edit link. The API rejects the change.

## 4. Regular user

1. Sign in as `user@example.com`.
2. The dashboard shows the personal card, not the admin stats.
3. The Administration group is absent.
4. Settings still works.

## 5. Approval mode

1. Sign in as a superadmin.
2. Open Users and enable registration approval.
3. Sign out.
4. Register a new account.
5. The login page says that an administrator must approve the account.
6. Sign in as the superadmin and verify the account.
7. Sign in with the new account.

## 6. Reset and verification email

1. Mailpit is at http://localhost:8025.
2. Register an account. Open the verification email and follow the link.
3. Sign out and use Forgot password.
4. Open the reset email and set a new password.

## 7. Account deletion

1. Sign in.
2. Open Settings, Profile.
3. Use Confirm password. The button becomes available.
4. Delete the account. The app returns to login with a message.
5. The last superadmin cannot complete this step.
