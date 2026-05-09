import { cn } from "@/shared/lib/utils"
import { Skeleton } from "../shadcn/skeleton"
import { useMemo } from "react"


interface PropsI {
  className?: string
  rows?: number
}

export function CardSkeleton({ className, rows = 3 }: PropsI) {
  const rowSkeletons = useMemo(
    () => Array.from({ length: rows }, () => crypto.randomUUID()),
    [rows]
  )
  return (
    <div className={cn("p-6 border rounded-xl bg-card shadow-sm", className)}>
      <Skeleton className="h-6 w-3/4 mb-4" />
      {rowSkeletons.map((id, i) => (
        <Skeleton
          key={id}
          className={cn("h-4 w-full", i !== rows - 1 && "mb-2")}
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

export function CardsGridSkeleton({ 
  count = 6, 
  className,
  columns = 3 
}: CardsGridSkeletonProps) {
  const columnsClass = {
    1: "grid-cols-1",
    2: "grid-cols-1 md:grid-cols-2",
    3: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
    4: "grid-cols-1 md:grid-cols-2 lg:grid-cols-4",
  }

  const skeletons = useMemo(
    () => Array.from({ length: count }, () => crypto.randomUUID()),
    [count]
  )

  return (
    <div className={cn("grid gap-4", columnsClass[columns], className)}>
      {skeletons.map(id => (
        <CardSkeleton key={id} />
      ))}
    </div>
  )
}