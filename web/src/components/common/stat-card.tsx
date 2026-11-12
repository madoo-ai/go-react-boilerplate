import { Card, CardAction, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { AppIcon } from "@/types/icon"

export function StatCard({
  label,
  value,
  hint,
  icon: Icon,
}: {
  label: string
  value: string | number
  hint?: string
  icon?: AppIcon
}) {
  return (
    <Card className="group/card relative h-full hover:ring-primary/20">
      <CardHeader>
        <CardTitle className="text-sm font-medium text-muted-foreground">{label}</CardTitle>
        {Icon ? (
          <CardAction>
            <span className="grid size-9 place-items-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/15 transition-colors group-hover/card:bg-primary/15">
              <Icon className="size-4" />
            </span>
          </CardAction>
        ) : null}
      </CardHeader>
      <CardContent className="flex flex-col gap-1">
        <div className="text-3xl font-light tracking-tight tabular-nums">{value}</div>
        {hint ? <p className="text-xs text-muted-foreground">{hint}</p> : null}
      </CardContent>
      <span className="pointer-events-none absolute inset-x-0 bottom-0 h-px origin-left scale-x-0 bg-gradient-to-r from-primary via-chart-2 to-transparent transition-transform duration-300 group-hover/card:scale-x-100" />
    </Card>
  )
}
