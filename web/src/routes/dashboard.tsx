import { Link } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import {
  CheckCircleIcon as CheckCircle2,
  ChevronRightIcon,
  ClockIcon as Clock3,
  Key01Icon as KeyRound,
  ShieldTickIcon as ShieldCheck,
  Users01Icon as Users,
} from "@x1nx3r/altics-react-icons"
import { api } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { useI18n, useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { initials } from "@/lib/utils"
import { DashboardCharts } from "@/components/dashboard/dashboard-charts"
import { PageHeader } from "@/components/common/page-header"
import { StatCard } from "@/components/common/stat-card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardAction, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

function statusVariant(status: string) {
  if (status === "active") return "default" as const
  if (status === "suspended") return "destructive" as const
  return "secondary" as const
}

export function DashboardPage() {
  const { user } = useAuth()
  const t = useTranslations("dashboard")
  const tCommon = useTranslations("common")
  const tUserColumns = useTranslations("users.columns")
  const { locale } = useI18n()
  usePageMeta(t("title"))

  const canViewAdmin = user?.permissions.includes("view_admin") ?? false
  const query = useQuery({ queryKey: ["dashboard"], queryFn: api.dashboard, enabled: canViewAdmin })
  const numberFormat = new Intl.NumberFormat(locale === "id" ? "id-ID" : "en-US")
  const dateFormat = new Intl.DateTimeFormat(locale === "id" ? "id-ID" : "en-US", { dateStyle: "medium" })

  if (!user) return null

  if (!canViewAdmin) {
    return (
      <div className="flex flex-col gap-6">
        <PageHeader
          title={t("title")}
          description={t("description")}
          actions={
            <Button asChild variant="outline">
              <Link to="/settings/profile">{tCommon("profile")}</Link>
            </Button>
          }
        />
        <Card className="relative gap-0 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_top_left,color-mix(in_oklab,var(--primary)_14%,transparent),transparent_62%)]" />
          <CardContent className="relative flex flex-col gap-6 p-6 sm:flex-row sm:items-center">
            <BusinessAvatar name={user.name} />
            <div className="flex min-w-0 flex-1 flex-col gap-1">
              <p className="text-xs text-muted-foreground">{tCommon("profile")}</p>
              <p className="font-display text-3xl leading-tight tracking-tight">{user.name}</p>
              <p className="truncate text-sm text-muted-foreground">{user.email}</p>
            </div>
            <div className="flex flex-wrap gap-2">
              <Badge variant="outline">{user.role?.label ?? tCommon("unknown")}</Badge>
              <Badge variant={statusVariant(user.status)}>{tCommon(user.status)}</Badge>
            </div>
          </CardContent>
          <div className="relative grid gap-4 border-t border-border/60 p-6 sm:grid-cols-3">
            <div className="flex flex-col gap-1 text-sm">
              <span className="text-muted-foreground">{tCommon("jobTitle")}</span>
              <span>{user.jobTitle ?? "—"}</span>
            </div>
            <div className="flex flex-col gap-1 text-sm">
              <span className="text-muted-foreground">{tCommon("department")}</span>
              <span>{user.department ?? "—"}</span>
            </div>
            <div className="flex flex-col gap-1 text-sm">
              <span className="text-muted-foreground">{tCommon("phone")}</span>
              <span>{user.phoneNumber ?? "—"}</span>
            </div>
          </div>
        </Card>
        <div className="grid gap-4 sm:grid-cols-3">
          <StatCard label={tCommon("role")} value={user.role?.label ?? tCommon("unknown")} icon={ShieldCheck} />
          <StatCard label={tCommon("status")} value={tCommon(user.status)} icon={CheckCircle2} />
          <StatCard label={tCommon("security")} value={user.twoFactorEnabled ? tCommon("enabled") : tCommon("disabled")} icon={KeyRound} />
        </div>
        <div className="flex flex-wrap gap-2">
          <Button asChild variant="outline"><Link to="/settings/profile">{tCommon("profile")}</Link></Button>
          <Button asChild variant="outline"><Link to="/settings/password">{tCommon("password")}</Link></Button>
          <Button asChild variant="outline"><Link to="/settings/security">{tCommon("security")}</Link></Button>
        </div>
      </div>
    )
  }

  const data = query.data
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title={t("title")}
        description={t("description")}
        actions={<Button asChild><Link to="/admin/users">{t("viewAll")}</Link></Button>}
      />
      {data ? (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <StatCard label={t("stats.totalUsers")} value={numberFormat.format(data.totalUsers)} icon={Users} />
            <StatCard label={t("stats.activeUsers")} value={numberFormat.format(data.activeUsers)} icon={CheckCircle2} />
            <StatCard label={t("stats.pendingUsers")} value={numberFormat.format(data.pendingUsers)} icon={Clock3} />
            <StatCard label={t("stats.roles")} value={numberFormat.format(data.rolesCount)} icon={ShieldCheck} />
          </div>
          <DashboardCharts roleData={data.roleData} signupData={data.signupData} />
          <Card>
            <CardHeader>
              <CardTitle>{t("recentUsers")}</CardTitle>
              <CardAction>
                <Button asChild variant="ghost" size="sm">
                  <Link to="/admin/users">
                    {t("viewAll")}
                    <ChevronRightIcon data-icon="inline-end" />
                  </Link>
                </Button>
              </CardAction>
            </CardHeader>
            <CardContent>
              {data.recentUsers.length === 0 ? (
                <p className="text-sm text-muted-foreground">{t("recentUsersEmpty")}</p>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{tCommon("name")}</TableHead>
                      <TableHead>{tCommon("role")}</TableHead>
                      <TableHead>{tCommon("status")}</TableHead>
                      <TableHead className="text-right">{tUserColumns("createdAt")}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data.recentUsers.map((row) => (
                      <TableRow key={row.id}>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <Avatar className="size-8">
                              <AvatarFallback className="bg-primary/10 text-xs text-primary">{initials(row.name)}</AvatarFallback>
                            </Avatar>
                            <div className="flex flex-col">
                              <span className="font-medium">{row.name}</span>
                              <span className="text-xs text-muted-foreground">{row.email}</span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell><Badge variant="outline">{row.role?.label ?? tCommon("unknown")}</Badge></TableCell>
                        <TableCell><Badge variant={statusVariant(row.status)}>{tCommon(row.status)}</Badge></TableCell>
                        <TableCell className="text-right text-muted-foreground">{dateFormat.format(new Date(row.createdAt))}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </>
      ) : (
        <p className="text-sm text-muted-foreground">{tCommon("loading")}</p>
      )}
    </div>
  )
}

function BusinessAvatar({ name }: { name: string }) {
  return (
    <Avatar className="size-16 rounded-2xl ring-1 ring-border/60">
      <AvatarFallback className="rounded-2xl bg-primary/10 text-xl text-primary">{initials(name)}</AvatarFallback>
    </Avatar>
  )
}
