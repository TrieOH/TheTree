import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { MemberAddToOrganizationI, OrganizationCreateI, OrganizationI, OrganizationMemberI } from "../model";
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

// Members

/**
 * Adds a new member to a organization on the server.
 * @param organization_id - The ID of the organization to add the member to.
 * @param memberData - The data for the new member.
 * @returns A promise that resolves to the API response containing the newly created member.
 */
export const addMemberToOrganizationFn = createClientOnlyFn((organization_id: string, memberData: MemberAddToOrganizationI) => {
  return authFetcher.post(`/organizations/${organization_id}/members`, memberData);
});

/**
 * Removes a member from a organization on the server.
 * @param organization_id - The ID of the organization to remove the member from.
 * @param actor_email - The email of the user to remove from the organization.
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromOrganizationFn = createClientOnlyFn((organization_id: string, actor_email: string) => {
  return authFetcher.delete(`/organizations/${organization_id}/members`, { actor_email });
});

/**
 * Fetches all organization members from the server.
 * @param organization_id - The ID of the organization to fetch members for.
 * @returns A promise that resolves to an array of members objects.
 */
export const getAllOrganizationsMemberFn = createClientOnlyFn((
  organization_id: string
) => {
  return tanstackQueryFetcher<OrganizationMemberI[]>(`/organizations/${organization_id}/members`);
});

/**
 * Query options for fetching all Members, using TanStack Query.
 * @param organization_id - The ID of the organization to fetch members for.
 * @returns An object containing the query key and query function for fetching all Members.
 */
export const allOrganizationsMembersQueryOptions = (organization_id: string) => {
  return queryOptions({
    queryKey: ['organizations', organization_id, 'members'],
    queryFn: () => getAllOrganizationsMemberFn(organization_id),
  })
}