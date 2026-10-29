import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { createMemo, Show } from "solid-js";
import { Sidebar as UISidebar, cn } from "@trieoh/ui-solid";
import LogOutIcon from "~icons/lucide/log-out";
import { useSessionActions } from "@/features/auths/hooks/use-session-actions";
import Logo from "@/shared/ui/Logo";
import { getAdminShellLabel, getAdminSidebarSections } from "./sidebar-menu";
import { useSidebar } from "./use-sidebar";

const LogOut = LogOutIcon as unknown as (props: { class?: string }) => JSX.Element;

/** Ported from the React item so `$param` routes match the same way later. */
function isActivePath(pathname: string, href: string, exact?: boolean): boolean {
  const pattern = `^${href
    .split("/")
    .map((segment) => {
      if (!segment) return "";
      if (segment.startsWith("$")) return "[^/]+";
      return segment.replace(/[.*+?^${}()|[\\]\\\\]/g, "\\$&");
    })
    .join("/")}${exact ? "$" : "(?:/.*)?$"}`;

  return new RegExp(pattern).test(pathname);
}

export function Sidebar(): JSX.Element {
  const { collapsed, toggleCollapsed, mobileOpen, setMobileOpen, desktop } = useSidebar();
  const { logout } = useSessionActions();
  const location = useLocation();
  const compact = () => desktop() && collapsed();
  const sections = createMemo(() => getAdminSidebarSections(location().pathname));
  const activeId = () => sections().flatMap((section) => section.items)
    .find((item) => isActivePath(location().pathname, item.to, item.exact))?.id;

  return (
    <>
      <button
        type="button"
        aria-label="Fechar menu"
        tabindex={mobileOpen() ? 0 : -1}
        onClick={() => setMobileOpen(false)}
        class={cn(
          "fixed inset-0 z-40 bg-black/25 transition-opacity duration-200 lg:hidden",
          mobileOpen() ? "opacity-100" : "pointer-events-none invisible opacity-0",
        )}
      />
      <div
        inert={!desktop() && !mobileOpen()}
        class={cn(
          "fixed inset-y-0 left-0 z-60 h-dvh print:hidden transition-transform duration-200 motion-reduce:transition-none lg:translate-x-0",
          mobileOpen() ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <UISidebar
          title="Univents"
          subtitle={getAdminShellLabel(location().pathname).title}
          brand={<div class="size-7 shrink-0"><Logo variant="icon" priority /></div>}
          ariaLabel="Navegação do admin"
          class="bg-card"
          collapsed={compact()}
          onCollapsedChange={() => desktop() ? toggleCollapsed() : setMobileOpen(false)}
          sections={sections()}
          activeId={activeId()}
          onSelect={() => setMobileOpen(false)}
          renderItem={(item, itemProps) => (
            <Link
              {...itemProps}
              to={item().to}
              params={item().params}
              onClick={(event: MouseEvent) => {
                if (isActivePath(location().pathname, item().to, true)) event.preventDefault();
                itemProps.onClick(event);
              }}
            />
          )}
          footer={
            <button
              type="button"
              aria-label="Sair"
              title={compact() ? "Sair" : undefined}
              onClick={() => void logout()}
              class={cn(
                "flex items-center text-[13px] text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                compact() ? "mx-auto size-9 justify-center rounded-full p-0" : "min-h-10 w-full gap-2.5 rounded-xl px-3",
              )}
            >
              <LogOut class="size-4 shrink-0" />
              <Show when={!compact()}><span>Sair da conta</span></Show>
            </button>
          }
        />
      </div>
    </>
  );
}
