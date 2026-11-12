import { Link, useParams } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import { api } from "@/api/client"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { PageHeader } from "@/components/common/page-header"
import { UserForm } from "@/components/users/user-form"
import { Button } from "@/components/ui/button"

export function UserCreatePage() {
  const t = useTranslations("users")
  const tCommon = useTranslations("common")
  usePageMeta(t("create"))
  const query = useQuery({ queryKey: ["users", "form-meta"], queryFn: () => api.users({ page: 1 }) })

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title={t("create")}
        actions={
          <Button asChild variant="outline">
            <Link to="/admin/users">{tCommon("back")}</Link>
          </Button>
        }
      />
      {query.data ? (
        <UserForm mode="create" roles={query.data.assignableRoles} />
      ) : (
        <p className="text-sm text-muted-foreground">{tCommon("loading")}</p>
      )}
    </div>
  )
}

export function UserEditPage() {
  const t = useTranslations("users")
  const tCommon = useTranslations("common")
  const { id } = useParams()
  usePageMeta(t("edit"))
  const query = useQuery({
    queryKey: ["user", id],
    queryFn: () => api.user(id!),
    enabled: !!id,
  })

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title={t("edit")}
        actions={
          <Button asChild variant="outline">
            <Link to="/admin/users">{tCommon("back")}</Link>
          </Button>
        }
      />
      {query.data ? (
        <UserForm mode="edit" user={query.data.user} roles={query.data.assignableRoles} />
      ) : (
        <p className="text-sm text-muted-foreground">{query.isError ? tCommon("noResults") : tCommon("loading")}</p>
      )}
    </div>
  )
}
