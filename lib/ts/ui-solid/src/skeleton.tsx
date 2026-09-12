import { cn } from "./lib/cn";

export interface SkeletonProps {
  class?: string;
}

export function Skeleton(props: SkeletonProps) {
  return (
    <div
      aria-hidden="true"
      class={cn("animate-pulse rounded-md bg-muted", props.class)}
    />
  );
}

/** Card-shaped placeholder that matches the admin grid's item height. */
export function CardSkeleton(props: { class?: string }) {
  return (
    <div
      class={cn(
        "flex flex-col gap-3 rounded-lg border border-border bg-card p-4",
        props.class,
      )}
    >
      <Skeleton class="h-28 w-full" />
      <Skeleton class="h-4 w-3/5" />
      <Skeleton class="h-3 w-2/5" />
    </div>
  );
}
