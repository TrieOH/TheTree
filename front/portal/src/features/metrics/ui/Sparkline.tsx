import { For } from "solid-js";

/** Requests per hour as small bars; an empty hour is a dot on the baseline. */
export function Sparkline(props: { values: number[]; label: string }) {
  const max = () => Math.max(1, ...props.values);
  return (
    <div role="img" aria-label={props.label} class="flex h-8 items-end gap-px">
      <For each={props.values}>
        {(value) => (
          <span
            class={value ? "w-1.5 rounded-sm bg-accent/70" : "h-px w-1.5 bg-border"}
            style={value ? { height: `${Math.max(8, (value / max()) * 100)}%` } : undefined}
          />
        )}
      </For>
    </div>
  );
}
