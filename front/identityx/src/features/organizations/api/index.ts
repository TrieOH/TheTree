import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { OrganizationCreateI, OrganizationI } from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Creates a new OrganizationI on the server.
 * @param orgData - The data for the new organization.
 * @returns A promise that resolves to the API response containing the newly created organization.
 */
export const createOrganizationFn = createClientOnlyFn((orgData: OrganizationCreateI) => {
  return authFetcher.post<OrganizationI>("/organizations", orgData);
});

/**
 * Fetches all organizations from the server.
 * @returns A promise that resolves to an array of organizations objects.
 */
export const getAllOrganizationsFn = createClientOnlyFn(() => {
  return tanstackQueryFetcher<OrganizationI[]>("/organizations");
});

/**
 * Query options for fetching all Organizations, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Organizations.
 */
export const allOrganizationsQueryOptions = () => {
  return queryOptions({
    queryKey: ['orgs'],
    queryFn: () => getAllOrganizationsFn(),
  })
}
