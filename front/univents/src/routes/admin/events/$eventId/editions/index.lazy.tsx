import { createLazyFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Calendar, Plus } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import { Button } from '@/shared/ui/shadcn/button'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import {
  useCreateEditionMutation,
  usePublishEditionMutation,
  useUpdateEditionMutation,
} from '@/features/editions/api/mutations'
import type { EditionI } from '@/features/editions/model'
import { AdminEditionCard } from '@/features/editions/ui/AdminEditionCard'
import { ManageEditionModal } from '@/features/editions/ui/ManageEditionModal'

const STATUS_SORT_ORDER: Record<EditionI['status'], number> = {
  draft: 0,
  announced: 1,
  open: 2,
  ongoing: 3,
  finished: 4,
  completed: 5,
  cancelled: 6,
  postponed: 7,
}

export const Route = createLazyFileRoute('/admin/events/$eventId/editions/')({
  component: EditionsRoute,
})

function EditionsRoute() {
  const { eventId } = Route.useParams()
  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const createEditionMutation = useCreateEditionMutation()
  const updateEditionMutation = useUpdateEditionMutation()
  const publishEditionMutation = usePublishEditionMutation()
  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<EditionI>>({
    field: 'starts_at',
    direction: 'desc',
  })
  const [modalState, setModalState] = useState<{ open: boolean; edition?: EditionI }>({
    open: false,
  })
  const [publishingEdition, setPublishingEdition] = useState<EditionI | null>(null)

  const filteredEditions = [...editions]
    .filter((edition) => {
      const search = filter.trim().toLowerCase()
      if (!search) return true

      return [
        edition.edition_name,
        edition.tagline ?? '',
        edition.location_name,
        edition.location_address,
        edition.type,
        edition.status,
      ].some((value) => value.toLowerCase().includes(search))
    })
    .sort((a, b) => {
      const direction = sort.direction === 'asc' ? 1 : -1

      if (sort.field === 'starts_at') {
        return (new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime()) * direction
      }

      if (sort.field === 'status') {
        return (STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status]) * direction
      }

      if (sort.field === 'type') {
        return String(a.type).localeCompare(String(b.type)) * direction
      }

      return String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
    })

  return (
    <div className="flex flex-wrap">
      <PaginatedContainer<EditionI>
        items={filteredEditions}
        layout="grid"
        minItemWidth="16rem"
        pageSize={4}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: 'starts_at',
            label: 'Início',
            comparator: (a, b) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
          },
          { key: 'edition_name', label: 'Nome' },
          { key: 'type', label: 'Frequência' },
          {
            key: 'status',
            label: 'Status',
            comparator: (a, b) => STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, local ou status..."
        itemLabel="edições"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalState({ open: true, edition: undefined })}
            className="h-9 gap-2"
          >
            <Plus className="size-4" />
            Nova edição
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Calendar}
            eyebrow="Edições"
            title="Nenhuma edição encontrada"
            description="Crie a primeira edição para começar a organizar esse evento."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((edition, idx) => (
            <AdminEditionCard
              key={edition.id}
              edition={edition}
              eventId={eventId}
              index={idx}
              onEdit={(currentEdition) => setModalState({ open: true, edition: currentEdition })}
              onPublish={edition.status === 'draft' ? () => setPublishingEdition(edition) : undefined}
            />
          ))
        }
      />

      <ManageEditionModal
        key={modalState.edition?.id ?? 'edition-create'}
        open={modalState.open}
        edition={modalState.edition}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }))
            return
          }

          setModalState({ open: false, edition: undefined })
        }}
        onCreate={async (values) => {
          const res = await createEditionMutation.mutateAsync({
            eventId,
            data: values,
          })

          return res.success ? res.data : false
        }}
        onUpdate={async (editionId, values) => {
          const res = await updateEditionMutation.mutateAsync({
            eventId,
            editionId,
            data: values,
          })

          return res.success ? res.data : false
        }}
      />

      <AlertModal
        open={Boolean(publishingEdition)}
        onOpenChange={() => setPublishingEdition(null)}
        title="Publicar edição?"
        description={
          publishingEdition
            ? `Ao publicar "${publishingEdition.edition_name}", ela ficará disponível como anunciada.`
            : undefined
        }
        confirmLabel="Publicar edição"
        variant="default"
        loading={publishEditionMutation.isPending}
        onConfirm={async () => {
          if (!publishingEdition) return
          publishEditionMutation.mutate({
            eventId,
            editionId: publishingEdition.id,
          })
          setPublishingEdition(null)
        }}
      />
    </div>
  )
}
