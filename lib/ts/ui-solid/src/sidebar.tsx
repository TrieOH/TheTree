import type { JSX } from "@solidjs/web";
import { For, Show, dynamicComponent, type Accessor } from "solid-js";
import MenuIcon from "~icons/lucide/menu";
import PanelLeftIcon from "~icons/lucide/panel-left";
import PanelLeftCloseIcon from "~icons/lucide/panel-left-close";
import { cn } from "./lib/cn";

export interface SidebarItem {
  id: string;
  label: string;
  icon?: (props: { class?: string }) => JSX.Element;
  badge?: string | number;
  disabled?: boolean;
}

export interface SidebarSection<T extends SidebarItem = SidebarItem> {
  title: string;
  items: readonly T[];
}

export interface SidebarItemRenderProps {
  class: string;
  children: JSX.Element;
  "aria-label": string;
  "aria-current": "page" | undefined;
  "aria-disabled": "true" | undefined;
  tabindex: -1 | undefined;
  title: string | undefined;
  onClick: (event: MouseEvent) => void;
}

export interface SidebarProps<T extends SidebarItem = SidebarItem> {
  title: string;
  subtitle?: string;
  brand?: JSX.Element;
  sections: readonly SidebarSection<T>[];
  activeId?: string;
  collapsed?: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  onSelect?: (item: T) => void;
  action?: JSX.Element;
  footer?: JSX.Element;
  ariaLabel?: string;
  class?: string;
  /** Integrate a router without coupling the sidebar to one. */
  renderItem?: (
    item: Accessor<T>,
    props: SidebarItemRenderProps,
  ) => JSX.Element;
}

/** Inline panel; the app owns positioning and the mobile drawer. */
export function Sidebar<T extends SidebarItem>(
  props: SidebarProps<T>,
): JSX.Element {
  return (
    <aside
      aria-label={props.ariaLabel ?? "Menu lateral"}
      class={cn(
        "flex h-full min-h-0 shrink-0 flex-col overflow-hidden border-r border-border/40 bg-muted/40 text-foreground",
        props.collapsed ? "w-16" : "w-64",
        props.class,
      )}
    >
      <header class={cn("flex shrink-0 items-center gap-2 px-3", props.collapsed ? "h-14 justify-center" : "h-16")}>
        <Show when={!props.collapsed}>
          <div class="flex min-w-0 flex-1 items-center gap-2.5">
            {props.brand}
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold tracking-tight">
                {props.title}
              </p>
              <Show when={props.subtitle}>
                <p class="mt-0.5 truncate text-xs text-muted-foreground">
                  {props.subtitle}
                </p>
              </Show>
            </div>
          </div>
        </Show>
        <button
          type="button"
          aria-label={props.collapsed ? "Expandir menu" : "Recolher menu"}
          aria-expanded={props.collapsed ? "false" : "true"}
          title={props.collapsed ? "Expandir menu" : "Recolher menu"}
          onClick={() => props.onCollapsedChange(!props.collapsed)}
          class={cn("group/sidebar-toggle relative inline-flex shrink-0 cursor-pointer items-center justify-center text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring", props.collapsed ? "size-9 rounded-full" : "ml-auto size-10 rounded-xl")}
        >
          <Show when={props.collapsed}>
            <span
              aria-hidden="true"
              class="absolute inset-0 flex items-center justify-center transition-opacity group-hover/sidebar-toggle:opacity-0 group-focus-visible/sidebar-toggle:opacity-0 motion-reduce:transition-none"
            >
              {props.brand ?? (
                <span class="text-lg font-semibold text-foreground">
                  {props.title.slice(0, 1)}
                </span>
              )}
            </span>
          </Show>
          <span
            aria-hidden="true"
            class={cn(
              "relative size-5 transition-opacity motion-reduce:transition-none",
              props.collapsed && "opacity-0 group-hover/sidebar-toggle:opacity-100 group-focus-visible/sidebar-toggle:opacity-100",
            )}
          >
            <Show when={props.collapsed} fallback={
              <>
                <PanelLeftIcon class="absolute inset-0 size-5 transition-opacity group-hover/sidebar-toggle:opacity-0 group-focus-visible/sidebar-toggle:opacity-0 motion-reduce:transition-none" />
                <PanelLeftCloseIcon class="absolute inset-0 size-5 opacity-0 transition-opacity group-hover/sidebar-toggle:opacity-100 group-focus-visible/sidebar-toggle:opacity-100 motion-reduce:transition-none" />
              </>
            }>
              <MenuIcon class="size-5" />
            </Show>
          </span>
        </button>

      </header>

      <Show when={props.action}>
        <div
          class={cn(
            "shrink-0 px-2 pb-3",
            props.collapsed && "flex justify-center",
          )}
        >
          {props.action}
        </div>
      </Show>

      <nav
        aria-label={props.ariaLabel ?? "Navegação"}
        class={cn("min-h-0 flex-1 overflow-y-auto overflow-x-hidden", props.collapsed ? "space-y-3 px-1.5 py-1" : "space-y-4 px-2 py-2")}
      >
        <For each={props.sections} keyed={(section) => section.title}>
          {(section) => (
            <div class="space-y-0.5">
              <Show when={!props.collapsed}>
                <h2 class="mb-1.5 truncate px-3 text-[11px] font-medium text-muted-foreground">
                  {section().title}
                </h2>
              </Show>
              <For each={section().items} keyed={(item) => item.id}>
                {(item) => {
                  const Icon = dynamicComponent(() => item().icon);
                  const active = () => props.activeId === item().id;
                  const itemProps: SidebarItemRenderProps = {
                    get class() {
                      return cn(
                        "relative flex items-center text-left text-[13px] outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring",
                        active()
                          ? "font-medium text-primary"
                          : "text-muted-foreground hover:bg-foreground/5 hover:text-foreground",
                        item().disabled && "pointer-events-none opacity-40",
                        props.collapsed
                          ? "mx-auto size-9 justify-center rounded-full p-0"
                          : "min-h-10 w-full gap-2.5 rounded-xl px-3",
                      );
                    },
                    get "aria-label"() {
                      return item().label;
                    },
                    get "aria-current"() {
                      return active() ? "page" : undefined;
                    },
                    get "aria-disabled"() {
                      return item().disabled ? "true" : undefined;
                    },
                    get tabindex() {
                      return item().disabled ? -1 : undefined;
                    },
                    get title() {
                      return props.collapsed ? item().label : undefined;
                    },
                    onClick: (event) => {
                      if (item().disabled) event.preventDefault();
                      else props.onSelect?.(item());
                    },
                    get children() {
                      return (
                        <>
                          <Show when={active()}>
                            <span
                              aria-hidden="true"
                              class="pointer-events-none absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-primary"
                            />
                          </Show>
                          <Show
                            when={item().icon}
                            fallback={
                              <span
                                aria-hidden="true"
                                class="inline-flex size-4 shrink-0 items-center justify-center text-xs font-medium"
                              >
                                {item().label.slice(0, 1)}
                              </span>
                            }
                          >
                            <Icon class="size-4 shrink-0" />
                          </Show>
                          <Show when={!props.collapsed}>
                            <span class="min-w-0 flex-1 truncate">
                              {item().label}
                            </span>
                            <Show when={item().badge !== undefined}>
                              <span class="rounded-full bg-background/70 px-2 py-0.5 text-[11px] tabular-nums text-muted-foreground">
                                {item().badge}
                              </span>
                            </Show>
                          </Show>
                        </>
                      );
                    },
                  };
                  return props.renderItem ? (
                    props.renderItem(item, itemProps)
                  ) : (
                    <button
                      type="button"
                      {...itemProps}
                      disabled={item().disabled}
                    />
                  );
                }}
              </For>
            </div>
          )}
        </For>
      </nav>

      <Show when={props.footer}>
        <footer class="shrink-0 p-2">{props.footer}</footer>
      </Show>
    </aside>
  );
}
