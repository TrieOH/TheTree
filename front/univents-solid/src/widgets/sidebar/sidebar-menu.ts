import type { JSX } from "@solidjs/web";

import LayoutDashboardIcon from "~icons/lucide/layout-dashboard";
import UploadCloudIcon from "~icons/lucide/upload-cloud";

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
const UploadCloud = UploadCloudIcon as unknown as (
  props: { class?: string },
) => JSX.Element;

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
export function getAdminShellLabel(pathname?: string): { title: string } {
  if (pathname?.startsWith("/admin/uploads")) {
    return { title: "Uploads" };
  }
  return { title: "Eventos" };
}
