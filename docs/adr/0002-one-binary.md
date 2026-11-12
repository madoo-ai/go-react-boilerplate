# 0002 — Go API plus React SPA, served as one binary

## Status

Accepted.

## Context

The source boilerplate is a Next.js full-stack app. A Go version needs a
clear split between the API and the UI. A second runtime would add a reverse
proxy, a second deploy unit and CORS.

## Decision

The Go server owns `/api/*` and the static React app on one origin. Vite
builds the app. Go embeds the build. Production runs one binary.

Development keeps Vite for HMR on port 5173. Vite proxies `/api` to the Go
server on port 8080. Node is not a production dependency.

## Consequences

- Cookies and CSRF checks stay simple because there is one origin.
- The final image is small and has one process.
- Developer experience keeps fast refresh.
- Server-side rendering for React is not available. Marketing and auth pages
  are prerendered at build time instead.
