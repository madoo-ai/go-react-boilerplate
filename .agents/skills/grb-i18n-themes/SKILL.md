---
name: grb-i18n-themes
description: Use when adding or changing a translation, locale behavior, a color palette, the theme picker, the light/dark/system mode or the no-flash boot script. Covers both message catalogues, cookie-backed preferences and generated theme assets.
---

# i18n and themes

Locale and appearance are cookie-backed. They are not database rows.

## i18n

- Locales: `id` (default) and `en`.
- Messages: `web/src/messages/id.json` and `web/src/messages/en.json`.
- Provider: `web/src/lib/i18n.tsx`.
- Cookie: `locale`.
- The Go server injects `APP_LOCALE` into `index.html`.
- The client uses the cookie first, then `window.__APP_LOCALE__`.

Add a message to both catalogues, then run:

```bash
cd web && npm run check:i18n
```

Use `useTranslations("namespace")` in components. Zod messages are keys such
as `validation.email`; translate them with `translateFieldError`. API errors
are also keys; render them with the global `t`.

## Themes

- Mode cookie: `theme` with `light`, `dark` or `system`.
- Palette cookie: `color-theme`.
- Palettes: `web/src/styles/themes.css`.
- Palette metadata: `web/src/lib/theme-options.ts`.
- Theme provider: `web/src/lib/theme.tsx`.
- The boot script in `web/index.html` applies both before first paint.

Do not hand-edit `theme-options.ts` or `themes.css`. The source is
`web/src/lib/opencode-themes.json`; regenerate the assets when the palette set
changes.

## Appearance rules

- Use semantic tokens (`bg-background`, `text-muted-foreground`, `border`).
- Do not use raw palette classes for app surfaces.
- Keep the reduced-motion path usable.
- Use `font-display` for marketing and auth editorial headings.
- Use Inter for product screens.

## Adding a locale

1. Add the code and label to `web/src/lib/i18n.tsx`.
2. Import the catalogue.
3. Update the server validation in `internal/xmodule/config.go`.
4. Update the prerender script and `index.html` only if the default changes.
5. Run `npm run check:i18n` and the Playwright locale test.
