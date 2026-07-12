import { createLazyFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { Calendar, Plus } from 'lucide-react'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import type { EventCreateSubmitI, EventI } from '@/features/events/model'
import { allOwnEventsQueryOptions } from '@/features/events/api'
import {
  useCreateEventMutation,
  usePatchEventMutation,
  usePublishEventMutation,
} from '@/features/events/api/mutations'
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
  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<EventI>>({
    field: 'created_at',
    direction: 'desc',
  })
  const [modalState, setModalState] = useState<{ open: boolean; event?: EventI }>({
    open: false,
  })
  const [publishingEvent, setPublishingEvent] = useState<EventI | null>(null)

  const { data: events = [] } = useQuery(allOwnEventsQueryOptions())
  const createMutation = useCreateEventMutation()
  const patchMutation = usePatchEventMutation()
  const publishEventMutation = usePublishEventMutation()

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
    <div className="flex flex-wrap p-6">
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
          <EmptyState
            icon={Calendar}
            eyebrow="Eventos"
            title="Nenhum evento encontrado"
            description="Crie um evento para começar a organizar o dashboard do admin."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
            action={
              <Button
                type="button"
                onClick={() => setModalState({ open: true, event: undefined })}
                size="sm"
                className="mt-0.5 h-9 rounded-sm gap-2 px-4 text-sm shadow-sm"
              >
                <Plus className="w-4 h-4" />
                Criar evento
              </Button>
            }
          />
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
          patchMutation.mutateAsync({ id, data: values as Partial<EventCreateSubmitI> }).then(
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
          publishEventMutation.mutate(publishingEvent.id)
        }}
        variant="success"
        loading={publishEventMutation.isPending}
      />
    </div>
  )
}
