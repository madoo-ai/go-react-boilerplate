import { Link, Navigate, Outlet, useLocation } from "react-router-dom"
import { ShieldOffIcon as ShieldOff } from "@x1nx3r/altics-react-icons"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { AppHeader } from "@/components/app-shell/app-header"
import { AppSidebar } from "@/components/app-shell/app-sidebar"
import { PageHeader } from "@/components/common/page-header"
import { SettingsNav } from "@/components/settings/settings-nav"
import { Button } from "@/components/ui/button"
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"

export function ProtectedLayout() {
  const { user, isLoading } = useAuth()
  const location = useLocation()
  const t = useTranslations("common")

  if (isLoading) {
    return <div className="grid min-h-svh place-items-center text-sm text-muted-foreground">{t("loading")}</div>
  }
  if (!user) {
    const next = encodeURIComponent(`${location.pathname}${location.search}`)
    return <Navigate to={`/login?next=${next}`} replace />
  }
  if (user.status !== "active") {
    return <Navigate to="/login?error=account_inactive" replace />
  }

  return (
    <SidebarProvider>
      <AppSidebar user={user} />
      <SidebarInset>
        <AppHeader />
        <div className="flex flex-1 flex-col gap-6 app-canvas p-4 md:p-6 lg:p-8">
          <div key={location.pathname} className="flex flex-col gap-6 motion-safe:animate-page-in">
            <Outlet />
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

export function RequirePermission({ permission, children }: { permission: string; children: React.ReactNode }) {
  const { user } = useAuth()
  const t = useTranslations("errorPages.forbidden")
  const tCommon = useTranslations("common")
  if (!user) return null
  if (!user.permissions.includes(permission)) {
    return (
      <div className="flex min-h-[60svh] items-center justify-center">
        <div className="flex max-w-md flex-col items-center gap-4 text-center">
          <ShieldOff className="size-8 text-muted-foreground" />
          <PageHeader title={t("title")} description={t("description")} />
          <Button asChild variant="outline">
            <Link to="/dashboard">{tCommon("backToDashboard")}</Link>
          </Button>
        </div>
      </div>
    )
  }
  return <>{children}</>
}

export function SettingsLayout() {
  const t = useTranslations("settings")
  usePageMeta(t("title"))

  return (
    <div className="flex flex-col gap-6">
      <PageHeader title={t("title")} description={t("description")} />
      <div className="grid gap-6 lg:grid-cols-[190px_minmax(0,1fr)]">
        <SettingsNav />
        <div className="min-w-0">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
