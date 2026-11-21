import { StatusDot } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import type { AttentionItem, MetricsTab } from "../attention";

/**
 * The page's verdict: everything worth a look, in plain sentences, each
 * opening the tab that explains it. When nothing is off it says so in one line.
 */
export function Attention(props: {
  items: readonly AttentionItem[];
  tabLabel: (tab: MetricsTab) => string;
  onOpen: (tab: MetricsTab) => void;
}) {
  return (
    <Show
      when={props.items.length > 0}
      fallback={
        <p class="mb-8 flex items-center gap-2.5 text-[13.5px]">
          <StatusDot tone="success" label="Healthy" size="md" />
          Nothing needs a look right now.
        </p>
      }
    >
      <section aria-labelledby="attention-title" class="mb-10">
        <h2 id="attention-title" class="mb-2 font-heading text-base font-semibold">
          Needs a look
        </h2>
        <ul class="divide-y divide-border border-y border-border">
          <For each={props.items}>
            {(item) => (
              <li>
                <button
                  type="button"
                  onClick={() => props.onOpen(item.tab)}
                  class="group flex w-full items-baseline gap-3 py-3 text-left text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <StatusDot
                    tone={item.tone}
                    label={item.tone === "destructive" ? "Serious" : "Worth a look"}
                    size="md"
                    class="relative top-[-1px] shrink-0"
                  />
                  <span class="min-w-0 flex-1">{item.text}</span>
                  <span class="shrink-0 text-[12.5px] text-muted-foreground underline-offset-4 group-hover:text-foreground group-hover:underline">
                    {props.tabLabel(item.tab)}
                  </span>
                </button>
              </li>
            )}
          </For>
        </ul>
      </section>
    </Show>
  );
}
