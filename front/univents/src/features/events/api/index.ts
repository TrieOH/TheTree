import { createClientOnlyFn, createServerFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { EventCreateI, EventI } from "../model";
import type { Permission } from "@soramux/node-perm-sdk"
import type { ImageURLUploadI } from "@/shared/model/generic";
import { authFetcher, simpleFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { serverPerm } from "@/features/auths/lib/server-auth";

/**
 * Creates a new Event on the server.
 * @param eventData - The data for the new event.
 * @returns A promise that resolves to the API response containing the newly created event.
 */
export const createEventFn = createClientOnlyFn((eventData: EventCreateI) => {
  return authFetcher.post<EventI>("/events", eventData);
});

/**
 * Update Event on the server.
 * @param eventData - The data of the event to update.
 * @returns A promise that resolves to the API response containing the updated event.
 */
export const patchEventFn = createClientOnlyFn((id: string, eventData: Partial<EventI>) => {
  return authFetcher.patch<EventI>(`/events/${id}`, eventData);
});

/**
 * Fetches a single own event from the server by filtering the list.
 * @param id - The event id
 * @returns A promise that resolves to the Event object.
 * @throws Error if not found.
 */
export const getOwnEventFn = createClientOnlyFn(async (id: string) => {
  const events = await getOwnEventsFn();
  const event = events.find(e => e.id === id);
  if (event) return event;
  throw new Error("Failed to find own event in list")
});

/**
 * Query options for fetching a single own event.
 */
export const ownEventQueryOptions = (id: string) => {
  return queryOptions({
    queryKey: ['events', 'own', id],
    queryFn: () => getOwnEventFn(id),
  })
}

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
 * Fetches a single public event from the server by filtering the list.
 * @param id - The event id
 * @returns A promise that resolves to the Event object.
 * @throws Error if not found.
 */
export const getEventFn = async (id: string) => {
  const events = await getEventsFn();
  const event = events.find(e => e.id === id);
  if (event) return event;
  throw new Error("Failed to find event in list")
};

/**
 * Query options for fetching a single public event.
 */
export const eventQueryOptions = (id: string) => {
  return queryOptions({
    queryKey: ['events', 'public', id],
    queryFn: () => getEventFn(id),
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
  return authFetcher.post<null>(`/events/${eventId}/publish`);
});


/**
 * Adds a MinIO URL to the event's gallery_urls array.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const addImageToTheEventGalleryFn = createClientOnlyFn((
  eventId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.post<EventI>(`/events/${eventId}/gallery`, urlData);
});

/**
 * Removes a URL from the event's gallery_urls array and deletes the object from MinIO.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const removeImageToTheEventGalleryFn = createClientOnlyFn((
  eventId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.delete<EventI>(`/events/${eventId}/gallery`, urlData);
});

/**
 * Sets the event banner URL. If the URL is not already in gallery_urls it is added automatically.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const setEventBannerFn = createClientOnlyFn((
  eventId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.put<EventI>(`/events/${eventId}/banner`, urlData);
});

/**
 * Clears the event banner. The image remains in gallery_urls.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const unsetEventBannerFn = createClientOnlyFn((
  eventId: string,
) => {
  return authFetcher.delete<EventI>(`/events/${eventId}/banner`);
});

/**
 * Sets the event logo URL. If the URL is not already in gallery_urls it is added automatically.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const setEventLogoFn = createClientOnlyFn((
  eventId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.put<EventI>(`/events/${eventId}/logo`, urlData);
});

/**
 * Clears the event logo. The image remains in gallery_urls.
 * @param eventId - The event id
 * @returns A promise that resolves to the API EventI response.
 */
export const unsetEventLogoFn = createClientOnlyFn((
  eventId: string,
) => {
  return authFetcher.delete<EventI>(`/events/${eventId}/logo`);
});

// Server
export const checkAdminPermissionFn = createServerFn({ method: "POST" })
  .inputValidator((data: Permission.CheckPermissionRequestI) => {
    return data
  })
  .handler(async ({ data }) => {
    const result = await serverPerm.check(data)
    if (result.success) return { success: true, data: result.data }
    return { success: false, message: result.message }
  });

