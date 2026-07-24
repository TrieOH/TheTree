import { Link } from '@tanstack/react-router'
import type { ReactNode } from 'react'
import { cn } from '@/shared/lib/utils'

type QuickActionVariant = 'default' | 'destructive'

type QuickActionLinkProps = {
  children: ReactNode
  disabled?: boolean
  variant?: QuickActionVariant
} & (
    | {
      to: string
      params?: Record<string, string>
      onClick?: never
    }
    | {
      to?: never
      params?: never
      onClick: () => void
    }
  )

export function QuickAction({
  children,
  disabled = false,
  variant = 'default',
  to,
  params,
  onClick,
}: QuickActionLinkProps) {
  const baseClassName = cn(
    'flex items-center justify-between rounded-2xl border border-dashed border-border/70 bg-muted/15 px-4 py-4 text-left',
    'transition-colors hover:border-border hover:bg-muted/30',
    disabled &&
    'cursor-not-allowed opacity-60 hover:border-border/70 hover:bg-muted/15',
    variant === 'destructive' &&
    'border-destructive/30 bg-destructive/5 text-destructive hover:border-destructive/50 hover:bg-destructive/10',
    variant === 'destructive' &&
    disabled &&
    'opacity-60 hover:border-destructive/30 hover:bg-destructive/5',
  )

  if (to) {
    return (
      <Link
        to={to}
        params={params as Record<string, string> | undefined}
        className={baseClassName}
      >
        {children}
      </Link>
    )
  }

  return (
    <button
      type="button"
      className={baseClassName}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  )
}