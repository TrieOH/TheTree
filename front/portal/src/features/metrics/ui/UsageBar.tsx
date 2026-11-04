import { cn } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import { PLAN, type Usage } from "../model";
import { formatPercent, formatUsage } from "../format";

/**
 * Used against included, with a tick where an even pace through the billing
 * cycle would be now; ahead of the tick means on course to run out.
 */
export function UsageBar(props: { usage: Usage; pace: number }) {
  const plan = () => PLAN[props.usage.resource];
  const share = () => props.usage.used / props.usage.included;
  const ahead = () => share() > props.pace;
  return (
    <div class="space-y-1.5">
      <div class="flex items-baseline justify-between gap-3 text-[13px]">
        <span class="font-medium">{plan().label}</span>
        <span class="tabular-nums text-muted-foreground">
          {formatUsage(props.usage.used, plan().unit)} of {formatUsage(props.usage.included, plan().unit)} ·{" "}
          <span class={cn(ahead() && "text-warn")}>{formatPercent(share(), 1)}</span>
        </span>
      </div>
      <div class="relative h-2 overflow-hidden rounded-full bg-muted">
        <div
          class={cn("h-full rounded-full", ahead() ? "bg-warn" : "bg-ok")}
          style={{ width: `${Math.min(100, share() * 100)}%` }}
        />
        <span
          aria-hidden="true"
          title="Even pace through the cycle"
          class="absolute top-0 h-full w-0.5 bg-foreground/60"
          style={{ left: `${Math.min(100, props.pace * 100)}%` }}
        />
        <Show when={props.usage.stopsAt}>
          {(stop) => (
            <span
              aria-hidden="true"
              title="The container guard stops containers here"
              class="absolute top-0 h-full w-0.5 bg-destructive-foreground"
              style={{ left: `${(stop() / props.usage.included) * 100}%` }}
            />
          )}
        </Show>
      </div>
      <Show when={props.usage.stopsAt}>
        {(stop) => (
          <p class="text-[12px] text-muted-foreground">
            The container guard stops every API container at {stop().toFixed(1)} GiB-h, until the cycle turns over.
          </p>
        )}
      </Show>
    </div>
  );
}
