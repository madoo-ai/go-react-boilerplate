import { CheckIcon as Check, Translate01Icon as Languages } from "@x1nx3r/altics-react-icons"
import { useI18n, useTranslations, isLocale, localeLabels, locales, type Locale } from "@/lib/i18n"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export function LocaleSwitcher() {
  const t = useTranslations("preferences")
  const { locale, setLocale } = useI18n()
  const activeLocale = isLocale(locale) ? locale : ("id" as const)

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-2" aria-label={t("language")}>
          <Languages />
          <span className="hidden sm:inline">{localeLabels[activeLocale]}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="min-w-40">
        {locales.map((option: Locale) => (
          <DropdownMenuItem key={option} onSelect={() => setLocale(option)}>
            <span className="flex-1">{localeLabels[option]}</span>
            {activeLocale === option ? <Check className="size-4" /> : null}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
