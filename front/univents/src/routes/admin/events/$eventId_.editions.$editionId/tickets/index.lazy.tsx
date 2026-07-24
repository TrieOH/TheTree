import { createLazyFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Ticket, Plus } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import { Button } from '@/shared/ui/shadcn/button'
import { allTicketsQueryOptions } from '@/features/tickets/api'
import {
  useCreateTicketMutation,
  useUpdateTicketMutation,
} from '@/features/tickets/api/mutations'
import type { TicketI } from '@/features/tickets/model'
import AdminTicketCard from '@/features/tickets/ui/AdminTicketCard'
import { ManageTicketModal } from '@/features/tickets/ui/ManageTicketModal'

export const Route = createLazyFileRoute('/admin/events/$eventId_/editions/$editionId/tickets/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { editionId } = Route.useParams()
  const { data: tickets = [] } = useQuery(allTicketsQueryOptions(editionId))
  const createTicketMutation = useCreateTicketMutation()
  const updateTicketMutation = useUpdateTicketMutation()
  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<TicketI>>({
    field: 'name',
    direction: 'asc',
  })
  const [modalState, setModalState] = useState<{ open: boolean; ticket?: TicketI }>({
    open: false,
  })

  const filteredTickets = [...tickets]
    .filter((ticket) => {
      const search = filter.trim().toLowerCase()
      if (!search) return true

      return [
        ticket.name,
        ticket.description ?? '',
        String(ticket.access_level),
        String(ticket.price_cents),
      ].some((value) => value.toLowerCase().includes(search))
    })
    .sort((a, b) => {
      const direction = sort.direction === 'asc' ? 1 : -1

      if (sort.field === 'price_cents') {
        return (a.price_cents - b.price_cents) * direction
      }

      if (sort.field === 'access_level') {
        return (a.access_level - b.access_level) * direction
      }

      return String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
    })

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<TicketI>
        items={filteredTickets}
        layout="grid"
        minItemWidth="16rem"
        pageSize={8}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: 'name', label: 'Nome' },
          { key: 'price_cents', label: 'Preço', comparator: (a, b) => a.price_cents - b.price_cents },
          { key: 'access_level', label: 'Nível de acesso', comparator: (a, b) => a.access_level - b.access_level },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, descrição ou nível..."
        itemLabel="tickets"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalState({ open: true, ticket: undefined })}
            className="h-9 gap-2"
          >
            <Plus className="size-4" />
            Novo ticket
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Ticket}
            eyebrow="Tickets"
            title="Nenhum ticket encontrado"
            description="Crie o primeiro ticket para começar a vender nessa edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((ticket, idx) => (
            <AdminTicketCard
              key={ticket.id}
              ticket={ticket}
              index={idx}
              onManage={(currentTicket) => setModalState({ open: true, ticket: currentTicket })}
            />
          ))
        }
      />

      <ManageTicketModal
        key={modalState.ticket?.id ?? 'ticket-create'}
        open={modalState.open}
        ticket={modalState.ticket}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }))
            return
          }

          setModalState({ open: false, ticket: undefined })
        }}
        onCreate={async (values) => {
          const res = await createTicketMutation.mutateAsync({
            editionId,
            data: values,
          })

          return res.success ? res.data : false
        }}
        onUpdate={async (ticketId, values) => {
          const res = await updateTicketMutation.mutateAsync({
            ticketId,
            data: values,
          })

          return res.success ? res.data : false
        }}
      />
    </div>
  )
}