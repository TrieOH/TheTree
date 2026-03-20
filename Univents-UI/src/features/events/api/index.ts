import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { EventCreateI, EventI } from "../model";
import { authFetcher, simpleFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Event on the server.
 * @param eventData - The data for the new event.
 * @returns A promise that resolves to the API response containing the newly created event.
 */
export const createEventFn = createClientOnlyFn((eventData: EventCreateI) => {
  return authFetcher.post<EventI>("/events", eventData);
});

/**
 * Fetches all own events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
export const getOwnEventsFn = createClientOnlyFn(async () => {
  try {
    return await tanstackQueryFetcher<EventI[]>("/events/own");
  } catch {
    return [];
  }
});

/**
 * Query options for fetching own events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching own events.
 */
export const ownEventsQueryOptions = () => {
  return queryOptions({
    queryKey: ['events', 'own'],
    queryFn: getOwnEventsFn,
  })
}

/**
 * Fetches all events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
export const getEventsFn = async () => {
  try {
    // FIXME: Use a alternative version like tanstackQuerySimpleFetcher
    const res = await simpleFetcher.get<EventI[]>("/events");
    if (res.success) return res.data
    return []
  } catch {
    return [];
  }
};

/**
 * Query options for fetching events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching events.
 */
export const eventsQueryOptions = () => {
  return queryOptions({
    queryKey: ['events', 'public'],
    queryFn: getEventsFn,
  })
}

/**
 * Publish a Event on the server.
 * @param eventId - The event id
 * @returns A promise that resolves to the API null response.
 */
export const publishEventFn = createClientOnlyFn((
  eventId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/publish`
  );
});

