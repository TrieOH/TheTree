import { Meter, type MeterMarker } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import { projectedShare } from "../attention";
import { PLAN, type Usage } from "../model";
import { formatPercent, formatUsage } from "../format";

/**
 * Used against included. The tick is where usage lands by the end of the cycle
 * if it keeps its pace so far; the red tick is where the container guard stops.
 */
export function UsageBar(props: { usage: Usage; pace: number }) {
  const plan = () => PLAN[props.usage.resource];
  const share = () => props.usage.used / props.usage.included;
  const projected = () => projectedShare(props.usage, props.pace);
  const limit = () => (props.usage.stopsAt ? props.usage.stopsAt / props.usage.included : 1);
  const tight = () => (projected() ?? share()) > Math.min(0.8, limit());
  const markers = (): MeterMarker[] => [
    ...(projected() !== null && projected()! <= 1
      ? [{ value: projected()!, label: "Where usage lands by the end of the cycle at this pace" }]
      : []),
    ...(props.usage.stopsAt
      ? [{ value: limit(), label: "The container guard stops containers here", tone: "destructive" as const }]
      : []),
  ];
  return (
    <div class="space-y-1.5">
      <Meter
        value={share()}
        label={plan().label}
        valueLabel={
          <>
            <span class="text-foreground">{formatPercent(share(), share() < 0.1 ? 1 : 0)}</span> now
            <Show when={projected() !== null}>
              , <span class={tight() ? "text-warn" : "text-foreground"}>{formatPercent(projected(), 0)}</span> by the
              end
            </Show>
          </>
        }
        tone={tight() ? "warning" : "success"}
        markers={markers()}
      />
      <p class="text-[12px] text-muted-foreground">
        {formatUsage(props.usage.used, plan().unit)} of {formatUsage(props.usage.included, plan().unit)} included
        <Show when={props.usage.stopsAt}>
          {(stop) => <>; the container guard stops every API container at {stop().toFixed(1)} GiB-h</>}
        </Show>
        .
      </p>
    </div>
  );
}
