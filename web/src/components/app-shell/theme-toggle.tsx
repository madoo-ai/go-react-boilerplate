import { CheckIcon as Check, Monitor01Icon as Monitor, Moon01Icon as Moon, SunIcon as Sun } from "@x1nx3r/altics-react-icons"
import { useTranslations } from "@/lib/i18n"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useTheme, themeModes, type Theme } from "@/lib/theme"

const themeIcons = { light: Sun, dark: Moon, system: Monitor } satisfies Record<Theme, typeof Sun>

export function ThemeToggle() {
  const t = useTranslations("preferences")
  const { theme, setTheme } = useTheme()
  const Icon = themeIcons[theme]

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" aria-label={t("theme")}>
          <Icon />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="min-w-40">
        {themeModes.map((option) => {
          const OptionIcon = themeIcons[option]
          return (
            <DropdownMenuItem key={option} onSelect={() => setTheme(option)}>
              <OptionIcon />
              <span className="flex-1">{t(option)}</span>
              {theme === option ? <Check className="size-4" /> : null}
            </DropdownMenuItem>
          )
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
