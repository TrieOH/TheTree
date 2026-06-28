import { cn } from "./cn"
import { useMemo } from "react"

function generateId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function")
    return crypto.randomUUID()

  return Math.random().toString(36).substring(2, 11)
}

interface CardSkeletonProps {
  className?: string
  rows?: number
}

/**
 * A single card skeleton used as a loading placeholder.
 */
export function CardSkeleton({ className, rows = 3 }: CardSkeletonProps) {
  const rowSkeletons = useMemo(
    () => Array.from({ length: rows }, () => generateId()),
    [rows],
  )

  return (
    <div className={cn("p-6 border rounded-xl bg-card shadow-sm", className)}>
      <div className="h-6 w-3/4 mb-4 bg-muted animate-pulse rounded-md" />
      {rowSkeletons.map((id, i) => (
        <div
          key={id}
          className={cn(
            "h-4 w-full bg-muted animate-pulse rounded-md",
            i !== rows - 1 && "mb-2",
          )}
        />
      ))}
    </div>
  )
}

interface CardsGridSkeletonProps {
  count?: number
  className?: string
  columns?: 1 | 2 | 3 | 4
}

/**
 * A responsive grid of card skeletons.
 */
export function CardsGridSkeleton({
  count = 6,
  className,
  columns = 3,
}: CardsGridSkeletonProps) {
  const columnsClass: Record<number, string> = {
    1: "grid-cols-1",
    2: "grid-cols-1 md:grid-cols-2",
    3: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
    4: "grid-cols-1 md:grid-cols-2 lg:grid-cols-4",
  }

  const skeletons = useMemo(
    () => Array.from({ length: count }, () => generateId()),
    [count],
  )

  return (
    <div className={cn("grid gap-4", columnsClass[columns], className)}>
      {skeletons.map((id) => (
        <CardSkeleton key={id} />
      ))}
    </div>
  )
}
