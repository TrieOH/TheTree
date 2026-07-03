import { cn } from '@/shared/lib/utils'
import { Fingerprint } from 'lucide-react'
import type { CapabilityI } from '../model'
import { timeAgo } from '@/shared/lib/date-utils'

interface PropsI {
  data: CapabilityI
}

export function CapabilityCard({ data }: PropsI) {
  return (
    <div
      className={cn(
        'bg-card rounded-sm w-full cursor-default',
        'ring-1 ring-foreground/10 shadow-xs',
        'flex items-center gap-3 px-4 py-3',
        'hover:ring-foreground/20 duration-150',
      )}
    >
      <div className="shrink-0 size-9 rounded-full bg-muted ring-1 ring-foreground/10 flex items-center justify-center">
        <Fingerprint className="size-4 text-muted-foreground" />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm font-semibold truncate">
            {data.resource}:{data.action}
          </span>
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          Created {timeAgo(data.created_at)}
        </p>
      </div>
    </div>
  )
}
