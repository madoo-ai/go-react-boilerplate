import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react"
import en from "@/messages/en.json"
import id from "@/messages/id.json"

export const locales = ["id", "en"] as const
export type Locale = (typeof locales)[number]

export const defaultLocale: Locale = "id"
export const localeCookieName = "locale"

export const localeLabels: Record<Locale, string> = {
  id: "Bahasa Indonesia",
  en: "English",
}

const catalogues: Record<Locale, unknown> = { en, id }

export function isLocale(value: string | null | undefined): value is Locale {
  return value === "en" || value === "id"
}

function readCookie(name: string): string | null {
  if (typeof document === "undefined") return null
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]+)`))
  return match ? decodeURIComponent(match[1]) : null
}

function lookup(locale: Locale, key: string): string | undefined {
  let current: unknown = catalogues[locale]
  for (const part of key.split(".")) {
    if (!current || typeof current !== "object") return undefined
    current = (current as Record<string, unknown>)[part]
  }
  return typeof current === "string" ? current : undefined
}

function interpolate(template: string, values?: Record<string, string | number>) {
  if (!values) return template
  return template.replace(/\{(\w+)\}/g, (match, name: string) =>
    name in values ? String(values[name]) : match,
  )
}

type I18nValue = {
  locale: Locale
  setLocale: (locale: Locale) => void
  t: (key: string, values?: Record<string, string | number>) => string
}

const I18nContext = createContext<I18nValue | null>(null)

export function I18nProvider({
  children,
  initialLocale,
}: {
  children: React.ReactNode
  initialLocale?: Locale
}) {
  const [locale, setLocaleState] = useState<Locale>(() => {
    const fromCookie = readCookie(localeCookieName)
    if (isLocale(fromCookie)) return fromCookie
    if (initialLocale) return initialLocale
    if (typeof window !== "undefined" && isLocale(window.__APP_LOCALE__)) {
      return window.__APP_LOCALE__
    }
    return defaultLocale
  })

  const setLocale = useCallback((next: Locale) => {
    setLocaleState(next)
    if (typeof document !== "undefined") {
      document.cookie = `${localeCookieName}=${next};path=/;max-age=31536000;samesite=lax`
      document.documentElement.lang = next
    }
  }, [])

  const t = useCallback(
    (key: string, values?: Record<string, string | number>) => {
      const message = lookup(locale, key) ?? lookup("en", key) ?? key
      return interpolate(message, values)
    },
    [locale],
  )

  useEffect(() => {
    if (typeof document !== "undefined") document.documentElement.lang = locale
  }, [locale])

  const value = useMemo(() => ({ locale, setLocale, t }), [locale, setLocale, t])
  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>
}

export function useI18n() {
  const context = useContext(I18nContext)
  if (!context) throw new Error("useI18n must be used inside I18nProvider")
  return context
}

export function useTranslations(namespace?: string) {
  const { t } = useI18n()
  return useCallback(
    (key: string, values?: Record<string, string | number>) =>
      namespace ? t(`${namespace}.${key}`, values) : t(key, values),
    [namespace, t],
  )
}

export function translateFieldError(
  message: string | undefined,
  t: (key: string) => string,
): string | undefined {
  if (!message) return undefined
  if (message.startsWith("validation.")) {
    return t(message.replace("validation.", ""))
  }
  return message
}
