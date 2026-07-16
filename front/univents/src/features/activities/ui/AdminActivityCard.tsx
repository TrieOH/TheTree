import type React from 'react'
import { motion } from 'motion/react'
import {
  CalendarDays,
  Coins,
  CheckCircle2,
  Eye,
  PencilLine,
  MapPin,
  MoreVertical,
  Users,
} from 'lucide-react'
import { formatTime, formatDuration, difficultyConfig } from './ActivityCard'
import type { ActivityI } from '@/features/activities/model'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/shared/ui/shadcn/context-menu'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/shadcn/dropdown-menu'
import { cn } from '@/shared/lib/utils'

export const statusConfig: Record<
  ActivityI['status'],
  { label: string; dot: string; pill: string }
> = {
  draft: {
    label: 'Rascunho',
    dot: 'bg-slate-400',
    pill: 'bg-slate-50 text-slate-600 border-slate-200',
  },
  published: {
    label: 'Publicado',
    dot: 'bg-blue-500',
    pill: 'bg-blue-50 text-blue-700 border-blue-200',
  },
  ongoing: {
    label: 'Em andamento',
    dot: 'bg-emerald-500',
    pill: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  },
  completed: {
    label: 'Concluído',
    dot: 'bg-slate-900',
    pill: 'bg-slate-100 text-slate-900 border-slate-200',
  },
  canceled: {
    label: 'Cancelado',
    dot: 'bg-red-500',
    pill: 'bg-red-50 text-red-700 border-red-200',
  },
}

interface AdminActivityCardProps {
  activity: ActivityI
  index?: number
  onManage: (activity: ActivityI) => void
  onPublish?: (activity: ActivityI) => void
  onComplete?: (activity: ActivityI) => void
}

function MenuItems({
  isContext = false,
  activity,
  onManage,
  onPublish,
  onComplete,
}: {
  isContext?: boolean
  activity: ActivityI
  onManage: () => void
  onPublish?: () => void
  onComplete?: () => void
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem
  const stop = (action?: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
    e.preventDefault()
    e.stopPropagation()
    action?.()
  }

  return (
    <>
      <Item onClick={stop(onManage)}>
        <PencilLine className="size-4" />
        <span>Editar atividade</span>
      </Item>

      {activity.status === 'draft' && onPublish && (
        <Item onClick={stop(onPublish)}>
          <Eye className="size-4" />
          <span>Publicar atividade</span>
        </Item>
      )}

      {onComplete && (
        <Item onClick={stop(onComplete)}>
          <CheckCircle2 className="size-4" />
          <span>Concluir atividade</span>
        </Item>
      )}
    </>
  )
}

export default function AdminActivityCard({
  activity,
  index = 0,
  onPublish,
  onComplete,
  onManage,
}: AdminActivityCardProps) {
  const status = statusConfig[activity.status]
  const difficulty = difficultyConfig[activity.difficulty]
  const hasCapacity = activity.has_capacity

  const handleEdit = () => onManage(activity)
  const handlePublish = () => onPublish?.(activity)
  const handleComplete = () => onComplete?.(activity)

  const Article = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.35, ease: [0.25, 0.1, 0.25, 1] }}
      className={cn(
        'group relative flex w-full min-w-62.5 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left',
        'ring-1 ring-foreground/10 shadow-xs',
        'transform-gpu will-change-transform',
        'transition-all duration-300 ease-out',
        'hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm',
        'focus:outline-none focus-visible:outline-none focus-visible:ring-0',
      )}
      role="button"
      tabIndex={0}
      onClick={handleEdit}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          handleEdit()
        }
      }}
    >
      <div className="relative overflow-hidden bg-linear-to-br from-muted via-background to-muted/40 px-4 py-4">
        <div className="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-primary/70 via-primary to-cyan-500/70" />

        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(e) => e.stopPropagation()}
                  className={cn(
                    'inline-flex size-9 items-center justify-center rounded-full',
                    'bg-background/85 text-foreground shadow-sm backdrop-blur-sm',
                    'transition-colors hover:bg-background',
                  )}
                  aria-label={`Abrir ações de ${activity.title}`}
                >
                  <MoreVertical className="size-4" />
                </button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems
                activity={activity}
                onManage={handleEdit}
                onPublish={handlePublish}
                onComplete={handleComplete}
              />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="flex items-start justify-between gap-3 pr-10">
          <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-1.5">
              <span className={cn(
                'inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm',
                status.pill,
              )}>
                <span className={cn('size-1.5 rounded-full', status.dot)} />
                <span className="max-w-28 truncate">{status.label}</span>
              </span>

              <span className={cn(
                'inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm',
                difficulty.pill,
                difficulty.pillText,
              )}>
                <span className="max-w-28 truncate">{difficulty.label}</span>
              </span>
            </div>

            <div className="min-w-0 space-y-1">
              <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
                {activity.title}
              </h3>
              {activity.presenter_name && (
                <p className="line-clamp-1 text-sm text-muted-foreground">
                  {activity.presenter_name}
                </p>
              )}
            </div>
          </div>

          <div className="shrink-0 text-right">
            <p className="text-sm font-semibold text-foreground">
              {formatTime(activity.starts_at)}
            </p>
            <p className="mt-0.5 text-[11px] text-muted-foreground">
              {formatDuration(activity.starts_at, activity.ends_at)}
            </p>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
        <div className="min-w-0 flex-1 space-y-1.5">
          <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            <span className="inline-flex min-w-0 items-center gap-1.5">
              <MapPin className="size-3.5 shrink-0" />
              <span className="truncate">{activity.location}</span>
            </span>
            <span className="inline-flex items-center gap-1.5">
              <CalendarDays className="size-3.5 shrink-0" />
              <span className="truncate">
                {formatTime(activity.starts_at)} - {formatTime(activity.ends_at)}
              </span>
            </span>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            {activity.token_cost > 0 && (
              <span className="inline-flex items-center gap-1.5 rounded-full border border-amber-500/20 bg-amber-500/10 px-2.5 py-1 text-[11px] font-medium text-amber-700">
                <Coins className="size-3.5" />
                {activity.token_cost} tokens
              </span>
            )}

            {hasCapacity && (
              <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/80 px-2.5 py-1 text-[11px] font-medium text-muted-foreground">
                <Users className="size-3.5" />
                {activity.remaining_capacity} / {activity.capacity} vagas
              </span>
            )}
          </div>
        </div>
      </div>
    </motion.article>
  )

  return (
    <ContextMenu>
      <ContextMenuTrigger render={Article} />
      <ContextMenuContent className="w-56">
        <MenuItems
          isContext
          activity={activity}
          onManage={handleEdit}
          onPublish={handlePublish}
          onComplete={handleComplete}
        />
      </ContextMenuContent>
    </ContextMenu>
  )
}
