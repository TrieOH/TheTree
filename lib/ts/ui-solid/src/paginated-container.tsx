import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal, onCleanup } from "solid-js";
import { cn } from "./lib/cn";

export type SortDirection = "asc" | "desc";

export interface SortField<T> {
  key: keyof T;
  label: string;
  comparator?: (a: T, b: T) => number;
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
  /** Full dataset — filter it yourself; this sorts and paginates. */
  items: T[];
  defaultPage?: number;
  /** Sortable fields. When present this component sorts; callers only filter. */
  sortFields?: Array<SortField<T>>;
  /** Controlled sort, paired with `onSortChange`. */
  sort?: SortState<T>;
  onSortChange?: (sort: SortState<T>) => void;
  filterValue?: string;
  onFilterChange?: (value: string) => void;
  filterPlaceholder?: string;
  /** Noun used in the footer counter, e.g. "eventos". */
  itemLabel?: string;
  /**
   * Renders one page. `animate` is false for the rows that appear because a
   * resize grew the page — the list re-sliced, nothing was navigated to.
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
   */
  maxRows?: number;
  /** Hard override for the page size; ignores `maxRows` and the measurement. */
  pageSize?: number;
}

function defaultComparator<T>(a: T, b: T, key: keyof T): number {
  const left = a[key];
  const right = b[key];
  if (typeof left === "number" && typeof right === "number") return left - right;
  return String(left).localeCompare(String(right));
}

/** Same windowing algorithm as the React container, ellipsis included. */
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
  const [columns, setColumns] = createSignal(1);
  const [resized, setResized] = createSignal(false);
  const [measured, setMeasured] = createSignal(false);

  let grid: HTMLDivElement | undefined;
  let resizeObserver: ResizeObserver | undefined;
  let settleTimer: ReturnType<typeof setTimeout> | undefined;
  let measuredOnce = false;

  const activeSort = () => props.sort ?? internalSort();

  /** Reads the resolved tracks off the grid itself: unit-agnostic, no parsing. */
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

    // The first measurement is immediate: waiting would show the wrong page size.
    measureColumns();
    if (typeof ResizeObserver === "undefined") return;

    resizeObserver = new ResizeObserver(scheduleMeasure);
    resizeObserver.observe(element);
  };

  onCleanup(() => {
    resizeObserver?.disconnect();
    clearTimeout(settleTimer);
  });

  const pageSize = () =>
    props.pageSize ?? Math.max(1, columns() * (props.maxRows ?? 2));

  const sorted = createMemo(() => {
    const sort = activeSort();
    const fields = props.sortFields;
    if (!sort || !fields?.length) return props.items;

    const field = fields.find((candidate) => candidate.key === sort.field);
    const direction = sort.direction === "desc" ? -1 : 1;

    return [...props.items].sort((a, b) => {
      const result = field?.comparator
        ? field.comparator(a, b)
        : defaultComparator(a, b, sort.field);
      return result * direction;
    });
  });

  const totalPages = createMemo(() => Math.max(1, Math.ceil(sorted().length / pageSize())));
  const currentPage = createMemo(() => Math.min(Math.max(page(), 1), totalPages()));
  const slice = createMemo(() => {
    const size = pageSize();
    const start = (currentPage() - 1) * size;
    return sorted().slice(start, start + size);
  });
  const range = createMemo(() => {
    const count = sorted().length;
    if (count === 0) return { start: 0, end: 0 };
    return {
      start: (currentPage() - 1) * pageSize() + 1,
      end: Math.min(currentPage() * pageSize(), count),
    };
  });

  const goTo = (next: number) => setPage(Math.min(Math.max(next, 1), totalPages()));
  const changeSort = (next: SortState<T>) => {
    if (!props.sort) setInternalSort(next);
    props.onSortChange?.(next);
  };
  const activeSortLabel = () => {
    const sort = activeSort();
    if (!sort) return null;
    return props.sortFields?.find((field) => field.key === sort.field)?.label ?? null;
  };

  return (
    <div class={cn("flex w-full flex-col gap-4", props.class)}>
      <header class="flex flex-wrap items-center gap-2">
        <Show when={props.onFilterChange}>
          <div class="relative min-w-56 flex-1 sm:max-w-xs">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
              class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
            >
              <circle cx="11" cy="11" r="8" />
              <path d="m21 21-4.3-4.3" stroke-linecap="round" />
            </svg>
            <input
              type="search"
              value={props.filterValue ?? ""}
              placeholder={props.filterPlaceholder ?? "Filtrar…"}
              aria-label={props.filterPlaceholder ?? "Filtrar"}
              onInput={(event) => props.onFilterChange?.(event.currentTarget.value)}
              class="h-9 w-full rounded-md border border-border bg-background py-1 pl-9 pr-3 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>
        </Show>

        <Show when={props.sortFields?.length}>
          <div class="relative">
            <button
              type="button"
              aria-expanded={sortOpen() ? "true" : "false"}
              aria-haspopup="dialog"
              onClick={() => setSortOpen((open) => !open)}
              class={cn(
                "inline-flex h-9 items-center gap-1.5 rounded-md border border-border px-3 text-sm transition-colors hover:bg-muted hover:text-foreground",
                activeSortLabel() ? "text-foreground" : "text-muted-foreground",
              )}
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                aria-hidden="true"
                class="size-4"
              >
                <path
                  d="m3 16 4 4 4-4M7 20V4M21 8l-4-4-4 4M17 4v16"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              <Show when={activeSortLabel()} fallback={<>Ordenar</>}>
                {(label) => (
                  <>
                    {label()}
                    <span aria-hidden="true">
                      {activeSort()?.direction === "desc" ? "↓" : "↑"}
                    </span>
                  </>
                )}
              </Show>
            </button>

            <Show when={sortOpen()}>
              <div
                role="dialog"
                aria-label="Opções de ordenação"
                class="absolute right-0 top-[calc(100%+6px)] z-50 w-52 rounded-md border border-border bg-card p-2 shadow-lg shadow-black/5"
              >
                <div class="mb-1 flex flex-col gap-0.5">
                  <For each={props.sortFields}>
                    {(field) => {
                      const isActive = () => activeSort()?.field === field.key;
                      return (
                        <button
                          type="button"
                          onClick={() =>
                            changeSort({
                              field: field.key,
                              direction:
                                isActive() && activeSort()?.direction === "asc" ? "desc" : "asc",
                            })}
                          class={cn(
                            "flex w-full cursor-pointer items-center justify-between rounded-md border px-2.5 py-2 text-left text-sm transition-colors select-none",
                            isActive()
                              ? "border-primary/20 bg-primary/10 font-medium text-primary"
                              : "border-transparent text-foreground hover:bg-muted",
                          )}
                        >
                          <span>{field.label}</span>
                          <span class={isActive() ? "text-primary" : "text-muted-foreground/40"}>
                            {isActive() ? (activeSort()?.direction === "asc" ? "↑" : "↓") : "—"}
                          </span>
                        </button>
                      );
                    }}
                  </For>
                </div>

                <div class="grid grid-cols-2 gap-1.5 border-t border-border pt-2">
                  <For each={["asc", "desc"] as const}>
                    {(direction) => (
                      <button
                        type="button"
                        onClick={() => {
                          const sort = activeSort();
                          if (sort) changeSort({ ...sort, direction });
                        }}
                        class={cn(
                          "flex items-center justify-center gap-1.5 rounded-md border py-1.5 text-xs font-medium transition-colors",
                          activeSort()?.direction === direction
                            ? "border-primary/30 bg-primary/10 text-primary"
                            : "border-border text-muted-foreground hover:bg-muted hover:text-foreground",
                        )}
                      >
                        {direction === "asc" ? "Asc" : "Desc"}
                      </button>
                    )}
                  </For>
                </div>
              </div>
            </Show>
          </div>
        </Show>

        <div class="ml-auto flex items-center gap-2">{props.headerActions}</div>
      </header>

      <Show when={sorted().length > 0} fallback={props.emptyState}>
        <Show when={props.layout !== "custom"}>
          <div
            ref={observeGrid}
            class={cn(
              props.layout === "list"
                ? "flex w-full flex-col"
                : "grid w-full",
              GAP_CLASS[props.gap ?? "3"],
            )}
            style={
              props.layout === "grid"
                ? {
                  "grid-template-columns": `repeat(auto-fill, minmax(${props.minItemWidth ?? "200px"}, 1fr))`,
                }
                : undefined
            }
          >
            {/*
              Items wait for the first measurement: rendering them at the
              default page size and correcting one frame later makes every card
              mount twice, and an entrance animation plays for nothing.
            */}
            <Show when={measured()}>
              {props.renderItems(slice(), { animate: !resized() })}
            </Show>
          </div>
        </Show>

        <Show when={props.layout === "custom"}>
          {props.renderItems(slice(), { animate: !resized() })}
        </Show>

        <footer class="flex flex-wrap items-center justify-between gap-3 border-t border-border pt-3">
          <p class="text-xs text-muted-foreground">
            Mostrando {range().start}–{range().end} de {sorted().length}{" "}
            {props.itemLabel ?? "itens"}
          </p>

          <div class="flex flex-wrap items-center gap-1">
            <PageButton label="Primeira página" disabled={currentPage() === 1} onClick={() => goTo(1)}>
              <path d="m11 17-5-5 5-5M18 17l-5-5 5-5" stroke-linecap="round" stroke-linejoin="round" />
            </PageButton>
            <PageButton
              label="Página anterior"
              disabled={currentPage() === 1}
              onClick={() => goTo(currentPage() - 1)}
            >
              <path d="m15 18-6-6 6-6" stroke-linecap="round" stroke-linejoin="round" />
            </PageButton>

            <div class="hidden items-center gap-1 sm:flex">
              <For each={buildPageNumbers(currentPage(), totalPages(), 7)}>
                {(entry) => (
                  <Show
                    when={typeof entry === "number"}
                    fallback={<span class="px-1 text-xs text-muted-foreground">…</span>}
                  >
                    <button
                      type="button"
                      aria-current={entry === currentPage() ? "page" : undefined}
                      aria-label={`Página ${entry}`}
                      onClick={() => goTo(Number(entry))}
                      class={cn(
                        "flex h-7 min-w-7 items-center justify-center rounded-sm border px-1.5 text-xs font-medium transition-colors select-none",
                        entry === currentPage()
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
              disabled={currentPage() === totalPages()}
              onClick={() => goTo(currentPage() + 1)}
            >
              <path d="m9 18 6-6-6-6" stroke-linecap="round" stroke-linejoin="round" />
            </PageButton>
            <PageButton
              label="Última página"
              disabled={currentPage() === totalPages()}
              onClick={() => goTo(totalPages())}
            >
              <path d="m13 17 5-5-5-5M6 17l5-5-5-5" stroke-linecap="round" stroke-linejoin="round" />
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
      onClick={() => props.onClick()}
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
