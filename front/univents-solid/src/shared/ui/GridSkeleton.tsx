import { For } from "solid-js";
import { cn } from "@trieoh/ui-solid";

/** Placeholder for an admin card grid; match `itemClass` height to the real cards. */
export function GridSkeleton(props: {
  /** CSS length for the grid's `minmax()` column width. */
  minItemWidth?: string;
  itemClass?: string;
  count?: number;
}) {
  return (
    <div
      class="grid gap-2"
      style={{
        "grid-template-columns": `repeat(auto-fill, minmax(${props.minItemWidth ?? "16rem"}, 1fr))`,
      }}
    >
      <For each={Array.from({ length: props.count ?? 4 }, (_, index) => index)}>
        {() => (
          <div class={cn("h-24 animate-pulse rounded-xl bg-muted", props.itemClass)} />
        )}
      </For>
    </div>
  );
}
