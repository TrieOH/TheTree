import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface CheckboxProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label: JSX.Element;
  description?: JSX.Element;
  disabled?: boolean;
  class?: string;
}

export function Checkbox(props: CheckboxProps) {
  return (
    <label
      class={cn(
        "flex items-center justify-between gap-4 rounded-md border bg-background p-3 text-sm transition-colors",
        props.checked && !props.disabled ? "border-accent/60" : "border-border",
        props.disabled
          ? "cursor-not-allowed text-muted-foreground"
          : "cursor-pointer text-foreground hover:border-muted-foreground",
        props.class,
      )}
    >
      <span class="min-w-0">
        <span class="block font-medium">{props.label}</span>
        <Show when={props.description}>
          <span class="mt-0.5 block text-xs leading-relaxed text-muted-foreground">
            {props.description}
          </span>
        </Show>
      </span>
      <span class="relative inline-flex size-5 shrink-0">
        <input
          type="checkbox"
          checked={props.checked}
          disabled={props.disabled}
          onChange={(event) => props.onChange(event.currentTarget.checked)}
          class="peer m-0 size-5 appearance-none rounded border-2 border-muted-foreground bg-background transition-colors checked:border-accent-foreground dark:checked:border-accent checked:bg-accent focus-visible:outline-2 focus-visible:outline-ring focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:border-muted-foreground dark:disabled:border-muted-foreground disabled:bg-muted forced-colors:appearance-auto"
        />
        <svg
          aria-hidden="true"
          viewBox="0 0 20 20"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="pointer-events-none absolute inset-0 size-5 p-0.5 text-accent-foreground opacity-0 peer-checked:opacity-100 peer-disabled:text-muted-foreground forced-colors:hidden"
        >
          <path d="m4 10 4 4 8-8" />
        </svg>
      </span>
    </label>
  );
}
