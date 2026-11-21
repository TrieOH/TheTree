import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "./lib/cn";

export interface TableColumn<T> {
  id: string;
  header: JSX.Element;
  cell: (item: T) => JSX.Element;
  align?: "left" | "right";
  /** Classes for this column's cells; a function gets the row. */
  class?: string | ((item: T) => string | false | undefined);
  /** Hover text for this column's cells. */
  title?: (item: T) => string | undefined;
}

export interface TableProps<T> {
  items: readonly T[];
  columns: readonly TableColumn<T>[];
  /** Read by screen readers; not shown. */
  caption: string;
  /** Shown instead of the table when there are no items. */
  emptyState?: JSX.Element;
  density?: "compact" | "default";
  class?: string;
}

/**
 * A plain read-only table: header, rows, nothing else. For sorting, selection
 * or pagination use DataTable.
 */
export function Table<T>(props: TableProps<T>) {
  const cellY = () => (props.density === "compact" ? "py-2" : "py-2.5");
  const edge = (index: number) =>
    index === props.columns.length - 1 ? undefined : "pr-4";
  const align = (column: TableColumn<T>) =>
    column.align === "right" ? "text-right" : undefined;
  return (
    <Show
      when={props.items.length > 0 || !props.emptyState}
      fallback={props.emptyState}
    >
      <div
        class={cn("w-full min-w-0 overflow-x-auto text-[13px]", props.class)}
      >
        <table class="w-full text-left">
          <caption class="sr-only">{props.caption}</caption>
          <thead class="text-[12px] text-muted-foreground">
            <tr class="border-b border-border">
              <For each={props.columns}>
                {(column, index) => (
                  <th
                    scope="col"
                    class={cn("py-2 font-normal", edge(index()), align(column))}
                  >
                    {column.header}
                  </th>
                )}
              </For>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <For each={props.items}>
              {(item) => (
                <tr>
                  <For each={props.columns}>
                    {(column, index) => (
                      <td
                        title={column.title?.(item)}
                        class={cn(
                          cellY(),
                          edge(index()),
                          align(column),
                          typeof column.class === "function"
                            ? column.class(item)
                            : column.class,
                        )}
                      >
                        {column.cell(item)}
                      </td>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </tbody>
        </table>
      </div>
    </Show>
  );
}
