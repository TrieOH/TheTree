import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";

export interface DashboardBarListItem {
  id: string;
  label: string;
  value: number;
  detail?: string;
  color?: string;
}

interface DashboardBarListProps {
  items: DashboardBarListItem[];
  maxValue?: number;
  total?: number;
  showPercentage?: boolean;
  emptyMessage?: string;
}

export function DashboardBarList(props: DashboardBarListProps): JSX.Element {
  const maximum = () =>
    props.maxValue ?? Math.max(...props.items.map((item) => item.value), 1);

  return (
    <Show
      when={props.items.length > 0}
      fallback={
        <p class="text-sm text-muted-foreground">
          {props.emptyMessage ?? "Nenhum dado disponível."}
        </p>
      }
    >
      <div class="space-y-4">
        <For each={props.items}>
          {(item) => {
            const percentage = () =>
              props.total ? (item.value / props.total) * 100 : 0;
            const width = () =>
              item.value ? Math.max((item.value / maximum()) * 100, 3) : 0;

            return (
              <div class="space-y-1.5">
                <div class="flex items-center justify-between gap-3 text-xs">
                  <span class="min-w-0 truncate font-medium text-foreground">
                    {item.label}
                  </span>
                  <span class="shrink-0 text-muted-foreground">
                    {item.detail ?? item.value}
                    {props.showPercentage
                      ? ` · ${percentage().toFixed(0)}%`
                      : ""}
                  </span>
                </div>
                <div class="h-2 overflow-hidden rounded-full bg-muted">
                  <div
                    class={`h-full rounded-full transition-all duration-500 ${
                      item.color ?? "bg-primary"
                    }`}
                    style={{ width: `${width()}%` }}
                  />
                </div>
              </div>
            );
          }}
        </For>
      </div>
    </Show>
  );
}
