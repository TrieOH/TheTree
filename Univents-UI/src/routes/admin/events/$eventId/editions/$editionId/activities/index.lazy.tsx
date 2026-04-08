import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import {
  Plus,
  Calendar,
  MoreVertical,
  CheckSquare,
  ShieldCheck,
} from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { ActivityCreateI, ActivityI } from '@/features/activities/model'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { cn } from '@/shared/lib/utils'
import { FormDrawer } from '@/widgets/form/ui/form-drawer'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { activityCreateSchema } from '@/features/activities/model'
import { getActivityFields } from '@/features/activities/model/field'
import {
  allAdminActivitiesQueryOptions,
  createActivityFn,
  publishActivityFn,
  updateActivityFn,
  allActivityAttendanceRecordsQueryOptions,
  allActivitiesQueryOptions,
} from '@/features/activities/api'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import AdminActivityCard from '@/features/activities/ui/AdminActivityCard'
import { getDirtyFields } from '@/shared/lib/diff'

export const Route = createLazyFileRoute(
  '/admin/events/$eventId/editions/$editionId/activities/'
)({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const { eventId, editionId } = Route.useParams()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [editingActivity, setEditingActivity] = useState<ActivityI | null>(null)
  const [deletingActivity, setDeletingActivity] = useState<ActivityI | null>(null)
  const [publishingActivity, setPublishingActivity] = useState<ActivityI | null>(null)
  const [viewingAttendance, setViewingAttendance] = useState<ActivityI | null>(null)
  const [isActionsOpen, setIsActionsOpen] = useState(false)

  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const edition = editions.find(e => e.id === editionId)

  const { data: activities = [] } = useQuery(allAdminActivitiesQueryOptions(eventId, editionId))

  const { data: attendanceRecords = [] } = useQuery(allActivityAttendanceRecordsQueryOptions(eventId, editionId, viewingAttendance?.id ?? ""))

  const createMutation = useMutation({
    mutationFn: (data: ActivityCreateI) => createActivityFn(data, eventId, editionId),
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData<ActivityI[]>(
          allAdminActivitiesQueryOptions(eventId, editionId).queryKey,
          (old) => [...(old ?? []), res.data]
        )
        setIsCreateOpen(false)
        setEditingActivity(null)
        toast.success('Atividade salva com sucesso!')
      } else toast.error(res.message || 'Erro ao salvar atividade')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  });

  const publishMutation = useMutation({
    mutationFn: (activityId: string) => publishActivityFn(eventId, editionId, activityId),
    onSuccess: (res, variables) => {
      if (res.success) {
        queryClient.setQueryData<ActivityI[]>(
          allAdminActivitiesQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((a) =>
            a.id === variables ? { ...a, status: 'published' as const } : a
          )
        )
        queryClient.setQueryData<ActivityI[]>(
          allActivitiesQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((a) =>
            a.id === variables ? { ...a, status: 'published' as const } : a
          )
        )
        setPublishingActivity(null)
        toast.success('Atividade publicada com sucesso!')
      } else toast.error(res.message || 'Erro ao publicar atividade')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: Partial<ActivityI> }) =>
      updateActivityFn(id, data, eventId, editionId),
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData<ActivityI[]>(
          allAdminActivitiesQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((a) => (a.id === res.data.id ? res.data : a))
        )
        queryClient.setQueryData<ActivityI[]>(
          allActivitiesQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((a) => (a.id === res.data.id ? res.data : a))
        )
        setEditingActivity(null) // Close the edit drawer
        toast.success('Atividade atualizada com sucesso!')
      } else toast.error(res.message || 'Erro ao atualizar atividade')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  });

  const duplicateMutation = useMutation({
    mutationFn: (activity: ActivityI) => {
      const duplicatedActivityData: ActivityCreateI = {
        ...activity,
        title: `${activity.title} (Cópia)`,
      };
      return createActivityFn(duplicatedActivityData, eventId, editionId);
    },
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData<ActivityI[]>(
          allAdminActivitiesQueryOptions(eventId, editionId).queryKey,
          (old) => [...(old ?? []), res.data]
        );
        toast.success('Atividade duplicada com sucesso!');
      } else toast.error(res.message || 'Erro ao duplicar atividade');
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  });

  const handleCreate = (data: ActivityCreateI) => {
    createMutation.mutate(data)
  }

  const handleUpdate = async (data: ActivityCreateI) => {
    if (!editingActivity) return

    const changes = getDirtyFields(data, editingActivity, [
      'title', 'capacity', 'description', 'difficulty', 'token_cost',
      'starts_at', 'ends_at', 'has_capacity', 'presenter_name', 'location'
    ])

    if (Object.keys(changes).length === 0) {
      toast.info('Nenhuma alteração detectada')
      setEditingActivity(null)
      return
    }

    await updateMutation.mutateAsync({ id: editingActivity.id, data })
  }

  const handlePublish = () => {
    if (!publishingActivity) return
    publishMutation.mutate(publishingActivity.id)
  }

  const loading = createMutation.isPending || publishMutation.isPending || updateMutation.isPending || duplicateMutation.isPending

  const groupedActivities = useMemo(() => {
    const sorted = [...activities].sort((a, b) =>
      new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime()
    )

    const groups: Record<string, ActivityI[]> = {}
    sorted.forEach(activity => {
      const date = new Date(activity.starts_at).toLocaleDateString('pt-BR', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
      })
      groups[date] ??= []
      groups[date].push(activity)
    })
    return groups
  }, [activities])

  return (
    <div className="min-h-screen bg-background relative pb-20 md:pb-0">
      <header className="sticky top-0 z-40 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <div className="flex items-center gap-2 shrink-0">
              <h1 className="text-lg md:text-xl font-semibold text-foreground">
                Atividades
                <span className="ml-2 text-sm font-normal text-muted-foreground">
                  ({activities.length})
                </span>
              </h1>
            </div>

            <div className="hidden sm:flex items-center gap-2 ml-auto">
              <button
                onClick={() => { setIsCreateOpen(true); }}
                className={cn(
                  "flex items-center gap-2 px-4 py-2 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium transition-all active:scale-95"
                )}
              >
                <Plus className="w-4 h-4" />
                Nova atividade
              </button>
            </div>

            <div className="sm:hidden flex items-center gap-1 ml-auto">
              <Drawer open={isActionsOpen} onOpenChange={setIsActionsOpen}>
                <DrawerTrigger asChild>
                  <button className={cn("flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted")}>
                    <MoreVertical className="w-5 h-5 text-foreground" />
                  </button>
                </DrawerTrigger>
                <DrawerContent className="z-60 rounded-t-2xl">
                  <DrawerHeader className="pb-4 border-b">
                    <DrawerTitle className="text-base font-semibold">Ações</DrawerTitle>
                  </DrawerHeader>
                  <div className="p-3 pb-8 space-y-1">
                    <button
                      onClick={() => { setIsActionsOpen(false); setIsCreateOpen(true) }}
                      className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                    >
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Plus className="w-4 h-4 text-primary" />
                      </div>
                      <span className="font-medium">Nova atividade</span>
                    </button>
                  </div>
                </DrawerContent>
              </Drawer>
            </div>

            <Link
              to="/events/$eventId/editions/$editionId/activities"
              params={{ eventId, editionId }}
              className={cn(
                "group relative flex items-center justify-center",
                "w-9 h-9 rounded-lg transition-all",
                "bg-primary text-primary-foreground",
                "hover:bg-primary/90",
                "shrink-0"
              )}
            >
              <ShieldCheck className="w-5 h-5" />
              <span
                className={cn(
                  "pointer-events-none absolute -bottom-9 right-0",
                  "whitespace-nowrap rounded-md px-2 py-1",
                  "bg-popover text-popover-foreground border text-xs shadow-md",
                  "opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0",
                  "transition-all"
                )}>
                Sair do admin
              </span>
            </Link>
          </div>
        </div>
      </header>

      <main className="max-w-3xl mx-auto px-4 py-4">
        <AnimatePresence mode="wait">
          {activities.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="flex flex-col items-center justify-center py-24 space-y-6"
            >
              <div className="w-20 h-20 rounded-3xl bg-muted flex items-center justify-center">
                <Calendar className="w-10 h-10 text-muted-foreground/30" />
              </div>
              <div className="text-center space-y-1">
                <h3 className="text-lg font-bold">Nenhuma atividade</h3>
                <p className="text-sm text-muted-foreground max-w-xs">
                  Comece criando a primeira atividade para {edition?.edition_name ?? 'esta edição'}.
                </p>
              </div>
              <button
                onClick={() => { setIsCreateOpen(true); }}
                className={cn(
                  "px-6 py-3 rounded-xl",
                  "bg-primary text-primary-foreground font-bold shadow-sm",
                  "active:scale-95 transition-all"
                )}
              >
                Criar Atividade
              </button>
            </motion.div>
          ) : (
            <div className="space-y-8">
              {Object.entries(groupedActivities).map(([date, dateActivities]) => (
                <div key={date} className="space-y-4">
                  <h2 className="sticky top-14 z-20 -mx-4 mb-2 bg-background/80 px-4 py-2 text-[10px] font-black uppercase tracking-widest text-muted-foreground backdrop-blur-sm">
                    {date}
                  </h2>
                  <div className="space-y-3">
                    {dateActivities.map((activity, idx) => (
                      <AdminActivityCard
                        key={activity.id}
                        activity={activity}
                        index={idx}
                        onEdit={setEditingActivity}
                        onPublish={setPublishingActivity}
                        onAttendance={setViewingAttendance}
                        onDelete={setDeletingActivity}
                        onDuplicate={(act) => {
                          duplicateMutation.mutate(act)
                        }}
                      />
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </AnimatePresence>
      </main>

      {/* Create Drawer */}
      <FormDrawer
        idPrefix="create-"
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        title="Nova atividade"
        fields={getActivityFields()}
        schema={activityCreateSchema}
        onSubmit={handleCreate}
        submitLabel="Criar atividade"
        loading={loading}
      />

      {/* Edit Drawer */}
      {editingActivity && (
        <FormDrawer
          idPrefix="edit-"
          open={!!editingActivity}
          onOpenChange={() => { setEditingActivity(null); }}
          title="Editar atividade"
          fields={getActivityFields()}
          defaultValues={editingActivity}
          schema={activityCreateSchema}
          onSubmit={handleUpdate}
          submitLabel="Salvar alterações"
          loading={loading}
        />
      )}

      {/* Delete Confirmation */}
      <AlertModal
        open={!!deletingActivity}
        onOpenChange={() => { setDeletingActivity(null); }}
        title="Excluir atividade?"
        description={`Tem certeza que deseja excluir "${deletingActivity?.title}"? Esta ação não pode ser desfeita.`}
        confirmLabel="Excluir"
        onConfirm={() => {
          toast.error('Funcionalidade não implementada ainda')
          setDeletingActivity(null)
        }}
        variant="destructive"
        loading={loading}
      />

      {/* Publish Confirmation */}
      <AlertModal
        open={!!publishingActivity}
        onOpenChange={() => { setPublishingActivity(null); }}
        title="Publicar atividade?"
        description={`Ao publicar "${publishingActivity?.title}", ela ficará visível para os participantes.`}
        confirmLabel="Publicar"
        onConfirm={handlePublish}
        variant="success"
        loading={loading}
      />

      {/* Attendance Drawer */}
      <Drawer open={!!viewingAttendance} onOpenChange={() => { setViewingAttendance(null); }}>
        <DrawerContent className="z-60 rounded-t-3xl max-h-[90vh]">
          {viewingAttendance && (
            <>
              <DrawerHeader className="pb-4 border-b">
                <DrawerTitle className="text-lg font-bold">
                  Lista de Presença
                </DrawerTitle>
                <p className="text-sm text-muted-foreground">
                  {viewingAttendance.title}
                </p>
              </DrawerHeader>
              <div className="p-4 overflow-y-auto">
                {attendanceRecords.length === 0 ? (
                  <div className="text-center py-12 text-muted-foreground space-y-4">
                    <div className="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center mx-auto opacity-50">
                      <CheckSquare className="w-8 h-8" />
                    </div>
                    <p className="font-medium">Nenhum registro de presença</p>
                  </div>
                ) : (
                  <div className="space-y-2">
                    {attendanceRecords.map((record) => (
                      <div key={record.id} className="flex items-center justify-between p-4 bg-muted/50 rounded-2xl border border-border/50">
                        <div className="flex flex-col">
                          <span className="text-sm font-bold">{record.user_id}</span>
                          <span className="text-[10px] text-muted-foreground font-medium uppercase tracking-wider">
                            {record.checked_in_at ? `Check-in: ${new Date(record.checked_in_at).toLocaleTimeString()}` : 'Aguardando'}
                          </span>
                        </div>
                        <span className={cn(
                          "text-[10px] font-black uppercase tracking-widest px-2.5 py-1 rounded-md",
                          record.status === 'checked_in' && "bg-emerald-50 text-emerald-600 border border-emerald-100",
                          record.status === 'registered' && "bg-blue-50 text-blue-600 border border-blue-100",
                          record.status === 'cancelled' && "bg-red-50 text-red-600 border border-red-100",
                        )}>
                          {record.status}
                        </span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </>
          )}
        </DrawerContent>
      </Drawer>
    </div>
  )
}