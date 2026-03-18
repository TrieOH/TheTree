import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { ActivityCreateI, ActivityI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Activity on the server.
 * @param activityData - The data for the new activity.
 * @returns A promise that resolves to the API response containing the newly created activity.
 */
export const createActivityFn = createClientOnlyFn((
  activityData: ActivityCreateI, eventId: string, editionId: string
) => {
  return authFetcher.post<ActivityI>(
    `/events/${eventId}/editions/${editionId}/activities`,
    activityData
  );
});


/**
 * Fetches all activities for a specific edition from the server.
 * @returns A promise that resolves to an array of Activity objects.
 */
export const getAllActivitiesFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<ActivityI[]>(`/events/${eventId}/editions/${editionId}/activities`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all activities for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all activities for a specific edition.
 */
export const allActivitiesQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['allActivities', eventId, editionId],
    queryFn: () => getAllActivitiesFn(eventId, editionId),
  })
}


/**
 * Fetches all admin activities for a specific edition from the server.
 * @returns A promise that resolves to an array of Activity objects.
 */
export const getAllAdminActivitiesFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<ActivityI[]>(`/events/${eventId}/editions/${editionId}/activities/admin`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all admin activities for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all admin activities for a specific edition.
 */
export const allAdminActivitiesQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['allAdminActivities', eventId, editionId],
    queryFn: () => getAllAdminActivitiesFn(eventId, editionId),
  })
};

/**
 * Publish a Activity on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param activityId - The activity id
 * @returns A promise that resolves to the API null response.
 */
export const publishActivityFn = createClientOnlyFn((
  eventId: string, editionId: string, activityId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/publish`
  );
});
