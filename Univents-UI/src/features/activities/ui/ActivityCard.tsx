import { motion } from 'motion/react'
import { MapPin, Users, Coins } from 'lucide-react'
import type { ActivityI } from '@/features/activities/model'
import { cn } from '@/shared/lib/utils'

type Difficulty = ActivityI['difficulty']

interface DifficultyConfig {
  label: string
  accent: string
  pill: string
  pillText: string
}

export const difficultyConfig: Record<Difficulty, DifficultyConfig> = {
  no_prerequisites: {
    label: 'Sem pré-requisitos',
    accent: 'bg-emerald-500',
    pill: 'bg-emerald-50',
    pillText: 'text-emerald-700',
  },
  beginner: {
    label: 'Iniciante',
    accent: 'bg-blue-500',
    pill: 'bg-blue-50',
    pillText: 'text-blue-700',
  },
  intermediate: {
    label: 'Intermediário',
    accent: 'bg-amber-400',
    pill: 'bg-amber-50',
    pillText: 'text-amber-700',
  },
  advanced: {
    label: 'Avançado',
    accent: 'bg-orange-500',
    pill: 'bg-orange-50',
    pillText: 'text-orange-700',
  },
  expert: {
    label: 'Especialista',
    accent: 'bg-red-500',
    pill: 'bg-red-50',
    pillText: 'text-red-700',
  },
}

export function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDuration(startsAt: string, endsAt: string) {
  const minutes = Math.round(
    (new Date(endsAt).getTime() - new Date(startsAt).getTime()) / 60000
  )
  if (minutes < 60) return `${minutes}min`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m > 0 ? `${h}h ${m}min` : `${h}h`
}

interface ActivityCardProps {
  activity: ActivityI
  registered?: boolean
  onClick?: () => void
  index?: number
}

export default function ActivityCard({
  activity,
  registered,
  onClick,
  index = 0,
}: ActivityCardProps) {
  const difficulty = difficultyConfig[activity.difficulty]
  const isFull = activity.has_capacity && activity.remaining_capacity <= 0

  return (
    <motion.button
      initial={{ opacity: 0, x: -10 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.05 }}
      onClick={onClick}
      className={cn(
        'grid w-full grid-cols-[3px_52px_1fr] gap-x-3 text-left',
        'rounded-lg border bg-card p-3.5 transition-colors',
        'hover:border-border/60 active:scale-[0.99]',
        registered && 'border-emerald-200 bg-emerald-50/60'
      )}
    >
      <span
        className={cn(
          'w-0.75 self-stretch rounded-full',
          registered ? 'bg-emerald-500' : isFull ? 'bg-red-400' : difficulty.accent
        )}
      />

      <div className="flex flex-col items-end border-r border-border/40 pr-3 pt-0.5">
        <span className="text-sm font-medium leading-tight text-foreground">
          {formatTime(activity.starts_at)}
        </span>
        <span className="mt-0.5 text-[11px] text-muted-foreground">
          {formatDuration(activity.starts_at, activity.ends_at)}
        </span>
      </div>

      <div className="flex min-w-0 flex-col gap-1.5">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <p className="truncate text-sm font-medium leading-snug text-foreground">
              {activity.title}
            </p>
            {activity.presenter_name && (
              <p className="mt-0.5 text-xs text-muted-foreground">
                {activity.presenter_name}
              </p>
            )}
          </div>

          {registered && (
            <span className="shrink-0 rounded bg-emerald-100 px-2 py-0.5 text-[11px] font-medium text-emerald-700">
              Inscrito
            </span>
          )}
          {isFull && !registered && (
            <span className="shrink-0 rounded bg-red-50 px-2 py-0.5 text-[11px] font-medium text-red-600">
              Esgotado
            </span>
          )}
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <span
            className={cn(
              'rounded text-[11px] font-medium',
              difficulty.pill,
              difficulty.pillText
            )}
          >
            {difficulty.label}
          </span>

          {activity.token_cost > 0 && (
            <span className="flex items-center gap-1 text-[11px] font-medium text-amber-700">
              <Coins className="h-3 w-3" />
              {activity.token_cost}
            </span>
          )}
        </div>

        <div className="flex items-center gap-3 text-[11px] text-muted-foreground">
          <span className="flex items-center gap-1">
            <MapPin className="h-3 w-3 shrink-0" />
            {activity.location}
          </span>
          {activity.has_capacity && (
            <span className="flex items-center gap-1">
              <Users className="h-3 w-3 shrink-0" />
              {activity.remaining_capacity} / {activity.capacity} vagas
            </span>
          )}
        </div>
      </div>
    </motion.button>
  )
}