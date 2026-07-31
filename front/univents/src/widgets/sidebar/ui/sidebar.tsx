import { useRouterState } from "@tanstack/react-router";
import { useAuthActions } from "@trieoh/front-core";
import { ChevronLeft, ChevronRight, LogOut } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { useSidebar } from "../hooks/use-sidebar";
import { getAdminShellLabel, getAdminSidebarSections } from "../sidebar-menu";
import { SidebarItem } from "./sidebar-item";

export function Sidebar() {
  const { collapsed, toggleCollapsed, mobileOpen, setMobileOpen } =
    useSidebar();
  const { handleLogout } = useAuthActions();
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const sections = getAdminSidebarSections(pathname);

  return (
    <>
      <div
        aria-hidden="true"
        onClick={() => setMobileOpen(false)}
        className={cn(
          "fixed inset-0 z-40 bg-primary/25 backdrop-blur-sm transition-opacity duration-300 lg:hidden",
          mobileOpen
            ? "pointer-events-auto opacity-100"
            : "pointer-events-none opacity-0",
        )}
      />

      <aside
        role="navigation"
        aria-label="Admin navigation"
        className={cn(
          "fixed inset-y-0 left-0 z-60 flex h-dvh w-72 flex-col border-r border-border/60 bg-card/95 shadow-xl shadow-black/5 backdrop-blur-xl",
          "transition-[width,transform] duration-300 ease-in-out",
          collapsed ? "lg:w-18" : "lg:w-[18rem]",
          mobileOpen ? "translate-x-0" : "-translate-x-full",
          "lg:translate-x-0",
        )}
      >
        <div className="relative flex h-16 shrink-0 items-center border-b border-border/60 px-4">
          <div
            className={cn(
              "min-w-0 flex-1 transition-all duration-200",
              collapsed ? "lg:pointer-events-none lg:opacity-0" : "opacity-100",
            )}
          >
            <p className="truncate text-sm font-semibold text-foreground">
              <span className="mr-2 text-xs uppercase tracking-[0.28em] text-primary">
                Univents
              </span>
              <span className="text-muted-foreground/70">·</span>
              <span className="ml-2">{getAdminShellLabel(pathname).title}</span>
            </p>
          </div>

          <button
            type="button"
            onClick={toggleCollapsed}
            className={cn(
              "hidden p-2 text-muted-foreground rounded-sm transition-all hover:bg-muted hover:text-foreground lg:flex lg:items-center lg:justify-center lg:absolute lg:top-1/2 lg:-translate-y-1/2",
              collapsed
                ? "lg:left-1/2 lg:-translate-x-1/2"
                : "lg:right-4 lg:translate-x-0",
            )}
            aria-label={collapsed ? "Expandir menu" : "Recolher menu"}
            title="Ctrl/Cmd + B"
          >
            {collapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </button>
        </div>

        <nav className="flex-1 space-y-5 overflow-y-auto overflow-x-hidden px-3 py-4">
          {sections.map((section) => (
            <div key={section.title} className="space-y-2">
              {!collapsed && (
                <p className="px-3 text-[10px] font-semibold uppercase tracking-[0.28em] text-muted-foreground/70">
                  {section.title}
                </p>
              )}
              <div className="space-y-1">
                {section.items.map((item) => (
                  <SidebarItem
                    key={item.id}
                    item={item}
                    collapsed={collapsed}
                    pathname={pathname}
                  />
                ))}
              </div>
            </div>
          ))}
        </nav>

        <div className="shrink-0 p-3">
          <div className="mx-1 mb-3 h-px bg-border/70" />
          <button
            type="button"
            className="flex w-full items-center gap-3 rounded-xl px-3 py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
            onClick={handleLogout}
          >
            <LogOut className="h-4 w-4 shrink-0" strokeWidth={2} />
            <span
              className={cn(
                "truncate transition-[opacity,width] duration-200",
                collapsed ? "lg:w-0 lg:opacity-0" : "w-auto opacity-100",
              )}
            >
              Sair
            </span>
          </button>
        </div>
      </aside>
    </>
  );
}
