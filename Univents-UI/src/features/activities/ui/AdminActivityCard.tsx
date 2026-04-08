import { motion } from 'motion/react'
import {
  MapPin,
  Users,
  Edit3,
  UserCheck,
  Eye,
  MoreVertical,
  Trash2,
  Copy,
  Coins,
} from 'lucide-react'
import { formatTime, formatDuration, difficultyConfig } from './ActivityCard'
import type { ActivityI } from '@/features/activities/model'
import { cn } from '@/shared/lib/utils'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'

export const statusConfig: Record<
  ActivityI['status'],
  { label: string; dot: string; bg: string; text: string }
> = {
  draft: {
    label: 'Rascunho',
    dot: 'bg-slate-400',
    bg: 'bg-slate-50',
    text: 'text-slate-600',
  },
  published: {
    label: 'Publicado',
    dot: 'bg-blue-500',
    bg: 'bg-blue-50',
    text: 'text-blue-700',
  },
  ongoing: {
    label: 'Em andamento',
    dot: 'bg-emerald-500',
    bg: 'bg-emerald-50',
    text: 'text-emerald-700',
  },
  completed: {
    label: 'Concluído',
    dot: 'bg-slate-900',
    bg: 'bg-slate-100',
    text: 'text-slate-900',
  },
  canceled: {
    label: 'Cancelado',
    dot: 'bg-red-500',
    bg: 'bg-red-50',
    text: 'text-red-700',
  },
}

interface AdminActivityCardProps {
  activity: ActivityI
  index?: number
  onEdit: (activity: ActivityI) => void
  onPublish: (activity: ActivityI) => void
  onAttendance: (activity: ActivityI) => void
  onDelete: (activity: ActivityI) => void
  onDuplicate: (activity: ActivityI) => void
}

export default function AdminActivityCard({
  activity,
  index = 0,
  onEdit,
  onPublish,
  onAttendance,
  onDelete,
  onDuplicate,
}: AdminActivityCardProps) {
  const difficulty = difficultyConfig[activity.difficulty]
  const status = statusConfig[activity.status]

  return (
    <motion.div
      initial={{ opacity: 0, x: -10 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.05 }}
      className={cn(
        'grid w-full grid-cols-[3px_52px_1fr] gap-x-3 text-left',
        'rounded-lg border bg-card p-3.5 transition-colors',
        'hover:border-border/60 active:scale-[0.99]',
      )}
    >
      <span
        className={cn(
          'w-0.75 self-stretch rounded-full',
          difficulty.accent
        )}
      />

      {/* Time Column */}
      <div className="flex flex-col items-end border-r border-border/40 pr-3 pt-0.5">
        <span className="text-sm font-medium leading-tight text-foreground">
          {formatTime(activity.starts_at)}
        </span>
        <span className="mt-0.5 text-[11px] text-muted-foreground">
          {formatDuration(activity.starts_at, activity.ends_at)}
        </span>
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-3">
          <div className="space-y-1 min-w-0">
            <div className="flex flex-wrap items-center gap-2">
              <span className={cn(
                "inline-flex items-center gap-1 px-1.5 py-0.5 rounded-md text-[9px] font-bold uppercase tracking-wider",
                status.bg,
                status.text
              )}>
                <span className={cn("w-1 h-1 rounded-full", status.dot)} />
                {status.label}
              </span>

              <span className={cn(
                "px-1.5 py-0.5 rounded-md text-[9px] font-bold uppercase tracking-wider",
                difficulty.pill,
                difficulty.pillText
              )}>
                {difficulty.label}
              </span>

              {activity.token_cost > 0 && (
                <span className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded-md text-[9px] font-bold uppercase tracking-wider bg-amber-50 text-amber-700">
                  <Coins className="w-2.5 h-2.5" />
                  {activity.token_cost}
                </span>
              )}
            </div>

            <h3 className="font-bold text-base text-foreground leading-tight truncate pr-4">
              {activity.title}
            </h3>

            {activity.presenter_name && (
              <p className="text-sm text-muted-foreground font-medium truncate">
                {activity.presenter_name}
              </p>
            )}
          </div>

          {/* Actions & Menu */}
          <div className="flex items-center gap-2 shrink-0">
            {/* Primary Action */}
            {activity.status === 'draft' && (
              <button
                onClick={() => { onPublish(activity); }}
                className={cn(
                  "hidden sm:flex px-3 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider",
                  "bg-primary/10 text-primary hover:bg-primary/20",
                  "transition-all active:scale-95"
                )}
              >
                Publicar
              </button>
            )}

            {/* Options Menu */}
            <Drawer>
              <DrawerTrigger asChild>
                <button className="flex items-center justify-center w-9 h-9 rounded-xl text-muted-foreground hover:bg-muted hover:text-foreground transition-all active:scale-90 outline-none">
                  <MoreVertical className="w-5 h-5" />
                </button>
              </DrawerTrigger>
              <DrawerContent className="z-60 rounded-t-3xl">
                <DrawerHeader className="pb-4 border-b">
                  <DrawerTitle className="text-base font-bold text-center">Gerenciar Atividade</DrawerTitle>
                </DrawerHeader>
                <div className="p-4 pb-12 space-y-2">
                  {activity.status === 'draft' && (
                    <button
                      onClick={() => { onPublish(activity); }}
                      className="w-full sm:hidden flex items-center gap-3 px-4 py-4 rounded-2xl bg-primary/10 text-primary active:bg-primary/20 transition-colors"
                    >
                      <div className="w-9 h-9 rounded-xl bg-primary/20 flex items-center justify-center">
                        <Eye className="w-5 h-5" />
                      </div>
                      <span className="font-bold text-sm">Publicar atividade</span>
                    </button>
                  )}

                  <button
                    onClick={() => { onEdit(activity); }}
                    className="w-full flex items-center gap-3 px-4 py-4 rounded-2xl hover:bg-muted active:bg-muted/80 transition-colors"
                  >
                    <div className="w-9 h-9 rounded-xl bg-primary/10 flex items-center justify-center">
                      <Edit3 className="w-5 h-5 text-primary" />
                    </div>
                    <span className="font-bold text-sm">Editar detalhes</span>
                  </button>

                  <button
                    onClick={() => { onAttendance(activity); }}
                    className="w-full flex items-center gap-3 px-4 py-4 rounded-2xl hover:bg-muted active:bg-muted/80 transition-colors"
                  >
                    <div className="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center">
                      <UserCheck className="w-5 h-5 text-blue-600" />
                    </div>
                    <span className="font-bold text-sm">Lista de presença</span>
                  </button>

                  <div className="grid grid-cols-2 gap-2 mt-2 pt-4 border-t border-border/50">
                    <button
                      onClick={() => { onDuplicate(activity); }}
                      className="flex flex-col items-center justify-center gap-2 p-4 rounded-2xl hover:bg-muted transition-colors"
                    >
                      <div className="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center">
                        <Copy className="w-5 h-5 text-slate-600" />
                      </div>
                      <span className="font-bold text-xs text-slate-600">Duplicar</span>
                    </button>

                    <button
                      onClick={() => { onDelete(activity); }}
                      className="flex flex-col items-center justify-center gap-2 p-4 rounded-2xl hover:bg-red-50 text-red-600 transition-colors"
                    >
                      <div className="w-10 h-10 rounded-xl bg-red-100 flex items-center justify-center">
                        <Trash2 className="w-5 h-5 text-destructive" />
                      </div>
                      <span className="font-bold text-xs text-destructive">Excluir</span>
                    </button>
                  </div>
                </div>
              </DrawerContent>
            </Drawer>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-3">
          <div className="flex items-center gap-1.5 text-[11px] font-medium text-muted-foreground/80">
            <MapPin className="w-3 h-3" />
            {activity.location}
          </div>

          {activity.has_capacity && (
            <div className="flex items-center gap-1.5 text-[11px] font-medium text-muted-foreground/80">
              <Users className="w-3 h-3" />
              {activity.remaining_capacity} / {activity.capacity} vagas
            </div>
          )}
        </div>
      </div>
    </motion.div>
  )
}