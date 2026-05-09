import { motion } from 'motion/react'
import { DoorOpen, MapPin, Coffee, Calendar, LogOut, Ticket, Users, Clock } from 'lucide-react'
import type { CheckpointI } from '@/features/checkpoints/model'
import { cn } from '@/shared/lib/utils'

interface CheckpointsTabProps {
  checkpoints: CheckpointI[]
}

const typeConfig = {
  entry: { icon: DoorOpen, label: 'Entrada', color: 'text-green-600 bg-green-500/10' },
  zone: { icon: MapPin, label: 'Zona', color: 'text-blue-600 bg-blue-500/10' },
  amenity: { icon: Coffee, label: 'Comodidade', color: 'text-amber-600 bg-amber-500/10' },
  session: { icon: Calendar, label: 'Sessão', color: 'text-purple-600 bg-purple-500/10' },
  exit: { icon: LogOut, label: 'Saída', color: 'text-red-600 bg-red-500/10' },
} as const

const accessConfig = {
  open: { label: 'Livre', icon: Users, color: 'text-green-600' },
  ticket: { label: 'Com Ticket', icon: Ticket, color: 'text-blue-600' },
  staff_only: { label: 'Staff Only', icon: Users, color: 'text-red-600' },
} as const

export function CheckpointsTab({ checkpoints }: CheckpointsTabProps) {
  const formatTime = (date: string | null) => {
    if (!date) return 'O tempo todo'
    return new Date(date).toLocaleString('pt-BR', {
      day: '2-digit',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <div className="space-y-6">
      {checkpoints.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          Nenhum checkpoint cadastrado.
        </div>
      ) : (
        <div className="grid gap-4">
          {checkpoints.map((checkpoint, idx) => {
            const TypeIcon = typeConfig[checkpoint.type].icon
            const AccessIcon = accessConfig[checkpoint.access_mode].icon

            return (
              <motion.div
                key={checkpoint.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: idx * 0.05 }}
                className="bg-card border border-border rounded-xl p-4 md:p-5 flex items-start gap-4 hover:border-primary/20 transition-colors"
              >
                <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center shrink-0", typeConfig[checkpoint.type].color)}>
                  <TypeIcon className="w-5 h-5" />
                </div>

                <div className="flex-1 space-y-2">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-2">
                    <h4 className="text-lg font-semibold text-foreground">{checkpoint.name}</h4>
                    <div className="flex items-center gap-2">
                      <span className={cn("text-xs font-medium px-2 py-1 rounded-full", typeConfig[checkpoint.type].color)}>
                        {typeConfig[checkpoint.type].label}
                      </span>
                      <span className={cn("text-xs font-medium px-2 py-1 rounded-full bg-muted flex items-center gap-1", accessConfig[checkpoint.access_mode].color)}>
                        <AccessIcon className="w-3 h-3" />
                        {accessConfig[checkpoint.access_mode].label}
                      </span>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1.5">
                      <Clock className="w-4 h-4" />
                      <span>Início: {formatTime(checkpoint.starts_at)}</span>
                    </div>
                    {checkpoint.ends_at && (
                      <div className="flex items-center gap-1.5">
                        <Clock className="w-4 h-4" />
                        <span>Término: {formatTime(checkpoint.ends_at)}</span>
                      </div>
                    )}
                  </div>
                </div>
              </motion.div>
            )
          })}
        </div>
      )}
    </div>
  )
}