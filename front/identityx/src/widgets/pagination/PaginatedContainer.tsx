import { ArrowDown, ArrowUp, ArrowUpDown, ChevronFirst, ChevronLast, ChevronLeft, ChevronRight, Minus, Search, X } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import type { KeyboardEvent, ReactNode } from "react"

export type SortDirection = "asc" | "desc"

export interface SortField<T> {
  key: keyof T
  label: string
  comparator?: (a: T, b: T) => number
}

export interface SortState<T> {
  field: keyof T
  direction: SortDirection
}

export interface PaginationInfo {
  page: number
  totalPages: number
  pageSize: number
  start: number
  end: number
  count: number
}

/**
 * Layout modes for the items wrapper.
 *
 * "grid"   - CSS grid with auto-fill columns. Control minimum item width via
 *            `minItemWidth` prop (default: "200px"). Items stretch to fill
 *            the row evenly - no gap on the right side.
 *
 * "list"   - `flex flex-col`, one item per row, full width.
 *
 * "custom" - No wrapper is rendered at all. Your `renderItems` receives the
 *            slice and you mount whatever container you want:
 *
 *            renderItems={(slice) => (
 *              <div className="grid grid-cols-4 gap-3">
 *                {slice.map(item => <MyCard key={item.id} item={item} />)}
 *              </div>
 *            )}
 */
export type LayoutMode = "grid" | "list" | "custom"

/**
 * Tailwind gap size applied to the items wrapper.
 * Maps to `gap-{n}` utility classes.
 */
export type GapSize = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "8" | "10"

export interface PaginatedContainerProps<T> {
  /** Full dataset - filter externally before passing. */
  items: T[]
  /** Items per page (default: 8). */
  pageSize?: number
  /** Starting page (default: 1). */
  defaultPage?: number

  /**
   * Sortable field definitions. When provided a Sort button appears in the
   * header that opens a dropdown panel (field picker + direction toggle).
   */
  sortFields?: SortField<T>[]
  /**
   * Controlled sort state. Pair with `onSortChange` for external control.
   * Leave unset to let the component manage sort internally.
   */
  sort?: SortState<T>
  onSortChange?: (sort: SortState<T>) => void

  /**
   * Controlled filter value. The component renders the search input and
   * calls `onFilterChange` on every keystroke - the actual filtering is
   * done by you so you can match any fields you want.
   */
  filterValue?: string
  onFilterChange?: (value: string) => void
  /** Placeholder for the filter input (default: "Filter…"). */
  filterPlaceholder?: string

  /**
   * Item noun used in the footer counter.
   * E.g. "organizations" → "Showing 1–4 of 24 organizations"
   */
  itemLabel?: string
  /** Render the current page slice. */
  renderItems: (slice: T[]) => ReactNode
  /** Shown when the item list is empty after filtering. */
  emptyState?: ReactNode
  /** Extra content placed in the header to the right of the Sort button. */
  headerActions?: ReactNode
  className?: string

  /**
   * Layout mode for the items wrapper (default: "grid").
   *
   * "grid"   - CSS grid auto-fill; items stretch to fill available width.
   * "list"   - flex-col; one item per row.
   * "custom" - no wrapper; you own the container inside `renderItems`.
   */
  layout?: LayoutMode
  /**
   * Gap between items (default: "3").
   * Maps to Tailwind's `gap-{n}` utility.
   */
  gap?: GapSize
  /**
   * Minimum width for each item in grid layout (default: "200px").
   * The grid will fit as many columns as possible at this minimum width,
   * then stretch items to fill the row - no empty space on the right.
   */
  minItemWidth?: string
}

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
}

interface ItemsWrapperProps {
  layout: LayoutMode
  gap: GapSize
  minItemWidth: string
  children: ReactNode
}

function ItemsWrapper({ layout, gap, minItemWidth, children }: ItemsWrapperProps) {
  const gapClass = GAP_CLASS[gap]

  if (layout === "custom") return <>{children}</>

  if (layout === "list") {
    return <div className={`flex flex-col w-full ${gapClass}`}>{children}</div>
  }

  // "grid" - auto-fill columns that stretch to fill the row.
  return (
    <div
      className={`w-full grid justify-center justify-items-center ${gapClass}`}
      style={{
        gridTemplateColumns: `repeat(auto-fill, minmax(${minItemWidth}, 1fr))`,
      }}
    >
      {children}
    </div>
  )
}

function buildPageNumbers(current: number, total: number): (number | "…")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  if (current <= 4) return [1, 2, 3, 4, 5, "…", total]
  if (current >= total - 3) return [1, "…", total - 4, total - 3, total - 2, total - 1, total]
  return [1, "…", current - 1, current, current + 1, "…", total]
}

function defaultComparator<T>(a: T, b: T, key: keyof T): number {
  const av = a[key]
  const bv = b[key]
  if (typeof av === "number" && typeof bv === "number") return av - bv
  return String(av).localeCompare(String(bv))
}

function applySorting<T>(
  items: T[],
  sort: SortState<T>,
  fields: SortField<T>[],
): T[] {
  const field = fields.find((f) => f.key === sort.field)
  const dir = sort.direction === "desc" ? -1 : 1
  return [...items].sort((a, b) => {
    const result = field?.comparator
      ? field.comparator(a, b)
      : defaultComparator(a, b, sort.field)
    return result * dir
  })
}

interface SortPanelProps<T> {
  fields: SortField<T>[]
  sort: SortState<T>
  onChange: (sort: SortState<T>) => void
  onClose: () => void
}

function SortPanel<T>({ fields, sort, onChange, onClose }: SortPanelProps<T>) {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) onClose()
    }
    document.addEventListener("mousedown", handler)
    return () => document.removeEventListener("mousedown", handler)
  }, [onClose])

  const handleFieldClick = (key: keyof T) => {
    onChange({
      field: key,
      direction:
        sort.field === key && sort.direction === "asc" ? "desc" : "asc",
    })
  }

  const handleFieldKeyDown = (
    e: KeyboardEvent<HTMLDivElement>,
    key: keyof T,
  ) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault()
      handleFieldClick(key)
    }
  }

  return (
    <div
      ref={ref}
      role="dialog"
      aria-label="Sort options"
      className="absolute right-0 top-[calc(100%+6px)] z-50 w-52 rounded-md border border-border bg-card shadow-lg shadow-black/5 p-2"
    >
      {/* Header */}
      <div className="flex items-center justify-between px-2 pb-2 mb-1">
        <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">
          Sort by
        </span>
        <button
          type="button"
          onClick={onClose}
          aria-label="Close sort panel"
          className="flex h-5 w-5 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <X size={12} />
        </button>
      </div>

      {/* Field list */}
      <div className="flex flex-col gap-0.5 mb-2">
        {fields.map((f) => {
          const isActive = sort.field === f.key
          return (
            <div
              key={String(f.key)}
              role="button"
              tabIndex={0}
              onClick={() => handleFieldClick(f.key)}
              onKeyDown={(e) => handleFieldKeyDown(e, f.key)}
              className={[
                "flex cursor-pointer items-center justify-between rounded-md px-2.5 py-2 text-sm transition-colors select-none border",
                isActive
                  ? "bg-primary/10 text-primary border-primary/20 font-medium"
                  : "text-foreground hover:bg-muted border-transparent",
              ].join(" ")}
            >
              <span>{f.label}</span>
              <span
                className={
                  isActive ? "text-primary" : "text-muted-foreground/40"
                }
              >
                {isActive ? (
                  sort.direction === "asc" ? (
                    <ArrowUp size={13} />
                  ) : (
                    <ArrowDown size={13} />
                  )
                ) : (
                  <Minus size={13} />
                )}
              </span>
            </div>
          )
        })}
      </div>

      {/* Direction toggle */}
      <div className="grid grid-cols-2 gap-1.5 border-t border-border pt-2">
        {(["asc", "desc"] as const).map((dir) => (
          <button
            type="button"
            key={dir}
            onClick={() => onChange({ ...sort, direction: dir })}
            className={[
              "flex items-center justify-center gap-1.5 rounded-md border py-1.5 text-xs font-medium transition-colors",
              sort.direction === dir
                ? "border-primary/30 bg-primary/10 text-primary"
                : "border-border text-muted-foreground hover:bg-muted hover:text-foreground",
            ].join(" ")}
          >
            {dir === "asc" ? <ArrowUp size={11} /> : <ArrowDown size={11} />}
            {dir === "asc" ? "Asc" : "Desc"}
          </button>
        ))}
      </div>
    </div>
  )
}

interface PaginationButtonProps {
  onClick: () => void
  disabled?: boolean
  active?: boolean
  title?: string
  "aria-current"?: "page" | undefined
  children: ReactNode
}

function PaginationButton({
  onClick,
  disabled = false,
  active = false,
  children,
  ...rest
}: PaginationButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={[
        "flex h-7 min-w-7 items-center justify-center rounded-sm border px-1.5 text-xs font-medium transition-colors select-none",
        "disabled:cursor-not-allowed disabled:opacity-35",
        active
          ? "border-primary/30 bg-primary/10 text-primary"
          : "border-border bg-transparent text-muted-foreground hover:bg-muted hover:text-foreground",
      ].join(" ")}
      {...rest}
    >
      {children}
    </button>
  )
}

export function PaginatedContainer<T>({
  items,
  pageSize = 8,
  defaultPage = 1,
  sortFields,
  sort: controlledSort,
  onSortChange,
  filterValue,
  onFilterChange,
  filterPlaceholder = "Filter…",
  itemLabel = "items",
  renderItems,
  emptyState,
  headerActions,
  className,
  layout = "grid",
  gap = "3",
  minItemWidth = "200px",
}: PaginatedContainerProps<T>) {
  const [page, setPage] = useState(defaultPage)
  const [internalSort, setInternalSort] = useState<SortState<T> | null>(
    sortFields ? { field: sortFields[0].key, direction: "asc" as const } : null,
  )
  const [sortOpen, setSortOpen] = useState(false)

  const activeSort = controlledSort ?? internalSort

  const handleSortChange = (next: SortState<T>) => {
    onSortChange?.(next)
    if (!controlledSort) setInternalSort(next)
  }

  const processedItems =
    activeSort && sortFields
      ? applySorting(items, activeSort, sortFields)
      : items

  const totalPages = Math.max(1, Math.ceil(processedItems.length / pageSize))
  const safePage = Math.min(page, totalPages)
  const start = (safePage - 1) * pageSize
  const slice = processedItems.slice(start, start + pageSize)
  const pageNums = buildPageNumbers(safePage, totalPages)

  const go = (p: number) => setPage(Math.max(1, Math.min(totalPages, p)))

  const activeSortLabel =
    activeSort && sortFields
      ? sortFields.find((f) => f.key === activeSort.field)?.label
      : null

  const hasHeader =
    onFilterChange !== undefined || sortFields || headerActions

  return (
    <div
      className={[
        "w-full overflow-visible rounded-md border border-border bg-card",
        className ?? "",
      ].join(" ")}
    >
      {/* Header */}
      {hasHeader && (
        <div className="flex flex-wrap items-center justify-between gap-2.5 border-b border-border px-4 py-3">
          {onFilterChange !== undefined && (
            <div className="relative w-full sm:flex-1 sm:max-w-sm">
              <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-muted-foreground flex items-center">
                <Search size={13} />
              </span>
              <input
                type="text"
                name="table-filter"
                value={filterValue ?? ""}
                onChange={(e) => {
                  onFilterChange(e.target.value)
                  setPage(1)
                }}
                placeholder={filterPlaceholder}
                aria-label={filterPlaceholder}
                className="h-9 w-full rounded-md border border-border bg-muted/50 pl-8 pr-3 text-sm text-foreground placeholder:text-muted-foreground outline-none transition-colors focus:border-ring focus:bg-background"
              />
            </div>
          )}

          <div className="flex flex-wrap items-center gap-2 sm:ml-auto">
            {headerActions}

            {sortFields && activeSort && (
              <div className="relative">
                <button
                  type="submit"
                  onClick={() => setSortOpen((v) => !v)}
                  aria-expanded={sortOpen}
                  aria-haspopup="dialog"
                  className={[
                    "flex h-9 items-center gap-1.5 rounded-md border px-3 text-sm transition-colors",
                    sortOpen
                      ? "border-ring/50 bg-muted text-foreground"
                      : "border-border text-muted-foreground hover:border-border/80 hover:bg-muted hover:text-foreground",
                  ].join(" ")}
                >
                  <ArrowUpDown size={13} />
                  <span>Sort</span>
                  {activeSortLabel && (
                    <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[11px] font-medium text-primary">
                      {activeSortLabel}
                      {" · "}
                      {activeSort.direction === "asc" ? "↑" : "↓"}
                    </span>
                  )}
                </button>

                {sortOpen && (
                  <SortPanel
                    fields={sortFields}
                    sort={activeSort}
                    onChange={(s) => {
                      handleSortChange(s)
                      setPage(1)
                    }}
                    onClose={() => setSortOpen(false)}
                  />
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Content */}
      <div className="p-4">
        {slice.length === 0 ? (
          emptyState ?? (
            <p className="py-10 text-center text-sm text-muted-foreground">
              No {itemLabel} found.
            </p>
          )
        ) : (
          <ItemsWrapper layout={layout} gap={gap} minItemWidth={minItemWidth}>
            {renderItems(slice)}
          </ItemsWrapper>
        )}
      </div>

      {/* Footer */}
      <div className="flex flex-wrap items-center justify-between gap-2.5 border-t border-border px-4 py-2.5">
        <span className="text-xs text-muted-foreground">
          {processedItems.length === 0 ? (
            `0 ${itemLabel}`
          ) : (
            <>
              Showing{" "}
              <span className="font-medium text-foreground">
                {start + 1}–{Math.min(start + pageSize, processedItems.length)}
              </span>{" "}
              of{" "}
              <span className="font-medium text-foreground">
                {processedItems.length}
              </span>{" "}
              {itemLabel}
            </>
          )}
        </span>

        <nav aria-label="Pagination" className="flex items-center gap-1">
          <PaginationButton
            onClick={() => go(1)}
            disabled={safePage === 1}
            title="First page"
          >
            <ChevronFirst size={14} />
          </PaginationButton>
          <PaginationButton
            onClick={() => go(safePage - 1)}
            disabled={safePage === 1}
            title="Previous page"
          >
            <ChevronLeft size={14} />
          </PaginationButton>

          {pageNums.map((p) =>
            p === "…" ? (
              <span
                key={p}
                className="flex h-7 w-7 items-center justify-center text-xs text-muted-foreground"
              >
                …
              </span>
            ) : (
              <PaginationButton
                key={p}
                onClick={() => go(p)}
                active={p === safePage}
                aria-current={p === safePage ? "page" : undefined}
              >
                {p}
              </PaginationButton>
            ),
          )}

          <PaginationButton
            onClick={() => go(safePage + 1)}
            disabled={safePage === totalPages}
            title="Next page"
          >
            <ChevronRight size={14} />
          </PaginationButton>
          <PaginationButton
            onClick={() => go(totalPages)}
            disabled={safePage === totalPages}
            title="Last page"
          >
            <ChevronLast size={14} />
          </PaginationButton>
        </nav>
      </div>
    </div>
  )
}