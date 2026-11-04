import { cn } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { findService } from "@/features/catalog/services";
import { formatCount, formatMs } from "../format";
import { DEGRADED_P95_MS, type EndpointStats, type ProblemStats } from "../model";

// Slowest endpoints and top error types, on Metrics and on each service page.
// Route templates and problem slugs only: never paths, messages or tenants.

const th = "py-2 pr-4 font-normal";
const td = "py-2.5 pr-4 tabular-nums";

export function EndpointTable(props: { endpoints: EndpointStats[]; showService?: boolean }) {
  return (
    <Show when={props.endpoints.length > 0} fallback={<p class="text-[13px] text-muted-foreground">No requests recorded.</p>}>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-[13px]">
          <thead class="text-[12px] text-muted-foreground">
            <tr class="border-b border-border">
              <Show when={props.showService}>
                <th class={th}>Service</th>
              </Show>
              <th class={th}>Endpoint</th>
              <th class={th}>p95</th>
              <th class="py-2 font-normal">Requests</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <For each={props.endpoints}>
              {(endpoint) => (
                <tr>
                  <Show when={props.showService}>
                    <td class="py-2.5 pr-4 font-medium">{findService(endpoint.service)?.name}</td>
                  </Show>
                  <td class="py-2.5 pr-4 font-mono text-[12.5px]">
                    <span class="text-muted-foreground">{endpoint.method}</span> {endpoint.route}
                  </td>
                  <td class={cn(td, (endpoint.p95Ms ?? 0) > DEGRADED_P95_MS && "text-warn")}>{formatMs(endpoint.p95Ms)}</td>
                  <td class="py-2.5 tabular-nums">{formatCount(endpoint.requests)}</td>
                </tr>
              )}
            </For>
          </tbody>
        </table>
      </div>
    </Show>
  );
}

export function ProblemTable(props: { problems: ProblemStats[]; showService?: boolean }) {
  return (
    <Show when={props.problems.length > 0} fallback={<p class="text-[13px] text-muted-foreground">No problem answers recorded.</p>}>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-[13px]">
          <thead class="text-[12px] text-muted-foreground">
            <tr class="border-b border-border">
              <Show when={props.showService}>
                <th class={th}>Service</th>
              </Show>
              <th class={th}>Problem</th>
              <th class={th}>Status</th>
              <th class="py-2 font-normal">Answers</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <For each={props.problems}>
              {(problem) => (
                <tr>
                  <Show when={props.showService}>
                    <td class="py-2.5 pr-4 font-medium">{findService(problem.service)?.name}</td>
                  </Show>
                  <td class="py-2.5 pr-4 font-mono text-[12.5px]">{problem.problem}</td>
                  <td class={cn(td, problem.status >= 500 && "text-warn")}>{problem.status || "—"}</td>
                  <td class="py-2.5 tabular-nums">{formatCount(problem.count)}</td>
                </tr>
              )}
            </For>
          </tbody>
        </table>
      </div>
    </Show>
  );
}
