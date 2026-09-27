import type React from "react"
import type { LucideIcon } from "lucide-react"
import { cn } from "./cn"

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  action?: React.ReactNode
  className?: string
  eyebrow?: string
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  className,
  eyebrow,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-2xl border border-border/60 bg-card px-6 py-8 text-center shadow-sm sm:px-8 sm:py-10",
        className,
      )}
    >
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.08),transparent_42%)]" />

      <div className="relative z-10 mx-auto flex max-w-md flex-col items-center">
        {eyebrow && (
          <p className="mb-3 text-xs uppercase tracking-[0.24em] text-muted-foreground">
            {eyebrow}
          </p>
        )}

        {Icon && (
          <div className="mb-5 flex h-14 w-14 items-center justify-center rounded-2xl border border-border/60 bg-muted/70 shadow-inner">
            <Icon className="h-6 w-6 text-muted-foreground" />
          </div>
        )}

        <h3 className="text-lg font-semibold tracking-tight text-foreground sm:text-xl">
          {title}
        </h3>

        {description && (
          <p className="mt-2 max-w-sm text-sm leading-relaxed text-muted-foreground sm:text-base">
            {description}
          </p>
        )}

        {action && (
          <div className="mt-5 flex justify-center">
            {action}
          </div>
        )}
      </div>
    </div>
  )
}
