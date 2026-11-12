import { useTheme } from "@/lib/theme"
import { useTranslations } from "@/lib/i18n"
import { ThemePill } from "@/components/theme/theme-pill"
import { COLOR_THEME_OPTIONS } from "@/lib/theme-options"

export function LandingThemeRail() {
  const t = useTranslations("landing.themes")
  const { colorTheme, setColorTheme } = useTheme()
  const active = COLOR_THEME_OPTIONS.find((option) => option.id === colorTheme)

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between gap-4 text-xs text-muted-foreground">
        <span className="flex items-center gap-2">
          <span className="size-1.5 rounded-full bg-chart-3" />
          {t("live")}
        </span>
        <span className="hidden sm:inline">
          {t("current")}: <span className="font-medium text-foreground">{active?.name}</span>
        </span>
      </div>
      <div className="flex [scrollbar-width:thin] gap-2 overflow-x-auto pb-1">
        {COLOR_THEME_OPTIONS.map((option) => (
          <ThemePill
            key={option.id}
            option={option}
            selected={colorTheme === option.id}
            onSelect={setColorTheme}
          />
        ))}
      </div>
    </div>
  )
}
