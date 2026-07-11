import type React from "react"
import { cn } from "./cn"
import { CardsGridSkeleton } from "./skeleton"

type GridColumns = 1 | 2 | 3 | 4

interface CardsStateGridProps<T> {
  items: T[]
  isLoading: boolean
  count?: number
  className?: string
  columns?: GridColumns
  renderItem: (item: T, index: number) => React.ReactNode
  renderEmpty?: React.ReactNode
  renderSkeleton?: (index: number) => React.ReactNode
}

const columnsClass: Record<GridColumns, string> = {
  1: "grid-cols-1",
  2: "grid-cols-1 md:grid-cols-2",
  3: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
  4: "grid-cols-1 md:grid-cols-2 lg:grid-cols-4",
}

export function CardsStateGrid<T>({
  items,
  isLoading,
  count = 6,
  className,
  columns = 3,
  renderItem,
  renderEmpty,
  renderSkeleton,
}: CardsStateGridProps<T>) {
  if (isLoading) {
    return (
      <CardsGridSkeleton
        count={count}
        columns={columns}
        className={className}
        renderItem={renderSkeleton}
      />
    )
  }

  if (items.length === 0) return renderEmpty ?? null

  return (
    <div className={cn("grid w-full gap-4", columnsClass[columns], className)}>
      {items.map((item, idx) => renderItem(item, idx))}
    </div>
  )
}
