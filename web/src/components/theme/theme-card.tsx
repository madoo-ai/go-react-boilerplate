import { COLOR_THEME_OPTIONS, type ColorThemeId } from "@/lib/theme-options"
import { CheckCircleIcon } from "@x1nx3r/altics-react-icons"
import { cn } from "@/lib/utils"

type ThemeOption = (typeof COLOR_THEME_OPTIONS)[number]

export function ThemeCard({
  option,
  selected,
  onSelect,
  className,
}: {
  option: ThemeOption
  selected: boolean
  onSelect: (theme: ColorThemeId) => void
  className?: string
}) {
  return (
    <button
      type="button"
      aria-pressed={selected}
      onClick={() => onSelect(option.id)}
      className={cn(
        "group flex items-center gap-3 rounded-lg border p-2.5 text-left transition-colors",
        selected
          ? "border-primary/40 bg-primary/5 ring-1 ring-primary/20"
          : "hover:border-primary/25 hover:bg-muted/40",
        className,
      )}
    >
      <span className="flex shrink-0 overflow-hidden rounded-md border transition-transform duration-200 group-hover:scale-[1.03]">
        {[option.light.neutral, option.light.primary, option.dark.neutral, option.dark.primary].map(
          (color, index) => (
            <span key={`${option.id}-${index}`} className="size-6" style={{ backgroundColor: color }} />
          ),
        )}
      </span>
      <span className="flex min-w-0 flex-1 flex-col">
        <span className="truncate text-sm font-medium">{option.name}</span>
        <span className="truncate text-xs text-muted-foreground">{option.id}</span>
      </span>
      {selected ? <CheckCircleIcon className="size-4 shrink-0 text-primary" /> : null}
    </button>
  )
}
