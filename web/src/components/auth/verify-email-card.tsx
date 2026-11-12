import { useEffect, useState } from "react"
import { Link, useSearchParams } from "react-router-dom"
import { api, ApiError } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { AuthCard } from "@/components/auth/auth-card"
import { FormAlert } from "@/components/auth/form-alert"
import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"

type State = "idle" | "verifying" | "success" | "error"

export function VerifyEmailCard() {
  const t = useTranslations("auth.verify")
  const tAll = useTranslations()
  const { user, refresh } = useAuth()
  const [params] = useSearchParams()
  const token = params.get("token")
  const [state, setState] = useState<State>(token ? "verifying" : "idle")
  const [resent, setResent] = useState(false)
  const [pending, setPending] = useState(false)

  usePageMeta(t("title"))

  useEffect(() => {
    if (!token) return
    let active = true
    api
      .verifyEmail(token)
      .then(() => {
        if (active) {
          setState("success")
          void refresh()
        }
      })
      .catch(() => {
        if (active) setState("error")
      })
    return () => {
      active = false
    }
  }, [token, refresh])

  async function resend() {
    if (!user) return
    setPending(true)
    try {
      await api.requestVerification()
      setResent(true)
    } catch (error) {
      if (error instanceof ApiError) setState("error")
    } finally {
      setPending(false)
    }
  }

  return (
    <AuthCard
      title={t("title")}
      description={state === "success" ? undefined : t("description")}
      footer={
        <Button asChild variant="outline" className="w-full">
          <Link to="/login">{t("back")}</Link>
        </Button>
      }
    >
      <div className="flex flex-col gap-4">
        {state === "verifying" ? (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Spinner /> {tAll("common.loading")}
          </div>
        ) : state === "error" ? (
          <FormAlert message={t("error")} />
        ) : state === "success" ? (
          <FormAlert message={t("success")} variant="default" />
        ) : (
          <FormAlert message={t("description")} variant="default" />
        )}
        {resent ? <FormAlert message={t("resent")} variant="default" /> : null}
        {user ? (
          <Button type="button" variant="outline" disabled={pending} onClick={() => void resend()}>
            {pending ? <Spinner data-icon="inline-start" /> : null}
            {t("resend")}
          </Button>
        ) : null}
      </div>
    </AuthCard>
  )
}
