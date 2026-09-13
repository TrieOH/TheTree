import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { Show } from "solid-js";

import ChevronLeftIcon from "~icons/lucide/chevron-left";
import MenuIcon from "~icons/lucide/menu";

import { cn } from "@trieoh/ui-solid";
import Logo from "@/shared/ui/Logo";

import { getAdminBackLink, getAdminShellLabel } from "./sidebar-menu";
import { useSidebar } from "./use-sidebar";

const ChevronLeft = ChevronLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Menu = MenuIcon as unknown as (props: { class?: string }) => JSX.Element;

/** Drawer trigger + title, shown only below `lg` where the rail collapses. */
export function MobileTopbar(): JSX.Element {
  const { setMobileOpen } = useSidebar();
  const location = useLocation();

  const label = () => getAdminShellLabel(location().pathname);
  const backLink = () => getAdminBackLink(location().pathname);

  return (
    <header class="sticky top-0 z-30 flex h-16 shrink-0 items-center border-b border-border/60 bg-card/95 px-3 shadow-sm shadow-black/5 backdrop-blur-xl lg:hidden!">
      <div class="flex min-w-0 flex-1 items-center gap-2">
        <Show
          when={backLink()}
          keyed
          fallback={
            <button
              type="button"
              onClick={() => setMobileOpen(true)}
              aria-label="Abrir menu"
              class="inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60"
            >
              <Menu class="h-5 w-5" />
            </button>
          }
        >
          {(link) => (
            <Link
              to={link.to}
              params={link.params}
              class="inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60"
              aria-label="Voltar"
            >
              <ChevronLeft class="h-5 w-5" />
            </Link>
          )}
        </Show>

        <div class="flex min-w-0 items-center gap-2">
          <div class="size-8 shrink-0">
            <Logo variant="icon" />
          </div>
          <h1 class="truncate text-sm font-semibold text-foreground">
            {label().title}
          </h1>
        </div>
      </div>

      <button
        type="button"
        onClick={() => setMobileOpen(true)}
        class={cn(
          "inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60",
          backLink() ? "opacity-100" : "hidden",
        )}
        aria-label="Abrir menu"
      >
        <Menu class="h-5 w-5" />
      </button>
    </header>
  );
}
