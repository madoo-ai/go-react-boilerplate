import { renderToString } from "react-dom/server"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { MemoryRouter } from "react-router-dom"
import { App } from "./app"
import { AuthProvider } from "@/lib/auth"
import { I18nProvider } from "@/lib/i18n"
import { ThemeProvider } from "@/lib/theme"
import { TooltipProvider } from "@/components/ui/tooltip"
import id from "@/messages/id.json"
import { APP_NAME } from "@/config/app"

export function render(path: string) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return renderToString(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <I18nProvider initialLocale="id">
          <ThemeProvider initialTheme="system" initialColorTheme="default">
            <AuthProvider>
              <TooltipProvider>
                <App />
              </TooltipProvider>
            </AuthProvider>
          </ThemeProvider>
        </I18nProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

function lookup(path: string): string {
  let current: unknown = id
  for (const part of path.split(".")) {
    if (!current || typeof current !== "object") return ""
    current = (current as Record<string, unknown>)[part]
  }
  return typeof current === "string" ? current : ""
}

const metaByRoute: Record<string, { titleKey: string; descriptionKey?: string }> = {
  "/": { titleKey: "landing.meta.title", descriptionKey: "landing.meta.description" },
  "/login": { titleKey: "auth.login.title" },
  "/register": { titleKey: "auth.register.title" },
  "/forgot-password": { titleKey: "auth.forgot.title" },
  "/reset-password": { titleKey: "auth.reset.title" },
  "/verify-email": { titleKey: "auth.verify.title" },
  "/two-factor": { titleKey: "auth.twoFactor.title" },
  "/403": { titleKey: "errorPages.forbidden.title" },
}

export function metaFor(path: string) {
  const entry = metaByRoute[path] ?? metaByRoute["/"]
  const title = lookup(entry.titleKey) || APP_NAME
  return {
    title: `${title} · __APP_NAME_TOKEN__`,
    description: entry.descriptionKey ? lookup(entry.descriptionKey) : lookup("landing.meta.description"),
  }
}
