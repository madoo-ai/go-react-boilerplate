import {
  ChevronRightIcon,
  Key01Icon,
  LayoutGrid01Icon,
  Settings01Icon,
  Users01Icon,
} from "@x1nx3r/altics-react-icons"

export type ProductPreviewLabels = {
  nav: { dashboard: string; users: string; roles: string; settings: string }
  stats: { totalUsers: string; activeUsers: string; pendingUsers: string; roles: string }
  recentUsers: string
  roleDistribution: string
  description: string
  themesBadge: string
}

const bars = [38, 54, 44, 68, 50, 74, 62, 82, 58, 70]
const recentPeople = [
  { initials: "SA", name: "Super Admin", status: "bg-success" },
  { initials: "HR", name: "HR Staff", status: "bg-success" },
  { initials: "RU", name: "Regular User", status: "bg-warning" },
]

export function ProductPreview({ labels }: { labels: ProductPreviewLabels }) {
  return (
    <div className="pointer-events-none relative mx-auto w-full max-w-[580px] select-none" aria-hidden>
      <div className="absolute -inset-6 rounded-[2.5rem] bg-gradient-to-br from-primary/20 via-chart-2/10 to-chart-3/20 blur-2xl" />
      <div className="relative overflow-hidden rounded-2xl border border-border/80 bg-card shadow-2xl ring-1 ring-foreground/5">
        <div className="flex items-center gap-2 border-b border-border/70 bg-muted/50 px-4 py-2.5">
          <span className="flex gap-1.5">
            <span className="size-2 rounded-full bg-destructive/60" />
            <span className="size-2 rounded-full bg-warning/60" />
            <span className="size-2 rounded-full bg-success/60" />
          </span>
          <span className="mx-auto h-5 w-40 rounded-full border border-border/70 bg-background/70 text-center text-[10px] leading-5 text-muted-foreground sm:w-52">
            app.local/dashboard
          </span>
          <span className="size-4" />
        </div>
        <div className="grid grid-cols-[100px_1fr] sm:grid-cols-[132px_1fr]">
          <aside className="border-r border-border/70 bg-sidebar p-3">
            <div className="mb-4 flex items-center gap-2">
              <span className="grid size-6 place-items-center rounded-md bg-primary text-[9px] font-semibold text-primary-foreground">G</span>
              <span className="truncate text-[10px] font-semibold">Boilerplate</span>
            </div>
            <div className="flex flex-col gap-1">
              {[
                { icon: LayoutGrid01Icon, label: labels.nav.dashboard, active: true },
                { icon: Users01Icon, label: labels.nav.users, active: false },
                { icon: Key01Icon, label: labels.nav.roles, active: false },
                { icon: Settings01Icon, label: labels.nav.settings, active: false },
              ].map((item) => {
                const Icon = item.icon
                return (
                  <div
                    key={item.label}
                    className={
                      item.active
                        ? "flex items-center gap-2 rounded-md bg-sidebar-accent px-2 py-1.5 text-[10px] font-medium text-sidebar-accent-foreground"
                        : "flex items-center gap-2 rounded-md px-2 py-1.5 text-[10px] text-muted-foreground"
                    }
                  >
                    <Icon className="size-3 shrink-0" />
                    <span className="truncate">{item.label}</span>
                  </div>
                )
              })}
            </div>
          </aside>
          <div className="min-w-0 bg-background p-3 sm:p-4">
            <div className="flex items-center justify-between text-[10px] text-muted-foreground">
              <span>{labels.nav.dashboard}</span>
              <span className="flex items-center gap-1">
                <span className="size-1.5 rounded-full bg-primary" />
                <span className="size-1.5 rounded-full bg-chart-2" />
              </span>
            </div>
            <div className="mt-3 flex items-start justify-between gap-2">
              <div className="min-w-0">
                <p className="truncate text-lg font-light tracking-tight">{labels.nav.dashboard}</p>
                <p className="mt-0.5 truncate text-[9px] text-muted-foreground">{labels.description}</p>
              </div>
              <span className="grid size-5 shrink-0 place-items-center rounded-full bg-primary text-primary-foreground">
                <ChevronRightIcon className="size-3" />
              </span>
            </div>
            <div className="mt-3 grid grid-cols-2 gap-1.5 sm:grid-cols-4">
              {[
                { label: labels.stats.totalUsers, value: "14" },
                { label: labels.stats.activeUsers, value: "10" },
                { label: labels.stats.pendingUsers, value: "3" },
                { label: labels.stats.roles, value: "3" },
              ].map((stat) => (
                <div key={stat.label} className="min-w-0 rounded-lg border border-border/70 bg-card p-1.5">
                  <p className="truncate text-[8px] leading-tight text-muted-foreground">{stat.label}</p>
                  <p className="mt-0.5 text-sm font-light tabular-nums">{stat.value}</p>
                </div>
              ))}
            </div>
            <div className="mt-2 grid grid-cols-1 gap-1.5 sm:grid-cols-[1.25fr_0.75fr]">
              <div className="rounded-lg border border-border/70 bg-card p-2">
                <p className="truncate text-[9px] text-muted-foreground">{labels.roleDistribution}</p>
                <div className="mt-2 flex h-16 items-end gap-1">
                  {bars.map((height, index) => (
                    <span
                      key={index}
                      className={index > 7 ? "w-2 rounded-t-sm bg-chart-2/80" : "w-2 rounded-t-sm bg-primary/80"}
                      style={{ height: `${height}%` }}
                    />
                  ))}
                </div>
              </div>
              <div className="rounded-lg border border-border/70 bg-card p-2">
                <p className="truncate text-[9px] text-muted-foreground">{labels.recentUsers}</p>
                <div className="mt-2 flex flex-col gap-1.5">
                  {recentPeople.map((person) => (
                    <div key={person.name} className="flex min-w-0 items-center gap-1.5">
                      <span className="grid size-4 shrink-0 place-items-center rounded-full bg-primary/10 text-[7px] font-medium text-primary">
                        {person.initials}
                      </span>
                      <span className="truncate text-[8px]">{person.name}</span>
                      <span className={`ml-auto size-1.5 shrink-0 rounded-full ${person.status}`} />
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className="absolute -bottom-4 -left-3 flex items-center gap-2 rounded-full border border-border/70 bg-background/90 px-3 py-1.5 text-[10px] font-medium shadow-lg backdrop-blur">
        <span className="size-1.5 rounded-full bg-chart-3" />
        {labels.themesBadge}
      </div>
    </div>
  )
}
