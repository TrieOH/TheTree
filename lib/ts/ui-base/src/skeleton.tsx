import { cn } from "./cn"

interface CardSkeletonProps {
  className?: string
  mediaClassName?: string
  badgeClassName?: string
  actionClassName?: string
  titleClassName?: string
  titleWidthClassName?: string
  descriptionClassName?: string
  footerClassName?: string
  showBadge?: boolean
  showAction?: boolean
  showDescription?: boolean
  showFooter?: boolean
  rows?: number
}

export function CardSkeleton({
  className,
  mediaClassName,
  badgeClassName,
  actionClassName,
  titleClassName,
  titleWidthClassName,
  descriptionClassName,
  footerClassName,
  showBadge = true,
  showAction = true,
  showDescription = true,
  showFooter = false,
  rows = 2,
}: CardSkeletonProps) {
  return (
    <article className={cn("group relative overflow-hidden rounded-2xl border border-border/60 bg-card shadow-sm", className)}>
      <div className={cn("relative aspect-4/3 bg-muted/70", mediaClassName)}>
        <div className="h-full w-full bg-[linear-gradient(110deg,transparent_18%,rgba(255,255,255,0.28)_36%,transparent_54%)] bg-[length:200%_100%] animate-pulse" />
        {showBadge && (
          <div className={cn("absolute top-3 left-3 h-5 w-24 rounded-full bg-background/90 backdrop-blur-sm animate-pulse", badgeClassName)} />
        )}
        {showAction && (
          <div className={cn("absolute top-3 right-3 h-8 w-8 rounded-full bg-background/90 backdrop-blur-sm animate-pulse", actionClassName)} />
        )}
      </div>

      <div className={cn("space-y-3 p-4 md:p-5", footerClassName)}>
        <div className="h-3 w-24 rounded-full bg-muted animate-pulse" />

        <div className={cn("space-y-2", titleClassName)}>
          <div className={cn("h-4 rounded-lg bg-muted animate-pulse", titleWidthClassName ?? "w-5/6")} />
          <div className="h-4 w-4/6 rounded-lg bg-muted animate-pulse" />
        </div>

        {showDescription && (
          <div className={cn("space-y-2", descriptionClassName)}>
            {Array.from({ length: rows }).map((_, idx) => (
              <div
                key={idx}
                className={cn(
                  "h-3.5 rounded-full bg-muted/80 animate-pulse",
                  idx === 0 && "w-full",
                  idx === 1 && "w-5/6",
                  idx >= 2 && "w-2/3",
                )}
              />
            ))}
          </div>
        )}

        {showFooter && <div className="h-3.5 w-3/5 rounded-full bg-muted/80 animate-pulse" />}
      </div>
    </article>
  )
}

interface CardsGridSkeletonProps {
  count?: number
  className?: string
  columns?: 1 | 2 | 3 | 4
  renderItem?: (index: number) => React.ReactNode
}

export function CardsGridSkeleton({
  count = 6,
  className,
  columns = 3,
  renderItem,
}: CardsGridSkeletonProps) {
  const columnsClass: Record<number, string> = {
    1: "grid-cols-1",
    2: "grid-cols-1 md:grid-cols-2",
    3: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
    4: "grid-cols-1 md:grid-cols-2 lg:grid-cols-4",
  }

  return (
    <div className={cn("grid w-full gap-4", columnsClass[columns], className)}>
      {Array.from({ length: count }).map((_, idx) => (
        <div key={idx} className="w-full">
          {renderItem ? renderItem(idx) : <CardSkeleton />}
        </div>
      ))}
    </div>
  )
}
