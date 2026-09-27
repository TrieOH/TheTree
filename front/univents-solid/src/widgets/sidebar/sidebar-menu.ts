import type { JSX } from "@solidjs/web";

import ActivityIcon from "~icons/lucide/activity";
import BadgeCheckIcon from "~icons/lucide/badge-check";
import BoxesIcon from "~icons/lucide/boxes";
import CalendarDaysIcon from "~icons/lucide/calendar-days";
import FileTextIcon from "~icons/lucide/file-text";
import LayoutDashboardIcon from "~icons/lucide/layout-dashboard";
import LayoutGridIcon from "~icons/lucide/layout-grid";
import PenLineIcon from "~icons/lucide/pen-line";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";
import TicketsIcon from "~icons/lucide/tickets";
import UploadCloudIcon from "~icons/lucide/upload-cloud";
import UsersIcon from "~icons/lucide/users";

export interface SidebarMenuItem {
  id: string;
  label: string;
  to: string;
  /** Rendered by the item, which owns the size and colour. */
  icon: (props: { class?: string }) => JSX.Element;
  params?: Record<string, string>;
  /** Match the path exactly instead of by prefix. */
  exact?: boolean;
}

export interface SidebarSection {
  title: string;
  items: SidebarMenuItem[];
}

export interface AdminRouteContext {
  eventId?: string;
  editionId?: string;
}

const Activity = ActivityIcon as unknown as (props: { class?: string }) => JSX.Element;
const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Boxes = BoxesIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileText = FileTextIcon as unknown as (props: { class?: string }) => JSX.Element;
const LayoutDashboard = LayoutDashboardIcon as unknown as (props: { class?: string }) => JSX.Element;
const LayoutGrid = LayoutGridIcon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShoppingBag = ShoppingBagIcon as unknown as (props: { class?: string }) => JSX.Element;
const Tickets = TicketsIcon as unknown as (props: { class?: string }) => JSX.Element;
const UploadCloud = UploadCloudIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

const EVENT_ROUTE_RE = /^\/admin\/events\/([^/]+)(?:\/editions(?:\/([^/]+))?)?/;

export function getAdminRouteContext(pathname: string = ""): AdminRouteContext {
  const match = pathname.match(EVENT_ROUTE_RE);

  if (!match) return {};

  const [, eventId, editionId] = match;
  return {
    eventId,
    editionId,
  };
}

export function getAdminSidebarSections(pathname: string = ""): SidebarSection[] {
  const { eventId, editionId } = getAdminRouteContext(pathname);

  if (eventId && editionId) {
    return [
      {
        title: "Edição",
        items: [
          {
            id: "edition-overview",
            label: "Visão geral",
            to: "/admin/events/$eventId/editions/$editionId",
            params: { eventId, editionId },
            icon: LayoutGrid,
            exact: true,
          },
          {
            id: "edition-programs",
            label: "Programação",
            to: "/admin/events/$eventId/editions/$editionId/programs",
            params: { eventId, editionId },
            icon: Activity,
            exact: false,
          },
          {
            id: "edition-products",
            label: "Produtos",
            to: "/admin/events/$eventId/editions/$editionId/products",
            params: { eventId, editionId },
            icon: Boxes,
            exact: false,
          },
          {
            id: "edition-purchases",
            label: "Compras",
            to: "/admin/events/$eventId/editions/$editionId/purchases",
            params: { eventId, editionId },
            icon: ShoppingBag,
            exact: false,
          },
          {
            id: "edition-badges",
            label: "Crachás",
            to: "/admin/events/$eventId/editions/$editionId/badges",
            params: { eventId, editionId },
            icon: BadgeCheck,
            exact: false,
          },
          {
            id: "edition-certifications",
            label: "Certificados",
            to: "/admin/events/$eventId/editions/$editionId/certifications",
            params: { eventId, editionId },
            icon: FileText,
            exact: false,
          },
          {
            id: "edition-signatures",
            label: "Assinaturas",
            to: "/admin/events/$eventId/editions/$editionId/signatures",
            params: { eventId, editionId },
            icon: PenLine,
            exact: false,
          },
          {
            id: "edition-tickets",
            label: "Tickets",
            to: "/admin/events/$eventId/editions/$editionId/tickets",
            params: { eventId, editionId },
            icon: Tickets,
            exact: false,
          },
        ],
      },
    ];
  }

  if (eventId) {
    return [
      {
        title: "Evento",
        items: [
          {
            id: "event-overview",
            label: "Visão geral",
            to: "/admin/events/$eventId",
            params: { eventId },
            icon: LayoutGrid,
            exact: true,
          },
          {
            id: "event-editions",
            label: "Edições",
            to: "/admin/events/$eventId/editions",
            params: { eventId },
            icon: CalendarDays,
            exact: false,
          },
          {
            id: "event-members",
            label: "Membros",
            to: "/admin/events/$eventId/members",
            params: { eventId },
            icon: Users,
            exact: false,
          },
        ],
      },
    ];
  }

  return [
    {
      title: "Admin",
      items: [
        {
          id: "events",
          label: "Eventos",
          to: "/admin/events",
          icon: LayoutDashboard,
          exact: true,
        },
        {
          id: "uploads",
          label: "Uploads",
          to: "/admin/uploads",
          icon: UploadCloud,
          exact: true,
        },
      ],
    },
  ];
}

/** Title shown in the sidebar header and the mobile topbar. */
export function getAdminShellLabel(pathname: string = ""): {
  eyebrow: string;
  title: string;
  subtitle: string;
} {
  if (pathname.startsWith("/admin/uploads")) {
    return {
      eyebrow: "Admin Univents",
      title: "Uploads",
      subtitle: "Processamento de mídia",
    };
  }

  const { eventId, editionId } = getAdminRouteContext(pathname);

  if (editionId) {
    return {
      eyebrow: "Admin Univents",
      title: "Edição",
      subtitle: "Área da edição",
    };
  }

  if (eventId) {
    return {
      eyebrow: "Admin Univents",
      title: "Evento",
      subtitle: "Área do evento",
    };
  }

  return {
    eyebrow: "Admin Univents",
    title: "Eventos",
    subtitle: "Painel administrativo",
  };
}

export function getAdminBackLink(pathname: string = ""): {
  to: string;
  params?: Record<string, string>;
} | null {
  const { eventId, editionId } = getAdminRouteContext(pathname);

  if (eventId && editionId) {
    const editionOverviewPath = `/admin/events/${eventId}/editions/${editionId}`;

    if (pathname === editionOverviewPath) {
      return {
        to: "/admin/events/$eventId",
        params: { eventId },
      };
    }

    return {
      to: "/admin/events/$eventId/editions/$editionId",
      params: { eventId, editionId },
    };
  }

  if (eventId) {
    return {
      to: "/admin/events",
      params: undefined,
    };
  }

  return null;
}
