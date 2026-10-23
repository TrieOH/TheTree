import type { JSX } from "@solidjs/web";
import { createEffect, Show } from "solid-js";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface CheckboxProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label: JSX.Element;
  description?: JSX.Element;
  disabled?: boolean;
  indeterminate?: boolean;
  required?: boolean;
  id?: string;
  name?: string;
  value?: string;
  position?: "left" | "right";
  alignment?: "start" | "center" | "end";
  size?: "sm" | "md" | "lg";
  variant?: "card" | "plain";
  shape?: Shape;
  controlShape?: Shape;
  class?: string;
}

const sizes = { sm: "size-4", md: "size-5", lg: "size-6" };
const alignments = {
  start: "items-start",
  center: "items-center",
  end: "items-end",
};

export function Checkbox(props: CheckboxProps) {
  let input: HTMLInputElement | undefined;
  createEffect(
    () => [!!props.indeterminate, props.checked] as const,
    ([indeterminate]) => {
      if (input) input.indeterminate = indeterminate;
    },
  );
  return (
    <label
      class={cn(
        "flex justify-between gap-3 text-sm transition-colors",
        alignments[props.alignment ?? "center"],
        props.position === "left" && "flex-row-reverse",
        props.variant !== "plain" && "rounded-md border bg-background p-3",
        props.checked && !props.disabled ? "border-accent/60" : "border-border",
        props.disabled
          ? "cursor-not-allowed text-muted-foreground"
          : "cursor-pointer text-foreground hover:border-muted-foreground",
        props.size === "sm" && "text-xs",
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <span class="min-w-0 flex-1">
        <span class="block font-medium">
          {props.label}
          <Show when={props.required}>
            <span
              aria-hidden="true"
              class="ml-1 text-destructive dark:text-destructive-foreground"
            >
              *
            </span>
          </Show>
        </span>
        <Show when={props.description}>
          <span class="mt-0.5 block text-xs leading-relaxed text-muted-foreground">
            {props.description}
          </span>
        </Show>
      </span>
      <span
        class={cn("relative inline-flex shrink-0", sizes[props.size ?? "md"])}
      >
        <input
          ref={(element) => {
            input = element;
          }}
          type="checkbox"
          id={props.id}
          name={props.name}
          value={props.value}
          checked={props.checked}
          disabled={props.disabled}
          required={props.required}
          aria-checked={
            props.indeterminate ? "mixed" : props.checked ? "true" : "false"
          }
          onChange={(event) => props.onChange(event.currentTarget.checked)}
          class={cn(
            "peer m-0 shrink-0 appearance-none rounded border-2 border-muted-foreground bg-background transition-colors checked:border-accent-foreground dark:checked:border-accent checked:bg-accent focus-visible:outline-2 focus-visible:outline-ring focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:border-muted-foreground dark:disabled:border-muted-foreground disabled:bg-muted forced-colors:appearance-auto",
            sizes[props.size ?? "md"],
            props.indeterminate &&
              "border-accent-foreground dark:border-accent bg-accent",
            shapeClass(props.controlShape),
          )}
        />
        <svg
          aria-hidden="true"
          viewBox="0 0 20 20"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class={cn(
            "pointer-events-none absolute inset-0 p-0.5 text-accent-foreground peer-checked:opacity-100 peer-disabled:text-muted-foreground forced-colors:hidden",
            sizes[props.size ?? "md"],
            props.indeterminate ? "opacity-100" : "opacity-0",
          )}
        >
          <path d={props.indeterminate ? "M4 10h12" : "m4 10 4 4 8-8"} />
        </svg>
      </span>
    </label>
  );
}
