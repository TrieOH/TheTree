import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { ActivityCreateI, ActivityI, AttendanceRecordI } from "../model";
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
 * Not Implemeted yet
 * Updates an existing Activity on the server.
 * @param activityId - The ID of the activity to update.
 * @param activityData - The updated data for the activity.
 * @returns A promise that resolves to the API response containing the updated activity.
 */
export const updateActivityFn = createClientOnlyFn((
  activityId: string, activityData: Partial<ActivityI>, eventId: string, editionId: string
) => {
  return authFetcher.patch<ActivityI>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}`,
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
    queryKey: ['activities', 'public', eventId, editionId],
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
    queryKey: ['activities', 'admin', eventId, editionId],
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

/**
 * Lists attendance records of the activity if you have activities:manage on the activity
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param activityId - The activity id
 * @returns A promise that resolves to the API AttendanceRecordI response.
 */
export const getAllActivityAttendanceRecordsFn = createClientOnlyFn((
  eventId: string, editionId: string, activityId: string
) => {
  return tanstackQueryFetcher<AttendanceRecordI[]>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/records`
  );
});

/**
 * Query options for fetching all admin activities for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all admin activities for a specific edition.
 */
export const allActivityAttendanceRecordsQueryOptions = (
  eventId: string, editionId: string, activityId: string
) => {
  return queryOptions({
    queryKey: ['activities', 'admin', eventId, editionId, 'records', activityId],
    queryFn: () => getAllActivityAttendanceRecordsFn(eventId, editionId, activityId),
  })
};


/**
 * If you have attendance:mark on the activity and the record status is registered, marks it as completed
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param activityId - The activity id
 * @returns A promise that resolves to the API null response.
 */
export const markAttendanceForUserInActivityFn = createClientOnlyFn((
  eventId: string, editionId: string, activityId: string, recordId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/records/${recordId}`
  );
});

/**
 * Registers the user to the specified activity if they have activities:attend permission on it
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param activityId - The activity id
 * @returns A promise that resolves to the API null response.
 */
export const registerUserInActivityFn = createClientOnlyFn((
  eventId: string, editionId: string, activityId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/register`
  );
});

/**
 * Unregisters the user from the specified activity if they are registered on it
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param activityId - The activity id
 * @returns A promise that resolves to the API null response.
 */
export const unregisterUserInActivityFn = createClientOnlyFn((
  eventId: string, editionId: string, activityId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/unregister`
  );
});