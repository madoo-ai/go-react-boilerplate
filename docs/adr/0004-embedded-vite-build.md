# 0004 — Embedded Vite build with build-time prerendering

## Status

Accepted.

## Context

A SPA sends the same HTML for every route. The public landing page and auth
pages need useful titles and descriptions for search engines and link
previews.

## Decision

The build has three steps:

1. Vite builds the client into `internal/handler/webui/dist`.
2. Vite builds `src/entry-server.tsx` into `web/dist-ssr`.
3. A Node script renders the public routes and writes real HTML files.

The Go handler serves a prerendered file when it exists. It falls back to the
SPA shell for app routes. The client uses `createRoot`, thus there is no
hydration contract to maintain.

## Consequences

- Public pages have real content and meta tags in the first response.
- The authenticated app stays a SPA.
- The build needs Node. The runtime does not.
- Prerendered routes must render without browser APIs. The shared providers
  guard `document` and `window`.
