import { Dynamic, type JSX } from "@solidjs/web";
import { cn } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import AlertIcon from "~icons/lucide/triangle-alert";
import InfoIcon from "~icons/lucide/info";
import LightbulbIcon from "~icons/lucide/lightbulb";
import OctagonIcon from "~icons/lucide/octagon-alert";

export type CalloutType = "note" | "tip" | "warning" | "danger";

const ICON = { note: InfoIcon, tip: LightbulbIcon, warning: AlertIcon, danger: OctagonIcon };

const TONE: Record<CalloutType, string> = {
  note: "border-border bg-card [&>svg]:text-muted-foreground",
  tip: "border-ok/40 bg-ok/5 [&>svg]:text-ok",
  warning: "border-warn/50 bg-warn/5 [&>svg]:text-warn",
  danger: "border-destructive-foreground/50 bg-destructive/10 [&>svg]:text-destructive dark:[&>svg]:text-destructive-foreground",
};

/** `<Callout type="warning" title="Wakes the API">…</Callout>`; type defaults to note. */
export function Callout(props: { type?: CalloutType; title?: string; children?: JSX.Element }) {
  const type = () => props.type ?? "note";
  return (
    <aside class={cn("not-prose my-5 flex gap-3 rounded-xl border px-4 py-3 text-[13.5px]", TONE[type()])}>
      <Dynamic component={ICON[type()]} aria-hidden="true" class="mt-0.5 size-4 shrink-0" />
      <div class="min-w-0 flex-1 [&>*:last-child]:mb-0 [&>p]:mb-2">
        <Show when={props.title}>
          <p class="mb-1 font-medium">{props.title}</p>
        </Show>
        {props.children}
      </div>
    </aside>
  );
}
