import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import {
  ChevronSelectorVerticalIcon as ChevronsUpDown,
  Key01Icon as KeyRound,
  LogOut01Icon as LogOut,
  PaletteIcon as Palette,
  ShieldTickIcon as ShieldCheck,
  User01Icon as UserRound,
} from "@x1nx3r/altics-react-icons"
import { toast } from "sonner"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { initials } from "@/lib/utils"
import type { User } from "@/api/types"

const accountLinks = [
  { href: "/settings/profile", labelKey: "profile", icon: UserRound },
  { href: "/settings/password", labelKey: "password", icon: KeyRound },
  { href: "/settings/appearance", labelKey: "appearance", icon: Palette },
  { href: "/settings/security", labelKey: "security", icon: ShieldCheck },
] as const

export function NavUser({ user }: { user: User }) {
  const t = useTranslations("common")
  const navigate = useNavigate()
  const { signOut } = useAuth()
  const [isSigningOut, setIsSigningOut] = useState(false)

  async function handleSignOut() {
    setIsSigningOut(true)
    try {
      await signOut()
      navigate("/login", { replace: true })
    } catch {
      toast.error(t("signOut"))
      setIsSigningOut(false)
    }
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="size-8 rounded-lg">
                <AvatarFallback className="rounded-lg bg-sidebar-primary text-xs text-sidebar-primary-foreground">
                  {initials(user.name)}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.name}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {user.role?.label ?? user.email}
                </span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side="right"
            align="end"
            sideOffset={8}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="size-8 rounded-lg">
                  <AvatarFallback className="rounded-lg bg-primary text-xs text-primary-foreground">
                    {initials(user.name)}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user.name}</span>
                  <span className="truncate text-xs text-muted-foreground">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              {accountLinks.map((link) => {
                const Icon = link.icon
                return (
                  <DropdownMenuItem key={link.href} asChild>
                    <Link to={link.href}>
                      <Icon />
                      {t(link.labelKey)}
                    </Link>
                  </DropdownMenuItem>
                )
              })}
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onSelect={(event) => {
                event.preventDefault()
                void handleSignOut()
              }}
              disabled={isSigningOut}
            >
              <LogOut />
              {t("signOut")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
