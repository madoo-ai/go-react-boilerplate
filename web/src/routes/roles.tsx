import { useQuery, useQueryClient } from "@tanstack/react-query"
import { api } from "@/api/client"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { PageHeader } from "@/components/common/page-header"
import { CreateRoleDialog, EditRoleDialog } from "@/components/roles/role-form-dialog"
import { RoleDeleteDialog } from "@/components/roles/role-delete-dialog"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

export function RolesPage() {
  const t = useTranslations("roles")
  const tCommon = useTranslations("common")
  const tPermissions = useTranslations("permissions")
  usePageMeta(t("title"))
  const queryClient = useQueryClient()
  const query = useQuery({ queryKey: ["roles"], queryFn: api.roles })

  function refresh() {
    void queryClient.invalidateQueries({ queryKey: ["roles"] })
    void queryClient.invalidateQueries({ queryKey: ["users"] })
    void queryClient.invalidateQueries({ queryKey: ["session"] })
  }

  if (!query.data) {
    return <p className="text-sm text-muted-foreground">{tCommon("loading")}</p>
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title={t("title")}
        description={t("description")}
        actions={<CreateRoleDialog permissions={query.data.permissions} onChanged={refresh} />}
      />
      <div className="overflow-x-auto rounded-lg border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("columns.name")}</TableHead>
              <TableHead>{t("columns.description")}</TableHead>
              <TableHead>{t("columns.permissions")}</TableHead>
              <TableHead>{t("columns.users")}</TableHead>
              <TableHead className="w-24 text-right"><span className="sr-only">{tCommon("actions")}</span></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {query.data.roles.map((role) => {
              const isProtected = role.name === "superadmin"
              const isDefault = role.name === "user"
              const canDelete = !isProtected && !isDefault && role.userCount === 0
              return (
                <TableRow key={role.id}>
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{role.label}</span>
                        {isProtected ? <Badge variant="secondary">{tCommon("protected")}</Badge> : null}
                      </div>
                      <span className="text-xs text-muted-foreground">{role.name}</span>
                    </div>
                  </TableCell>
                  <TableCell className="max-w-xs text-sm text-muted-foreground">{role.description ?? "—"}</TableCell>
                  <TableCell>
                    {role.permissions.length === 0 ? (
                      <span className="text-sm text-muted-foreground">—</span>
                    ) : (
                      <div className="flex flex-wrap gap-1">
                        {role.permissions.map((key) => (
                          <Badge key={key} variant="outline">{tPermissions(key)}</Badge>
                        ))}
                      </div>
                    )}
                  </TableCell>
                  <TableCell className="tabular-nums">{role.userCount}</TableCell>
                  <TableCell>
                    <div className="flex items-center justify-end gap-1">
                      {isProtected ? (
                        <span className="text-xs text-muted-foreground">{tCommon("protected")}</span>
                      ) : (
                        <>
                          <EditRoleDialog role={role} permissions={query.data.permissions} onChanged={refresh} />
                          <RoleDeleteDialog role={{ id: role.id, label: role.label }} disabled={!canDelete} onChanged={refresh} />
                        </>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
