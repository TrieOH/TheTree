import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "node_modules/@tanstack/react-query/build/modern/queryOptions";
import type { EditionCreateI, EditionI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

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
export const getAllEditionsFn = createClientOnlyFn(async (eventId: string) => {
  try {
    return await tanstackQueryFetcher<EditionI[]>(`/events/${eventId}/editions`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all event editions, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all event editions.
 */
export const allEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: ['allEditions', eventId],
    queryFn: () => getAllEditionsFn(eventId),
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
    queryKey: ['allAdminEditions', eventId],
    queryFn: () => getAllAdminEditionsFn(eventId),
  })
};