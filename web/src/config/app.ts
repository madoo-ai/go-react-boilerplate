export const APP_NAME =
  (typeof window !== "undefined" && window.__APP_NAME__) ||
  import.meta.env.VITE_APP_NAME ||
  "Go React Boilerplate"

export const APP_DESCRIPTION =
  "A Go + React boilerplate with auth, RBAC, i18n and a themeable app shell."
