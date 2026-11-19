import ArrowUp from "~icons/lucide/arrow-up";
import ArrowDown from "~icons/lucide/arrow-down";
import ArrowUpDown from "~icons/lucide/arrow-up-down";
import { BulkActions, type BulkAction } from "./bulk-actions";
import { Loading } from "./loading";
import type { JSX } from "@solidjs/web";
import { For, Show, createMemo } from "solid-js";
import { Button } from "./button";
import { Checkbox } from "./checkbox";
import { Pagination, type PaginationProps } from "./pagination";
import { cn } from "./lib/cn";

export interface DataTableColumn<T> {
  id: string;
  header: string;
  cell: (item: T) => JSX.Element;
  sortable?: boolean;
  class?: string;
  /** Fixed column width, e.g. "12rem" or "40%". */
  width?: string;
}
export interface DataTableSort {
  column: string;
  direction: "asc" | "desc";
}
export interface DataTableProps<T> {
  items: readonly T[];
  columns: readonly DataTableColumn<T>[];
  getRowKey: (item: T) => string | number;
  /** Accessible description of the table. */
  caption: string;
  /** Selection persists across pages; select-all affects only visible rows. */
  selectedKeys?: readonly (string | number)[];
  onSelectionChange?: (keys: (string | number)[]) => void;
  sort?: DataTableSort;
  /** Caller applies sorting, locally or on the server. */
  onSortChange?: (sort: DataTableSort) => void;
  rowActions?: (item: T) => JSX.Element;
  pagination?: PaginationProps;
  /** Persistent controls; remain available while rows are selected. */
  toolbar?: JSX.Element;
  bulkActions?: readonly BulkAction[];
  bulkActionsIconOnly?: boolean;
  loading?: boolean;
  error?: string;
  onRetry?: () => void;
  emptyState?: JSX.Element;
  /** Minimum table width before scrolling internally on narrow screens. */
  minWidth?: string;
  class?: string;
}
export function DataTable<T>(props: DataTableProps<T>) {
  const selected = createMemo(() => new Set(props.selectedKeys ?? []));
  const selectable = () => !!props.onSelectionChange;
  const allSelected = () =>
    props.items.length > 0 &&
    props.items.every((item) => selected().has(props.getRowKey(item)));
  const someSelected = () =>
    props.items.some((item) => selected().has(props.getRowKey(item)));
  const selectAll = (checked: boolean) => {
    const next = new Set(selected());
    for (const item of props.items) {
      const key = props.getRowKey(item);
      if (checked) next.add(key);
      else next.delete(key);
    }
    props.onSelectionChange?.([...next]);
  };
  const selectRow = (item: T, checked: boolean) => {
    const next = new Set(selected());
    const key = props.getRowKey(item);
    if (checked) next.add(key);
    else next.delete(key);
    props.onSelectionChange?.([...next]);
  };
  const sortBy = (id: string) =>
    props.onSortChange?.({
      column: id,
      direction:
        props.sort?.column === id && props.sort.direction === "asc"
          ? "desc"
          : "asc",
    });
  const span = () =>
    props.columns.length + (selectable() ? 1 : 0) + (props.rowActions ? 1 : 0);
  return (
    <div class={cn("w-full min-w-0 space-y-3", props.class)}>
      <Show when={props.toolbar || (selectable() && props.bulkActions?.length)}>
        <div class="min-w-0 space-y-2">
          {props.toolbar}
          <Show when={selectable() && props.bulkActions?.length}>
            <div class="grid min-w-0 grid-cols-[minmax(0,1fr)] items-center">
              <p
                class={cn(
                  "col-start-1 row-start-1 text-xs text-muted-foreground tabular-nums",
                  selected().size > 0 && "invisible",
                )}
                aria-hidden={selected().size > 0 ? "true" : undefined}
              >
                {props.pagination?.totalItems ?? props.items.length}{" "}
                {props.pagination?.itemLabel ?? "itens"}
                {props.pagination?.totalItems === undefined
                  ? " nesta página"
                  : ""}
              </p>
              <BulkActions
                reserveSpace
                class="col-start-1 row-start-1"
                selectedCount={selected().size}
                actions={props.bulkActions ?? []}
                iconOnly={props.bulkActionsIconOnly}
                busy={props.loading}
                onClearSelection={() => props.onSelectionChange?.([])}
              />
            </div>
          </Show>
        </div>
      </Show>
      <div
        class="w-full min-w-0 overflow-x-auto rounded-md border border-border"
        tabindex={0}
        role="region"
        aria-label={props.caption}
        aria-busy={props.loading ? "true" : "false"}
      >
        <table
          class="w-full table-fixed border-collapse text-left text-sm"
          style={{
            "min-width":
              props.minWidth ??
              `${props.columns.length * 8 + (selectable() ? 3 : 0) + (props.rowActions ? 7 : 0)}rem`,
          }}
        >
          <colgroup>
            <Show when={selectable()}>
              <col style={{ width: "3rem" }} />
            </Show>
            <For each={props.columns}>
              {(column) => <col style={{ width: column.width }} />}
            </For>
            <Show when={props.rowActions}>
              <col style={{ width: "7rem" }} />
            </Show>
          </colgroup>
          <caption class="sr-only">{props.caption}</caption>
          <thead class="border-b border-border bg-muted/40 text-xs text-muted-foreground">
            <tr>
              <Show when={selectable()}>
                <th scope="col" class="w-12 p-3">
                  <Checkbox
                    size="sm"
                    variant="plain"
                    label={<span class="sr-only">Selecionar página</span>}
                    checked={allSelected()}
                    indeterminate={someSelected() && !allSelected()}
                    disabled={props.loading || !props.items.length}
                    onChange={selectAll}
                  />
                </th>
              </Show>
              <For each={props.columns}>
                {(column) => (
                  <th
                    scope="col"
                    class={cn("p-3 font-medium", column.class)}
                    aria-sort={
                      props.sort?.column === column.id
                        ? props.sort.direction === "asc"
                          ? "ascending"
                          : "descending"
                        : undefined
                    }
                  >
                    <Show
                      when={column.sortable && props.onSortChange}
                      fallback={column.header}
                    >
                      <Button
                        size="sm"
                        variant="ghost"
                        class="-ml-3"
                        disabled={props.loading}
                        onClick={() => sortBy(column.id)}
                      >
                        {column.header}
                        <Show
                          when={props.sort?.column === column.id}
                          fallback={
                            <ArrowUpDown
                              class="size-3.5 shrink-0"
                              aria-hidden="true"
                            />
                          }
                        >
                          <Show
                            when={props.sort?.direction === "asc"}
                            fallback={
                              <ArrowDown
                                class="size-3.5 shrink-0"
                                aria-hidden="true"
                              />
                            }
                          >
                            <ArrowUp
                              class="size-3.5 shrink-0"
                              aria-hidden="true"
                            />
                          </Show>
                        </Show>
                      </Button>
                    </Show>
                  </th>
                )}
              </For>
              <Show when={props.rowActions}>
                <th scope="col" class="p-3 text-right font-medium">
                  Ações
                </th>
              </Show>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <Show
              when={!props.error && props.items.length > 0}
              fallback={
                <tr>
                  <td
                    colspan={span()}
                    class="p-8 text-center text-muted-foreground"
                  >
                    <Show
                      when={props.error}
                      fallback={
                        <Show
                          when={props.loading}
                          fallback={
                            props.emptyState ?? "Nenhum resultado encontrado."
                          }
                        >
                          <Loading size="sm" tone="muted" label="Carregando…" />
                        </Show>
                      }
                    >
                      <p role="alert">{props.error}</p>
                      <Show when={props.onRetry}>
                        <Button
                          size="sm"
                          variant="outline"
                          class="mt-3"
                          onClick={() => props.onRetry?.()}
                        >
                          Tentar novamente
                        </Button>
                      </Show>
                    </Show>
                  </td>
                </tr>
              }
            >
              <For each={props.items}>
                {(item) => (
                  <tr
                    class="hover:bg-muted/30"
                    aria-selected={
                      selectable()
                        ? selected().has(props.getRowKey(item))
                          ? "true"
                          : "false"
                        : undefined
                    }
                  >
                    <Show when={selectable()}>
                      <td class="p-3">
                        <Checkbox
                          size="sm"
                          variant="plain"
                          label={
                            <span class="sr-only">
                              Selecionar {String(props.getRowKey(item))}
                            </span>
                          }
                          checked={selected().has(props.getRowKey(item))}
                          disabled={props.loading}
                          onChange={(checked) => selectRow(item, checked)}
                        />
                      </td>
                    </Show>
                    <For each={props.columns}>
                      {(column) => (
                        <td class={cn("p-3 wrap-anywhere", column.class)}>
                          {column.cell(item)}
                        </td>
                      )}
                    </For>
                    <Show when={props.rowActions}>
                      <td class="p-3 text-right">{props.rowActions?.(item)}</td>
                    </Show>
                  </tr>
                )}
              </For>
            </Show>
          </tbody>
        </table>
      </div>
      <Show when={props.loading && props.items.length > 0}>
        <Loading size="sm" tone="muted" label="Atualizando…" />
      </Show>
      <Show when={props.pagination}>
        {(pagination) => (
          <Pagination
            {...pagination()}
            disabled={props.loading || pagination().disabled}
          />
        )}
      </Show>
    </div>
  );
}
