import { describe, expect, it } from "vitest"
import { act, renderHook } from "@testing-library/react"
import type { ReactNode } from "react"
import { I18nProvider, useI18n, useTranslations } from "./i18n"
import { isColorTheme, isTheme } from "./theme"

function wrapper({ children }: { children: ReactNode }) {
  return <I18nProvider initialLocale="en">{children}</I18nProvider>
}

describe("i18n", () => {
  it("translates namespaced keys and interpolates values", () => {
    const { result } = renderHook(
      () => ({
        login: useTranslations("auth.login"),
        common: useTranslations("common"),
      }),
      { wrapper },
    )
    expect(result.current.login("title")).toBe("Sign in")
    expect(result.current.common("page", { page: 2, pages: 5 })).toBe("Page 2 of 5")
  })

  it("falls back to the key when a message is missing", () => {
    const { result } = renderHook(() => useTranslations(), { wrapper })
    expect(result.current("does.not.exist")).toBe("does.not.exist")
  })

  it("switches locale at runtime", () => {
    const { result } = renderHook(() => useI18n(), { wrapper })
    expect(result.current.t("auth.login.title")).toBe("Sign in")
    act(() => result.current.setLocale("id"))
    expect(result.current.t("auth.login.title")).toBe("Masuk")
  })
})

describe("theme guards", () => {
  it("accepts only known modes", () => {
    expect(isTheme("light")).toBe(true)
    expect(isTheme("dark")).toBe(true)
    expect(isTheme("system")).toBe(true)
    expect(isTheme("sepia")).toBe(false)
  })

  it("accepts only known palettes", () => {
    expect(isColorTheme("default")).toBe(true)
    expect(isColorTheme("dracula")).toBe(true)
    expect(isColorTheme("not-a-theme")).toBe(false)
  })
})
