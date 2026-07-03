import { createClientOnlyFn } from "@tanstack/react-start";
import type { CapabilityCreateI, CapabilityI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { queryOptions } from "@tanstack/react-query";

/**
 * Creates a new capability on the server.
 * @param project_id - The ID of the project.
 * @param capabilityData - The data for the new capability.
 * @returns A promise that resolves to the API response containing the newly created capability.
 */
export const createCapabilityFn = createClientOnlyFn((project_id: string, capabilityData: CapabilityCreateI) => {
  return authFetcher.post<CapabilityI>(`projects/${project_id}/capabilities`, capabilityData);
});

/**
 * Fetches all capabilities from the server.
 * @param project_id - The ID of the project.
 * @returns A promise that resolves to an array of CapabilityI objects.
 */
export const getCapabilitiesFn = createClientOnlyFn(async (project_id: string) => {
  return await tanstackQueryFetcher<CapabilityI[]>(`projects/${project_id}/capabilities`);
});

/**
 * Query options for fetching capabilities.
 * @param project_id - The ID of the project to fetch API keys for.
 * @returns An object containing the query key and query function for fetching capabilities.
 */
export const allCapabilitiesQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ["projects", project_id, "capabilities"],
    queryFn: () => getCapabilitiesFn(project_id),
  });
};