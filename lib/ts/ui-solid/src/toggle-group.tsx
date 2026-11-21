import type { JSX } from "@solidjs/web";
import { For } from "solid-js";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface ToggleGroupItem<T extends string = string> {
  value: T;
  label: JSX.Element;
  disabled?: boolean;
}

export interface ToggleGroupProps<T extends string = string> {
  items: readonly ToggleGroupItem<T>[];
  value: T;
  onChange: (value: T) => void;
  ariaLabel: string;
  size?: "sm" | "md";
  shape?: Shape;
  class?: string;
}

const sizes = {
  sm: "h-7 px-3 text-[12.5px]",
  md: "h-8 px-3.5 text-[13px]",
};

/** A row of chips where one is picked, like a filter. Each chip is a pressed or unpressed button. */
export function ToggleGroup<T extends string = string>(
  props: ToggleGroupProps<T>,
) {
  return (
    <div
      role="group"
      aria-label={props.ariaLabel}
      class={cn("flex flex-wrap gap-1.5", props.class)}
    >
      <For each={props.items}>
        {(item) => {
          const active = () => item.value === props.value;
          return (
            <button
              type="button"
              aria-pressed={active() ? "true" : "false"}
              disabled={item.disabled}
              onClick={() => props.onChange(item.value)}
              class={cn(
                "rounded-full border transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
                sizes[props.size ?? "sm"],
                active()
                  ? "border-foreground/30 bg-foreground/10 text-foreground"
                  : "border-border text-muted-foreground hover:bg-foreground/5",
                shapeClass(props.shape),
              )}
            >
              {item.label}
            </button>
          );
        }}
      </For>
    </div>
  );
}
