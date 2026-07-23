import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { EventI } from '../model'
import { createEventFn, publishEventFn } from './index'
import { eventKeys } from './query-keys'

function upsertById(events: EventI[] | undefined, event: EventI) {
  const list = events ?? []
  const index = list.findIndex((item) => item.id === event.id)

  if (index === -1) return [...list, event]

  const next = [...list]
  next[index] = event
  return next
}

function removeById(events: EventI[] | undefined, eventId: string) {
  return (events ?? []).filter((item) => item.id !== eventId)
}

function shouldBePublic(event: Pick<EventI, 'status'>) {
  return event.status !== 'draft'
}

function syncEventCaches(queryClient: QueryClient, event: EventI) {
  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) =>
    upsertById(old, event),
  )

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) => {
    if (shouldBePublic(event)) return upsertById(old, event)
    return removeById(old, event.id)
  })
}

function syncEventStatusInCaches(
  queryClient: QueryClient,
  eventId: string,
  status: EventI['status'],
  updatedAt: string,
) {
  const ownEvent = queryClient
    .getQueryData<EventI[]>(eventKeys.ownLists())
    ?.find((event) => event.id === eventId)

  if (!ownEvent) return

  const nextEvent = { ...ownEvent, status, updated_at: updatedAt }

  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) =>
    upsertById(old, nextEvent),
  )

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) =>
    upsertById(old, nextEvent),
  )
}

export function useCreateEventMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: createEventFn,
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao criar evento')
        return
      }

      syncEventCaches(queryClient, res.data)
      toast.success('Evento criado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function usePublishEventMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (eventId: string) => publishEventFn(eventId),
    onSuccess: (res, eventId) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar evento')
        return
      }

      syncEventStatusInCaches(
        queryClient,
        eventId,
        'active',
        new Date().toISOString(),
      )
      toast.success('Evento publicado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}
