import { useState } from "react"
import { SearchLgIcon as Search, XIcon as X } from "@x1nx3r/altics-react-icons"
import { useTranslations } from "@/lib/i18n"
import { Button } from "@/components/ui/button"
import { InputGroup, InputGroupAddon, InputGroupInput } from "@/components/ui/input-group"
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

type RoleOption = { name: string; label: string }
export type UserFilterValues = { search?: string; role?: string; status?: string }

export function UsersToolbar({
  roles,
  initialSearch,
  initialRole,
  initialStatus,
  onApply,
}: {
  roles: RoleOption[]
  initialSearch?: string
  initialRole?: string
  initialStatus?: string
  onApply: (values: UserFilterValues) => void
}) {
  const t = useTranslations("users")
  const tCommon = useTranslations("common")
  const [search, setSearch] = useState(initialSearch ?? "")
  const [role, setRole] = useState(initialRole)
  const [status, setStatus] = useState(initialStatus)
  const hasFilters = !!search || !!role || !!status

  function apply(overrides: UserFilterValues = {}) {
    onApply({ search: search.trim() || undefined, role, status, ...overrides })
  }

  return (
    <form
      className="flex flex-col gap-3 lg:flex-row lg:items-center"
      onSubmit={(event) => {
        event.preventDefault()
        apply()
      }}
    >
      <InputGroup className="lg:max-w-sm">
        <InputGroupAddon>
          <Search />
        </InputGroupAddon>
        <InputGroupInput
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder={t("searchPlaceholder")}
          aria-label={tCommon("search")}
        />
      </InputGroup>
      <div className="flex flex-col gap-2 sm:flex-row">
        <Select
          value={role ?? "all"}
          onValueChange={(value) => {
            const next = value === "all" ? undefined : value
            setRole(next)
            apply({ role: next })
          }}
        >
          <SelectTrigger className="w-full sm:w-44" aria-label={t("filterRole")}>
            <SelectValue placeholder={t("filterRole")} />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="all">{tCommon("all")}</SelectItem>
              {roles.map((option) => (
                <SelectItem key={option.name} value={option.name}>{option.label}</SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select
          value={status ?? "all"}
          onValueChange={(value) => {
            const next = value === "all" ? undefined : value
            setStatus(next)
            apply({ status: next })
          }}
        >
          <SelectTrigger className="w-full sm:w-44" aria-label={t("filterStatus")}>
            <SelectValue placeholder={t("filterStatus")} />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="all">{tCommon("all")}</SelectItem>
              <SelectItem value="active">{tCommon("active")}</SelectItem>
              <SelectItem value="pending">{tCommon("pending")}</SelectItem>
              <SelectItem value="suspended">{tCommon("suspended")}</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>
      <div className="flex items-center gap-2">
        <Button type="submit" variant="secondary">
          <Search data-icon="inline-start" />
          {tCommon("search")}
        </Button>
        {hasFilters ? (
          <Button
            type="button"
            variant="ghost"
            onClick={() => {
              setSearch("")
              setRole(undefined)
              setStatus(undefined)
              onApply({})
            }}
          >
            <X data-icon="inline-start" />
            {tCommon("remove")}
          </Button>
        ) : null}
      </div>
    </form>
  )
}
