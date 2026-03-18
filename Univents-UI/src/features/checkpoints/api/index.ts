import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { CheckpointCreateI, CheckpointI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Checkpoint on the server.
 * @param checkpointData - The data for the new checkpoint.
 * @returns A promise that resolves to the API response containing the newly created checkpoint.
 */
export const createCheckpointFn = createClientOnlyFn((
  checkpointData: CheckpointCreateI, eventId: string, editionId: string
) => {
  return authFetcher.post<CheckpointI>(
    `/events/${eventId}/editions/${editionId}/checkpoints`,
    checkpointData
  );
});

/**
 * Fetches all checkpoints for a specific edition from the server.
 * @returns A promise that resolves to an array of Checkpoint objects.
 */
export const getAllCheckpointsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<CheckpointI[]>(`/events/${eventId}/editions/${editionId}/checkpoints`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all checkpoints for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all checkpoints for a specific edition.
 */
export const allCheckpointsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['checkpoints', 'public', eventId, editionId],
    queryFn: () => getAllCheckpointsFn(eventId, editionId),
  })
}
