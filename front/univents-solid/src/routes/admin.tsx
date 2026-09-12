import type { JSX } from "@solidjs/web";
import { createFileRoute, Outlet } from "@tanstack/solid-router";

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
 */
function AdminShell(): JSX.Element {
  const { collapsed } = useSidebar();

  return (
    <div class="min-h-dvh bg-background">
      <Sidebar />

      <div
        class={cn(
          "flex min-h-dvh flex-col transition-[padding] duration-300 ease-in-out",
          collapsed() ? "lg:pl-18" : "lg:pl-72",
        )}
      >
        <div class="print:hidden">
          <MobileTopbar />
        </div>

        <div class="sticky top-0 z-30 hidden bg-card/95 shadow-sm shadow-black/5 lg:block">
          <Breadcrumb />
        </div>

        <main class="flex-1 px-6 py-6 pb-28">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
