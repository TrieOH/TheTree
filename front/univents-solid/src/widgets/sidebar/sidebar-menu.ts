import type { JSX } from "@solidjs/web";

import LayoutDashboardIcon from "~icons/lucide/layout-dashboard";

export interface SidebarMenuItem {
  id: string;
  label: string;
  to: string;
  /** Rendered by the item, which owns the size and colour. */
  icon: (props: { class?: string }) => JSX.Element;
  /** Match the path exactly instead of by prefix. */
  exact?: boolean;
}

export interface SidebarSection {
  title: string;
  items: SidebarMenuItem[];
}

const LayoutDashboard = LayoutDashboardIcon as unknown as (
  props: { class?: string },
) => JSX.Element;

/**
 * Only the sections whose routes exist in this app. The React admin also has
 * uploads, and the whole event/edition tree; each entry lands here with its
 * route — a menu link to a missing route is worse than a missing link.
 */
export function getAdminSidebarSections(): SidebarSection[] {
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
      ],
    },
  ];
}

/** Title shown in the sidebar header and the mobile topbar. */
export function getAdminShellLabel(): { title: string } {
  return { title: "Eventos" };
}
