import { lazy, Suspense } from "react"
import { Navigate, Outlet, Route, Routes, useLocation } from "react-router-dom"
import { useAuth } from "@/lib/auth"
import { AuthLayout } from "@/components/auth/auth-layout"
import { MarketingFooter } from "@/components/marketing/marketing-footer"
import { MarketingHeader } from "@/components/marketing/marketing-header"
import { ProtectedLayout, RequirePermission, SettingsLayout } from "@/components/app-shell/app-layout"
import { ErrorPage } from "@/routes/errors"
import { LandingPage } from "@/routes/landing"
import {
  ConfirmPasswordPage,
  ForgotPasswordPage,
  LoginPage,
  RegisterPage,
  ResetPasswordPage,
  TwoFactorPage,
  VerifyEmailPage,
} from "@/routes/auth"

const DashboardPage = lazy(() => import("@/routes/dashboard").then((module) => ({ default: module.DashboardPage })))
const UsersPage = lazy(() => import("@/routes/users").then((module) => ({ default: module.UsersPage })))
const UserCreatePage = lazy(() => import("@/routes/user-form").then((module) => ({ default: module.UserCreatePage })))
const UserEditPage = lazy(() => import("@/routes/user-form").then((module) => ({ default: module.UserEditPage })))
const RolesPage = lazy(() => import("@/routes/roles").then((module) => ({ default: module.RolesPage })))
const ProfileSettingsPage = lazy(() =>
  import("@/routes/settings").then((module) => ({ default: module.ProfileSettingsPage })),
)
const PasswordSettingsPage = lazy(() =>
  import("@/routes/settings").then((module) => ({ default: module.PasswordSettingsPage })),
)
const AppearanceSettingsPage = lazy(() =>
  import("@/routes/settings").then((module) => ({ default: module.AppearanceSettingsPage })),
)
const SecuritySettingsPage = lazy(() =>
  import("@/routes/settings").then((module) => ({ default: module.SecuritySettingsPage })),
)

function MarketingLayout() {
  return (
    <div className="flex min-h-svh flex-col bg-background">
      <MarketingHeader />
      <main className="flex-1">
        <Outlet />
      </main>
      <MarketingFooter />
    </div>
  )
}

function GuestOnly({ children }: { children: React.ReactNode }) {
  const { user } = useAuth()
  if (user) return <Navigate to="/dashboard" replace />
  return <>{children}</>
}

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  const location = useLocation()
  if (isLoading) {
    return <div className="grid min-h-svh place-items-center text-sm text-muted-foreground">Loading…</div>
  }
  if (!user) {
    const next = encodeURIComponent(`${location.pathname}${location.search}`)
    return <Navigate to={`/login?next=${next}`} replace />
  }
  return <>{children}</>
}

function AuthLayoutRoute() {
  return (
    <AuthLayout>
      <Outlet />
    </AuthLayout>
  )
}

export function App() {
  return (
    <Suspense fallback={<div className="grid min-h-svh place-items-center text-sm text-muted-foreground">Loading…</div>}>
      <Routes>
      <Route element={<MarketingLayout />}>
        <Route index element={<LandingPage />} />
      </Route>

      <Route element={<AuthLayoutRoute />}>
        <Route path="/login" element={<GuestOnly><LoginPage /></GuestOnly>} />
        <Route path="/register" element={<GuestOnly><RegisterPage /></GuestOnly>} />
        <Route path="/forgot-password" element={<GuestOnly><ForgotPasswordPage /></GuestOnly>} />
        <Route path="/reset-password" element={<GuestOnly><ResetPasswordPage /></GuestOnly>} />
        <Route path="/verify-email" element={<VerifyEmailPage />} />
        <Route path="/two-factor" element={<TwoFactorPage />} />
        <Route path="/confirm-password" element={<RequireAuth><ConfirmPasswordPage /></RequireAuth>} />
      </Route>

      <Route element={<ProtectedLayout />}>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route
          path="/admin/users"
          element={<RequirePermission permission="manage_users"><UsersPage /></RequirePermission>}
        />
        <Route
          path="/admin/users/create"
          element={<RequirePermission permission="manage_users"><UserCreatePage /></RequirePermission>}
        />
        <Route
          path="/admin/users/:id/edit"
          element={<RequirePermission permission="manage_users"><UserEditPage /></RequirePermission>}
        />
        <Route
          path="/admin/roles"
          element={<RequirePermission permission="manage_roles"><RolesPage /></RequirePermission>}
        />
        <Route path="/settings" element={<SettingsLayout />}>
          <Route index element={<Navigate to="/settings/profile" replace />} />
          <Route path="profile" element={<ProfileSettingsPage />} />
          <Route path="password" element={<PasswordSettingsPage />} />
          <Route path="appearance" element={<AppearanceSettingsPage />} />
          <Route path="security" element={<SecuritySettingsPage />} />
        </Route>
      </Route>

      <Route path="/403" element={<ErrorPage variant="forbidden" />} />
      <Route path="*" element={<ErrorPage variant="notFound" />} />
      </Routes>
    </Suspense>
  )
}
