import {
  Beaker01Icon,
  CursorClick01Icon,
  Database01Icon,
  ShieldTickIcon as ShieldIcon,
} from "@x1nx3r/altics-react-icons"
import {
  DockerIcon,
  GoIcon,
  PostgreSQLIcon,
  ReactIcon,
  TailwindCSSIcon,
  TypeScriptIcon,
} from "@x1nx3r/altics-react-icons/tech"
import { useTranslations } from "@/lib/i18n"
import type { AppIcon } from "@/types/icon"

export function StackList() {
  const t = useTranslations("landing.stack")
  const itemsT = useTranslations("landing.stack.items")

  const items: Array<{ icon: AppIcon; name: string; role: string }> = [
    { icon: GoIcon, name: itemsT("next.name"), role: itemsT("next.role") },
    { icon: ReactIcon, name: itemsT("react.name"), role: itemsT("react.role") },
    { icon: TypeScriptIcon, name: itemsT("typescript.name"), role: itemsT("typescript.role") },
    { icon: PostgreSQLIcon, name: itemsT("prisma.name"), role: itemsT("prisma.role") },
    { icon: Database01Icon, name: itemsT("postgres.name"), role: itemsT("postgres.role") },
    { icon: TailwindCSSIcon, name: itemsT("tailwind.name"), role: itemsT("tailwind.role") },
    { icon: ShieldIcon, name: itemsT("auth.name"), role: itemsT("auth.role") },
    { icon: DockerIcon, name: itemsT("docker.name"), role: itemsT("docker.role") },
    { icon: Beaker01Icon, name: itemsT("vitest.name"), role: itemsT("vitest.role") },
    { icon: CursorClick01Icon, name: itemsT("playwright.name"), role: itemsT("playwright.role") },
  ]

  return (
    <div className="flex flex-col gap-8">
      <div className="grid gap-x-10 md:grid-cols-2">
        {items.map((item) => {
          const Icon = item.icon
          return (
            <div key={item.name} className="group flex items-start gap-3 border-t border-border/60 py-4">
              <span className="grid size-9 shrink-0 place-items-center rounded-lg bg-muted/60 text-muted-foreground ring-1 ring-border/60 transition-colors group-hover:text-foreground">
                <Icon className="size-4" />
              </span>
              <div className="flex min-w-0 flex-col gap-0.5">
                <span className="truncate text-sm font-medium">{item.name}</span>
                <span className="truncate text-xs text-muted-foreground">{item.role}</span>
              </div>
            </div>
          )
        })}
      </div>
      <p className="text-xs text-muted-foreground">{t("note")}</p>
    </div>
  )
}
