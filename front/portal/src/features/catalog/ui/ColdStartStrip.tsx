import { For } from "solid-js";
import type { ColdStartBar } from "../overview";

const label = (bar: ColdStartBar) => {
  const date = new Date(`${bar.day}T00:00:00Z`).toLocaleDateString("en", { month: "short", day: "numeric", timeZone: "UTC" });
  if (bar.count === 0) return `${date}: no cold starts`;
  return `${date}: ${bar.count} ${bar.count === 1 ? "start" : "starts"}${bar.failed ? `, ${bar.failed} failed` : ""}`;
};

/**
 * Cold starts per day. Each bar is a day's starts; the part that failed sits at
 * its base in the destructive colour, so a bad day reads as a red bar.
 */
export function ColdStartStrip(props: { bars: readonly ColdStartBar[]; height?: number }) {
  const peak = () => Math.max(1, ...props.bars.map((b) => b.count));
  return (
    <div>
      <div class="flex items-end gap-[3px]" style={{ height: `${props.height ?? 56}px` }} role="list" aria-label="Cold starts per day">
        <For each={props.bars}>
          {(bar) => (
            <span
              role="listitem"
              aria-label={label(bar)}
              title={label(bar)}
              class="flex min-w-0 flex-1 flex-col-reverse overflow-hidden rounded-[2px]"
              style={{ height: bar.count ? `${Math.max(8, (bar.count / peak()) * 100)}%` : "2px" }}
            >
              <span
                class="w-full bg-destructive dark:bg-destructive-foreground"
                style={{ height: bar.count ? `${(bar.failed / bar.count) * 100}%` : "0" }}
              />
              <span class={bar.count ? "w-full flex-1 bg-muted-foreground/45" : "w-full flex-1 bg-border"} />
            </span>
          )}
        </For>
      </div>
      <div class="mt-1.5 flex justify-between text-[11.5px] text-muted-foreground" aria-hidden="true">
        <span>{props.bars.length} days ago</span>
        <span>today</span>
      </div>
    </div>
  );
}
