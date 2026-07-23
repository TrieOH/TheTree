import { createClientOnlyFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { EventCreateSubmitI, EventI } from '../model'
import {
  authFetcher,
  publicQueryFetcher,
  authQueryFetcher,
} from '@/shared/lib/api/fetch'
import { eventKeys } from './query-keys'

/**
 * Creates a new Event on the server.
 * @param eventData - The data for the new event.
 * @returns A promise that resolves to the API response containing the newly created event.
 */
export const createEventFn = createClientOnlyFn(
  (eventData: EventCreateSubmitI) => {
    return authFetcher.post<EventI>('/events', eventData)
  },
)

/**
 * Publish a Event on the server.
 * @param eventId - The event id
 * @returns A promise that resolves to the API null response.
 */
export const publishEventFn = createClientOnlyFn((eventId: string) => {
  return authFetcher.post<null>(`/events/${eventId}/publish`)
})

/**
 * Fetches all public events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
export const getPublicEventsFn = async () => {
  return publicQueryFetcher<EventI[]>('/events')
}

/**
 * Query options for fetching events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching events.
 */
export const allPublicEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.publicLists(),
    queryFn: getPublicEventsFn,
  })
}

/**
 * Fetches all own events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
export const getOwnEventsFn = createClientOnlyFn(async () => {
  return authQueryFetcher<EventI[]>('/events/owned')
})

/**
 * Query options for fetching own events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching own events.
 */
export const allOwnEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.ownLists(),
    queryFn: getOwnEventsFn,
  })
}

/**
 * Fetches a single public event from the server by filtering the list.
 * @param id - The event id
 * @returns A promise that resolves to the Event object.
 * @throws Error if not found.
 */
export const getPublicEventFn = async (id: string) => {
  const events = await getPublicEventsFn()
  const event = events.find((e) => e.id === id)
  if (event) return event
  throw new Error('Failed to find event in list')
}

/**
 * Query options for fetching a single public event.
 */
export const publicEventQueryOptions = (id: string) => {
  return queryOptions({
    queryKey: ['events', 'public', id],
    queryFn: () => getPublicEventFn(id),
  })
}
