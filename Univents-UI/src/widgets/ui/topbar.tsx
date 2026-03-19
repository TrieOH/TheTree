import { Link, useRouterState } from "@tanstack/react-router"
import { Fragment } from "react/jsx-runtime"
import { SidebarTrigger } from "@/shared/ui/shadcn/sidebar"
import { Separator } from "@/shared/ui/shadcn/separator"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/shared/ui/shadcn/breadcrumb"

const LABELS: Record<string, string> = {
  admin: "Admin",
  events: "Eventos",
  editions: "Edições",
  activities: "Atividades",
  tickets: "Ingressos",
  products: "Produtos",
  checkpoints: "Checkpoints",
  participants: "Participantes",
  auth: "Autenticação",
  profile: "Perfil",
}

function useCrumbs() {
  const { location } = useRouterState()
  const segments = location.pathname.split("/").filter(Boolean)
  return segments.map((seg, i) => ({
    label: LABELS[seg] ?? (seg.length > 20 ? `#${seg.slice(0, 6)}…` : seg),
    href: "/" + segments.slice(0, i + 1).join("/"),
    isLast: i === segments.length - 1,
  }))
}

export function AppTopbar() {
  const crumbs = useCrumbs()

  return (
    <header className="flex h-12 shrink-0 items-center gap-3 border-b border-sidebar-border bg-background px-4">
      <SidebarTrigger />
      <Separator orientation="vertical" className="h-4" />
      <Breadcrumb>
        <BreadcrumbList>
          {crumbs.map((crumb, i) => (
            <Fragment key={crumb.href}>
              {i > 0 && <BreadcrumbSeparator />}
              <BreadcrumbItem>
                {crumb.isLast ? (
                  <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
                ) : (
                  <BreadcrumbLink render={<Link to={crumb.href} />}>
                    {crumb.label}
                  </BreadcrumbLink>
                )}
              </BreadcrumbItem>
            </Fragment>
          ))}
        </BreadcrumbList>
      </Breadcrumb>
    </header>
  )
}