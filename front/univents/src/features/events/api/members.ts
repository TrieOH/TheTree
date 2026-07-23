import { createClientOnlyFn } from '@tanstack/react-start'
import {
  queryOptions,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import { toast } from 'sonner'
import { authFetcher, authQueryFetcher } from '@/shared/lib/api/fetch'
import { eventKeys } from './query-keys'
import type { EventMemberRole } from '../model/member'

export type { EventMemberRole } from '../model/member'

export interface EventMemberI {
  id: string
  event_id: string
  user_id: string
  role: EventMemberRole
  created_at: string
  updated_at: string | null
  deleted_at: string | null
  /** Available locally for members added during the current session. */
  email?: string
}

export interface AddEventMemberInput {
  eventId: string
  email: string
  role: EventMemberRole
}

export interface RemoveEventMemberInput {
  eventId: string
  userId: string
  email: string
}

export const getEventMembersFn = createClientOnlyFn((eventId: string) => {
  return authQueryFetcher<EventMemberI[]>(`/events/${eventId}/members`)
})

export const allEventMembersQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: eventKeys.members(eventId),
    queryFn: () => getEventMembersFn(eventId),
  })

export const addEventMemberFn = createClientOnlyFn(
  ({ eventId, email, role }: AddEventMemberInput) => {
    return authFetcher.post<EventMemberI>(`/events/${eventId}/members`, {
      email,
      role,
    })
  },
)

export const removeEventMemberFn = createClientOnlyFn(
  ({ eventId, userId, email }: RemoveEventMemberInput) => {
    return authFetcher.delete<null>(`/events/${eventId}/members/${userId}`, {
      email,
    })
  },
)

export function useAddEventMemberMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: addEventMemberFn,
    onSuccess: (res, input) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao adicionar membro')
        return
      }

      queryClient.setQueryData<EventMemberI[]>(
        eventKeys.members(input.eventId),
        (old = []) => [
          ...old.filter((member) => member.id !== res.data.id),
          { ...res.data, email: input.email },
        ],
      )
      toast.success('Membro adicionado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useRemoveEventMemberMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: removeEventMemberFn,
    onSuccess: (res, input) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao remover membro')
        return
      }

      queryClient.setQueryData<EventMemberI[]>(
        eventKeys.members(input.eventId),
        (old = []) => old.filter((member) => member.user_id !== input.userId),
      )
      toast.success('Membro removido com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}
