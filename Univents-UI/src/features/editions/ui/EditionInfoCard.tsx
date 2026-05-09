import { cn } from "@/shared/lib/utils"

export default function EditionInfoCard({
  icon,
  label,
  value,
  sub,
  iconClass = 'bg-primary/8 text-primary',
  footer,
}: {
  icon: React.ReactNode
  label: string
  value: string
  sub?: string
  iconClass?: string
  footer?: React.ReactNode
}) {
  return (
    <div className="bg-card border border-border rounded-2xl p-4">
      <div className="flex items-start gap-3">
        <div className={cn('w-9 h-9 rounded-xl flex items-center justify-center shrink-0', iconClass)}>
          {icon}
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-[10px] font-semibold tracking-widest uppercase text-muted-foreground mb-0.5">
            {label}
          </p>
          <p className="text-sm font-semibold text-foreground leading-snug">{value}</p>
          {sub && <p className="text-xs text-muted-foreground mt-0.5 font-normal">{sub}</p>}
        </div>
      </div>
      {footer && <div className="mt-3">{footer}</div>}
    </div>
  )
}