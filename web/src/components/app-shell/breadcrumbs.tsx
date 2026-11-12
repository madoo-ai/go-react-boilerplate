import { Fragment } from "react"
import { Link, useLocation } from "react-router-dom"
import { useTranslations } from "@/lib/i18n"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"

export function Breadcrumbs() {
  const { pathname } = useLocation()
  const tNav = useTranslations("nav")
  const tCommon = useTranslations("common")

  const labels: Record<string, () => string> = {
    admin: () => tNav("administration"),
    users: () => tNav("users"),
    roles: () => tNav("roles"),
    settings: () => tNav("settings"),
    profile: () => tNav("profile"),
    password: () => tNav("password"),
    appearance: () => tNav("appearance"),
    security: () => tNav("security"),
    create: () => tCommon("create"),
    edit: () => tCommon("edit"),
  }

  const segments = pathname.split("/").filter(Boolean)
  const items = segments
    .map((segment, index) => {
      const label = labels[segment]?.()
      if (!label) return null
      return { href: `/${segments.slice(0, index + 1).join("/")}`, label }
    })
    .filter((item): item is { href: string; label: string } => !!item)

  if (items.length === 0 || pathname === "/dashboard") {
    return (
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbPage>{tNav("dashboard")}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
    )
  }

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink asChild>
            <Link to="/dashboard">{tNav("dashboard")}</Link>
          </BreadcrumbLink>
        </BreadcrumbItem>
        {items.map((item, index) => (
          <Fragment key={item.href}>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              {index === items.length - 1 ? (
                <BreadcrumbPage>{item.label}</BreadcrumbPage>
              ) : (
                <BreadcrumbLink asChild>
                  <Link to={item.href}>{item.label}</Link>
                </BreadcrumbLink>
              )}
            </BreadcrumbItem>
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
