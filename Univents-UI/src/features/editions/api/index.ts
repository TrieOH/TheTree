import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { EditionCreateI, EditionI } from "../model";
import { authFetcher, simpleFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Edition on the server.
 * @param editionData - The data for the new edition.
 * @returns A promise that resolves to the API response containing the newly created edition.
 */
export const createEditionFn = createClientOnlyFn((editionData: EditionCreateI, eventId: string) => {
  return authFetcher.post<EditionI>(`/events/${eventId}/editions`, editionData);
});

/**
 * Fetches all event editions from the server.
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getAllEditionsFn = async (eventId: string) => {
  try {
    // FIXME: Use a alternative version like tanstackQuerySimpleFetcher
    const res = await simpleFetcher.get<EditionI[]>(`/events/${eventId}/editions`);
    if (res.success) return res.data
    return []
  } catch {
    return [];
  }
};

/**
 * Query options for fetching all event editions, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all event editions.
 */
export const allEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: ['editions', 'public', eventId],
    queryFn: () => getAllEditionsFn(eventId),
  })
}

/**
 * Query options for fetching a specific event edition, using TanStack Query.
 * If the list of all editions is already in cache, it uses that data.
 * Otherwise, it fetches the list and filters for the specific ID.
 * @returns An object containing the query key and query function for fetching a specific event edition.
 */
export const editionQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['editions', 'public', eventId, editionId],
    queryFn: async () => {
      const editions = await getAllEditionsFn(eventId);
      return editions.find(e => e.id === editionId) ?? null;
    },
  })
}


/**
 * Fetches all admin event editions from the server.
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getAllAdminEditionsFn = createClientOnlyFn(async (eventId: string) => {
  try {
    return await tanstackQueryFetcher<EditionI[]>(`/events/${eventId}/editions/admin`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all admin event editions, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all admin event editions.
 */
export const allAdminEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: ['editions', 'admin', eventId],
    queryFn: () => getAllAdminEditionsFn(eventId),
  })
};

/**
 * Publish a Edition on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to the API null response.
 */
export const publishEditionFn = createClientOnlyFn((
  eventId: string, editionId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/announce`
  );
});

/**
 * Connect Payment Account a Edition on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param credentialId - The Credential id
 * @param provider - Payment Provider
 * @returns A promise that resolves to the API null response.
 */
export const connectPaymentAccountToEditionFn = createClientOnlyFn((
  eventId: string, editionId: string, credentialId: string, provider: string, public_key: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/payments/connect?credential_id=${credentialId}&provider=${provider}&public_key=${public_key}`
  );
});

/**
 * Connect Payment Account a Edition on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to the API null response.
 */
export const disconnectPaymentAccountToEditionFn = createClientOnlyFn((
  eventId: string, editionId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/payments/disconnect`
  );
});
