import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { NamespaceCreateI, NamespaceI } from "../model";
import { queryOptions } from "@tanstack/react-query";

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
