# Users and roles

## User list

`GET /api/users` requires `manage_users`.

Query parameters:

| Parameter | Meaning                                   |
| --------- | ----------------------------------------- |
| `page`    | 1-based page number, 10 users per page    |
| `search`  | Case-insensitive name or email match      |
| `role`    | Role name filter                          |
| `status`  | `pending`, `active` or `suspended`        |

The response includes:

- the page of users
- `total`, `page`, `pageCount`, `pageSize`
- all roles for the filter
- assignable roles for the actor
- the registration approval flag
- the superadmin count
- `canManage` and `canDelete` for each row

## Create a user

`POST /api/users` requires `manage_users`.

The actor must be allowed to assign the selected role. The email must be
available. The API creates the user and the credential in one transaction.
The email starts unverified.

## Edit a user

`PATCH /api/users/{id}` requires `manage_users`.

Rules:

- The actor must pass `CanManageUser`.
- An actor cannot change its own role or status.
- Only a superadmin can edit a superadmin.
- The last superadmin cannot be demoted.
- An email change resets `email_verified`.
- An optional password field replaces the credential.
- A status other than `active` deletes the user sessions.

## Delete a user

`DELETE /api/users/{id}` requires `manage_users`.

Rules:

- An actor cannot delete itself.
- An actor must pass `CanManageUser`.
- The last superadmin cannot be deleted.

The database cascades the delete to credentials, sessions and TOTP rows.

## Verify a pending user

`POST /api/users/{id}/verify` requires `manage_users`.

The endpoint changes `pending` to `active`. It does nothing for other
statuses. This is the approval step for registration approval mode.

## Registration approval

`PATCH /api/admin/registration-approval` requires `manage_users`.

The endpoint writes `settings.registration_approval` as `true` or `false`.
The register use case reads this value on the server. A client cannot override
it.

## Roles

`GET /api/roles` requires `manage_roles`.

The response contains every role with its permissions and user count, plus
the full permission registry.

## Create a role

`POST /api/roles` requires `manage_roles`.

Rules:

- The name uses lowercase letters, numbers and underscores.
- The name must be unique.
- The label is required.
- Unknown permission keys are ignored.
- The new role is not a system role.

## Edit a role

`PATCH /api/roles/{id}` requires `manage_roles`.

Rules:

- The `superadmin` role is immutable.
- The `user` role can change its label and permissions, but not its name.
- A duplicate name is rejected.
- The API replaces the permission matrix in one transaction.

## Delete a role

`DELETE /api/roles/{id}` requires `manage_roles`.

Rules:

- The `superadmin` role cannot be deleted.
- The `user` role cannot be deleted.
- A role with users cannot be deleted.

## Protected entities

| Entity            | Rule                                              |
| ----------------- | ------------------------------------------------- |
| `superadmin` role | Immutable. Only a role manager can assign it.     |
| `user` role       | Default role. The name cannot change.             |
| Last superadmin   | Cannot be deleted or demoted.                     |
| Self              | Cannot be deleted from the admin screen.          |
