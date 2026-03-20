import { motion } from 'motion/react'
import { Clock, MapPin, User, Coins, Users } from 'lucide-react'
import type { ActivityI } from '@/features/activities/model'
import { cn } from '@/shared/lib/utils'

interface ActivitiesTabProps {
  activities: ActivityI[]
}

const difficultyLabels: Record<string, string> = {
  no_prerequisites: 'Sem pré-requisitos',
  beginner: 'Iniciante',
  intermediate: 'Intermediário',
  advanced: 'Avançado',
  expert: 'Especialista',
}

const difficultyColors: Record<string, string> = {
  no_prerequisites: 'text-green-600 bg-green-500/10',
  beginner: 'text-blue-600 bg-blue-500/10',
  intermediate: 'text-yellow-600 bg-yellow-500/10',
  advanced: 'text-orange-600 bg-orange-500/10',
  expert: 'text-red-600 bg-red-500/10',
}

export function ActivitiesTab({ activities }: ActivitiesTabProps) {
  const formatDateTime = (date: string) => {
    return new Date(date).toLocaleString('pt-BR', {
      day: '2-digit',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const publishedActivities = activities.filter(a => a.status === 'published')
  const draftActivities = activities.filter(a => a.status === 'draft')

  return (
    <div className="space-y-6">
      {publishedActivities.length === 0 && draftActivities.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          Nenhuma atividade cadastrada.
        </div>
      ) : (
        <>
          {publishedActivities.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                Atividades Confirmadas
              </h3>
              <div className="grid gap-4">
                {publishedActivities.map((activity, idx) => (
                  <motion.div
                    key={activity.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.05 }}
                    className="bg-card border border-border rounded-xl p-4 md:p-5 space-y-4 hover:border-primary/20 transition-colors"
                  >
                    <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
                      <div className="space-y-2">
                        <h4 className="text-lg font-semibold text-foreground">{activity.title}</h4>
                        {activity.description && (
                          <p className="text-sm text-muted-foreground line-clamp-2">{activity.description}</p>
                        )}
                        <div className="flex flex-wrap gap-2">
                          <span className={cn("px-2 py-1 rounded-full text-xs font-medium", difficultyColors[activity.difficulty])}>
                            {difficultyLabels[activity.difficulty]}
                          </span>
                          {activity.token_cost > 0 && (
                            <span className="px-2 py-1 rounded-full text-xs font-medium bg-amber-500/10 text-amber-600 flex items-center gap-1">
                              <Coins className="w-3 h-3" />
                              {activity.token_cost} tokens
                            </span>
                          )}
                          {activity.token_cost === 0 && (
                            <span className="px-2 py-1 rounded-full text-xs font-medium bg-green-500/10 text-green-600">
                              Gratuito
                            </span>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="flex flex-wrap gap-4 text-sm text-muted-foreground pt-2 border-t border-border">
                      <div className="flex items-center gap-1.5">
                        <Clock className="w-4 h-4" />
                        <span>{formatDateTime(activity.starts_at)}</span>
                      </div>
                      <div className="flex items-center gap-1.5">
                        <MapPin className="w-4 h-4" />
                        <span>{activity.location}</span>
                      </div>
                      {activity.presenter_name && (
                        <div className="flex items-center gap-1.5">
                          <User className="w-4 h-4" />
                          <span>{activity.presenter_name}</span>
                        </div>
                      )}
                      {activity.has_capacity && (
                        <div className="flex items-center gap-1.5">
                          <Users className="w-4 h-4" />
                          <span>{activity.remaining_capacity} de {activity.capacity} vagas</span>
                          <div className="w-16 h-1.5 bg-muted rounded-full overflow-hidden">
                            <div
                              className="h-full bg-primary rounded-full"
                              style={{ width: `${((activity.capacity - activity.remaining_capacity) / activity.capacity) * 100}%` }}
                            />
                          </div>
                        </div>
                      )}
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          )}

          {draftActivities.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                Em Breve
              </h3>
              <div className="grid gap-4 opacity-60">
                {draftActivities.map((activity) => (
                  <div
                    key={activity.id}
                    className="bg-muted/50 border border-border rounded-xl p-4 md:p-5"
                  >
                    <h4 className="text-lg font-medium text-foreground">{activity.title}</h4>
                    <p className="text-sm text-muted-foreground">Em breve mais informações</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  )
}