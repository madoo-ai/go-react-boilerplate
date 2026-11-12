# Internationalization

## Locales

- `id` — Bahasa Indonesia (default)
- `en` — English

The two catalogues live in `web/src/messages/en.json` and
`web/src/messages/id.json`.

## Source of truth

The `locale` cookie is the only source of truth. There is no `/[locale]` URL
segment. When the cookie is missing, the provider uses the value that the Go
server injects into the HTML shell from `APP_LOCALE`. The frontend default is
`id`.

The provider reads the cookie during render. On the server it uses the
`initialLocale` prop. `setLocale` writes the cookie and updates
`document.documentElement.lang`.

## Using messages

```tsx
const t = useTranslations("auth.login")
t("title") // auth.login.title

const tAll = useTranslations()
tAll("errors.forbidden")
```

`useTranslations` returns a function. Interpolation uses braces:

```json
{ "page": "Page {page} of {pages}" }
```

```tsx
t("page", { page: 2, pages: 5 })
```

If a key is missing in the active locale, the provider tries English. If the
key is still missing, it returns the key. This makes a missing string visible
in the UI and in tests.

## API errors

The API sends translation keys, not sentences:

```json
{ "error": "errors.forbidden" }
```

Form validation uses `validation.*` keys. `translateFieldError` maps a Zod
message such as `validation.email` to the `validation` namespace.

## Parity check

```bash
cd web && npm run check:i18n
```

The command flattens both catalogues and fails when one key is missing from
the other. The Makefile runs the same check from `make lint`.

## Prerendering and locale

The prerender script always renders the Indonesian catalogue. It sets
`<html lang="id">`. The client reads the cookie or the injected `APP_LOCALE`
value and re-renders in the active locale.

## Adding a locale

1. Add the locale to `locales` in `web/src/lib/i18n.tsx`.
2. Add a catalogue import.
3. Add the label to `localeLabels`.
4. Add the locale to the server validation in `internal/xmodule/config.go`.
5. Add the prerender lang if the default changes.
