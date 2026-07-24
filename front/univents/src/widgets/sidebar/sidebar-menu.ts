import {
  Activity,
  CalendarDays,
  FileText,
  LayoutDashboard,
  LayoutGrid,
  PenLine,
  UploadCloud,
  Users,
  type LucideIcon,
  Boxes,
  Tickets,
} from 'lucide-react'

export interface SidebarMenuItem {
  id: string
  label: string
  to: string
  icon: LucideIcon
  params?: Record<string, string>
  exact?: boolean
}

export interface SidebarSection {
  title: string
  items: SidebarMenuItem[]
}

export interface AdminRouteContext {
  eventId?: string
  editionId?: string
}

const EVENT_ROUTE_RE = /^\/admin\/events\/([^/]+)(?:\/editions(?:\/([^/]+))?)?/

export function getAdminRouteContext(pathname: string): AdminRouteContext {
  const match = pathname.match(EVENT_ROUTE_RE)

  if (!match) return {}

  const [, eventId, editionId] = match
  return {
    eventId,
    editionId,
  }
}

export function getAdminSidebarSections(pathname: string): SidebarSection[] {
  const { eventId, editionId } = getAdminRouteContext(pathname)

  if (eventId && editionId) {
    return [
      {
        title: 'Edição',
        items: [
          {
            id: 'edition-overview',
            label: 'Visão geral',
            to: '/admin/events/$eventId/editions/$editionId',
            params: { eventId, editionId },
            icon: LayoutGrid,
            exact: true,
          },
          {
            id: 'edition-activities',
            label: 'Atividades',
            to: '/admin/events/$eventId/editions/$editionId/activities',
            params: { eventId, editionId },
            icon: Activity,
            exact: false,
          },
          {
            id: 'edition-products',
            label: 'Produtos',
            to: '/admin/events/$eventId/editions/$editionId/products',
            params: { eventId, editionId },
            icon: Boxes,
            exact: false,
          },
          {
            id: 'edition-certifications',
            label: 'Certificados',
            to: '/admin/events/$eventId/editions/$editionId/certifications',
            params: { eventId, editionId },
            icon: FileText,
            exact: false,
          },
          {
            id: 'edition-signatures',
            label: 'Assinaturas',
            to: '/admin/events/$eventId/editions/$editionId/signatures',
            params: { eventId, editionId },
            icon: PenLine,
            exact: false,
          },
          {
            id: 'edition-tickets',
            label: 'Tickets',
            to: '/admin/events/$eventId/editions/$editionId/tickets',
            params: { eventId, editionId },
            icon: Tickets,
            exact: false,
          },
        ],
      },
    ]
  }

  if (eventId) {
    return [
      {
        title: 'Evento',
        items: [
          {
            id: 'event-overview',
            label: 'Visão geral',
            to: '/admin/events/$eventId',
            params: { eventId },
            icon: LayoutGrid,
            exact: true,
          },
          {
            id: 'event-editions',
            label: 'Edições',
            to: '/admin/events/$eventId/editions',
            params: { eventId },
            icon: CalendarDays,
            exact: false,
          },
          {
            id: 'event-members',
            label: 'Membros',
            to: '/admin/events/$eventId/members',
            params: { eventId },
            icon: Users,
            exact: false,
          },
        ],
      },
    ]
  }

  return [
    {
      title: 'Admin',
      items: [
        {
          id: 'events',
          label: 'Eventos',
          to: '/admin/events',
          icon: LayoutDashboard,
          exact: true,
        },
        {
          id: 'uploads',
          label: 'Uploads',
          to: '/admin/uploads',
          icon: UploadCloud,
          exact: true,
        },
      ],
    },
  ]
}

export function getAdminShellLabel(pathname: string) {
  if (pathname.startsWith('/admin/uploads')) {
    return {
      eyebrow: 'Admin Univents',
      title: 'Uploads',
      subtitle: 'Processamento de mídia',
    }
  }

  const { eventId, editionId } = getAdminRouteContext(pathname)

  if (editionId) {
    return {
      eyebrow: 'Admin Univents',
      title: 'Edição',
      subtitle: 'Área da edição',
    }
  }

  if (eventId) {
    return {
      eyebrow: 'Admin Univents',
      title: 'Evento',
      subtitle: 'Área do evento',
    }
  }

  return {
    eyebrow: 'Admin Univents',
    title: 'Eventos',
    subtitle: 'Painel administrativo',
  }
}

export function getAdminBackLink(pathname: string) {
  const { eventId, editionId } = getAdminRouteContext(pathname)

  if (eventId && editionId) {
    const editionOverviewPath = `/admin/events/${eventId}/editions/${editionId}`

    if (pathname === editionOverviewPath) {
      return {
        to: '/admin/events/$eventId',
        params: { eventId },
      }
    }

    return {
      to: '/admin/events/$eventId/editions/$editionId',
      params: { eventId, editionId },
    }
  }

  if (eventId) {
    return {
      to: '/admin/events',
      params: undefined,
    }
  }

  return null
}
