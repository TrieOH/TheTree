import { Dynamic, type JSX } from "@solidjs/web";
import { Show } from "solid-js";
import InfoIcon from "~icons/lucide/info";
import LightbulbIcon from "~icons/lucide/lightbulb";
import OctagonIcon from "~icons/lucide/octagon-alert";
import AlertIcon from "~icons/lucide/triangle-alert";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export type CalloutTone = "info" | "success" | "warning" | "destructive";

export interface CalloutProps {
  tone?: CalloutTone;
  title?: JSX.Element;
  /** Replaces the tone's icon; pass `null` to hide it. */
  icon?: JSX.Element | null;
  shape?: Shape;
  children?: JSX.Element;
  class?: string;
}

const icons = {
  info: InfoIcon,
  success: LightbulbIcon,
  warning: AlertIcon,
  destructive: OctagonIcon,
};

const tones: Record<CalloutTone, string> = {
  info: "border-border bg-card [&>[data-callout-icon]]:text-muted-foreground",
  success:
    "border-success/40 bg-success/5 [&>[data-callout-icon]]:text-success",
  warning:
    "border-warning/50 bg-warning/5 [&>[data-callout-icon]]:text-warning",
  destructive:
    "border-destructive-foreground/50 bg-destructive/10 [&>[data-callout-icon]]:text-destructive dark:[&>[data-callout-icon]]:text-destructive-foreground",
};

/** A highlighted note inside a page: context, a tip, a warning or a danger. */
export function Callout(props: CalloutProps) {
  const tone = () => props.tone ?? "info";
  return (
    <aside
      class={cn(
        "flex gap-3 rounded-xl border px-4 py-3 text-[13.5px] text-foreground",
        tones[tone()],
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <Show when={props.icon !== null}>
        <span
          data-callout-icon
          aria-hidden="true"
          class="mt-0.5 shrink-0 [&_svg]:size-4"
        >
          <Show
            when={props.icon}
            fallback={<Dynamic component={icons[tone()]} class="size-4" />}
          >
            {props.icon}
          </Show>
        </span>
      </Show>
      <div class="min-w-0 flex-1 [&>*:last-child]:mb-0 [&>p]:mb-2">
        <Show when={props.title}>
          <p class="mb-1 font-medium">{props.title}</p>
        </Show>
        {props.children}
      </div>
    </aside>
  );
}
