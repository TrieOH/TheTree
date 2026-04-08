import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { motion } from 'motion/react'
import {
  CalendarDays,
  Clock,
  MapPin,
  Users,
  Filter,
  CheckCircle2,
  Lock,
  ArrowRight,
  ShieldCheck,
  Coins,
} from 'lucide-react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { ActivityI, AttendanceRecordI } from '@/features/activities/model'
import { cn } from '@/shared/lib/utils'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import {
  allActivitiesQueryOptions,
  registerUserInActivityFn,
  unregisterUserInActivityFn,
} from '@/features/activities/api'
import { usePermissions } from '@/features/auths/hooks/use-permissions'
import {
  canAnnounceActivity,
  canCreateActivity,
  canManageActivity,
  canPublishActivity,
} from '@/features/activities/model/permissions'
import ActivityCard, {
  difficultyConfig,
  formatDuration,
  formatTime,
} from '@/features/activities/ui/ActivityCard'

export const Route = createLazyFileRoute(
  '/events/$eventId/editions/$editionId/activities/'
)({
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const isAuthenticated = Route.useRouteContext().auth?.isAuthenticated
  const auth = Route.useRouteContext().auth?.auth
  const userProfile = auth?.profile()

  const [difficultyFilter, setDifficultyFilter] = useState<ActivityI['difficulty'] | 'all'>('all')
  const [selectedActivity, setSelectedActivity] = useState<ActivityI | null>(null)
  const [isFilterOpen, setIsFilterOpen] = useState(false)

  const { data: activities = [] } = useQuery(allActivitiesQueryOptions(eventId, editionId))

  // mock
  const myRegistrations: AttendanceRecordI[] = []

  const { some: somePerms } = usePermissions(
    { canCreateActivity, canManageActivity, canPublishActivity, canAnnounceActivity },
    userProfile?.id
  )

  const isAdmin = somePerms('canCreateActivity', 'canManageActivity', 'canPublishActivity', 'canAnnounceActivity')

  const registerMutation = useMutation({
    mutationFn: (activityId: string) => registerUserInActivityFn(eventId, editionId, activityId),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Inscrição realizada com sucesso!')
        setSelectedActivity(null)
      } else {
        toast.error(res.message || 'Erro ao realizar inscrição')
      }
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const unregisterMutation = useMutation({
    mutationFn: (activityId: string) => unregisterUserInActivityFn(eventId, editionId, activityId),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Inscrição cancelada com sucesso!')
        setSelectedActivity(null)
      } else toast.error(res.message || 'Erro ao cancelar inscrição')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const filteredActivities = useMemo(() => {
    return activities
      .filter(a => a.status === 'published' || a.status === 'ongoing')
      .filter(a => {
        if (difficultyFilter === 'all') return true
        return a.difficulty === difficultyFilter
      })
      .sort((a, b) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime())
  }, [activities, difficultyFilter])

  const groupedActivities = useMemo(() => {
    const groups: Record<string, ActivityI[]> = {}
    filteredActivities.forEach(activity => {
      const date = new Date(activity.starts_at).toLocaleDateString('pt-BR', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
      })
      groups[date] ??= []
      groups[date].push(activity)
    })
    return groups
  }, [filteredActivities])

  const isRegistered = (activityId: string) => {
    return myRegistrations.some(r => r.activity_id === activityId)
  }

  const canRegister = (activity: ActivityI) => {
    if (!isAuthenticated) return false
    if (activity.status !== 'published') return false
    if (activity.has_capacity && activity.remaining_capacity <= 0) return false
    return true
  }

  const handleRegister = () => {
    if (!selectedActivity) return
    registerMutation.mutate(selectedActivity.id)
  }

  const handleUnregister = () => {
    if (!selectedActivity) return
    unregisterMutation.mutate(selectedActivity.id)
  }

  const loading = registerMutation.isPending || unregisterMutation.isPending

  return (
    <div className="min-h-screen bg-background pb-24">
      {/* Header */}
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <h1 className="text-lg md:text-xl font-semibold text-foreground">
              Programação
              <span className="ml-2 text-sm font-normal text-muted-foreground">
                ({filteredActivities.length})
              </span>
            </h1>

            <div className="flex items-center gap-2">
              <Drawer open={isFilterOpen} onOpenChange={setIsFilterOpen}>
                <DrawerTrigger asChild>
                  <button className="flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted transition-colors relative">
                    <Filter className="w-4 h-4 text-muted-foreground" />
                    {difficultyFilter !== 'all' && (
                      <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-primary" />
                    )}
                  </button>
                </DrawerTrigger>
                <DrawerContent className="z-60 rounded-t-2xl">
                  <DrawerHeader className="pb-4 border-b">
                    <DrawerTitle className="text-base font-semibold">Filtrar por nível</DrawerTitle>
                  </DrawerHeader>
                  <div className="p-4 pb-8 space-y-2">
                    <button
                      onClick={() => { setDifficultyFilter('all'); setIsFilterOpen(false) }}
                      className={cn(
                        "w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-colors",
                        difficultyFilter === 'all' ? "bg-primary/10 text-primary" : "hover:bg-muted"
                      )}
                    >
                      <span className="font-medium">Todos os níveis</span>
                      {difficultyFilter === 'all' && <CheckCircle2 className="w-4 h-4 ml-auto" />}
                    </button>
                    {Object.entries(difficultyConfig).map(([key, config]) => (
                      <button
                        key={key}
                        onClick={() => { setDifficultyFilter(key as ActivityI['difficulty']); setIsFilterOpen(false) }}
                        className={cn(
                          "w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-colors",
                          difficultyFilter === key
                            ? "bg-primary/10 text-primary"
                            : "hover:bg-muted"
                        )}
                      >
                        <span className={cn("w-2 h-2 rounded-full", config.accent)} />
                        <span className="font-medium">{config.label}</span>
                        {difficultyFilter === key && <CheckCircle2 className="w-4 h-4 ml-auto" />}
                      </button>
                    ))}
                  </div>
                </DrawerContent>
              </Drawer>
              {isAdmin && (
                <Link
                  to="/admin/events/$eventId/editions/$editionId/activities"
                  params={{ eventId, editionId }}
                  className={cn(
                    "group relative flex items-center justify-center",
                    "w-9 h-9 rounded-lg transition-all duration-200",
                    "text-muted-foreground hover:text-foreground",
                    "hover:bg-muted active:bg-muted/60",
                    "shrink-0"
                  )}
                  aria-label="Painel administrativo"
                >
                  <ShieldCheck className="w-5 h-5 transition-transform duration-200 group-hover:scale-110" />
                  <span className={cn(
                    "pointer-events-none absolute -bottom-9 right-0",
                    "whitespace-nowrap rounded-md px-2 py-1",
                    "bg-popover text-popover-foreground border border-border",
                    "text-xs shadow-md",
                    "opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0",
                    "transition-all duration-150"
                  )}>
                    Modo admin
                  </span>
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <main className="max-w-3xl mx-auto px-4 py-4">
        {filteredActivities.length === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="flex flex-col items-center justify-center py-16 space-y-4"
          >
            <div className="w-16 h-16 rounded-2xl bg-muted flex items-center justify-center">
              <CalendarDays className="w-8 h-8 text-muted-foreground/30" />
            </div>
            <div className="text-center">
              <h3 className="font-medium text-foreground">Nenhuma atividade encontrada</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Tente ajustar os filtros ou buscar por outro termo
              </p>
            </div>
          </motion.div>
        ) : (
          <div className="space-y-6">
            {Object.entries(groupedActivities).map(([date, dateActivities], groupIdx) => (
              <motion.div
                key={date}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: groupIdx * 0.1 }}
              >
                <h2 className="sticky top-14 z-20 -mx-4 mb-2 bg-background/80 px-4 py-2 text-[10px] font-black uppercase tracking-widest text-muted-foreground backdrop-blur-sm">
                  {date}
                </h2>

                <div className="space-y-3">
                  {dateActivities.map((activity, idx) => (
                    <ActivityCard
                      key={activity.id}
                      activity={activity}
                      index={idx}
                      registered={isRegistered(activity.id)}
                      onClick={() => { setSelectedActivity(activity); }}
                    />
                  ))}
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </main>

      {/* Activity Detail Drawer */}
      <Drawer open={!!selectedActivity} onOpenChange={() => { setSelectedActivity(null); }}>
        <DrawerContent className="z-60 rounded-t-3xl max-h-[90vh]">
          {selectedActivity && (
            <>
              <DrawerHeader className="pb-4 border-b">
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <span className={cn(
                      "inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-medium mb-2",
                      difficultyConfig[selectedActivity.difficulty].pill,
                      difficultyConfig[selectedActivity.difficulty].pillText
                    )}>
                      {difficultyConfig[selectedActivity.difficulty].label}
                    </span>
                    <DrawerTitle className="text-lg font-bold text-left">
                      {selectedActivity.title}
                    </DrawerTitle>
                    {selectedActivity.presenter_name && (
                      <p className="text-sm text-muted-foreground mt-1 text-left">
                        por {selectedActivity.presenter_name}
                      </p>
                    )}
                  </div>

                  {isRegistered(selectedActivity.id) && (
                    <span className="inline-flex items-center gap-1 px-3 py-1.5 rounded-full text-xs font-medium bg-primary/10 text-primary border border-primary/20 shrink-0">
                      <CheckCircle2 className="w-3.5 h-3.5" />
                      Inscrito
                    </span>
                  )}
                </div>
              </DrawerHeader>

              <div className="p-4 space-y-4 overflow-y-auto">
                {/* Info cards */}
                <div className="grid grid-cols-2 gap-3">
                  <div className="bg-muted rounded-xl p-3">
                    <div className="flex items-center gap-2 text-muted-foreground mb-1">
                      <Clock className="w-4 h-4" />
                      <span className="text-xs font-medium">Horário</span>
                    </div>
                    <p className="text-sm font-semibold">
                      {formatTime(selectedActivity.starts_at)} – {formatTime(selectedActivity.ends_at)}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {formatDuration(selectedActivity.starts_at, selectedActivity.ends_at)}
                    </p>
                  </div>

                  <div className="bg-muted rounded-xl p-3">
                    <div className="flex items-center gap-2 text-muted-foreground mb-1">
                      <MapPin className="w-4 h-4" />
                      <span className="text-xs font-medium">Local</span>
                    </div>
                    <p className="text-sm font-semibold truncate">
                      {selectedActivity.location}
                    </p>
                  </div>
                </div>

                {/* Capacity */}
                {selectedActivity.has_capacity && (
                  <div className="bg-muted rounded-xl p-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Users className="w-4 h-4" />
                        <span className="text-xs font-medium">Vagas disponíveis</span>
                      </div>
                      <span className="text-sm font-semibold">
                        {selectedActivity.remaining_capacity} / {selectedActivity.capacity}
                      </span>
                    </div>
                    <div className="w-full h-2 bg-background rounded-full mt-2 overflow-hidden">
                      <div
                        className={cn(
                          "h-full rounded-full transition-all",
                          selectedActivity.remaining_capacity / selectedActivity.capacity > 0.3
                            ? "bg-emerald-500"
                            : selectedActivity.remaining_capacity > 0
                              ? "bg-amber-500"
                              : "bg-red-500"
                        )}
                        style={{
                          width: `${(selectedActivity.remaining_capacity / selectedActivity.capacity) * 100}%`
                        }}
                      />
                    </div>
                  </div>
                )}

                {/* Cost */}
                {selectedActivity.token_cost > 0 && (
                  <div className="flex items-center gap-3 p-3 bg-amber-500/10 rounded-xl border border-amber-500/20">
                    <div className="w-10 h-10 rounded-lg bg-amber-500/20 flex items-center justify-center">
                      <Coins className="w-5 h-5 text-amber-600" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-amber-800">Custo em tokens</p>
                      <p className="text-lg font-bold text-amber-900">
                        {selectedActivity.token_cost} tokens
                      </p>
                    </div>
                  </div>
                )}

                {/* Description */}
                {selectedActivity.description && (
                  <div>
                    <h4 className="text-xs font-semibold tracking-wider uppercase text-muted-foreground mb-2">
                      Sobre a atividade
                    </h4>
                    <p className="text-sm text-foreground/80 leading-relaxed whitespace-pre-wrap">
                      {selectedActivity.description}
                    </p>
                  </div>
                )}

                {/* Action button */}
                <div className="pt-2">
                  {!isAuthenticated ? (
                    <Link
                      to="/auth"
                      search={{ redirect: `/events/${eventId}/editions/${editionId}/activities` }}
                      className={cn(
                        "flex items-center justify-center gap-2 w-full py-3.5 rounded-xl",
                        "bg-primary text-primary-foreground font-medium",
                        "hover:bg-primary/90 transition-colors"
                      )}
                    >
                      <Lock className="w-4 h-4" />
                      Faça login para se inscrever
                    </Link>
                  ) : isRegistered(selectedActivity.id) ? (
                    <button
                      onClick={handleUnregister}
                      disabled={loading}
                      className={cn(
                        "flex items-center justify-center gap-2 w-full py-3.5 rounded-xl",
                        "bg-red-500/10 text-red-600 border border-red-500/20 font-medium",
                        "hover:bg-red-500/20 transition-colors",
                        "disabled:opacity-50 disabled:cursor-not-allowed"
                      )}
                    >
                      {loading ? (
                        <div className="w-5 h-5 border-2 border-red-600/30 border-t-red-600 rounded-full animate-spin" />
                      ) : (
                        "Cancelar inscrição"
                      )}
                    </button>
                  ) : selectedActivity.has_capacity && selectedActivity.remaining_capacity <= 0 ? (
                    <button
                      disabled
                      className="flex items-center justify-center gap-2 w-full py-3.5 rounded-xl bg-muted text-muted-foreground font-medium cursor-not-allowed"
                    >
                      Vagas esgotadas
                    </button>
                  ) : (
                    <button
                      onClick={handleRegister}
                      disabled={loading || !canRegister(selectedActivity)}
                      className={cn(
                        "flex items-center justify-center gap-2 w-full py-3.5 rounded-xl",
                        "bg-primary text-primary-foreground font-medium",
                        "hover:bg-primary/90 transition-colors",
                        "disabled:opacity-50 disabled:cursor-not-allowed"
                      )}
                    >
                      {loading ? (
                        <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      ) : selectedActivity.status !== 'published' ? (
                        "Inscrições fechadas"
                      ) : (
                        <>
                          Realizar inscrição
                          <ArrowRight className="w-4 h-4" />
                        </>
                      )}
                    </button>
                  )}
                </div>
              </div>
            </>
          )}
        </DrawerContent>
      </Drawer>
    </div>
  )
}