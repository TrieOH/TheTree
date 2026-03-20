import { Link, useRouterState } from "@tanstack/react-router"
import { useAuth } from "@soramux/node-auth-sdk/react"
import {
  Calendar,
  ChevronRight,
  Globe,
  Home,
  LogIn,
  LogOut,
  MapPin,
  Package,
  Ticket,
  Zap,
} from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
  SidebarSeparator,
} from "@/shared/ui/shadcn/sidebar"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/shared/ui/shadcn/collapsible"
import { cn } from "@/shared/lib/utils"
import { useAuthActions } from "@/features/auths/hooks/use-auth-actions"


const publicNav = [
  { label: "Início", icon: Home, to: "/", exact: true },
]

const adminNav = [
  {
    exact: false,
    label: "Eventos",
    icon: Globe,
    to: "/admin/events",
    children: [
      { label: "Edições", icon: Calendar, toPattern: "/admin/events/$eventId/editions" },
      { label: "Atividades", icon: Zap, toPattern: "/admin/events/$eventId/editions/$editionId/activities" },
      { label: "Ingressos", icon: Ticket, toPattern: "/admin/events/$eventId/editions/$editionId/tickets" },
      { label: "Produtos", icon: Package, toPattern: "/admin/events/$eventId/editions/$editionId/products" },
      { label: "Checkpoints", icon: MapPin, toPattern: "/admin/events/$eventId/editions/$editionId/checkpoints" },
    ],
  },
]


function useIsActive(to: string, exact = false) {
  const { location } = useRouterState()
  if (to === "/" || to.startsWith("/?")) return location.pathname === "/"
  if (exact) return location.pathname === to
  return location.pathname.startsWith(to)
}

function useAdminEventsCtx() {
  const { location } = useRouterState()
  const m = /\/admin\/events\/([^/]+)(?:\/editions\/([^/]+))?/.exec(location.pathname)
  return {
    eventId: m?.[1],
    editionId: m?.[2],
    isInsideEvents: !!m,
    isInsideEdition: !!m?.[2],
  }
}

const activeBtn = "bg-primary! text-primary-foreground! [&_svg]:text-primary-foreground!"
const activeSub = "text-accent! [&_svg]:text-accent!"


function PublicNavItem({ item }: { item: (typeof publicNav)[number] }) {
  const active = useIsActive(item.to)
  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        isActive={active}
        tooltip={item.label}
        render={<Link to={item.to} />}
        className={cn(active && activeBtn)}
      >
        <item.icon className="h-4 w-4 shrink-0" />
        <span>{item.label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}


// ✅ Componente separado para cada child do adminNav — hooks no topo do componente
function NavSubItem({
  child,
  eventId,
  editionId,
  isInsideEdition,
}: {
  child: (typeof adminNav)[number]["children"][number]
  eventId: string | undefined
  editionId: string | undefined
  isInsideEdition: boolean
}) {
  const resolved = child.toPattern
    .replace("$eventId", eventId ?? "")
    .replace("$editionId", editionId ?? "")

  const childActive = useIsActive(resolved)

  if (child.toPattern.includes("$editionId") && !isInsideEdition) return null

  return (
    <SidebarMenuSubItem className="mt-1">
      <SidebarMenuSubButton
        isActive={childActive}
        render={<Link to={resolved} className="rounded-xs h-8" />}
        className={cn(childActive && activeSub)}
      >
        <child.icon className="h-3.5 w-3.5 shrink-0" />
        <span>{child.label}</span>
      </SidebarMenuSubButton>
    </SidebarMenuSubItem>
  )
}


function NavItem({ item }: { item: (typeof adminNav)[number] }) {
  const isActive = useIsActive(item.to, item.exact)
  const { isInsideEvents, eventId, editionId, isInsideEdition } = useAdminEventsCtx()

  if (!item.children.length) {
    return (
      <SidebarMenuItem>
        <SidebarMenuButton
          isActive={isActive}
          tooltip={item.label}
          render={<Link to={item.to} />}
          className={cn(isActive && activeBtn)}
        >
          <item.icon className="h-4 w-4 shrink-0" />
          <span>{item.label}</span>
        </SidebarMenuButton>
      </SidebarMenuItem>
    )
  }

  return (
    <Collapsible defaultOpen={isActive} className="group/collapsible">
      <SidebarMenuItem>
        <CollapsibleTrigger
          render={
            <SidebarMenuButton
              isActive={isActive}
              tooltip={item.label}
              className={cn(isActive && activeBtn)}
            />
          }
        >
          <Link
            to={item.to}
            className="flex items-center gap-2 flex-1 min-w-0"
            onClick={(e) => { e.stopPropagation(); }}
          >
            <item.icon className="h-4 w-4 shrink-0" />
            <span className={cn(isActive && activeBtn)}>{item.label}</span>
          </Link>
          <ChevronRight className="ml-auto h-3.5 w-3.5 shrink-0 opacity-40 transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
        </CollapsibleTrigger>

        {isInsideEvents && (
          <CollapsibleContent>
            <SidebarMenuSub>
              {item.children.map((child) => (
                <NavSubItem
                  key={child.label}
                  child={child}
                  eventId={eventId}
                  editionId={editionId}
                  isInsideEdition={isInsideEdition}
                />
              ))}
            </SidebarMenuSub>
          </CollapsibleContent>
        )}
      </SidebarMenuItem>
    </Collapsible>
  )
}

export function AppSidebar() {
  const { isAuthenticated } = useAuth()
  const { handleLogout } = useAuthActions();

  return (
    <Sidebar collapsible="icon" variant="sidebar">
      <SidebarHeader className="h-12 justify-center items-center border-b border-sidebar-border">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              tooltip="Univents"
              render={<Link to="/" />}
              className="bg-transparent! text-primary!"
            >
              <span
                className={cn(
                  "flex h-8 w-8 shrink-0 items-center justify-center",
                  "rounded-none font-bold text-xl"
                )}
              >
                UV
              </span>
              <span className="font-semibold tracking-tight">Univents</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="overflow-hidden">
        {!isAuthenticated && (
          <>
            <SidebarGroup>
              <SidebarGroupLabel>Navegar</SidebarGroupLabel>
              <SidebarMenu>
                {publicNav.map((item) => (
                  <PublicNavItem key={item.label} item={item} />
                ))}
              </SidebarMenu>
            </SidebarGroup>

            <SidebarSeparator />

            <SidebarGroup>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    tooltip="Entrar"
                    render={<Link to="/auth" />}
                    className={cn(useIsActive("/auth", true) && activeBtn)}
                    isActive={useIsActive("/auth", true)}
                  >
                    <LogIn className="h-4 w-4 shrink-0" />
                    <span>Entrar</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroup>
          </>
        )}

        {isAuthenticated && (
          <SidebarGroup>
            <SidebarGroupLabel>Admin</SidebarGroupLabel>
            <SidebarMenu>
              {adminNav.map((item) => (
                <NavItem key={item.label} item={item} />
              ))}
            </SidebarMenu>
          </SidebarGroup>
        )}
      </SidebarContent>

      <SidebarFooter>
        {isAuthenticated && (
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                tooltip="Sair"
                onClick={() => { handleLogout() }}
                className="text-destructive hover:bg-destructive/10 hover:text-destructive [&_svg]:text-destructive"
              >
                <LogOut className="h-4 w-4 shrink-0" />
                <span>Sair</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        )}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}