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
        "flex items-center justify-between gap-4 rounded-md border border-border bg-background p-3 text-sm cursor-pointer",
        props.disabled && "opacity-50 cursor-not-allowed",
        props.class,
      )}
    >
      <span>
        <span class="block">{props.label}</span>
        <Show when={props.description}>
          <span class="block text-xs text-muted-foreground">
            {props.description}
          </span>
        </Show>
      </span>
      <input
        type="checkbox"
        checked={props.checked}
        disabled={props.disabled}
        onChange={(event) => props.onChange(event.currentTarget.checked)}
        class="size-4 shrink-0 rounded border-input accent-primary focus-visible:outline-2 focus-visible:outline-ring focus-visible:outline-offset-2"
      />
    </label>
  );
}
