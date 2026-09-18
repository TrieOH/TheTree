import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createMemo,
  createSignal,
  onCleanup,
  onSettled,
} from "solid-js";

import { cn } from "./lib/cn";

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
   * Full dataset.
   *
   * Filtering, sorting and pagination happen inside this component.
   */
  items: T[];

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
  if (typeof left === "number" && typeof right === "number") return left - right;


  return String(left ?? "").localeCompare(
    String(right ?? ""),
    undefined,
    {
      numeric: true,
      sensitivity: "base",
    },
  );
}

function matchesSearch(value: unknown, search: string, seen = new WeakSet<object>()): boolean {
  if (value === null || value === undefined) return false;

  if (
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean" ||
    typeof value === "bigint"
  ) {
    return String(value)
      .toLocaleLowerCase()
      .includes(search);
  }

  if (value instanceof Date) {
    return value
      .toISOString()
      .toLocaleLowerCase()
      .includes(search);
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

function buildPageNumbers(
  current: number,
  total: number,
  maxVisible: number,
): Array<number | "…"> {
  const visible = Math.max(3, Math.min(maxVisible, total));
  if (total <= visible) return Array.from({ length: total }, (_, index) => index + 1);

  const middleCount = visible - 2;
  let start = current - Math.floor(middleCount / 2);
  let end = start + middleCount - 1;

  if (start < 2) {
    end += 2 - start;
    start = 2;
  }
  if (end > total - 1) {
    start -= end - (total - 1);
    end = total - 1;
  }

  start = Math.max(2, start);
  end = Math.min(total - 1, end);

  const pages: Array<number | "…"> = [1];
  if (start > 2) pages.push("…");
  for (let page = start; page <= end; page += 1) pages.push(page);
  if (end < total - 1) pages.push("…");
  pages.push(total);
  return pages;
}

/**
 * Toolbar + page + footer around a list, mirroring `@trieoh/ui-react`'s
 * PaginatedContainer: same sorting rules, same pagination windowing, same
 * footer wording, same pill/active colours. Icons are inline SVG — this package
 * carries no icon dependency on purpose.
 *
 * The page size is measured, not configured: it observes how many columns the
 * grid actually renders at the current width and multiplies that by `maxRows`.
 * Cells are equal because the tracks are equal, and a narrow viewport simply gets
 * fewer columns instead of a page that overflows the screen.
 */
export function PaginatedContainer<T>(props: PaginatedContainerProps<T>) {
  const [page, setPage] = createSignal(props.defaultPage ?? 1);
  const [internalSort, setInternalSort] = createSignal<SortState<T> | undefined>(
    props.sortFields?.length
      ? { field: props.sortFields[0]!.key, direction: "asc" }
      : undefined,
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
    const usable = tracks.length > 0 && tracks.every((track) => /^[\d.]+px$/.test(track));

    // Same count: leave the signal alone, so the page is not re-sliced.
    if (tracks.length !== columns() && usable) {
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
    const handlePointerDown = (
      event: PointerEvent,
    ) => {
      if (
        sortOpen() &&
        sortRoot &&
        !sortRoot.contains(
          event.target as Node,
        )
      ) {
        setSortOpen(false);
      }
    };

    const handleKeyDown = (
      event: KeyboardEvent,
    ) => {
      if (
        event.key === "Escape" &&
        sortOpen()
      ) {
        setSortOpen(false);
      }
    };

    document.addEventListener(
      "pointerdown",
      handlePointerDown,
    );

    document.addEventListener(
      "keydown",
      handleKeyDown,
    );

    const handleResize = () => {
      if (sortOpen()) {
        clampSortMenu();
      }
    };
    window.addEventListener("resize", handleResize);

    return () => {
      document.removeEventListener(
        "pointerdown",
        handlePointerDown,
      );

      document.removeEventListener(
        "keydown",
        handleKeyDown,
      );
      window.removeEventListener(
        "resize",
        handleResize,
      );

      resizeObserver?.disconnect();

      clearTimeout(settleTimer);
    };
  });

  const rowsFor = (
    columnCount: number,
  ) =>
    typeof props.maxRows === "function"
      ? props.maxRows(columnCount)
      : (props.maxRows ?? 2);

  const pageSize = () =>
    props.pageSize ??
    Math.max(
      1,
      columns() *
      rowsFor(columns()),
    );

  /*
   * Search all searchable data by default.
   */
  const filtered = createMemo(() => {
    const search = (
      props.filterValue ?? ""
    )
      .trim()
      .toLocaleLowerCase();

    if (!search) {
      return props.items;
    }

    return props.items.filter(
      (item) => {
        if (
          props.filterFields?.length
        ) {
          return props.filterFields.some(
            (key) =>
              matchesSearch(
                item[key],
                search,
              ),
          );
        }

        return matchesSearch(
          item,
          search,
        );
      },
    );
  });

  const sorted = createMemo(() => {
    const sort = activeSort();

    const fields =
      props.sortFields;

    if (
      !sort ||
      !fields?.length
    ) {
      return filtered();
    }

    const field = fields.find(
      (candidate) =>
        candidate.key ===
        sort.field,
    );

    const direction =
      sort.direction === "desc"
        ? -1
        : 1;

    return [...filtered()].sort(
      (a, b) => {
        const result =
          field?.comparator
            ? field.comparator(a, b)
            : defaultComparator(
              a,
              b,
              sort.field,
            );

        return result * direction;
      },
    );
  });

  const totalPages =
    createMemo(() =>
      Math.max(
        1,
        Math.ceil(
          sorted().length /
          pageSize(),
        ),
      ),
    );

  const currentPage =
    createMemo(() =>
      Math.min(
        Math.max(page(), 1),
        totalPages(),
      ),
    );

  const slice = createMemo(() => {
    const size = pageSize();

    const start =
      (currentPage() - 1) *
      size;

    return sorted().slice(
      start,
      start + size,
    );
  });

  const range = createMemo(() => {
    const count =
      sorted().length;

    if (count === 0) {
      return {
        start: 0,
        end: 0,
      };
    }

    return {
      start:
        (currentPage() - 1) *
        pageSize() +
        1,

      end: Math.min(
        currentPage() *
        pageSize(),
        count,
      ),
    };
  });

  const goTo = (
    next: number,
  ) => {
    setResized(false);

    setPage(
      Math.min(
        Math.max(next, 1),
        totalPages(),
      ),
    );
  };

  const changeSort = (
    next: SortState<T>,
  ) => {
    if (!props.sort) {
      setInternalSort(next);
    }

    props.onSortChange?.(next);

    setPage(1);
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
  const selectSortField = (
    field: SortField<T>,
  ) => {
    const current = activeSort();

    if (
      current?.field ===
      field.key
    ) {
      changeSort({
        field: field.key,

        direction:
          current.direction === "asc"
            ? "desc"
            : "asc",
      });

      return;
    }

    changeSort({
      field: field.key,
      direction: "asc",
    });
  };

  const directionLabel = (
    direction: SortDirection,
  ) => {
    const sort = activeSort();

    const field =
      props.sortFields?.find(
        (candidate) =>
          candidate.key ===
          sort?.field,
      );

    if (direction === "asc") {
      return (
        field?.ascLabel ??
        "Ascendente"
      );
    }

    return (
      field?.descLabel ??
      "Descendente"
    );
  };

  return (
    <div
      class={cn(
        "flex w-full flex-col gap-4",
        props.class,
      )}
    >
      <header class="flex w-full flex-wrap items-center gap-2">
        {/* SEARCH */}
        <Show when={props.onFilterChange}>
          <div class="relative min-w-56 flex-1">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
              class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
            >
              <circle
                cx="11"
                cy="11"
                r="8"
              />

              <path
                d="m21 21-4.3-4.3"
                stroke-linecap="round"
              />
            </svg>

            <input
              type="search"
              value={
                props.filterValue ??
                ""
              }
              placeholder={
                props.filterPlaceholder ??
                "Buscar…"
              }
              aria-label={
                props.filterPlaceholder ??
                "Buscar"
              }
              onInput={(event) => {
                setPage(1);
                setResized(false);

                props.onFilterChange?.(
                  event.currentTarget
                    .value,
                );
              }}
              class="h-9 w-full rounded-md border border-border bg-background py-1 pl-9 pr-3 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>
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
                  aria-expanded={
                    sortOpen()
                      ? "true"
                      : "false"
                  }
                  aria-haspopup="dialog"
                  onClick={() =>
                    setSortOpen(
                      (open) => !open,
                    )
                  }
                  class="inline-flex size-9 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  {/* FILTER / SORT ICON */}
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    aria-hidden="true"
                    class="size-4"
                  >
                    <path
                      d="M3 6h18"
                      stroke-linecap="round"
                    />

                    <path
                      d="M6 12h12"
                      stroke-linecap="round"
                    />

                    <path
                      d="M10 18h4"
                      stroke-linecap="round"
                    />
                  </svg>
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
                            activeSort()?.field ===
                            field.key;

                          return (
                            <button
                              type="button"
                              onClick={() =>
                                selectSortField(
                                  field,
                                )
                              }
                              class={cn(
                                "flex w-full items-center gap-2.5 rounded-md px-2.5 py-2 text-left transition-colors",
                                isActive()
                                  ? "bg-muted text-foreground"
                                  : "text-foreground hover:bg-muted/70",
                              )}
                            >
                              {/* FIELD ICON */}
                              <svg
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                aria-hidden="true"
                                class={cn(
                                  "size-4 shrink-0",
                                  isActive()
                                    ? "text-foreground"
                                    : "text-muted-foreground",
                                )}
                              >
                                <path
                                  d="M7 7h10M7 12h7M7 17h4"
                                  stroke-linecap="round"
                                />
                              </svg>

                              {/* LABEL + DESCRIPTION */}
                              <span class="min-w-0 flex-1">
                                <span
                                  class={cn(
                                    "block truncate text-sm",
                                    isActive() &&
                                    "font-medium",
                                  )}
                                >
                                  {field.label}
                                </span>

                                <Show when={isActive()}>
                                  <span class="mt-0.5 block truncate text-[11px] leading-tight text-muted-foreground">
                                    {activeSort()
                                      ?.direction ===
                                      "asc"
                                      ? (
                                        field.ascLabel ??
                                        "Ascendente"
                                      )
                                      : (
                                        field.descLabel ??
                                        "Descendente"
                                      )}
                                  </span>
                                </Show>
                              </span>

                              {/* CURRENT DIRECTION ARROW */}
                              <Show when={isActive()}>
                                <svg
                                  viewBox="0 0 24 24"
                                  fill="none"
                                  stroke="currentColor"
                                  stroke-width="2"
                                  aria-hidden="true"
                                  class="size-4 shrink-0 text-primary"
                                >
                                  <Show
                                    when={
                                      activeSort()
                                        ?.direction ===
                                      "asc"
                                    }
                                    fallback={
                                      <path
                                        d="M12 5v14M18 13l-6 6-6-6"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                      />
                                    }
                                  >
                                    <path
                                      d="M12 19V5M6 11l6-6 6 6"
                                      stroke-linecap="round"
                                      stroke-linejoin="round"
                                    />
                                  </Show>
                                </svg>
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
                          <For
                            each={
                              [
                                "asc",
                                "desc",
                              ] as const
                            }
                          >
                            {(direction) => {
                              const isActive = () =>
                                activeSort()
                                  ?.direction ===
                                direction;

                              return (
                                <button
                                  type="button"
                                  onClick={() => {
                                    const sort =
                                      activeSort();

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
                                  <svg
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    aria-hidden="true"
                                    class="size-3.5 shrink-0"
                                  >
                                    <Show
                                      when={
                                        direction ===
                                        "asc"
                                      }
                                      fallback={
                                        <path
                                          d="M12 5v14M18 13l-6 6-6-6"
                                          stroke-linecap="round"
                                          stroke-linejoin="round"
                                        />
                                      }
                                    >
                                      <path
                                        d="M12 19V5M6 11l6-6 6 6"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                      />
                                    </Show>
                                  </svg>

                                  <span class="truncate">
                                    {directionLabel(
                                      direction,
                                    )}
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

      <Show
        when={sorted().length > 0}
        fallback={props.emptyState}
      >
        <Show
          when={
            props.layout !==
            "custom"
          }
        >
          <div
            ref={observeGrid}
            class={cn(
              props.layout === "list"
                ? "flex w-full flex-col"
                : "grid w-full",

              GAP_CLASS[
              props.gap ?? "3"
              ],
            )}
            style={
              props.layout === "grid"
                ? {
                  "grid-template-columns": `repeat(auto-fill, minmax(${props.minItemWidth ??
                    "200px"
                    }, 1fr))`,
                }
                : undefined
            }
          >
            <Show when={measured()}>
              {props.renderItems(
                slice(),
                {
                  animate:
                    !resized(),
                },
              )}
            </Show>
          </div>
        </Show>

        <Show
          when={
            props.layout ===
            "custom"
          }
        >
          {props.renderItems(
            slice(),
            {
              animate:
                !resized(),
            },
          )}
        </Show>

        <footer class="flex flex-wrap items-center justify-between gap-3 border-t border-border pt-3">
          <p class="text-xs text-muted-foreground">
            Mostrando{" "}
            {range().start}–
            {range().end} de{" "}
            {sorted().length}{" "}
            {props.itemLabel ?? "itens"}
          </p>

          <div class="flex flex-wrap items-center gap-1">
            <PageButton
              label="Primeira página"
              disabled={
                currentPage() === 1
              }
              onClick={() =>
                goTo(1)
              }
            >
              <path
                d="m11 17-5-5 5-5M18 17l-5-5 5-5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </PageButton>

            <PageButton
              label="Página anterior"
              disabled={
                currentPage() === 1
              }
              onClick={() =>
                goTo(
                  currentPage() - 1,
                )
              }
            >
              <path
                d="m15 18-6-6 6-6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </PageButton>

            <div class="hidden items-center gap-1 sm:flex">
              <For
                each={buildPageNumbers(
                  currentPage(),
                  totalPages(),
                  7,
                )}
              >
                {(entry) => (
                  <Show
                    when={
                      typeof entry ===
                      "number"
                    }
                    fallback={
                      <span class="px-1 text-xs text-muted-foreground">
                        …
                      </span>
                    }
                  >
                    <button
                      type="button"
                      aria-current={
                        entry ===
                          currentPage()
                          ? "page"
                          : undefined
                      }
                      aria-label={`Página ${entry}`}
                      onClick={() =>
                        goTo(
                          Number(entry),
                        )
                      }
                      class={cn(
                        "flex h-7 min-w-7 items-center justify-center rounded-sm border px-1.5 text-xs font-medium transition-colors select-none",
                        entry ===
                          currentPage()
                          ? "border-primary/30 bg-primary/10 text-primary"
                          : "border-border text-muted-foreground hover:bg-muted hover:text-foreground",
                      )}
                    >
                      {entry}
                    </button>
                  </Show>
                )}
              </For>
            </div>

            <PageButton
              label="Próxima página"
              disabled={
                currentPage() ===
                totalPages()
              }
              onClick={() =>
                goTo(
                  currentPage() + 1,
                )
              }
            >
              <path
                d="m9 18 6-6-6-6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </PageButton>

            <PageButton
              label="Última página"
              disabled={
                currentPage() ===
                totalPages()
              }
              onClick={() =>
                goTo(
                  totalPages(),
                )
              }
            >
              <path
                d="m13 17 5-5-5-5M6 17l5-5-5-5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </PageButton>
          </div>
        </footer>
      </Show>
    </div>
  );
}

function PageButton(props: {
  label: string;
  disabled: boolean;
  onClick: () => void;
  children: JSX.Element;
}) {
  return (
    <button
      type="button"
      aria-label={props.label}
      disabled={props.disabled}
      onClick={() =>
        props.onClick()
      }
      class="flex size-7 items-center justify-center rounded-sm border border-border text-muted-foreground transition-colors hover:bg-muted hover:text-foreground disabled:cursor-not-allowed disabled:opacity-35"
    >
      <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        aria-hidden="true"
        class="size-3.5"
      >
        {props.children}
      </svg>
    </button>
  );
}