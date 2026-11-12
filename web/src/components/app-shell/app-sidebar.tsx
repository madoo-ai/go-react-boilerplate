import { Link } from "react-router-dom"
import { useTranslations } from "@/lib/i18n"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"
import { APP_NAME } from "@/config/app"
import {
  adminGroup,
  adminNavItems,
  filterNavItems,
  mainNavItems,
  shouldShowAdminGroup,
} from "@/config/navigation"
import type { User } from "@/api/types"
import { AppLogoMark } from "./app-logo"
import { NavMain } from "./nav-main"
import { NavUser } from "./nav-user"

export function AppSidebar({ user }: { user: User }) {
  const t = useTranslations("nav")
  const navUser = {
    permissions: user.permissions,
    role: user.role?.name ?? null,
  }
  const mainItems = filterNavItems(mainNavItems, navUser)
  const adminItems = shouldShowAdminGroup(navUser) ? filterNavItems(adminNavItems, navUser) : []

  return (
    <Sidebar collapsible="icon" variant="sidebar">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/dashboard">
                <AppLogoMark />
                <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                  <span className="truncate font-semibold">{APP_NAME}</span>
                  <span className="truncate text-xs text-muted-foreground">
                    {user.role?.label ?? t("dashboard")}
                  </span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={mainItems} />
        <NavMain items={adminItems} labelKey={adminGroup.labelKey} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
