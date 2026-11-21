import { Link } from "@tanstack/solid-router";
import { Sparkline, cn } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { findService } from "@/features/catalog/services";
import { timeAgo } from "@/features/delivery/stats";
import { StatusDot, statusLabel } from "@/widgets/shell/StatusDot";
import { formatCount, formatMs, formatPercent } from "../format";
import { DEGRADED_P95_MS, type ServiceMetrics } from "../model";

/** Under this, 30-day success is worth a look. */
const SUCCESS_TARGET = 0.99;

const grid =
  "grid grid-cols-[minmax(0,1fr)_repeat(4,minmax(0,1fr))] gap-x-4 gap-y-3 md:grid-cols-[9rem_minmax(0,1fr)_4.5rem_3rem_4.5rem_6rem] md:items-center";

/**
 * One row per API on the same 24-hour clock, so a spike in one service lines up
 * with the others. Each strip is scaled to its own busiest hour, written beside it.
 */
export function ServiceLanes(props: { services: readonly ServiceMetrics[] }) {
  return (
    <div role="table" aria-label="APIs, last 24 hours" class="text-[13px]">
      <div role="row" class={cn(grid, "hidden border-b border-border pb-2 text-[12px] text-muted-foreground md:grid")}>
        <span role="columnheader">API</span>
        <span role="columnheader">Requests per hour</span>
        <span role="columnheader" class="text-right">
          Requests
        </span>
        <span role="columnheader" class="text-right">
          5xx
        </span>
        <span role="columnheader" class="text-right">
          p95
        </span>
        <span role="columnheader" class="text-right">
          30-day success
        </span>
      </div>
      <For each={props.services}>
        {(service) => {
          const peak = () => Math.max(0, ...service.hourly);
          return (
            <Link
              role="row"
              to="/services/$service"
              params={{ service: service.id }}
              class={cn(
                grid,
                "border-b border-border py-3.5 transition-colors last:border-b-0 hover:bg-foreground/[0.03] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              )}
            >
              <span role="cell" class="col-span-full flex min-w-0 flex-col gap-0.5 md:col-span-1">
                <span class="flex items-center gap-2 font-medium">
                  <StatusDot status={service.status} />
                  {findService(service.id)?.name}
                </span>
                <span class="pl-3.5 text-[12px] text-muted-foreground">
                  <Show
                    when={service.status === "operational"}
                    fallback={
                      <span class={service.status === "degraded" ? "text-warn" : undefined}>
                        {statusLabel(service.status)}
                      </span>
                    }
                  >
                    <Show when={service.lastRequestAt} fallback="no calls yet">
                      {(at) => <>called {timeAgo(at())}</>}
                    </Show>
                  </Show>
                </span>
              </span>
              <span role="cell" class="col-span-full flex min-w-0 items-end gap-2 md:col-span-1">
                <Sparkline
                  type="bar"
                  fill
                  data={service.hourly}
                  height={30}
                  color="var(--accent)"
                  label={`${findService(service.id)?.name}: requests per hour, last 24 hours`}
                  valueFormatter={(value) => `${value} ${value === 1 ? "request" : "requests"}`}
                  class="flex-1"
                />
                <span class="w-12 shrink-0 text-right text-[11.5px] leading-none tabular-nums text-muted-foreground">
                  <Show when={peak() > 0} fallback="idle">
                    {formatCount(peak())}/h
                  </Show>
                </span>
              </span>
              <Stat label="Requests" value={service.requests24h === null ? "—" : formatCount(service.requests24h)} />
              <Stat label="5xx" value={service.errors24h ?? "—"} warn={(service.errors24h ?? 0) > 0} />
              <Stat label="p95" value={formatMs(service.p95Ms24h)} warn={(service.p95Ms24h ?? 0) > DEGRADED_P95_MS} />
              <Stat
                label="30-day success"
                value={formatPercent(service.successRate30d, 1)}
                warn={service.successRate30d !== null && service.successRate30d < SUCCESS_TARGET}
              />
            </Link>
          );
        }}
      </For>
      <div class={cn(grid, "hidden pt-1.5 text-[11.5px] text-muted-foreground md:grid")} aria-hidden="true">
        <span />
        <span class="flex justify-between pr-14">
          <span>24 h ago</span>
          <span>12 h ago</span>
          <span>now</span>
        </span>
      </div>
    </div>
  );
}

/** A number in a lane; on narrow screens it carries its own label instead of the header. */
function Stat(props: { label: string; value: string | number; warn?: boolean }) {
  return (
    <span role="cell" class="flex min-w-0 flex-col md:block md:text-right">
      <span class="text-[11.5px] text-muted-foreground md:hidden">{props.label}</span>
      <span class={cn("tabular-nums", props.warn ? "font-medium text-warn" : undefined)}>{props.value}</span>
    </span>
  );
}
