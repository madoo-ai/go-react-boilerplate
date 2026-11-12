import { useTheme } from "@/lib/theme"
import { ThemeCard } from "@/components/theme/theme-card"
import { COLOR_THEME_OPTIONS } from "@/lib/theme-options"

const featuredIds = [
  "default",
  "dracula",
  "nord",
  "catppuccin",
  "tokyonight",
  "vercel",
  "gruvbox",
  "rosepine",
  "opencode",
] as const

export function LandingThemeGrid() {
  const { colorTheme, setColorTheme } = useTheme()
  const featured = featuredIds
    .map((id) => COLOR_THEME_OPTIONS.find((option) => option.id === id))
    .filter((option): option is (typeof COLOR_THEME_OPTIONS)[number] => !!option)

  return (
    <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
      {featured.map((option) => (
        <ThemeCard
          key={option.id}
          option={option}
          selected={colorTheme === option.id}
          onSelect={setColorTheme}
          className="bg-background/70"
        />
      ))}
    </div>
  )
}
