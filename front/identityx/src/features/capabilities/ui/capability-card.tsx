import { cn } from '@/shared/lib/utils'
import { Layers3 } from 'lucide-react'
import type { CapabilityI } from '../model'
import { timeAgo } from '@trieoh/shared-utils'

interface PropsI {
  data: CapabilityI
}

export function CapabilityCard({ data }: PropsI) {
  return (
    <div
      className={cn(
        'bg-card rounded-md w-full min-w-0 cursor-default',
        'ring-1 ring-foreground/10 shadow-xs',
        'flex items-start gap-3 px-3 py-3 sm:px-4',
        'hover:ring-foreground/20 duration-150',
      )}
    >
      <div className="shrink-0 size-10 rounded-md bg-muted ring-1 ring-foreground/10 flex items-center justify-center">
        <Layers3 className="size-4 text-muted-foreground" />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <span className="block text-sm font-semibold truncate font-mono">
              {data.resource}:{data.action}
            </span>
            <p className="mt-1 text-xs text-muted-foreground truncate">
              Capability
            </p>
          </div>
        </div>

        <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px] text-muted-foreground">
          <span>
            Created {timeAgo(data.created_at)}
          </span>
        </div>
      </div>
    </div>
  )
}
