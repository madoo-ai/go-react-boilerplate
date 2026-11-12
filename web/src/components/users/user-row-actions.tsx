import { useState } from "react"
import { Link } from "react-router-dom"
import { DotsHorizontalIcon as MoreHorizontal, Pencil01Icon as Pencil, Trash02Icon as Trash2, UserCheck01Icon as UserCheck } from "@x1nx3r/altics-react-icons"
import { toast } from "sonner"
import { api, ApiError } from "@/api/client"
import type { ListUser } from "@/api/types"
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
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Spinner } from "@/components/ui/spinner"

export function UserRowActions({
  user,
  canManage,
  canDelete,
  onChanged,
}: {
  user: ListUser
  canManage: boolean
  canDelete: boolean
  onChanged: () => void
}) {
  const t = useTranslations("users")
  const tAll = useTranslations()
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [pending, setPending] = useState<"verify" | "delete" | null>(null)

  if (!canManage) return null

  async function verify() {
    setPending("verify")
    try {
      await api.verifyUser(user.id)
      toast.success(t("messages.verified"))
      onChanged()
    } catch (error) {
      toast.error(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    } finally {
      setPending(null)
    }
  }

  async function remove() {
    setPending("delete")
    try {
      await api.deleteUser(user.id)
      setDeleteOpen(false)
      toast.success(t("messages.deleted"))
      onChanged()
    } catch (error) {
      toast.error(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    } finally {
      setPending(null)
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon-sm" aria-label={tAll("common.openMenu")}>
            <MoreHorizontal />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="min-w-44">
          <DropdownMenuGroup>
            <DropdownMenuItem asChild>
              <Link to={`/admin/users/${user.id}/edit`}>
                <Pencil />
                {t("actions.edit")}
              </Link>
            </DropdownMenuItem>
            {user.status === "pending" ? (
              <DropdownMenuItem onSelect={() => void verify()} disabled={pending !== null}>
                <UserCheck />
                {t("actions.verify")}
              </DropdownMenuItem>
            ) : null}
          </DropdownMenuGroup>
          {canDelete ? (
            <>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" onSelect={() => setDeleteOpen(true)}>
                <Trash2 />
                {t("actions.delete")}
              </DropdownMenuItem>
            </>
          ) : null}
        </DropdownMenuContent>
      </DropdownMenu>

      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("delete.title")}</AlertDialogTitle>
            <AlertDialogDescription>{t("delete.description", { name: user.name })}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={pending !== null}>{tAll("common.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              disabled={pending !== null}
              onClick={(event) => {
                event.preventDefault()
                void remove()
              }}
            >
              {pending === "delete" ? <Spinner /> : null}
              {t("delete.confirm")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
