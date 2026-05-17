import { Link } from '@tanstack/react-router';
import type { ReactNode } from 'react'
import { cn } from '@/shared/lib/utils'

export function SectionCard({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="bg-muted/40 rounded-xl border border-border/50 px-4 py-3">
      <h3 className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/70 mb-2.5">
        {label}
      </h3>
      {children}
    </div>
  )
}

export function InfoRow({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex flex-col gap-0.5">
      <span className="text-[11px] text-muted-foreground">{label}</span>
      <span className={cn('text-sm text-foreground font-medium', mono && 'font-mono text-xs')}>{value}</span>
    </div>
  )
}

export function SocialChip({ href, label, icon }: { href: string; label: string; icon: ReactNode }) {
  return (
    <Link
      to={href}
      target="_blank"
      rel="noopener noreferrer"
      className={cn(
        'flex items-center gap-2 px-3 py-2 rounded-lg',
        'bg-background border border-border',
        'text-sm text-foreground/80 hover:text-foreground',
        'hover:bg-muted transition-colors',
      )}
    >
      <span className="text-muted-foreground">{icon}</span>
      <span>{label}</span>
    </Link>
  )
}
