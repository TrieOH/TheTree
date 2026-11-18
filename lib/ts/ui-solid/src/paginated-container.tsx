import SortIcon from "~icons/lucide/list-filter";
import CheckIcon from "~icons/lucide/check";
import ArrowDownIcon from "~icons/lucide/arrow-down";
import ArrowUpIcon from "~icons/lucide/arrow-up";
import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createMemo,
  createSignal,
  onSettled,
  untrack,
} from "solid-js";

import { cn } from "./lib/cn";
import { Loading } from "./loading";
import { Pagination } from "./pagination";
import { Button } from "./button";
import { SearchInput } from "./search-input";

export type SortDirection = "asc" | "desc";

export interface SortField<T> {
  key: keyof T;
  label: string;
  comparator?: (a: T, b: T) => number;

  /**
   * Optional labels describing each direction.
   *
   * Examples:
   * ascLabel: "A → Z"
   * descLabel: "Z → A"
   *
   * ascLabel: "Mais antigos"
   * descLabel: "Mais recentes"
   */
  ascLabel?: string;
  descLabel?: string;
}

export interface SortState<T> {
  field: keyof T;
  direction: SortDirection;
}

export type LayoutMode = "grid" | "list" | "custom";
export type GapSize = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "8" | "10";

const GAP_CLASS: Record<GapSize, string> = {
  "0": "gap-0",
  "1": "gap-1",
  "2": "gap-2",
  "3": "gap-3",
  "4": "gap-4",
  "5": "gap-5",
  "6": "gap-6",
  "8": "gap-8",
  "10": "gap-10",
};

export interface PaginatedContainerProps<T> {
  /**
   * Full dataset in local mode; only the current page in server mode.
   *
   * Local mode filters, sorts and slices. With `pagination`, these operations
   * belong to the caller and items are rendered unchanged.
   */
  items: T[];
  /** Server mode: items already represent one filtered/sorted page. */
  pagination?: {
    page: number;
    pageSize: number;
    totalItems?: number;
    hasNextPage?: boolean;
    hasPreviousPage?: boolean;
    onPageChange: (page: number) => void;
    onPageSizeChange?: (size: number) => void;
  };
  onPageSizeChange?: (size: number) => void;
  /** Hide the optional selector even when a page-size handler is supplied. */
  showPageSizeSelector?: boolean;
  pageSizeOptions?: readonly number[];
  loading?: boolean;
  error?: string;
  onRetry?: () => void;

  defaultPage?: number;

  /** Sortable fields. */
  sortFields?: Array<SortField<T>>;

  /** Controlled sort, paired with `onSortChange`. */
  sort?: SortState<T>;

  onSortChange?: (sort: SortState<T>) => void;

  filterValue?: string;

  onFilterChange?: (value: string) => void;

  filterPlaceholder?: string;

  /**
   * Restricts which fields are searched.
   *
   * When omitted, every searchable field in the item is searched.
   */
  filterFields?: Array<keyof T>;

  /** Noun used in the footer counter, e.g. "eventos". */
  itemLabel?: string;

  /**
   * Renders one page.
   *
   * `animate` is false for rows that appeared because a resize
   * increased the available page size.
   */
  renderItems: (slice: T[], options: { animate: boolean }) => JSX.Element;
  emptyState?: JSX.Element;
  headerActions?: JSX.Element;
  class?: string;
  layout?: LayoutMode;
  gap?: GapSize;
  /** Grid minimum item width, e.g. "16rem". The column count follows the width. */
  minItemWidth?: string;
  /**
   * Rows per page. With the measured column count this decides how many items
   * fit on screen, so callers do not have to guess a page size.
   *
   * Pass a function when one row count cannot serve every width: a single
   * column needs more rows than a three-column grid to fill the same screen.
   */
  maxRows?: number | ((columns: number) => number);
  /** Hard override for the page size; ignores `maxRows` and the measurement. */
  pageSize?: number;
}

function defaultComparator<T>(a: T, b: T, key: keyof T): number {
  const left = a[key];
  const right = b[key];
  if (typeof left === "number" && typeof right === "number")
    return left - right;

  return String(left ?? "").localeCompare(String(right ?? ""), undefined, {
    numeric: true,
    sensitivity: "base",
  });
}

function matchesSearch(
  value: unknown,
  search: string,
  seen = new WeakSet<object>(),
): boolean {
  if (value === null || value === undefined) return false;

  if (
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean" ||
    typeof value === "bigint"
  ) {
    return String(value).toLocaleLowerCase().includes(search);
  }

  if (value instanceof Date) {
    return value.toISOString().toLocaleLowerCase().includes(search);
  }

  if (Array.isArray(value)) {
    return value.some((item) => matchesSearch(item, search, seen));
  }

  if (typeof value === "object") {
    if (seen.has(value)) return false;

    seen.add(value);

    return Object.values(value as Record<string, unknown>).some((item) =>
      matchesSearch(item, search, seen),
    );
  }

  return false;
}

/**
 * Toolbar + page + footer around a list, mirroring `@trieoh/ui-react`'s
 * PaginatedContainer: same sorting rules, same pagination windowing, same
 * footer wording, same pill/active colours. Icons use the shared Lucide collection.
 *
 * The page size is measured, not configured: it observes how many columns the
 * grid actually renders at the current width and multiplies that by `maxRows`.
 * Cells are equal because the tracks are equal, and a narrow viewport simply gets
 * fewer columns instead of a page that overflows the screen.
 */
export function PaginatedContainer<T>(props: PaginatedContainerProps<T>) {
  const [page, setPage] = createSignal(untrack(() => props.defaultPage ?? 1));
  const [internalSort, setInternalSort] = createSignal<
    SortState<T> | undefined
  >(
    untrack(() =>
      props.sortFields?.length
        ? { field: props.sortFields[0]!.key, direction: "asc" }
        : undefined,
    ),
  );
  const [sortOpen, setSortOpen] = createSignal(false);
  let sortMenuEl: HTMLElement | undefined;

  const clampSortMenu = (el?: HTMLElement | null) => {
    const target = el ?? sortMenuEl;
    if (!target || typeof window === "undefined") return;
    requestAnimationFrame(() => {
      target.style.transform = "";
      const rect = target.getBoundingClientRect();
      const vw = window.innerWidth;
      if (rect.left < 8) {
        target.style.transform = `translateX(${8 - rect.left}px)`;
      } else if (rect.right > vw - 8) {
        target.style.transform = `translateX(${vw - 8 - rect.right}px)`;
      }
    });
  };
  const [columns, setColumns] = createSignal(1);
  const [resized, setResized] = createSignal(false);
  const [measured, setMeasured] = createSignal(false);

  let grid: HTMLDivElement | undefined;
  let sortRoot: HTMLDivElement | undefined;
  let resizeObserver: ResizeObserver | undefined;
  let settleTimer: ReturnType<typeof setTimeout> | undefined;
  let measuredOnce = false;

  const activeSort = () => props.sort ?? internalSort();

  const measureColumns = () => {
    if (!grid) return;
    const tracks = getComputedStyle(grid)
      .gridTemplateColumns.split(" ")
      .filter((track) => track !== "" && track !== "none");

    // Only trust a resolved list: without layout (jsdom, an unsupported value)
    // this returns the specified `repeat(...)`, which is not a column count.
    const usable =
      tracks.length > 0 && tracks.every((track) => /^[\d.]+px$/.test(track));

    // Same count: leave the signal alone, so the page is not re-sliced.
    if (tracks.length !== untrack(columns) && usable) {
      setColumns(tracks.length);
      // A later change is the viewport moving, not the user navigating: the
      // rows it adds are not a page entrance.
      if (measuredOnce) setResized(true);
    }

    measuredOnce = true;
    setMeasured(true);
  };

  /**
   * Resizes here are usually an animation — the admin rail takes 300ms to
   * collapse — and every frame that changes the column count re-slices the page,
   * which reads as the items flickering. Wait for the width to settle instead.
   */
  const scheduleMeasure = () => {
    clearTimeout(settleTimer);
    settleTimer = setTimeout(measureColumns, 150);
  };

  /** Called by Solid with the element on mount, and with `undefined` on removal. */
  const observeGrid = (element: HTMLDivElement | undefined) => {
    grid = element;
    resizeObserver?.disconnect();
    resizeObserver = undefined;
    if (!element) return;

    measureColumns();
    if (typeof ResizeObserver === "undefined") return;

    resizeObserver = new ResizeObserver(scheduleMeasure);
    resizeObserver.observe(element);
  };

  onSettled(() => {
    const handlePointerDown = (event: PointerEvent) => {
      if (sortOpen() && sortRoot && !sortRoot.contains(event.target as Node)) {
        setSortOpen(false);
      }
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && sortOpen()) {
        setSortOpen(false);
      }
    };

    document.addEventListener("pointerdown", handlePointerDown);

    document.addEventListener("keydown", handleKeyDown);

    const handleResize = () => {
      if (sortOpen()) {
        clampSortMenu();
      }
    };
    window.addEventListener("resize", handleResize);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);

      document.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("resize", handleResize);

      resizeObserver?.disconnect();

      clearTimeout(settleTimer);
    };
  });

  const rowsFor = (columnCount: number) =>
    typeof props.maxRows === "function"
      ? props.maxRows(columnCount)
      : (props.maxRows ?? 2);

  const pageSize = () =>
    props.pagination?.pageSize ??
    props.pageSize ??
    Math.max(1, columns() * rowsFor(columns()));

  /*
   * Search all searchable data by default.
   */
  const filtered = createMemo(() => {
    if (props.pagination) return props.items;
    const search = (props.filterValue ?? "").trim().toLocaleLowerCase();

    if (!search) {
      return props.items;
    }

    return props.items.filter((item) => {
      if (props.filterFields?.length) {
        return props.filterFields.some((key) =>
          matchesSearch(item[key], search),
        );
      }

      return matchesSearch(item, search);
    });
  });

  const sorted = createMemo(() => {
    if (props.pagination) return props.items;
    const sort = activeSort();

    const fields = props.sortFields;

    if (!sort || !fields?.length) {
      return filtered();
    }

    const field = fields.find((candidate) => candidate.key === sort.field);

    const direction = sort.direction === "desc" ? -1 : 1;

    return [...filtered()].sort((a, b) => {
      const result = field?.comparator
        ? field.comparator(a, b)
        : defaultComparator(a, b, sort.field);

      return result * direction;
    });
  });

  const totalItems = () =>
    props.pagination ? props.pagination.totalItems : sorted().length;
  const totalPages = () =>
    totalItems() === undefined
      ? Infinity
      : Math.max(1, Math.ceil(totalItems()! / pageSize()));
  const currentPage = () =>
    Math.min(Math.max(props.pagination?.page ?? page(), 1), totalPages());

  const slice = createMemo(() => {
    if (props.pagination) return props.items;
    const size = pageSize();

    const start = (currentPage() - 1) * size;

    return sorted().slice(start, start + size);
  });

  const goTo = (next: number) => {
    setResized(false);
    const target = Math.min(Math.max(next, 1), totalPages());
    if (props.pagination) props.pagination.onPageChange(target);
    else setPage(target);
  };

  const changeSort = (next: SortState<T>) => {
    if (!props.sort) {
      setInternalSort(next);
    }

    props.onSortChange?.(next);

    goTo(1);
    setResized(false);
  };

  /**
   * First click on another field:
   * selects it as ascending.
   *
   * Clicking the already active field:
   * toggles asc ↔ desc.
   *
   * The direction buttons at the bottom still work independently.
   */
  const selectSortField = (field: SortField<T>) => {
    const current = activeSort();

    if (current?.field === field.key) {
      changeSort({
        field: field.key,

        direction: current.direction === "asc" ? "desc" : "asc",
      });

      return;
    }

    changeSort({
      field: field.key,
      direction: "asc",
    });
  };

  const directionLabel = (direction: SortDirection) => {
    const sort = activeSort();

    const field = props.sortFields?.find(
      (candidate) => candidate.key === sort?.field,
    );

    if (direction === "asc") {
      return field?.ascLabel ?? "Ascendente";
    }

    return field?.descLabel ?? "Descendente";
  };

  return (
    <div class={cn("flex w-full min-w-0 flex-col gap-4", props.class)}>
      <header class="flex w-full flex-wrap items-center gap-2">
        {/* SEARCH */}
        <Show when={props.onFilterChange}>
          <SearchInput
            class="min-w-0 flex-1 basis-40"
            value={props.filterValue ?? ""}
            placeholder={props.filterPlaceholder ?? "Buscar…"}
            aria-label={props.filterPlaceholder ?? "Buscar"}
            disabled={props.loading}
            onInput={(event) => {
              goTo(1);
              props.onFilterChange?.(event.currentTarget.value);
            }}
          />
        </Show>

        {/* HEADER ACTIONS & SORT */}
        <Show when={props.headerActions || props.sortFields?.length}>
          <div class="flex flex-wrap items-center gap-2 sm:ml-auto">
            {props.headerActions}

            {/* SORT */}
            <Show when={props.sortFields?.length}>
              <div
                ref={(element) => {
                  sortRoot = element;
                }}
                class="relative shrink-0"
              >
                <button
                  type="button"
                  aria-label="Ordenar"
                  title="Ordenar"
                  aria-expanded={sortOpen() ? "true" : "false"}
                  aria-haspopup="dialog"
                  onClick={() => setSortOpen((open) => !open)}
                  class="inline-flex size-9 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  {/* FILTER / SORT ICON */}
                  <SortIcon class="size-4" aria-hidden="true" />
                </button>

                <Show when={sortOpen()}>
                  <div
                    ref={(el) => {
                      sortMenuEl = el;
                      clampSortMenu(el);
                    }}
                    role="dialog"
                    aria-label="Opções de ordenação"
                    class="
                  absolute
                  right-0
                  top-[calc(100%+6px)]
                  z-50
                  w-64
                  max-w-[calc(100vw-1rem)]
                  overflow-hidden
                  rounded-lg
                  border
                  border-border
                  bg-popover
                  text-popover-foreground
                  shadow-lg
                  shadow-black/5
                "
                  >
                    {/* HEADER */}
                    <div class="border-b border-border px-3 py-2">
                      <p class="text-xs font-medium text-muted-foreground">
                        Ordenar por
                      </p>
                    </div>

                    {/* SORT FIELDS */}
                    <div class="p-1.5">
                      <For each={props.sortFields}>
                        {(field) => {
                          const isActive = () =>
                            activeSort()?.field === field.key;

                          return (
                            <button
                              type="button"
                              onClick={() => selectSortField(field)}
                              class={cn(
                                "flex w-full items-center gap-2.5 rounded-md px-2.5 py-2 text-left transition-colors",
                                isActive()
                                  ? "bg-muted text-foreground"
                                  : "text-foreground hover:bg-muted/70",
                              )}
                            >
                              {/* FIELD ICON */}
                              <CheckIcon
                                class={cn(
                                  "size-4 shrink-0",
                                  isActive()
                                    ? "text-foreground"
                                    : "text-muted-foreground",
                                )}
                                aria-hidden="true"
                              />

                              {/* LABEL + DESCRIPTION */}
                              <span class="min-w-0 flex-1">
                                <span
                                  class={cn(
                                    "block truncate text-sm",
                                    isActive() && "font-medium",
                                  )}
                                >
                                  {field.label}
                                </span>

                                <Show when={isActive()}>
                                  <span class="mt-0.5 block truncate text-[11px] leading-tight text-muted-foreground">
                                    {activeSort()?.direction === "asc"
                                      ? (field.ascLabel ?? "Ascendente")
                                      : (field.descLabel ?? "Descendente")}
                                  </span>
                                </Show>
                              </span>

                              {/* CURRENT DIRECTION ARROW */}
                              <Show when={isActive()}>
                                <Show
                                  when={activeSort()?.direction === "asc"}
                                  fallback={
                                    <ArrowDownIcon
                                      class="size-4 shrink-0 text-primary"
                                      aria-hidden="true"
                                    />
                                  }
                                >
                                  <ArrowUpIcon
                                    class="size-4 shrink-0 text-primary"
                                    aria-hidden="true"
                                  />
                                </Show>
                              </Show>
                            </button>
                          );
                        }}
                      </For>
                    </div>

                    {/* DIRECTION */}
                    <Show when={activeSort()}>
                      <div class="border-t border-border p-2">
                        <div class="grid grid-cols-2 gap-1.5">
                          <For each={["asc", "desc"] as const}>
                            {(direction) => {
                              const isActive = () =>
                                activeSort()?.direction === direction;

                              return (
                                <button
                                  type="button"
                                  onClick={() => {
                                    const sort = activeSort();

                                    if (!sort) {
                                      return;
                                    }

                                    changeSort({
                                      ...sort,
                                      direction,
                                    });
                                  }}
                                  class={cn(
                                    "flex min-w-0 items-center justify-center gap-1.5 rounded-md border px-2 py-1.5 text-xs font-medium transition-colors",
                                    isActive()
                                      ? "border-primary/30 bg-primary/10 text-primary"
                                      : "border-border text-muted-foreground hover:bg-muted hover:text-foreground",
                                  )}
                                >
                                  <Show
                                    when={direction === "asc"}
                                    fallback={
                                      <ArrowDownIcon
                                        class="size-3.5 shrink-0"
                                        aria-hidden="true"
                                      />
                                    }
                                  >
                                    <ArrowUpIcon
                                      class="size-3.5 shrink-0"
                                      aria-hidden="true"
                                    />
                                  </Show>

                                  <span class="truncate">
                                    {directionLabel(direction)}
                                  </span>
                                </button>
                              );
                            }}
                          </For>
                        </div>
                      </div>
                    </Show>
                  </div>
                </Show>
              </div>
            </Show>
          </div>
        </Show>
      </header>

      <Show when={props.loading && sorted().length === 0}>
        <Loading size="sm" tone="muted" class="py-6" label="Carregando…" />
      </Show>
      <Show when={props.error}>
        <div role="alert" class="py-6 text-center text-sm">
          <p>{props.error}</p>
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
        </div>
      </Show>
      <Show
        when={!props.error && (sorted().length > 0 || !!props.pagination)}
        fallback={props.loading || props.error ? undefined : props.emptyState}
      >
        <Show when={props.layout !== "custom"}>
          <div
            ref={observeGrid}
            class={cn(
              props.layout === "list" ? "flex w-full flex-col" : "grid w-full",

              GAP_CLASS[props.gap ?? "3"],
            )}
            style={
              props.layout === "grid"
                ? {
                    "grid-template-columns": `repeat(auto-fill, minmax(min(100%, ${
                      props.minItemWidth ?? "200px"
                    }), 1fr))`,
                  }
                : undefined
            }
          >
            <Show when={measured()}>
              {props.renderItems(slice(), {
                animate: !resized(),
              })}
            </Show>
          </div>
        </Show>

        <Show when={props.layout === "custom"}>
          {props.renderItems(slice(), {
            animate: !resized(),
          })}
        </Show>

        <Pagination
          class="border-t border-border pt-3"
          page={currentPage()}
          pageSize={pageSize()}
          totalItems={totalItems()}
          itemLabel={props.itemLabel}
          disabled={props.loading || !!props.error}
          hasNextPage={props.pagination?.hasNextPage}
          hasPreviousPage={props.pagination?.hasPreviousPage}
          onPageChange={goTo}
          showPageSizeSelector={props.showPageSizeSelector}
          pageSizeOptions={props.pageSizeOptions}
          onPageSizeChange={
            props.pagination?.onPageSizeChange ??
            (props.onPageSizeChange
              ? (size) => {
                  props.onPageSizeChange?.(size);
                  goTo(1);
                }
              : undefined)
          }
        />
      </Show>
    </div>
  );
}
