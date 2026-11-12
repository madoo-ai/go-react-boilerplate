import { Link, useSearchParams } from "react-router-dom"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { PlusIcon as Plus } from "@x1nx3r/altics-react-icons"
import { api } from "@/api/client"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { DataPagination } from "@/components/common/data-pagination"
import { PageHeader } from "@/components/common/page-header"
import { RegistrationApprovalCard } from "@/components/users/registration-approval-card"
import { UsersTable } from "@/components/users/users-table"
import { UsersToolbar, type UserFilterValues } from "@/components/users/users-toolbar"
import { Button } from "@/components/ui/button"

export function UsersPage() {
  const t = useTranslations("users")
  const tCommon = useTranslations("common")
  usePageMeta(t("title"))
  const [params, setParams] = useSearchParams()
  const queryClient = useQueryClient()

  const page = Number(params.get("page") ?? "1") || 1
  const search = params.get("search") ?? undefined
  const role = params.get("role") ?? undefined
  const status = params.get("status") ?? undefined

  const query = useQuery({
    queryKey: ["users", page, search, role, status],
    queryFn: () => api.users({ page, search, role, status }),
  })

  function applyFilters(values: UserFilterValues) {
    const next = new URLSearchParams()
    if (values.search) next.set("search", values.search)
    if (values.role) next.set("role", values.role)
    if (values.status) next.set("status", values.status)
    setParams(next)
  }

  function changePage(target: number) {
    const next = new URLSearchParams(params)
    if (target > 1) next.set("page", String(target))
    else next.delete("page")
    setParams(next)
  }

  function refresh() {
    void queryClient.invalidateQueries({ queryKey: ["users"] })
    void queryClient.invalidateQueries({ queryKey: ["dashboard"] })
  }

  const data = query.data

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title={t("title")}
        description={t("description")}
        actions={
          <Button asChild>
            <Link to="/admin/users/create">
              <Plus data-icon="inline-start" />
              {t("create")}
            </Link>
          </Button>
        }
      />
      {data ? (
        <>
          <RegistrationApprovalCard key={data.registrationApproval ? "on" : "off"} enabled={data.registrationApproval} onChanged={refresh} />
          <div className="flex flex-col gap-3">
            <UsersToolbar
              key={`${search ?? ""}-${role ?? ""}-${status ?? ""}`}
              roles={data.roles.map((item) => ({ name: item.name, label: item.label }))}
              initialSearch={search}
              initialRole={role}
              initialStatus={status}
              onApply={applyFilters}
            />
            <UsersTable users={data.users} onChanged={refresh} />
            <DataPagination
              page={data.page}
              pageCount={data.pageCount}
              total={data.total}
              pageSize={data.pageSize}
              onPageChange={changePage}
            />
          </div>
        </>
      ) : (
        <p className="text-sm text-muted-foreground">{tCommon("loading")}</p>
      )}
    </div>
  )
}
