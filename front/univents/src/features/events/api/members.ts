import { createClientOnlyFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
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
