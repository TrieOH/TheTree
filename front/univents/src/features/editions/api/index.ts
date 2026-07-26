import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import {
  authFetcher,
  authQueryFetcher,
  publicQueryFetcher,
} from "@/shared/lib/api/fetch";
import {
  type EditionApiI,
  type EditionCreateOutputI,
  type EditionPatchOutputI,
  normalizeEdition,
} from "../model";
import { editionKeys } from "./query-keys";

/**
 * Creates a new Edition on the server.
 * @param editionData - The data for the new edition.
 * @returns A promise that resolves to the API response containing the newly created edition.
 */
export const createEditionFn = createClientOnlyFn(
  async (editionData: EditionCreateOutputI, eventId: string) => {
    const res = await authFetcher.post<EditionApiI>(
      `/events/${eventId}/editions`,
      editionData,
    );
    return res.success ? { ...res, data: normalizeEdition(res.data) } : res;
  },
);

export const patchEditionFn = createClientOnlyFn(
  async (
    eventId: string,
    editionId: string,
    editionData: EditionPatchOutputI,
  ) => {
    const res = await authFetcher.patch<EditionApiI>(
      `/events/${eventId}/editions/${editionId}`,
      editionData,
    );
    return res.success ? { ...res, data: normalizeEdition(res.data) } : res;
  },
);

export const publishEditionFn = createClientOnlyFn(
  (eventId: string, editionId: string) =>
    authFetcher.post<null>(`/events/${eventId}/editions/${editionId}/publish`),
);

/**
 * Fetches all event editions from the server.
 * @param eventId - The event id
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getAllPublicEditionsFn = async (eventId: string) => {
  const editions = await publicQueryFetcher<EditionApiI[]>(
    `/events/${eventId}/editions`,
  );
  return editions.filter((edition) => !edition.is_draft).map(normalizeEdition);
};

const getPublicEdition = async (path: string) => {
  try {
    return normalizeEdition(await publicQueryFetcher<EditionApiI>(path));
  } catch {
    return null;
  }
};

const getPublicEditions = async (path: string) =>
  (await publicQueryFetcher<EditionApiI[]>(path)).map(normalizeEdition);

export const activeEditionQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.activeByEvent(eventId),
    queryFn: () => getPublicEdition(`/events/${eventId}/editions/active`),
  });

export const pastEditionsQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.pastByEvent(eventId),
    queryFn: () => getPublicEditions(`/events/${eventId}/editions/past`),
  });

export const upcomingEditionsQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.upcomingByEvent(eventId),
    queryFn: () => getPublicEditions(`/events/${eventId}/editions/upcoming`),
  });

/**
 * Query options for fetching all event editions, using TanStack Query.
 * @param eventId - The event id
 * @returns An object containing the query key and query function for fetching all event editions.
 */
export const allPublicEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.publicListByEvent(eventId),
    queryFn: () => getAllPublicEditionsFn(eventId),
  });
};

/**
 * Fetches all admin event editions from the server.
 * @param eventId - The event id
 * @returns A promise that resolves to an array of Edition objects.
 */
export const getDraftEditionsFn = createClientOnlyFn(
  async (eventId: string) => {
    const editions = await authQueryFetcher<EditionApiI[]>(
      `/events/${eventId}/editions/draft`,
    );
    return editions.map(normalizeEdition);
  },
);

export const getAllAdminEditionsFn = createClientOnlyFn(
  async (eventId: string) => {
    const [publicEditions, draftEditions] = await Promise.all([
      getAllPublicEditionsFn(eventId),
      getDraftEditionsFn(eventId),
    ]);

    return [...publicEditions, ...draftEditions].filter(
      (edition, index, editions) =>
        editions.findIndex((candidate) => candidate.id === edition.id) ===
        index,
    );
  },
);

/**
 * Query options for fetching all admin event editions, using TanStack Query.
 * @param eventId - The event id
 * @returns An object containing the query key and query function for fetching all admin event editions.
 */
export const allAdminEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.adminListByEvent(eventId),
    queryFn: () => getAllAdminEditionsFn(eventId),
  });
};
