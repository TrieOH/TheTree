import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { EditionCreateOutputI, EditionI } from "../model";
import { authFetcher, publicQueryFetcher, authQueryFetcher } from "@/shared/lib/api/fetch";
import { editionKeys } from "./query-keys";

/**
 * Creates a new Edition on the server.
 * @param editionData - The data for the new edition.
 * @returns A promise that resolves to the API response containing the newly created edition.
 */
export const createEditionFn = createClientOnlyFn((editionData: EditionCreateOutputI, eventId: string) => {
  return authFetcher.post<EditionI>(`/events/${eventId}/editions`, editionData);
});

/**
 * Update an Edition on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param editionData - The updated edition data.
 * @returns A promise that resolves to the API response containing the updated edition.
 */
export const patchEditionFn = createClientOnlyFn((eventId: string, editionId: string, editionData: EditionCreateOutputI) => {
  return authFetcher.patch<EditionI>(`/events/${eventId}/editions/${editionId}`, editionData);
});

/**
 * Publish a Edition on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to the API null response.
 */
export const publishEditionFn = createClientOnlyFn((eventId: string, editionId: string) => {
  return authFetcher.post<null>(`/events/${eventId}/editions/${editionId}/announce`);
});

/**
 * Fetches all event editions from the server.
 * @param eventId - The event id
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getAllPublicEditionsFn = async (eventId: string) => {
  return publicQueryFetcher<EditionI[]>(`/events/${eventId}/editions`);
};

/**
 * Query options for fetching all event editions, using TanStack Query.
 * @param eventId - The event id
 * @returns An object containing the query key and query function for fetching all event editions.
 */
export const allPublicEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.publicListByEvent(eventId),
    queryFn: () => getAllPublicEditionsFn(eventId),
  })
}

/**
 * Fetches all admin event editions from the server.
 * @param eventId - The event id
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getAllAdminEditionsFn = createClientOnlyFn(async (eventId: string) => {
  return await authQueryFetcher<EditionI[]>(`/events/${eventId}/editions/admin`);
});

/**
 * Query options for fetching all admin event editions, using TanStack Query.
 * @param eventId - The event id
 * @returns An object containing the query key and query function for fetching all admin event editions.
 */
export const allAdminEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.adminListByEvent(eventId),
    queryFn: () => getAllAdminEditionsFn(eventId),
  })
};

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


// FIXME: I NEED TO DELETE EVERYTHING BELOW THIS LINE AND REPLACE IT


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
      const editions = await getAllPublicEditionsFn(eventId);
      return editions.find(e => e.id === editionId) ?? null;
    },
  })
}
