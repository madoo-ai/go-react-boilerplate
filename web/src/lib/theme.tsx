import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react"
import { COLOR_THEME_OPTIONS, type ColorThemeId } from "@/lib/theme-options"

export const themeModes = ["light", "dark", "system"] as const
export type Theme = (typeof themeModes)[number]

export const themeCookieName = "theme"
export const colorThemeCookieName = "color-theme"
export const defaultTheme: Theme = "system"
export const defaultColorTheme: ColorThemeId = "default"

export function isTheme(value: string | null | undefined): value is Theme {
  return value === "light" || value === "dark" || value === "system"
}

export function isColorTheme(value: string | null | undefined): value is ColorThemeId {
  return !!value && COLOR_THEME_OPTIONS.some((option) => option.id === value)
}

function readCookie(name: string): string | null {
  if (typeof document === "undefined") return null
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]+)`))
  return match ? decodeURIComponent(match[1]) : null
}

function resolveTheme(theme: Theme): "light" | "dark" {
  if (theme === "system") {
    if (typeof window === "undefined") return "light"
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  }
  return theme
}

function applyTheme(theme: Theme) {
  if (typeof document === "undefined") return
  const resolved = resolveTheme(theme)
  const root = document.documentElement
  root.classList.toggle("dark", resolved === "dark")
  root.style.colorScheme = resolved
}

function applyColorTheme(theme: ColorThemeId) {
  if (typeof document === "undefined") return
  if (theme === "default") {
    delete document.documentElement.dataset.colorTheme
    return
  }
  document.documentElement.dataset.colorTheme = theme
}

type ThemeValue = {
  theme: Theme
  colorTheme: ColorThemeId
  setTheme: (theme: Theme) => void
  setColorTheme: (theme: ColorThemeId) => void
}

const ThemeContext = createContext<ThemeValue | null>(null)

export function ThemeProvider({
  children,
  initialTheme,
  initialColorTheme,
}: {
  children: React.ReactNode
  initialTheme?: Theme
  initialColorTheme?: ColorThemeId
}) {
  const [theme, setThemeState] = useState<Theme>(() => {
    const fromCookie = readCookie(themeCookieName)
    return initialTheme ?? (isTheme(fromCookie) ? fromCookie : defaultTheme)
  })
  const [colorTheme, setColorThemeState] = useState<ColorThemeId>(() => {
    const fromCookie = readCookie(colorThemeCookieName)
    return initialColorTheme ?? (isColorTheme(fromCookie) ? fromCookie : defaultColorTheme)
  })

  useEffect(() => {
    applyTheme(theme)
    applyColorTheme(colorTheme)
    if (theme !== "system" || typeof window === "undefined") return
    const media = window.matchMedia("(prefers-color-scheme: dark)")
    const onChange = () => applyTheme("system")
    media.addEventListener("change", onChange)
    return () => media.removeEventListener("change", onChange)
  }, [theme, colorTheme])

  const setTheme = useCallback((next: Theme) => {
    setThemeState(next)
    if (typeof document !== "undefined") {
      document.cookie = `${themeCookieName}=${next};path=/;max-age=31536000;samesite=lax`
    }
    applyTheme(next)
  }, [])

  const setColorTheme = useCallback((next: ColorThemeId) => {
    setColorThemeState(next)
    if (typeof document !== "undefined") {
      document.cookie = `${colorThemeCookieName}=${next};path=/;max-age=31536000;samesite=lax`
    }
    applyColorTheme(next)
  }, [])

  const value = useMemo(
    () => ({ theme, colorTheme, setTheme, setColorTheme }),
    [theme, colorTheme, setTheme, setColorTheme],
  )
  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
}

export function useTheme() {
  const context = useContext(ThemeContext)
  if (!context) throw new Error("useTheme must be used inside ThemeProvider")
  return context
}
