import type { JSX } from "@solidjs/web";
import { For, Show, createSignal, dynamicComponent } from "solid-js";
import { cn } from "./lib/cn";
import {
  tabClass,
  tabCountClass,
  tabIconClass,
  tabListClass,
  type TabsOrientation,
  type TabsSize,
  type TabsVariant,
} from "./lib/tab-styles";

export type { TabsOrientation, TabsSize, TabsVariant };
export type TabsActivation = "automatic" | "manual";

export interface TabsItem<T extends string = string> {
  id: T;
  label: string;
  icon?: (props: { class?: string }) => JSX.Element;
  count?: number;
  badge?: string | number;
  disabled?: boolean;
  /**
   * Panel shown while this tab is selected. Leave it out on every item to use
   * the tabs as a bare switch for content you render yourself (a page section
   * driven by `value`/`onChange`): no panel is mounted then.
   */
  content?: JSX.Element;
}

export interface TabsProps<T extends string = string> {
  items: readonly TabsItem<T>[];
  /** Selected tab. Pass it (with `onChange`) to control the component. */
  value?: T;
  /** Initial tab when uncontrolled. Defaults to the first enabled one. */
  defaultValue?: T;
  onChange?: (id: T) => void;
  /** Visual style of the triggers. */
  variant?: TabsVariant;
  /** Trigger density. */
  size?: TabsSize;
  /** Which way the triggers run. Arrow keys follow it. */
  orientation?: TabsOrientation;
  /** Stretch the triggers to share the full width (horizontal only). */
  fullWidth?: boolean;
  /**
   * `automatic` selects a tab as soon as it receives focus from the arrow
   * keys; `manual` only moves focus and waits for Enter or Space.
   */
  activation?: TabsActivation;
  /** Render every panel and hide the inactive ones, instead of mounting only the selected one. */
  keepMounted?: boolean;
  ariaLabel?: string;
  class?: string;
  /** Extra classes for the row of triggers. */
  listClass?: string;
  /** Extra classes for each panel. */
  panelClass?: string;
}

/**
 * Accessible tabs: a `tablist` of triggers that each control a `tabpanel`,
 * with roving focus and Arrow/Home/End navigation. Give each item `content` and
 * the component renders the panels, or omit it and drive your own content from
 * `value`/`onChange`.
 */
export function Tabs<T extends string = string>(props: TabsProps<T>): JSX.Element {
  const variant = () => props.variant ?? "underline";
  const size = () => props.size ?? "md";
  const orientation = () => props.orientation ?? "horizontal";
  const fullWidth = () => props.fullWidth ?? false;
  const activation = () => props.activation ?? "automatic";
  const keepMounted = () => props.keepMounted ?? false;
  const vertical = () => orientation() === "vertical";
  const hasPanels = () => props.items.some((item) => item.content !== undefined);

  const baseId = `tabs-${Math.random().toString(36).slice(2, 9)}`;
  const tabId = (id: string) => `${baseId}-tab-${id}`;
  const panelId = (id: string) => `${baseId}-panel-${id}`;

  const firstEnabled = () =>
    props.items.find((item) => !item.disabled)?.id ?? props.items[0]?.id;
  const [inner, setInner] = createSignal<string | undefined>(props.defaultValue);
  const active = () => props.value ?? inner() ?? firstEnabled();

  let list!: HTMLDivElement;

  const select = (id: T) => {
    if (props.value === undefined) setInner(() => id as string);
    props.onChange?.(id);
  };

  const enabledTabs = () =>
    Array.from(
      list.querySelectorAll<HTMLButtonElement>('[role="tab"]:not([disabled])'),
    );

  const onKeyDown = (event: KeyboardEvent) => {
    const next = vertical() ? "ArrowDown" : "ArrowRight";
    const previous = vertical() ? "ArrowUp" : "ArrowLeft";
    if (![next, previous, "Home", "End"].includes(event.key)) return;

    const tabs = enabledTabs();
    const current = tabs.indexOf(event.target as HTMLButtonElement);
    if (current === -1) return;

    event.preventDefault();
    const target =
      event.key === "Home"
        ? 0
        : event.key === "End"
          ? tabs.length - 1
          : (current + (event.key === next ? 1 : -1) + tabs.length) % tabs.length;

    tabs[target].focus();
    if (activation() === "automatic") tabs[target].click();
  };

  return (
    <div
      class={cn(
        "flex",
        hasPanels() && "gap-4",
        vertical() ? "flex-row" : "flex-col",
        props.class,
      )}
    >
      <div
        ref={(element) => {
          list = element;
        }}
        role="tablist"
        aria-label={props.ariaLabel}
        aria-orientation={orientation()}
        onKeyDown={onKeyDown}
        class={cn(
          tabListClass({
            variant: variant(),
            orientation: orientation(),
            fullWidth: fullWidth(),
          }),
          vertical() && "w-auto shrink-0",
          props.listClass,
        )}
      >
        <For each={props.items} keyed={(item) => item.id}>
          {(item) => {
            const Icon = dynamicComponent(() => item().icon);
            const isActive = () => active() === item().id;
            const displayCount = () => item().count ?? item().badge;
            const rendered = () =>
              item().content !== undefined && (isActive() || keepMounted());

            return (
              <button
                type="button"
                role="tab"
                id={tabId(item().id)}
                disabled={item().disabled}
                aria-selected={isActive() ? "true" : "false"}
                aria-controls={rendered() ? panelId(item().id) : undefined}
                tabindex={isActive() ? 0 : -1}
                onClick={() => select(item().id)}
                class={tabClass({
                  variant: variant(),
                  orientation: orientation(),
                  size: size(),
                  active: isActive(),
                  fullWidth: fullWidth(),
                })}
              >
                <Show when={item().icon}>
                  <Icon class={tabIconClass(isActive(), variant())} />
                </Show>

                <span class="whitespace-nowrap">{item().label}</span>

                <Show when={displayCount() !== undefined}>
                  <span class={tabCountClass(isActive(), variant())}>
                    {displayCount()}
                  </span>
                </Show>
              </button>
            );
          }}
        </For>
      </div>

      <Show when={hasPanels()}>
        <div class="min-w-0 flex-1">
          <For each={props.items} keyed={(item) => item.id}>
            {(item) => {
              const isActive = () => active() === item().id;

              return (
                <Show
                  when={
                    item().content !== undefined && (isActive() || keepMounted())
                  }
                >
                  <div
                    role="tabpanel"
                    id={panelId(item().id)}
                    aria-labelledby={tabId(item().id)}
                    tabindex="0"
                    hidden={!isActive()}
                    class={cn(
                      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                      props.panelClass,
                    )}
                  >
                    {item().content}
                  </div>
                </Show>
              );
            }}
          </For>
        </div>
      </Show>
    </div>
  );
}
