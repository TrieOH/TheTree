import type { JSX } from "@solidjs/web";
import XIcon from "~icons/lucide/x";
import { For, Show, createEffect, createSignal } from "solid-js";
import { Button, type ButtonVariant } from "./button";
import { cn } from "./lib/cn";

export interface BulkAction {
  id: string;
  label: string;
  onClick: () => void;
  variant?: ButtonVariant;
  disabled?: boolean;
  icon?: JSX.Element;
  iconOnly?: boolean;
}
export interface BulkActionsProps {
  selectedCount: number;
  actions: readonly BulkAction[];
  onClearSelection: () => void;
  busy?: boolean;
  /** Force icons only. Otherwise labels adapt to available space and action count. */
  iconOnly?: boolean;
  /** Keep the action row in layout when empty, useful when sharing a toolbar slot. */
  reserveSpace?: boolean;
  class?: string;
}
/** Selection and action side effects belong to the caller. */
export function BulkActions(props: BulkActionsProps) {
  let region!: HTMLDivElement;
  let counter!: HTMLSpanElement;
  let actions!: HTMLDivElement;
  const [compact, setCompact] = createSignal(true);
  const actionCompact = (action: BulkAction) =>
    !!action.icon &&
    ((action.iconOnly ?? props.iconOnly) === true || compact());
  createEffect(
    () => ({
      options: props.actions.map((action) => ({
        label: action.label,
        icon: !!action.icon,
        forced: (action.iconOnly ?? props.iconOnly) === true,
      })),
      forceClear: props.iconOnly === true,
      count: props.selectedCount,
      visible: props.selectedCount > 0 || props.reserveSpace,
    }),
    ({ options, forceClear, visible }) => {
      if (!visible || !region || !counter || !actions) return;
      let frame = 0;
      const measure = () => {
        frame = 0;
        const labels = [
          ...actions.querySelectorAll<HTMLElement>("[data-bulk-label]"),
        ];
        if (labels.length !== options.length + 1) return;
        // Label scroll widths remain measurable when visually hidden.
        const widths = labels.map((label, index) => {
          const option =
            index < options.length
              ? options[index]
              : { icon: true, forced: forceClear };
          if (option.icon && option.forced) return 32;
          return label.scrollWidth + 24 + (option.icon ? 24 : 0);
        });
        setCompact(
          widths.reduce((sum, width) => sum + width, 0) +
            Math.max(0, widths.length - 1) * 4 >
            actions.clientWidth,
        );
      };
      const schedule = () => {
        if (!frame) frame = requestAnimationFrame(measure);
      };
      measure();
      if (typeof ResizeObserver === "undefined") return;
      const observer = new ResizeObserver(schedule);
      observer.observe(region);
      observer.observe(counter);
      for (const label of actions.querySelectorAll<HTMLElement>(
        "[data-bulk-label]",
      ))
        observer.observe(label);
      return () => {
        observer.disconnect();
        cancelAnimationFrame(frame);
      };
    },
  );
  return (
    <Show when={props.selectedCount > 0 || props.reserveSpace}>
      <div
        ref={(element) => {
          region = element;
        }}
        role="region"
        aria-label="Ações em lote"
        aria-hidden={props.selectedCount === 0 ? "true" : undefined}
        inert={props.selectedCount === 0}
        aria-busy={props.busy ? "true" : "false"}
        class={cn(
          "flex w-full min-h-9 min-w-0 items-center gap-4 text-muted-foreground",
          props.selectedCount === 0 && "invisible pointer-events-none",
          props.class,
        )}
      >
        <span
          ref={(element) => {
            counter = element;
          }}
          class="max-w-[45%] shrink-0 text-xs tabular-nums"
          aria-live="polite"
        >
          {props.selectedCount} selecionado
          {props.selectedCount === 1 ? "" : "s"}
        </span>
        <div
          ref={(element) => {
            actions = element;
          }}
          class="ml-auto flex min-w-0 flex-1 flex-wrap items-center justify-end gap-1"
        >
          <For each={props.actions}>
            {(action) => (
              <Button
                size="sm"
                variant={action.variant ?? "ghost"}
                disabled={props.busy || action.disabled}
                aria-label={action.label}
                title={action.label}
                class={
                  actionCompact(action)
                    ? "size-8 px-0"
                    : "h-auto min-h-8 max-w-full whitespace-normal py-1"
                }
                onClick={() => action.onClick()}
              >
                <Show when={action.icon}>
                  <span class="size-4 shrink-0 [&>svg]:size-full">
                    {action.icon}
                  </span>
                </Show>
                <span
                  data-bulk-label
                  class={
                    actionCompact(action) ? "sr-only" : "min-w-0 wrap-anywhere"
                  }
                >
                  {action.label}
                </span>
              </Button>
            )}
          </For>
          <Button
            size="sm"
            variant="ghost"
            disabled={props.busy}
            aria-label="Limpar seleção"
            title="Limpar seleção"
            class={
              props.iconOnly || compact()
                ? "size-8 px-0"
                : "h-auto min-h-8 max-w-full whitespace-normal py-1"
            }
            onClick={() => props.onClearSelection()}
          >
            <XIcon class="size-4 shrink-0" aria-hidden="true" />
            <span
              data-bulk-label
              class={props.iconOnly || compact() ? "sr-only" : undefined}
            >
              Limpar seleção
            </span>
          </Button>
        </div>
      </div>
    </Show>
  );
}
