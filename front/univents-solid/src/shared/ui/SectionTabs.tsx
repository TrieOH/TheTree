import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "@trieoh/ui-solid";

export interface SectionTabItem<T extends string = string> {
  id: T;
  label: string;
  icon?: (props: { class?: string }) => JSX.Element;
  count?: number;
  badge?: string | number;
}

export interface SectionTabsProps<T extends string = string> {
  items: readonly SectionTabItem<T>[];
  active: T;
  onChange: (id: T) => void;
  ariaLabel?: string;
  class?: string;
}

export function SectionTabs<T extends string = string>(
  props: SectionTabsProps<T>,
): JSX.Element {
  return (
    <nav
      class={cn(
        "flex w-full min-w-0 items-center gap-1 overflow-x-auto border-b border-border scrollbar-none",
        props.class,
      )}
      aria-label={props.ariaLabel}
    >
      <For each={props.items}>
        {({ id, label, icon: Icon, count, badge }) => {
          const isActive = () => props.active === id;
          const displayCount = () => count ?? badge;

          return (
            <button
              type="button"
              onClick={() => props.onChange(id)}
              class={cn(
                "group relative inline-flex shrink-0 items-center gap-2 border-b-2 px-3.5 py-2.5 text-sm font-medium transition-colors cursor-pointer",
                isActive()
                  ? "border-primary text-foreground font-semibold"
                  : "border-transparent text-muted-foreground hover:text-foreground",
              )}
              aria-current={isActive() ? "page" : undefined}
            >
              <Show when={Icon}>
                {Icon!({
                  class: cn(
                    "size-4 shrink-0 transition-colors",
                    isActive()
                      ? "text-primary"
                      : "text-muted-foreground group-hover:text-foreground",
                  ),
                })}
              </Show>

              <span class="whitespace-nowrap">{label}</span>

              <Show when={displayCount() !== undefined}>
                <span
                  class={cn(
                    "ml-1 inline-flex shrink-0 items-center rounded-full px-2 py-0.5 text-[11px] font-mono leading-none transition-colors",
                    isActive()
                      ? "bg-primary/15 text-primary font-medium"
                      : "bg-muted text-muted-foreground",
                  )}
                >
                  {displayCount()}
                </span>
              </Show>
            </button>
          );
        }}
      </For>
    </nav>
  );
}
