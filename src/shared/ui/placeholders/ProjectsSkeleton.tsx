import { CardsGridSkeleton } from "./CardSkeleton"

interface PropsI {
  count?: number
}

export function ProjectsSkeleton({ count = 6 }: PropsI) {
  return (
    <main className="w-full bg-background flex flex-col items-center my-4">
      {/* Header skeleton */}
      <div className="text-center space-y-2 mb-7">
        <div className="h-9 bg-muted rounded-lg w-48 mx-auto animate-pulse" />
        <div className="h-4 bg-muted rounded w-64 mx-auto animate-pulse" />
      </div>
      
      {/* Grid skeleton */}
      <div className="max-w-7xl w-full xs:px-4">
        <CardsGridSkeleton count={count} columns={3} />
      </div>
    </main>
  )
}