import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  addOrganizationMember,
  createOrganization,
  listOrganizationMembers,
  listOrganizations,
  removeOrganizationMember,
} from "@trieoh/identityx-api";
import type {
  MemberAddToOrganizationI,
  OrganizationCreateI,
  OrganizationI,
  OrganizationMemberI,
} from "../model";

/**
 * Creates a new OrganizationI on the server.
 * @param orgData - The data for the new organization.
 * @returns A promise that resolves to the API response containing the newly created organization.
 */
export const createOrganizationFn = createClientOnlyFn(
  (orgData: OrganizationCreateI) => {
    return createOrganization(orgData).then(orvalData<OrganizationI>);
  },
);

/**
 * Fetches all organizations from the server.
 * @returns A promise that resolves to an array of organizations objects.
 */
export const getAllOrganizationsFn = createClientOnlyFn(() => {
  return listOrganizations().then(orvalData<OrganizationI[]>);
});

/**
 * Query options for fetching all Organizations, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Organizations.
 */
export const allOrganizationsQueryOptions = () => {
  return queryOptions({
    queryKey: ["orgs"],
    queryFn: () => getAllOrganizationsFn(),
  });
};

// Members

/**
 * Adds a new member to a organization on the server.
 * @param organization_id - The ID of the organization to add the member to.
 * @param memberData - The data for the new member.
 * @returns A promise that resolves to the API response containing the newly created member.
 */
export const addMemberToOrganizationFn = createClientOnlyFn(
  (organization_id: string, memberData: MemberAddToOrganizationI) => {
    return addOrganizationMember(organization_id, memberData).then(
      orvalData<void>,
    );
  },
);

/**
 * Removes a member from a organization on the server.
 * @param organization_id - The ID of the organization to remove the member from.
 * @param actor_email - The email of the user to remove from the organization.
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromOrganizationFn = createClientOnlyFn(
  (organization_id: string, actor_email: string) => {
    return removeOrganizationMember(organization_id, { actor_email }).then(
      orvalData<void>,
    );
  },
);

/**
 * Fetches all organization members from the server.
 * @param organization_id - The ID of the organization to fetch members for.
 * @returns A promise that resolves to an array of members objects.
 */
export const getAllOrganizationsMemberFn = createClientOnlyFn(
  (organization_id: string) => {
    return listOrganizationMembers(organization_id).then(
      orvalData<OrganizationMemberI[]>,
    );
  },
);

/**
 * Query options for fetching all Members, using TanStack Query.
 * @param organization_id - The ID of the organization to fetch members for.
 * @returns An object containing the query key and query function for fetching all Members.
 */
export const allOrganizationsMembersQueryOptions = (
  organization_id: string,
) => {
  return queryOptions({
    queryKey: ["organizations", organization_id, "members"],
    queryFn: () => getAllOrganizationsMemberFn(organization_id),
  });
};
