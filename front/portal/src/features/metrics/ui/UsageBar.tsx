import { Meter } from "@trieoh/ui-solid";
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
      <Meter
        value={share()}
        label={plan().label}
        valueLabel={
          <>
            {formatUsage(props.usage.used, plan().unit)} of {formatUsage(props.usage.included, plan().unit)} ·{" "}
            <span class={ahead() ? "text-warn" : undefined}>{formatPercent(share(), 1)}</span>
          </>
        }
        tone={ahead() ? "warning" : "success"}
        markers={[
          { value: props.pace, label: "Even pace through the cycle" },
          ...(props.usage.stopsAt
            ? [
                {
                  value: props.usage.stopsAt / props.usage.included,
                  label: "The container guard stops containers here",
                  tone: "destructive" as const,
                },
              ]
            : []),
        ]}
      />
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
