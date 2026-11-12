import {
  Key01Icon as KeyRound,
  LayoutGrid01Icon as LayoutGrid,
  PaletteIcon as Palette,
  ShieldTickIcon as ShieldCheck,
  Users01Icon as Users,
} from "@x1nx3r/altics-react-icons"
import type { AppIcon } from "@/types/icon"

export type NavTitleKey =
  | "dashboard"
  | "users"
  | "roles"
  | "administration"
  | "profile"
  | "password"
  | "appearance"
  | "security"

export type NavItem = {
  titleKey: NavTitleKey
  href: string
  icon: AppIcon
  permissions?: string[]
  roles?: string[]
  children?: NavItem[]
}

export type NavUser = {
  permissions: string[]
  role: string | null
}

export const mainNavItems: NavItem[] = [
  { titleKey: "dashboard", href: "/dashboard", icon: LayoutGrid },
]

export const adminNavItems: NavItem[] = [
  {
    titleKey: "users",
    href: "/admin/users",
    icon: Users,
    permissions: ["manage_users"],
  },
  {
    titleKey: "roles",
    href: "/admin/roles",
    icon: KeyRound,
    permissions: ["manage_roles"],
  },
]

export const adminGroup = {
  labelKey: "administration" as const,
  icon: ShieldCheck,
  permissions: ["view_admin"],
}

export function canSeeNavItem(item: NavItem, user: NavUser): boolean {
  if (item.permissions?.length) {
    if (!item.permissions.some((permission) => user.permissions.includes(permission))) {
      return false
    }
  }
  if (item.roles?.length) {
    if (!user.role || !item.roles.includes(user.role)) return false
  }
  return true
}

export function filterNavItems(items: NavItem[], user: NavUser): NavItem[] {
  return items
    .filter((item) => canSeeNavItem(item, user))
    .map((item) => ({
      ...item,
      children: item.children ? filterNavItems(item.children, user) : undefined,
    }))
    .filter((item) => !item.children || item.children.length > 0)
}

export function shouldShowAdminGroup(user: NavUser): boolean {
  return (
    adminGroup.permissions.some((permission) => user.permissions.includes(permission)) &&
    filterNavItems(adminNavItems, user).length > 0
  )
}

export const settingsNavItems: NavItem[] = [
  { titleKey: "profile", href: "/settings/profile", icon: Users },
  { titleKey: "password", href: "/settings/password", icon: KeyRound },
  { titleKey: "appearance", href: "/settings/appearance", icon: Palette },
  { titleKey: "security", href: "/settings/security", icon: ShieldCheck },
]
