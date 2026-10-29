import type { JSX } from "@solidjs/web";
import { For, Show, type Accessor } from "solid-js";
import { cn } from "./lib/cn";

export interface BreadcrumbItem {
  path: string;
  label: string;
}

export interface BreadcrumbProps {
  items: readonly BreadcrumbItem[];
  ariaLabel?: string;
  class?: string;
  /** Supply a router link to keep client-side navigation. */
  renderLink?: (item: Accessor<BreadcrumbItem>, props: { class: string }) => JSX.Element;
}

export function Breadcrumb(props: BreadcrumbProps): JSX.Element {
  const linkClass = "max-w-37.5 truncate transition-colors hover:text-primary";

  return (
    <nav
      aria-label={props.ariaLabel ?? "Trilha de navegação"}
      class={cn(
        "flex items-center space-x-2 text-muted-foreground",
        "font-bold uppercase tracking-[0.2em] text-[10px]",
        "h-16 border-b border-border/60 px-6",
        "bg-background/95 backdrop-blur-md",
        "overflow-x-auto whitespace-nowrap min-w-0 max-w-full",
        props.class,
      )}
    >
      <For each={props.items} keyed={(item) => item.path}>
        {(item, index) => (
          <>
            <Show when={index() > 0}>
              <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3 shrink-0 text-muted-foreground/40">
                <path d="m9 18 6-6-6-6" />
              </svg>
            </Show>
            <Show
              when={index() < props.items.length - 1}
              fallback={<span aria-current="page" class="max-w-37.5 truncate text-foreground">{item().label}</span>}
            >
              {props.renderLink
                ? props.renderLink(item, { class: linkClass })
                : <a href={item().path} class={linkClass}>{item().label}</a>}
            </Show>
          </>
        )}
      </For>
    </nav>
  );
}
