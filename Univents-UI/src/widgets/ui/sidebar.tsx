import { Link, useLocation } from "@tanstack/react-router"
import { useAuth } from "@soramux/node-auth-sdk/react"
import {
  BarChart3,
  BookOpen,
  Calendar,
  ChevronRight,
  Globe,
  Home,
  Info,
  LogIn,
  LogOut,
  MapPin,
  Package,
  Ticket,
  Zap,
} from "lucide-react"
import type {
  LucideIcon} from "lucide-react";

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


interface NavChild {
  label: string
  icon: LucideIcon
  toPattern: string
}

interface NavItemConfig {
  label: string
  icon: LucideIcon
  to: string
  exact: boolean
  children?: NavChild[]
}


const publicNav: NavItemConfig[] = [
  { label: "Início", icon: Home, to: "/", exact: true },
  {
    label: "Eventos",
    icon: Calendar,
    to: "/events",
    exact: false,
    children: [
      { label: "Edições", icon: Calendar, toPattern: "/events/$eventId/editions" },
      { label: "Detalhes", icon: Info, toPattern: "/events/$eventId/editions/$editionId" },
      { label: "Produtos", icon: Package, toPattern: "/events/$eventId/editions/$editionId/products" },
    ],
  },
  { label: "Blog", icon: BookOpen, to: "/blog", exact: false },
  { label: "Comparativo", icon: BarChart3, to: "/comparative", exact: true },
  { label: "Sobre", icon: Info, to: "/about", exact: true },
]

const adminNav: NavItemConfig[] = [
  {
    exact: false,
    label: "Eventos (Adm)",
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
  const location = useLocation()
  if (to === "/" || to.startsWith("/?")) return location.pathname === "/"
  if (exact) return location.pathname === to
  return location.pathname.startsWith(to)
}

function useAdminEventsCtx() {
  const location = useLocation()
  const m = /\/admin\/events\/([^/]+)(?:\/editions\/([^/]+))?/.exec(location.pathname)
  return {
    eventId: m?.[1],
    editionId: m?.[2],
    isInsideEvents: !!m,
    isInsideEdition: !!m?.[2],
  }
}

function useParticipantEventsCtx(): NavContext {
  const location = useLocation()
  if (location.pathname.startsWith("/admin")) {
    return {
      isInsideEvents: false,
      isInsideEdition: false,
      eventId: undefined,
      editionId: undefined,
    }
  }
  const m = /\/events\/([^/]+)(?:\/editions\/([^/]+))?/.exec(location.pathname)
  return {
    eventId: m?.[1],
    editionId: m?.[2],
    isInsideEvents: !!m,
    isInsideEdition: !!m?.[2],
  }
}

const activeBtn = "bg-primary! text-primary-foreground! [&_svg]:text-primary-foreground!"
const activeSub = "text-accent! [&_svg]:text-accent!"


interface NavContext {
  eventId: string | undefined
  editionId: string | undefined
  isInsideEvents: boolean
  isInsideEdition: boolean
}

function NavSubItem({
  child,
  ctx,
}: {
  child: NavChild
  ctx: NavContext
}) {
  const resolved = child.toPattern
    .replace("$eventId", ctx.eventId ?? "")
    .replace("$editionId", ctx.editionId ?? "")

  const childActive = useIsActive(resolved)

  if (child.toPattern.includes("$editionId") && !ctx.isInsideEdition) return null
  if (child.toPattern.includes("$eventId") && !ctx.isInsideEvents) return null

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


function NavItem({
  item,
  ctx,
}: {
  item: NavItemConfig
  ctx: NavContext
}) {
  const isActive = useIsActive(item.to, item.exact)

  if (!item.children?.length) {
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

        {ctx.isInsideEvents && (
          <CollapsibleContent>
            <SidebarMenuSub>
              {item.children.map((child) => (
                <NavSubItem
                  key={child.label}
                  child={child}
                  ctx={ctx}
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
  const { handleLogout } = useAuthActions()
  const adminCtx = useAdminEventsCtx()
  const participantCtx = useParticipantEventsCtx()

  return (
    <Sidebar collapsible="icon" variant="sidebar">
      <SidebarHeader
        className={cn(
          "h-12 justify-center items-center border-b",
          "border-sidebar-border"
        )}
      >
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
        <SidebarGroup>
          <SidebarGroupLabel>Explorar</SidebarGroupLabel>
          <SidebarMenu>
            {publicNav.map((item) => (
              <NavItem key={item.label} item={item} ctx={participantCtx} />
            ))}
          </SidebarMenu>
        </SidebarGroup>

        <SidebarSeparator />

        <SidebarGroup>
          <SidebarGroupLabel>Gerenciamento</SidebarGroupLabel>
          <SidebarMenu>
            {adminNav.map((item) => (
              <NavItem key={item.label} item={item} ctx={adminCtx} />
            ))}
          </SidebarMenu>
        </SidebarGroup>

        {!isAuthenticated && (
          <>
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
      </SidebarContent>

      <SidebarFooter>
        {isAuthenticated && (
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                tooltip="Sair"
                onClick={() => { void handleLogout() }}
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