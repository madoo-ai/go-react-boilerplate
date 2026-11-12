---
name: grb-api-contract
description: Use when adding or changing an API endpoint or a React data fetch. Covers DTOs, translation-key errors, validation fields, the X-UI-Request mutation guard, TanStack Query patterns and the session cache.
---

# API and React contract

The API lives in `internal/handler/httpapi`. The client lives in
`web/src/api`. The contract has three rules: translation keys, same-origin
mutations, and query-key invalidation.

## Errors

Every error is a translation key:

```json
{ "error": "errors.forbidden" }
```

Validation errors add field keys:

```json
{
  "error": "validation.failed",
  "fields": { "email": "validation.email" }
}
```

- Use `entity.NewError("errors.x")` or `entity.ValidationError(fields)`.
- `statusFor` maps known codes to HTTP statuses.
- Do not return English prose from the API.
- The client renders the key with `useTranslations()`.

## Mutations

Every `POST`, `PATCH` and `DELETE` must send:

```
Content-Type: application/json
X-UI-Request: 1
```

If an `Origin` header is present, it must match the request host, `APP_URL`
or `TRUSTED_ORIGINS`. The guard is in `internal/handler/httpapi/server.go`.
The client sets the header in `web/src/api/client.ts`.

## DTOs

- Internal entities use `time.Time` and empty strings.
- HTTP DTOs use JSON tags and `null` for empty optional strings.
- Map entities to DTOs in `internal/handler/httpapi/dto.go`.
- Never return gorm rows or internal structs directly.

## Adding an endpoint

1. Add the use case method and its validation.
2. Add the handler and register the route with the right permission.
3. Add the request and response DTO.
4. Add the method to `web/src/api/client.ts` and the types to
   `web/src/api/types.ts`.
5. Use `useQuery` for reads and `useMutation` for writes.
6. Invalidate every affected query key after a write.
7. Add the error keys to both message catalogues.
8. Add an E2E test when the flow crosses auth or permissions.

## Query keys

| Data             | Key                              |
| ---------------- | -------------------------------- |
| Session          | `["session"]`                    |
| Dashboard        | `["dashboard"]`                  |
| Users list       | `["users", page, search, ...]`   |
| User edit        | `["user", id]`                   |
| Role list        | `["roles"]`                      |
| Account overview | `["settings", "account"]`        |

Writes that change users or roles should also invalidate `["dashboard"]` and
`["session"]` when the current user could be affected.

## Session cache

`AuthProvider` loads `/api/auth/session` and stores the result under
`["session"]`. The login, profile update and 2FA responses return the fresh
session payload. Call `setSession(response)` instead of refetching.

## Status codes

| Code                         | Status |
| ---------------------------- | ------ |
| `errors.unauthorized`        | 401    |
| `errors.forbidden`           | 403    |
| `errors.userNotFound`        | 404    |
| `errors.emailInUse`          | 409    |
| `errors.tooManyRequests`     | 429    |
| `validation.failed`          | 400    |
| `errors.unexpected`          | 500    |

## Checklist

- [ ] The route has a permission gate if it is not self-service.
- [ ] The request is validated on the server.
- [ ] Errors are translation keys with correct statuses.
- [ ] The DTO has no internal types.
- [ ] The client invalidates the affected query keys.
- [ ] Both message catalogues have the new keys.
