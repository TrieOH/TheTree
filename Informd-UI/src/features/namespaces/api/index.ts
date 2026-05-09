import { authFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn, createServerFn } from "@tanstack/react-start";
import type { NamespaceCreateI, NamespaceI } from "../model";
import { queryOptions } from "@tanstack/react-query";
import { lookupResources } from "@soramux/node-perm-sdk";
import { serverPerm } from "#/shared/lib/api/server-auth";


const getNamespaceIds = createServerFn({ method: 'GET' })
  .inputValidator((userId: string) => userId)
  .handler(async ({ data: userId }) => {
    const request = lookupResources().subject("user", userId)
      .permission("view").resourceType("namespace").build()
    const userIds = [];
    const stream = serverPerm.lookupResources(request);
    for await (const response of stream) {
      if(response.result) userIds.push(response.result.resourceObjectId)
    }
    return userIds
  })

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
 * @param userId if the user that want to see the namespaces
 * @returns A promise that resolves to an array of namespaces objects.
 */
export const getAllNamespacesFn = createClientOnlyFn(async (userId: string) => {
  const ids = await getNamespaceIds({ data: userId })
  const res = await authFetcher.post<NamespaceI[]>("/namespaces/bulk", {ids});
  return res.success ? res.data : []
});

/**
 * Query options for fetching all Namespaces, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Namespaces.
 */
export const allNamespacesQueryOptions = (userId: string) => {
  return queryOptions({
    queryKey: ['namespaces'],
    queryFn: () => getAllNamespacesFn(userId),
  })
}
