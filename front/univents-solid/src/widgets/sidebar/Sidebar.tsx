import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { For, Show } from "solid-js";

import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import LogOutIcon from "~icons/lucide/log-out";

import { cn } from "@trieoh/ui-solid";

import { useSessionActions } from "@/features/auths/hooks/use-session-actions";
import Logo from "@/shared/ui/Logo";

import { getAdminShellLabel, getAdminSidebarSections } from "./sidebar-menu";
import { useSidebar } from "./use-sidebar";

const ChevronLeft = ChevronLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronRight = ChevronRightIcon as unknown as (props: { class?: string }) => JSX.Element;
const LogOut = LogOutIcon as unknown as (props: { class?: string }) => JSX.Element;

/** Ported from the React item so `$param` routes match the same way later. */
function isActivePath(pathname: string, href: string, exact?: boolean): boolean {
  const pattern = `^${href
    .split("/")
    .map((segment) => {
      if (!segment) return "";
      if (segment.startsWith("$")) return "[^/]+";
      return segment.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    })
    .join("/")}${exact ? "$" : "(?:/.*)?$"}`;

  return new RegExp(pattern).test(pathname);
}

export function Sidebar(): JSX.Element {
  const { collapsed, toggleCollapsed, mobileOpen, setMobileOpen } = useSidebar();
  const { logout } = useSessionActions();
  const location = useLocation();
  const sections = getAdminSidebarSections();

  return (
    <>
      <div
        aria-hidden="true"
        onClick={() => setMobileOpen(false)}
        class={cn(
          "fixed inset-0 z-40 bg-primary/25 backdrop-blur-sm transition-opacity duration-300 lg:hidden",
          mobileOpen() ? "pointer-events-auto opacity-100" : "pointer-events-none opacity-0",
        )}
      />

      <aside
        role="navigation"
        aria-label="Navegação do admin"
        class={cn(
          "fixed inset-y-0 left-0 z-60 flex h-dvh w-72 flex-col border-r border-border/60 bg-card/95 shadow-xl shadow-black/5 backdrop-blur-xl",
          "print:hidden",
          "transition-[width,transform] duration-300 ease-in-out",
          collapsed() ? "lg:w-18" : "lg:w-[18rem]",
          mobileOpen() ? "translate-x-0" : "-translate-x-full",
          "lg:translate-x-0",
        )}
      >
        <div class="relative flex h-16 shrink-0 items-center border-b border-border/60 px-4">
          <div
            class={cn(
              "min-w-0 flex-1 transition-all duration-200",
              collapsed() ? "lg:pointer-events-none lg:opacity-0" : "opacity-100",
            )}
          >
            <div class="flex items-center gap-2">
              <div class="size-8 shrink-0">
                <Logo variant="icon" imgClassName="object-left" priority />
              </div>
              <span class={cn("text-muted-foreground/70", collapsed() && "hidden")}>·</span>
              <span class={cn("ml-2 truncate text-sm font-semibold", collapsed() && "hidden")}>
                {getAdminShellLabel().title}
              </span>
            </div>
          </div>

          <button
            type="button"
            onClick={() => toggleCollapsed()}
            title="Ctrl/Cmd + B"
            aria-label={collapsed() ? "Expandir menu" : "Recolher menu"}
            class={cn(
              "hidden p-2 text-muted-foreground transition-all hover:bg-muted hover:text-foreground lg:absolute lg:top-1/2 lg:flex lg:-translate-y-1/2 lg:items-center lg:justify-center",
              collapsed() ? "lg:left-1/2 lg:-translate-x-1/2" : "lg:right-4",
            )}
          >
            <Show when={collapsed()} fallback={<ChevronLeft class="h-4 w-4" />}>
              <ChevronRight class="h-4 w-4" />
            </Show>
          </button>
        </div>

        <nav class="flex-1 space-y-5 overflow-y-auto overflow-x-hidden px-3 py-4">
          <For each={sections}>
            {(section) => (
              <div class="space-y-2">
                <Show when={!collapsed()}>
                  <p class="px-3 text-[10px] font-semibold uppercase tracking-[0.28em] text-muted-foreground/70">
                    {section.title}
                  </p>
                </Show>

                <div class="space-y-1">
                  <For each={section.items}>
                    {(item) => {
                      const active = () =>
                        isActivePath(location().pathname, item.to, item.exact);

                      return (
                        <Link
                          to={item.to}
                          onClick={() => setMobileOpen(false)}
                          aria-current={active() ? "page" : undefined}
                          class={cn(
                            "group relative flex items-center gap-3 rounded-xl px-3 py-3 text-sm font-medium outline-none transition-colors duration-150 focus-visible:ring-2 focus-visible:ring-ring",
                            active()
                              ? "text-primary"
                              : "text-muted-foreground hover:bg-muted/70 hover:text-foreground",
                          )}
                        >
                          <span
                            class={cn(
                              "absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-primary transition-all duration-200 ease-out",
                              active() ? "scale-y-100 opacity-100" : "scale-y-75 opacity-0",
                            )}
                          />

                          {item.icon({
                            class: cn(
                              "h-4 w-4 shrink-0 transition-colors",
                              active()
                                ? "text-primary"
                                : "text-muted-foreground group-hover:text-foreground",
                            ),
                          })}

                          <span
                            class={cn(
                              "truncate transition-[opacity,width] duration-200",
                              collapsed() ? "lg:w-0 lg:opacity-0" : "w-auto opacity-100",
                            )}
                          >
                            {item.label}
                          </span>

                          <Show when={collapsed()}>
                            <span class="pointer-events-none absolute left-full top-1/2 z-50 ml-3 hidden -translate-y-1/2 whitespace-nowrap rounded-md bg-primary px-2.5 py-1.5 text-xs font-medium text-primary-foreground opacity-0 shadow-md transition-opacity duration-150 group-hover:opacity-100 lg:block">
                              {item.label}
                            </span>
                          </Show>
                        </Link>
                      );
                    }}
                  </For>
                </div>
              </div>
            )}
          </For>
        </nav>

        <div class="shrink-0 p-3">
          <div class="mx-1 mb-3 h-px bg-border/70" />
          <button
            type="button"
            onClick={() => void logout()}
            class={cn(
              "flex w-full items-center gap-3 rounded-xl px-3 py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive",
              collapsed() && "lg:justify-center lg:px-0",
            )}
          >
            <LogOut class="h-4 w-4 shrink-0" />
            <span
              class={cn(
                "truncate transition-[opacity,width] duration-200",
                collapsed() ? "lg:w-0 lg:opacity-0" : "w-auto opacity-100",
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
