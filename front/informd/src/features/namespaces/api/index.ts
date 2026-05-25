import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type {
  MemberAddToNamespaceI,
  NamespaceCreateI,
  NamespaceI,
  NamespaceMemberI
} from "../model";
import { queryOptions } from "@tanstack/react-query";
import type { FormCreateI, FormI } from "#/features/forms/model";

/**
 * Creates a new NamespaceI on the server.
 * @param namespaceData - The data for the new namespace.
 * @returns A promise that resolves to the API response containing the newly created namespace.
 */
export const createNamespaceFn = createClientOnlyFn((namespaceData: NamespaceCreateI) => {
  return authFetcher.post<NamespaceI>("/namespaces", namespaceData);
});

/**
 * Fetches all namespaces from the server.
 * @returns A promise that resolves to an array of namespaces objects.
 */
export const getAllNamespacesFn = createClientOnlyFn(() => {
  return tanstackQueryFetcher<NamespaceI[]>("/namespaces");
});

/**
 * Query options for fetching all Namespaces, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Namespaces.
 */
export const allNamespacesQueryOptions = () => {
  return queryOptions({
    queryKey: ['namespaces'],
    queryFn: () => getAllNamespacesFn(),
  })
}

// Members

/**
 * Adds a new member to a namespace on the server.
 * @param namespace_id - The ID of the namespace to add the member to.
 * @param memberData - The data for the new member.
 * @returns A promise that resolves to the API response containing the newly created member.
 */
export const addMemberToNamespaceFn = createClientOnlyFn((namespace_id: string, memberData: MemberAddToNamespaceI) => {
  return authFetcher.post(`/namespaces/${namespace_id}/members`, memberData);
});

/**
 * Removes a member from a namespace on the server.
 * @param namespace_id - The ID of the namespace to remove the member from.
 * @param user_id - The ID of the user to remove from the namespace.
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromNamespaceFn = createClientOnlyFn((namespace_id: string, user_id: string) => {
  return authFetcher.delete(`/namespaces/${namespace_id}/members`, { user_id });
});

/**
 * Fetches all namespace members from the server.
 * @param namespace_id - The ID of the namespace to fetch members for.
 * @returns A promise that resolves to an array of members objects.
 */
export const getAllNamespacesMemberFn = createClientOnlyFn((
  namespace_id: string
) => {
  return tanstackQueryFetcher<NamespaceMemberI[]>(`/namespaces/${namespace_id}/members`);
});

/**
 * Query options for fetching all Members, using TanStack Query.
 * @param namespace_id - The ID of the namespace to fetch members for.
 * @returns An object containing the query key and query function for fetching all Members.
 */
export const allNamespacesMembersQueryOptions = (namespace_id: string) => {
  return queryOptions({
    queryKey: ['namespaces', namespace_id, 'members'],
    queryFn: () => getAllNamespacesMemberFn(namespace_id),
  })
}

// Form

/**
 * Creates a new FormI on the server.
 * @param namespace_id - The ID of the namespace to create the form in.
 * @param formData - The data for the new form.
 * @returns A promise that resolves to the API response containing the newly created form.
 */
export const createFormOnNamespaceFn = createClientOnlyFn((
  namespace_id: string, formData: FormCreateI
) => {
  return authFetcher.post<FormI>(`/namespaces/${namespace_id}/forms`, formData);
});

/**
 * Fetches all forms for a specific namespace from the server.
 * @param namespace_id - The ID of the namespace to fetch forms for.
 * @returns A promise that resolves to an array of form objects.
 */
export const getAllNamespacesFormsFn = createClientOnlyFn((namespace_id: string) => {
  return tanstackQueryFetcher<FormI[]>(`/namespaces/${namespace_id}/forms`);
});

/**
 * Query options for fetching all Forms, using TanStack Query.
 * @param namespace_id - The ID of the namespace to fetch forms for.
 * @returns An object containing the query key and query function for fetching all Forms.
 */
export const allNamespacesFormsQueryOptions = (namespace_id: string) => {
  return queryOptions({
    queryKey: ['namespaces', namespace_id, 'forms'],
    queryFn: () => getAllNamespacesFormsFn(namespace_id),
  })
}
