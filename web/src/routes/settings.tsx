import { useQuery } from "@tanstack/react-query"
import { api } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { DeleteAccountCard } from "@/components/settings/delete-account-card"
import { ProfileForm } from "@/components/settings/profile-form"
import { PasswordForm } from "@/components/settings/password-form"
import { AppearanceSettings } from "@/components/settings/appearance-settings"
import { SecuritySettings } from "@/components/settings/security-settings"

export function ProfileSettingsPage() {
  const t = useTranslations("settings.profile")
  const tCommon = useTranslations("common")
  const { user } = useAuth()
  usePageMeta(t("title"))
  const query = useQuery({ queryKey: ["settings", "account"], queryFn: api.accountOverview, enabled: !!user })

  if (!user) return null
  return (
    <div className="flex flex-col gap-6">
      <ProfileForm user={user} />
      {query.data ? (
        <DeleteAccountCard canDelete={query.data.canDelete} confirmed={query.data.confirmed} />
      ) : (
        <p className="text-sm text-muted-foreground">{tCommon("loading")}</p>
      )}
    </div>
  )
}

export function PasswordSettingsPage() {
  const t = useTranslations("settings.password")
  usePageMeta(t("title"))
  return <PasswordForm />
}

export function AppearanceSettingsPage() {
  const t = useTranslations("settings.appearance")
  usePageMeta(t("title"))
  return <AppearanceSettings />
}

export function SecuritySettingsPage() {
  const t = useTranslations("settings.security")
  const { user } = useAuth()
  usePageMeta(t("title"))
  if (!user) return null
  return <SecuritySettings twoFactorEnabled={user.twoFactorEnabled} />
}
