import { Link, useLocation } from "react-router-dom"
import { useTranslations } from "@/lib/i18n"
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import type { NavItem, NavTitleKey } from "@/config/navigation"

function isActivePath(pathname: string, href: string) {
  if (href === "/") return pathname === "/"
  return pathname === href || pathname.startsWith(`${href}/`)
}

export function NavMain({ items, labelKey }: { items: NavItem[]; labelKey?: NavTitleKey }) {
  const { pathname } = useLocation()
  const t = useTranslations("nav")

  if (items.length === 0) return null

  return (
    <SidebarGroup>
      {labelKey ? <SidebarGroupLabel>{t(labelKey)}</SidebarGroupLabel> : null}
      <SidebarGroupContent>
        <SidebarMenu>
          {items.map((item) => {
            const Icon = item.icon
            const active = isActivePath(pathname, item.href)
            return (
              <SidebarMenuItem key={item.href}>
                <SidebarMenuButton asChild isActive={active} tooltip={t(item.titleKey)}>
                  <Link to={item.href} className="relative">
                    {active ? (
                      <span className="absolute top-1/2 left-0 h-4 w-0.5 -translate-y-1/2 rounded-full bg-sidebar-primary" />
                    ) : null}
                    <Icon />
                    <span>{t(item.titleKey)}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
