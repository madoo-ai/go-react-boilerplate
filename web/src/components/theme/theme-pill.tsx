import { COLOR_THEME_OPTIONS, type ColorThemeId } from "@/lib/theme-options"
import { CheckCircleIcon } from "@x1nx3r/altics-react-icons"
import { cn } from "@/lib/utils"

type ThemeOption = (typeof COLOR_THEME_OPTIONS)[number]

export function ThemePill({
  option,
  selected,
  onSelect,
}: {
  option: ThemeOption
  selected: boolean
  onSelect: (theme: ColorThemeId) => void
}) {
  return (
    <button
      type="button"
      aria-label={option.name}
      aria-pressed={selected}
      onClick={() => onSelect(option.id)}
      className={cn(
        "group flex shrink-0 items-center gap-2 rounded-full border px-2.5 py-1.5 text-xs transition-colors",
        selected
          ? "border-primary/40 bg-primary/10 text-foreground"
          : "bg-background/60 text-muted-foreground hover:border-primary/25 hover:text-foreground",
      )}
    >
      <span className="flex overflow-hidden rounded-full border border-border/60">
        <span className="size-3" style={{ backgroundColor: option.light.primary }} />
        <span className="size-3" style={{ backgroundColor: option.dark.primary }} />
      </span>
      <span className="max-w-28 truncate">{option.name}</span>
      {selected ? <CheckCircleIcon className="size-3.5 text-primary" /> : null}
    </button>
  )
}
