import { useState } from "react"
import { Trash02Icon as Trash2 } from "@x1nx3r/altics-react-icons"
import { toast } from "sonner"
import { api, ApiError } from "@/api/client"
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
import { Spinner } from "@/components/ui/spinner"

export function RoleDeleteDialog({
  role,
  disabled,
  onChanged,
}: {
  role: { id: string; label: string }
  disabled: boolean
  onChanged: () => void
}) {
  const t = useTranslations("roles")
  const tCommon = useTranslations("common")
  const tAll = useTranslations()
  const [open, setOpen] = useState(false)
  const [pending, setPending] = useState(false)

  async function remove() {
    setPending(true)
    try {
      await api.deleteRole(role.id)
      setOpen(false)
      toast.success(t("messages.deleted"))
      onChanged()
    } catch (error) {
      toast.error(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    } finally {
      setPending(false)
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="ghost" size="icon-sm" aria-label={tCommon("delete")} disabled={disabled}>
          <Trash2 />
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t("delete.title")}</AlertDialogTitle>
          <AlertDialogDescription>{t("delete.description", { name: role.label })}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={pending}>{tCommon("cancel")}</AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            disabled={pending}
            onClick={(event) => {
              event.preventDefault()
              void remove()
            }}
          >
            {pending ? <Spinner /> : null}
            {t("delete.confirm")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
