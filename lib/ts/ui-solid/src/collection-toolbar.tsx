import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { Button } from "./button";
import { Combobox } from "./combobox";
import { Pagination, type PaginationProps } from "./pagination";
import { FilterBar, type FilterBarProps } from "./filter-bar";
import { SearchInput } from "./search-input";
import ArrowUpDown from "~icons/lucide/arrow-up-down";
import List from "~icons/lucide/list";
import LayoutGrid from "~icons/lucide/layout-grid";
import { cn } from "./lib/cn";

export interface CollectionToolbarProps {
  search?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;
  totalItems?: number;
  itemLabel?: string;
  sortOptions?: readonly { value: string; label: string }[];
  sort?: string;
  onSortChange?: (value: string) => void;
  layout?: "list" | "grid";
  onLayoutChange?: (layout: "list" | "grid") => void;
  /** Filter controls composed from Combobox, date inputs, etc. */
  filters?: JSX.Element;
  actions?: JSX.Element;
  /** Optional active filter chips, composed with FilterBar. */
  filterBar?: FilterBarProps;
  /** Controlled navigation; the caller still supplies the displayed items. */
  pagination?: PaginationProps;
  disabled?: boolean;
  class?: string;
}

/** Controlled toolbar shared by paged and continuously loaded collections. */
export function CollectionToolbar(props: CollectionToolbarProps) {
  return (
    <div class={cn("w-full min-w-0 space-y-3", props.class)}>
      <div
        role="group"
        aria-label="Controles da coleção"
        class={cn("flex w-full min-w-0 flex-wrap items-center gap-2")}
      >
        <Show when={props.onSearchChange}>
          <SearchInput
            class="min-w-0 flex-1 basis-40"
            aria-label={props.searchPlaceholder ?? "Buscar"}
            placeholder={props.searchPlaceholder ?? "Buscar…"}
            value={props.search ?? ""}
            disabled={props.disabled}
            onInput={(event) =>
              props.onSearchChange?.(event.currentTarget.value)
            }
          />
        </Show>
        {props.filters}
        <Show when={props.totalItems !== undefined}>
          <span class="text-xs text-muted-foreground" aria-live="polite">
            {props.totalItems} {props.itemLabel ?? "itens"}
          </span>
        </Show>
        <Show when={props.sortOptions?.length && props.onSortChange}>
          <Combobox
            class="w-8 shrink-0"
            shape="soft"
            iconOnly
            icon={<ArrowUpDown class="size-4" />}
            triggerClass="border-input shadow-xs"
            size="sm"
            searchable={false}
            aria-label="Ordenação"
            placeholder="Ordenar por"
            options={props.sortOptions ?? []}
            value={props.sort}
            disabled={props.disabled}
            onChange={(value) => props.onSortChange?.(value)}
          />
        </Show>
        <Show when={props.onLayoutChange}>
          <div role="group" aria-label="Visualização" class="flex gap-1">
            <For each={["list", "grid"] as const}>
              {(layout) => (
                <Button
                  size="sm"
                  variant="ghost"
                  class={cn(
                    "size-8 px-0 text-muted-foreground",
                    props.layout === layout && "bg-muted text-foreground",
                  )}
                  shape="soft"
                  aria-label={
                    layout === "list"
                      ? "Visualização em lista"
                      : "Visualização em grid"
                  }
                  aria-pressed={props.layout === layout ? "true" : "false"}
                  disabled={props.disabled}
                  onClick={() => props.onLayoutChange?.(layout)}
                >
                  {layout === "list" ? (
                    <List class="size-4" />
                  ) : (
                    <LayoutGrid class="size-4" />
                  )}
                </Button>
              )}
            </For>
          </div>
        </Show>
        {props.actions}
      </div>
      <Show when={props.filterBar}>
        {(filters) => (
          <FilterBar
            {...filters()}
            disabled={props.disabled || filters().disabled}
          />
        )}
      </Show>
      <Show when={props.pagination}>
        {(pagination) => (
          <Pagination
            {...pagination()}
            disabled={props.disabled || pagination().disabled}
          />
        )}
      </Show>
    </div>
  );
}
