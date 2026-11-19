import { For, Show, createMemo } from "solid-js";
import { Loading } from "./loading";
import { cn } from "./lib/cn";
export interface DashboardBarListItem {
  id: string;
  label: string;
  value: number;
  detail?: string;
  /** Tailwind background class, matching existing Univents dashboards. */
  color?: string;
}
export interface DashboardBarListProps {
  items: readonly DashboardBarListItem[];
  maxValue?: number;
  total?: number;
  showPercentage?: boolean;
  valueFormatter?: (value: number) => string;
  emptyMessage?: string;
  loading?: boolean;
  class?: string;
}
export function DashboardBarList(props: DashboardBarListProps) {
  const values = createMemo(() =>
    props.items.map((item) => (Number.isFinite(item.value) ? item.value : 0)),
  );
  const maximum = () =>
    props.maxValue !== undefined &&
    Number.isFinite(props.maxValue) &&
    props.maxValue > 0
      ? props.maxValue
      : Math.max(1, ...values());
  const total = () =>
    props.total !== undefined && Number.isFinite(props.total) && props.total > 0
      ? props.total
      : 0;
  return (
    <div class={cn("min-w-0", props.class)}>
      <Show
        when={!props.loading}
        fallback={<Loading size="sm" tone="muted" class="min-h-24" />}
      >
        <Show
          when={props.items.length > 0}
          fallback={
            <p class="text-sm text-muted-foreground">
              {props.emptyMessage ?? "Nenhum dado disponível."}
            </p>
          }
        >
          <ul class="space-y-4">
            <For each={props.items}>
              {(item, index) => {
                const value = () => values()[index()] ?? 0;
                const width = () =>
                  Math.max(0, Math.min(100, (value() / maximum()) * 100));
                return (
                  <li class="min-w-0 space-y-1.5">
                    <div class="flex min-w-0 flex-wrap items-start justify-between gap-x-3 gap-y-1 text-xs">
                      <span class="min-w-0 font-medium wrap-anywhere">
                        {item.label}
                      </span>
                      <span class="ml-auto shrink-0 text-muted-foreground tabular-nums">
                        {item.detail ??
                          props.valueFormatter?.(value()) ??
                          value()}
                        {props.showPercentage
                          ? ` · ${total() ? ((value() / total()) * 100).toFixed(0) : 0}%`
                          : ""}
                      </span>
                    </div>
                    <div
                      role="meter"
                      aria-label={item.label}
                      aria-valuemin={0}
                      aria-valuemax={maximum()}
                      aria-valuenow={Math.min(maximum(), Math.max(0, value()))}
                      aria-valuetext={item.detail ?? String(value())}
                      class="h-2 overflow-hidden rounded-full bg-muted"
                    >
                      <div
                        class={cn(
                          "h-full rounded-full transition-[width] duration-300 motion-reduce:transition-none",
                          item.color ?? "bg-primary",
                        )}
                        style={{ width: `${width()}%` }}
                      />
                    </div>
                  </li>
                );
              }}
            </For>
          </ul>
        </Show>
      </Show>
    </div>
  );
}
