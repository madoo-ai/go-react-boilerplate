import { Link } from "react-router-dom"
import {
  ChevronDownIcon,
  Database01Icon,
  Key01Icon,
  PaletteIcon,
  Rocket01Icon,
  ShieldTickIcon,
  Star01Icon,
} from "@x1nx3r/altics-react-icons"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { LandingThemeGrid } from "@/components/marketing/landing-theme-grid"
import { LandingThemeRail } from "@/components/marketing/landing-theme-rail"
import { ProductPreview } from "@/components/marketing/product-preview"
import { StackList } from "@/components/marketing/stack-list"
import { Button } from "@/components/ui/button"

const demoAccounts = [
  { role: "Super Admin", email: "superadmin@example.com" },
  { role: "HR Staff", email: "hr@example.com" },
  { role: "User", email: "user@example.com" },
]

export function LandingPage() {
  const { user } = useAuth()
  const tHero = useTranslations("landing.hero")
  const tFeatures = useTranslations("landing.features")
  const tFeatureItems = useTranslations("landing.features.items")
  const tThemes = useTranslations("landing.themes")
  const tStack = useTranslations("landing.stack")
  const tDemo = useTranslations("landing.demo")
  const tDashboard = useTranslations("dashboard")
  const tNav = useTranslations("nav")
  const tMeta = useTranslations("landing.meta")

  usePageMeta(tMeta("title"), tMeta("description"))

  const features = [
    { icon: ShieldTickIcon, title: tFeatureItems("auth.title"), body: tFeatureItems("auth.body") },
    { icon: Key01Icon, title: tFeatureItems("rbac.title"), body: tFeatureItems("rbac.body") },
    { icon: PaletteIcon, title: tFeatureItems("shell.title"), body: tFeatureItems("shell.body") },
    { icon: Database01Icon, title: tFeatureItems("data.title"), body: tFeatureItems("data.body") },
    { icon: Star01Icon, title: tFeatureItems("design.title"), body: tFeatureItems("design.body") },
    { icon: Rocket01Icon, title: tFeatureItems("ops.title"), body: tFeatureItems("ops.body") },
  ]

  const previewLabels = {
    nav: {
      dashboard: tNav("dashboard"),
      users: tNav("users"),
      roles: tNav("roles"),
      settings: tNav("settings"),
    },
    stats: {
      totalUsers: tDashboard("stats.totalUsers"),
      activeUsers: tDashboard("stats.activeUsers"),
      pendingUsers: tDashboard("stats.pendingUsers"),
      roles: tDashboard("stats.roles"),
    },
    recentUsers: tDashboard("recentUsers"),
    roleDistribution: tDashboard("charts.roleDistribution"),
    description: tDashboard("description"),
    themesBadge: tHero("themesBadge"),
  }

  return (
    <>
      <section className="relative isolate overflow-hidden border-b border-border/60">
        <div className="absolute inset-0 -z-10 bg-mesh" />
        <div className="absolute inset-0 -z-10 bg-grid opacity-[0.14] [mask-image:linear-gradient(to_bottom,black,transparent_88%)]" />
        <div className="absolute top-8 -left-40 -z-10 size-[34rem] rounded-full bg-[radial-gradient(circle_at_center,color-mix(in_oklab,var(--primary)_32%,transparent),transparent_70%)] opacity-70 motion-safe:animate-mesh" />
        <div className="mx-auto max-w-7xl px-6 pt-20 pb-16 lg:pt-28">
          <div className="grid items-center gap-14 lg:grid-cols-[1.05fr_0.95fr]">
            <div className="max-w-2xl">
              <h1 className="font-display text-5xl leading-[0.98] tracking-[-0.03em] text-balance motion-safe:animate-rise sm:text-6xl lg:text-7xl">
                {tHero("title")}
              </h1>
              <p className="mt-6 max-w-xl text-base leading-relaxed text-muted-foreground [animation-delay:80ms] motion-safe:animate-rise">
                {tHero("description")}
              </p>
              <div className="mt-8 flex flex-wrap items-center gap-3 [animation-delay:160ms] motion-safe:animate-rise">
                <Button asChild size="lg">
                  <Link to={user ? "/dashboard" : "/login"}>
                    {user ? tHero("primaryCtaAuthed") : tHero("primaryCta")}
                  </Link>
                </Button>
                <Button asChild variant="outline" size="lg">
                  <a href="#features">
                    {tHero("secondaryCta")}
                    <ChevronDownIcon data-icon="inline-end" />
                  </a>
                </Button>
              </div>
              <p className="mt-5 max-w-xl text-xs leading-relaxed text-muted-foreground [animation-delay:240ms] motion-safe:animate-rise">
                {tHero("note")}
              </p>
            </div>
            <div className="[animation-delay:120ms] motion-safe:animate-rise">
              <ProductPreview labels={previewLabels} />
            </div>
          </div>
          <div className="mt-14 rounded-2xl border border-border/70 bg-background/70 p-4 shadow-sm backdrop-blur-md [animation-delay:280ms] motion-safe:animate-rise">
            <LandingThemeRail />
          </div>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-7xl scroll-mt-24 px-6 py-24">
        <div className="grid gap-12 lg:grid-cols-[0.8fr_1.2fr] lg:gap-20">
          <div className="lg:sticky lg:top-28 lg:self-start">
            <h2 className="font-display text-4xl leading-[1.02] tracking-[-0.02em] text-balance sm:text-5xl">
              {tFeatures("title")}
            </h2>
            <p className="mt-4 max-w-md text-sm leading-relaxed text-muted-foreground">
              {tFeatures("description")}
            </p>
          </div>
          <ul className="grid gap-x-10 sm:grid-cols-2">
            {features.map((feature) => {
              const Icon = feature.icon
              return (
                <li key={feature.title} className="group border-t border-border/60 py-7">
                  <span className="mb-4 grid size-9 place-items-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/15 transition-colors group-hover:bg-primary/15">
                    <Icon className="size-4" />
                  </span>
                  <h3 className="text-base font-medium">{feature.title}</h3>
                  <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{feature.body}</p>
                </li>
              )
            })}
          </ul>
        </div>
      </section>

      <section id="themes" className="scroll-mt-24 border-y border-border/60 bg-muted/30">
        <div className="mx-auto max-w-7xl px-6 py-24">
          <div className="grid gap-10 lg:grid-cols-[0.8fr_1.2fr] lg:gap-20">
            <div className="lg:sticky lg:top-28 lg:self-start">
              <h2 className="font-display text-4xl leading-[1.02] tracking-[-0.02em] text-balance sm:text-5xl">
                {tThemes("title")}
              </h2>
              <p className="mt-4 max-w-md text-sm leading-relaxed text-muted-foreground">
                {tThemes("description")}
              </p>
              <Button asChild variant="outline" className="mt-6">
                <Link to="/settings/appearance">{tThemes("all")}</Link>
              </Button>
            </div>
            <div className="flex flex-col gap-3">
              <span className="text-xs font-medium text-muted-foreground">{tThemes("featured")}</span>
              <LandingThemeGrid />
            </div>
          </div>
        </div>
      </section>

      <section id="stack" className="mx-auto max-w-7xl scroll-mt-24 px-6 py-24">
        <div className="grid gap-12 lg:grid-cols-[0.8fr_1.2fr] lg:gap-20">
          <div className="lg:sticky lg:top-28 lg:self-start">
            <h2 className="font-display text-4xl leading-[1.02] tracking-[-0.02em] text-balance sm:text-5xl">
              {tStack("title")}
            </h2>
            <p className="mt-4 max-w-md text-sm leading-relaxed text-muted-foreground">
              {tStack("description")}
            </p>
          </div>
          <StackList />
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-6 pb-24">
        <div className="relative overflow-hidden rounded-3xl bg-foreground px-6 py-12 text-background sm:px-10 lg:px-14 lg:py-16">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_right,color-mix(in_oklab,var(--primary)_35%,transparent),transparent_55%)]" />
          <div className="relative grid gap-12 lg:grid-cols-[1.05fr_0.95fr]">
            <div className="max-w-lg">
              <h2 className="font-display text-4xl leading-[1.02] tracking-[-0.02em] text-balance sm:text-5xl">
                {tDemo("title")}
              </h2>
              <p className="mt-4 text-sm leading-relaxed text-background/70">{tDemo("description")}</p>
              <div className="mt-8 flex flex-wrap gap-3">
                <Button asChild size="lg">
                  <Link to="/login">{tDemo("signIn")}</Link>
                </Button>
                <Button
                  asChild
                  size="lg"
                  variant="outline"
                  className="border-background/20 bg-transparent text-background hover:bg-background/10 hover:text-background"
                >
                  <Link to="/register">{tDemo("create")}</Link>
                </Button>
              </div>
            </div>
            <div className="flex flex-col gap-6">
              <div className="rounded-2xl border border-background/10 bg-background/5 p-4">
                <p className="text-xs font-medium text-background/70">{tDemo("accounts")}</p>
                <div className="mt-3 flex flex-col">
                  {demoAccounts.map((account) => (
                    <div key={account.email} className="flex items-center justify-between gap-4 border-t border-background/10 py-2.5 first:border-t-0">
                      <span className="text-xs font-medium">{account.role}</span>
                      <span className="truncate font-mono text-[11px] text-background/60">{account.email}</span>
                    </div>
                  ))}
                </div>
                <p className="mt-3 text-xs text-background/60">
                  {tDemo("password")}: <code className="rounded bg-background/10 px-1.5 py-0.5 font-mono text-[11px] text-background">password</code>
                </p>
              </div>
              <div className="rounded-2xl border border-background/10 bg-background/5 p-4">
                <p className="text-xs font-medium text-background/70">{tDemo("command")}</p>
                <pre className="mt-3 overflow-x-auto font-mono text-[11px] leading-relaxed text-background/80">
                  <span className="text-background/50">$</span> make setup{"\n"}
                  <span className="text-background/50">$</span> make dev
                </pre>
              </div>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}
