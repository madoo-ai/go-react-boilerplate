# 0005 — Cookie-backed locale and theme

## Status

Accepted.

## Context

The Next.js app stores locale and appearance in cookies. This avoids a
`/[locale]` prefix and keeps user preferences across tabs. The Go port must
preserve the behavior.

## Decision

- `locale` cookie: `id` or `en`.
- `theme` cookie: `light`, `dark` or `system`.
- `color-theme` cookie: palette ID.
- Other values: `sidebar_state` for the sidebar.

The client writes these cookies. The server does not store them in Postgres.
An inline script in `index.html` applies the mode and palette before first
paint.

## Consequences

- No database writes for appearance changes.
- Preferences apply before the app renders, so there is no flash.
- The server cannot render user-specific appearance. Prerendered pages use
  the default and the client updates after load.
- Clearing cookies resets preferences to the defaults.
