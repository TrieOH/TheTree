import type { JSX } from "@solidjs/web";
import { createFileRoute, Outlet, useLocation } from "@tanstack/solid-router";
import { Show } from "solid-js";

import { cn } from "@trieoh/ui-solid";

import { requireAuth } from "@/features/auths/lib/route-guard";
import { Breadcrumb } from "@/shared/ui/Breadcrumb";
import { MobileTopbar } from "@/widgets/sidebar/MobileTopbar";
import { Sidebar } from "@/widgets/sidebar/Sidebar";
import { SidebarProvider, useSidebar } from "@/widgets/sidebar/use-sidebar";

export const Route = createFileRoute("/admin")({
  beforeLoad: requireAuth,
  head: () => ({
    meta: [{ title: "Admin - Univents" }],
  }),
  component: AdminLayout,
});

function AdminLayout(): JSX.Element {
  return (
    <SidebarProvider>
      <AdminShell />
    </SidebarProvider>
  );
}

/**
 * Same shell as the React admin: rail above `lg`, drawer below it, a sticky
 * breadcrumb on desktop and a topbar with the drawer trigger on mobile. The
 * content padding follows the rail width so nothing hides behind it.
 *
 * In fullscreen editor routes (certifications/editor, badges/editor, signatures/editor,
 * programs/calendar, occurrences draw), the sidebar, topbar, and breadcrumb are hidden
 * so the editor fills the full screen.
 */
function AdminShell(): JSX.Element {
  const { collapsed } = useSidebar();
  const location = useLocation();

  const isFullScreenEditor = () => {
    const p = location().pathname.replace(/\/+$/, "");
    return (
      p.endsWith("/certifications/editor") ||
      p.endsWith("/badges/editor") ||
      p.endsWith("/signatures/editor") ||
      p.endsWith("/programs/calendar") ||
      (p.includes("/occurrences/") && p.endsWith("/draw"))
    );
  };

  return (
    <Show
      when={!isFullScreenEditor()}
      fallback={
        <div class="h-dvh overflow-hidden bg-background">
          <Outlet />
        </div>
      }
    >
      <div class="min-h-dvh min-w-0 max-w-full bg-background overflow-x-hidden">
        <Sidebar />

        <div
          class={cn(
            "flex min-h-dvh min-w-0 max-w-full flex-col transition-[padding] duration-300 ease-in-out",
            collapsed() ? "lg:pl-18" : "lg:pl-72",
          )}
        >
          <div class="print:hidden shrink-0 min-w-0 max-w-full">
            <MobileTopbar />
          </div>

          <div class="sticky top-0 z-30 hidden min-w-0 max-w-full shrink-0 bg-card/95 shadow-sm shadow-black/5 lg:block">
            <Breadcrumb />
          </div>

          <main class="flex-1 min-w-0 max-w-full px-4 sm:px-6 py-6 pb-28">
            <Outlet />
          </main>
        </div>
      </div>
    </Show>
  );
}
