import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { For, Show, createMemo } from "solid-js";

import ChevronRightIcon from "~icons/lucide/chevron-right";

import { cn } from "@trieoh/ui-solid";

const ChevronRight = ChevronRightIcon as unknown as (props: { class?: string }) => JSX.Element;

/** Portuguese labels for the segments we know; anything else keeps its own text. */
const LABELS: Record<string, string> = {
  admin: "Admin",
  events: "Eventos",
  editions: "Edições",
  products: "Produtos",
  purchases: "Compras",
  programs: "Programação",
};

const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/**
 * Path-derived crumbs, styled like the React one: uppercase micro-type,
 * chevron separators and id/long segments abbreviated so the strip never grows
 * past the viewport.
 */
export function Breadcrumb(): JSX.Element {
  const location = useLocation();

  const crumbs = createMemo(() => {
    const segments = location().pathname.split("/").filter(Boolean);

    return segments.map((segment, index) => {
      const label = LABELS[segment] ?? segment.charAt(0).toUpperCase() + segment.slice(1);
      const display =
        UUID.test(segment) || label.length > 20
          ? `${label.slice(0, 4)}...${label.slice(-2)}`
          : label;

      return {
        path: `/${segments.slice(0, index + 1).join("/")}`,
        label: display,
        last: index === segments.length - 1,
      };
    });
  });

  return (
    <nav
      aria-label="Trilha de navegação"
      class={cn(
        "flex items-center space-x-2 text-muted-foreground",
        "font-bold uppercase tracking-[0.2em] text-[10px]",
        "h-16 border-b border-border/60 px-6",
        "bg-background/95 backdrop-blur-md",
        "overflow-x-auto whitespace-nowrap",
      )}
    >
      <For each={crumbs()}>
        {(crumb, index) => (
          <>
            <Show when={index() > 0}>
              <ChevronRight class="h-3 w-3 shrink-0 text-muted-foreground/40" />
            </Show>

            <Show
              when={!crumb.last}
              fallback={
                <span class="max-w-37.5 truncate text-foreground">{crumb.label}</span>
              }
            >
              <Link
                to={crumb.path}
                class="max-w-37.5 truncate transition-colors hover:text-primary"
              >
                {crumb.label}
              </Link>
            </Show>
          </>
        )}
      </For>
    </nav>
  );
}

