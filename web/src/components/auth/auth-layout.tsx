import { AppLogo } from "@/components/app-shell/app-logo"
import { LocaleSwitcher } from "@/components/app-shell/locale-switcher"
import { ThemeToggle } from "@/components/app-shell/theme-toggle"
import { useTranslations } from "@/lib/i18n"
import { COLOR_THEME_OPTIONS } from "@/lib/theme-options"

const swatchIds = ["default", "dracula", "nord", "catppuccin", "tokyonight", "vercel"] as const

export function AuthLayout({ children }: { children: React.ReactNode }) {
  const t = useTranslations("auth.brand")
  const swatches = swatchIds
    .map((id) => COLOR_THEME_OPTIONS.find((option) => option.id === id))
    .filter((option): option is (typeof COLOR_THEME_OPTIONS)[number] => !!option)

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-8 p-6 md:p-10">
        <div className="flex items-center justify-between gap-4">
          <AppLogo />
          <div className="flex items-center gap-1">
            <LocaleSwitcher />
            <ThemeToggle />
          </div>
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-sm motion-safe:animate-rise">{children}</div>
        </div>
        <p className="text-center text-xs text-muted-foreground md:text-left">{t("stack")}</p>
      </div>
      <div className="relative hidden overflow-hidden border-l border-border/60 bg-mesh lg:block">
        <div className="absolute inset-0 bg-grid opacity-[0.18] [mask-image:linear-gradient(to_bottom,black,transparent_85%)]" />
        <div className="absolute inset-0 bg-gradient-to-br from-transparent via-transparent to-primary/10" />
        <div className="relative flex h-full flex-col justify-between p-10">
          <div className="flex flex-wrap gap-2" aria-hidden>
            {swatches.map((swatch) => (
              <span key={swatch.id} className="flex overflow-hidden rounded-full border border-border/60 bg-background/50 shadow-sm">
                <span className="size-4" style={{ backgroundColor: swatch.light.primary }} />
                <span className="size-4" style={{ backgroundColor: swatch.dark.primary }} />
              </span>
            ))}
          </div>
          <div className="flex flex-col gap-4">
            <blockquote className="max-w-md font-display text-4xl leading-[1.02] tracking-[-0.02em] text-foreground">
              {t("tagline")}
            </blockquote>
            <p className="text-sm text-muted-foreground">{t("stack")}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
