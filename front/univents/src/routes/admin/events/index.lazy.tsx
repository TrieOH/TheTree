import { createLazyFileRoute } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Calendar, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import type { EventCreateSubmitI, EventI } from '@/features/events/model'
import {
  createEventFn,
  eventsQueryOptions,
  ownEventsQueryOptions,
  patchEventFn,
  publishEventFn,
} from '@/features/events/api'
import AdminEventCard from '@/features/events/ui/AdminEventCard'
import { ManageEventModal } from '@/features/events/ui/ManageEventModal'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { Button } from '@/shared/ui/shadcn/button'

export const Route = createLazyFileRoute('/admin/events/')({
  component: RouteComponent,
})

const STATUS_SORT_ORDER: Record<EventI['status'], number> = {
  draft: 0,
  active: 1,
  archived: 2,
  discontinued: 3,
}

function RouteComponent() {
  const queryClient = useQueryClient()
  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<EventI>>({
    field: 'created_at',
    direction: 'desc',
  })
  const [modalState, setModalState] = useState<{ open: boolean; event?: EventI }>({
    open: false,
  })
  const [publishingEvent, setPublishingEvent] = useState<EventI | null>(null)

  const { data: events = [] } = useQuery(ownEventsQueryOptions())

  const createMutation = useMutation({
    mutationFn: createEventFn,
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao criar evento')
        return
      }

      queryClient.setQueryData<EventI[]>(
        ownEventsQueryOptions().queryKey,
        (old = []) => [res.data, ...old],
      )
      void queryClient.invalidateQueries({ queryKey: eventsQueryOptions().queryKey })
      toast.success('Evento criado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const patchMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<EventCreateSubmitI> }) =>
      patchEventFn(id, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao atualizar evento')
        return
      }

      queryClient.setQueryData<EventI[]>(
        ownEventsQueryOptions().queryKey,
        (old = []) => old.map((event) => (event.id === res.data.id ? res.data : event)),
      )
      void queryClient.invalidateQueries({ queryKey: eventsQueryOptions().queryKey })
      toast.success('Evento atualizado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const publishMutation = useMutation({
    mutationFn: publishEventFn,
    onSuccess: (res, eventId) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar evento')
        return
      }

      queryClient.setQueryData<EventI[]>(
        ownEventsQueryOptions().queryKey,
        (old = []) =>
          old.map((event) =>
            event.id === eventId ? { ...event, status: 'active' as const } : event,
          ),
      )
      void queryClient.invalidateQueries({ queryKey: eventsQueryOptions().queryKey })
      setPublishingEvent(null)
      toast.success('Evento publicado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const filteredEvents = [...events]
    .filter((event) => {
      const search = filter.trim().toLowerCase()
      if (!search) return true

      return [
        event.name,
        event.slug,
        event.acronym ?? '',
        event.contact_email,
        event.status,
      ].some((value) => value.toLowerCase().includes(search))
    })
    .sort((a, b) => {
      const direction = sort.direction === 'asc' ? 1 : -1

      if (sort.field === 'created_at') {
        return (new Date(a.created_at).getTime() - new Date(b.created_at).getTime()) * direction
      }

      if (sort.field === 'editions_count') {
        return (a.editions_count - b.editions_count) * direction
      }

      if (sort.field === 'status') {
        return (STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status]) * direction
      }

      return String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
    })

  return (
    <div className="flex flex-wrap p-4">
      <PaginatedContainer<EventI>
        items={filteredEvents}
        layout="grid"
        minItemWidth="16rem"
        pageSize={4}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: 'created_at',
            label: 'Criado em',
            comparator: (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
          },
          { key: 'name', label: 'Nome' },
          { key: 'slug', label: 'Slug' },
          {
            key: 'status',
            label: 'Status',
            comparator: (a, b) => STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
          },
          {
            key: 'editions_count',
            label: 'Edições',
            comparator: (a, b) => a.editions_count - b.editions_count,
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, slug, sigla ou e-mail..."
        itemLabel="eventos"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalState({ open: true, event: undefined })}
            size="sm"
            className="rounded-sm py-4 gap-2"
          >
            <Plus className="w-4 h-4" />
            Novo evento
          </Button>
        }
        emptyState={
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <div className="mb-4 flex size-16 items-center justify-center rounded-2xl border border-border bg-muted/40">
              <Calendar className="size-8 text-muted-foreground/40" />
            </div>
            <h3 className="text-base font-medium text-foreground">
              Nenhum evento encontrado
            </h3>
            <p className="mt-1 max-w-sm text-sm text-muted-foreground">
              Crie um evento para começar a organizar o dashboard do admin.
            </p>
            <Button
              type="button"
              onClick={() => setModalState({ open: true, event: undefined })}
              size="sm"
              className="mt-5 rounded-sm gap-2"
            >
              <Plus className="w-4 h-4" />
              Criar evento
            </Button>
          </div>
        }
        renderItems={(slice) =>
          slice.map((event, idx) => (
            <AdminEventCard
              key={event.id}
              event={event}
              index={idx}
              onEdit={(selectedEvent) => setModalState({ open: true, event: selectedEvent })}
              onPublish={setPublishingEvent}
            />
          ))
        }
      />

      <ManageEventModal
        key={modalState.event?.id ?? 'event-create'}
        open={modalState.open}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }))
            return
          }

          setModalState({ open: false, event: undefined })
        }}
        event={modalState.event}
        onCreate={(values) =>
          createMutation.mutateAsync(values).then(
            (res) => (res.success ? res.data : false),
            () => false,
          )
        }
        onUpdate={(id, values) =>
          patchMutation.mutateAsync({ id, data: values }).then(
            (res) => (res.success ? res.data : false),
            () => false,
          )
        }
      />

      <AlertModal
        open={!!publishingEvent}
        onOpenChange={() => setPublishingEvent(null)}
        title="Publicar evento?"
        description={`Ao publicar "${publishingEvent?.name}", ele ficará visível para o público.`}
        confirmLabel="Publicar"
        onConfirm={() => {
          if (!publishingEvent) return
          publishMutation.mutate(publishingEvent.id)
        }}
        variant="success"
        loading={publishMutation.isPending}
      />
    </div>
  )
}
