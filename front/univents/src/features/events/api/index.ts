import { createClientOnlyFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { EventCreateOutputI, EventI } from '../model'
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
  (eventData: EventCreateOutputI) => {
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

export const discontinueEventFn = createClientOnlyFn((eventId: string) => {
  return authFetcher.post<null>(`/events/${eventId}/discontinue`)
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
 * Fetches all joined events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
export const getJoinedEventsFn = createClientOnlyFn(async () => {
  return authQueryFetcher<EventI[]>('/events/joined')
})

/**
 * Query options for fetching joined events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching joined events.
 */
export const allJoinedEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.joinedLists(),
    queryFn: getJoinedEventsFn,
  })
}

/**
 * Fetches all public events from the server.
 * @param slug - the event slug
 * @returns A promise that resolves to an Event object.
 */
export const getPublicEventBySlugFn = async (slug: string) => {
  return publicQueryFetcher<EventI | null>(`/events/${slug}:by-slug`).catch(() => null)
}

/**
 * Query options for fetching event, using TanStack Query.
 * @param slug - the event slug
 * @returns An object containing the query key and query function for fetching the event.
 */
export const publicEventBySlugQueryOptions = (slug: string) => {
  return queryOptions({
    queryKey: eventKeys.detail.publicBySlug(slug),
    queryFn: () => getPublicEventBySlugFn(slug),
  })
}