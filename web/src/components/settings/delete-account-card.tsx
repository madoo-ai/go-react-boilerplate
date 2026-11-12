import { Link, useNavigate } from "react-router-dom"
import { Trash02Icon as Trash2 } from "@x1nx3r/altics-react-icons"
import { toast } from "sonner"
import { useState } from "react"
import { api, ApiError } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Spinner } from "@/components/ui/spinner"

export function DeleteAccountCard({ canDelete, confirmed }: { canDelete: boolean; confirmed: boolean }) {
  const t = useTranslations("settings.danger")
  const tCommon = useTranslations("common")
  const tAll = useTranslations()
  const navigate = useNavigate()
  const { signOut } = useAuth()
  const [pending, setPending] = useState(false)

  async function handleDelete() {
    setPending(true)
    try {
      await api.deleteAccount()
      toast.success(t("deleted"))
      await signOut()
      navigate("/login?deleted=1", { replace: true })
    } catch (error) {
      toast.error(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
      setPending(false)
    }
  }

  return (
    <Card className="border-destructive/30">
      <CardHeader>
        <CardTitle className="text-destructive">{t("title")}</CardTitle>
        <CardDescription>{t("description")}</CardDescription>
      </CardHeader>
      <CardContent>
        {!canDelete ? (
          <p className="text-sm text-muted-foreground">{tAll("errors.lastSuperadmin")}</p>
        ) : !confirmed ? (
          <Button asChild variant="outline">
            <Link to="/confirm-password?next=/settings/profile">{t("confirmPassword")}</Link>
          </Button>
        ) : (
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive">
                <Trash2 data-icon="inline-start" />
                {t("delete")}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{t("confirmTitle")}</AlertDialogTitle>
                <AlertDialogDescription>{t("confirmDescription")}</AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel disabled={pending}>{tCommon("cancel")}</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  disabled={pending}
                  onClick={(event) => {
                    event.preventDefault()
                    void handleDelete()
                  }}
                >
                  {pending ? <Spinner /> : null}
                  {t("confirm")}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}
      </CardContent>
    </Card>
  )
}
