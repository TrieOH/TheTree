import { authFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn, createServerFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from "../model";
import { queryOptions } from "@tanstack/react-query";
import { lookupResources } from "@soramux/node-perm-sdk";
import { serverPerm } from "#/shared/lib/api/server-auth";

const getApiKeyIds = createServerFn({ method: 'GET' })
  .inputValidator((userId: string) => userId)
  .handler(async ({ data: userId }) => {
    const request = lookupResources().subject("user", userId)
      .permission("view").resourceType("api_key").build()
    const userIds = [];
    const stream = serverPerm.lookupResources(request);
    for await (const response of stream) {
      if (response.result) userIds.push(response.result.resourceObjectId)
    }
    return userIds
  })

/**
 * Create a new API key for the specified namespace on the server.
 * @param apiKeyData - The data for the new API key.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to the API response containing the newly created API key.
 */
export const createApiKeyOnNamespaceFn = createClientOnlyFn((apiKeyData: ApiKeyCreateI) => {
  return authFetcher.post<ApiKeyCreateResponseI>('/api-keys', apiKeyData);
});


/**
 * Fetches all API keys for the specified namespace from the server.
 * @param userId if the user that want to see the namespace forms
 * @returns A promise that resolves to an array of API key objects.
 */
export const getAllNamespaceApiKeysFn = createClientOnlyFn(async (userId: string) => {
  const ids = await getApiKeyIds({ data: userId })
  const res = await authFetcher.post<ApiKeyI[]>('/api-keys/bulk', { ids });
  return res.success ? res.data : []
});

/**
 * Query options for fetching all API keys for a specific namespace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allNamespaceApiKeysQueryOptions = (userId: string) => {
  return queryOptions({
    queryKey: ['namespaces', userId, "keys"],
    queryFn: () => getAllNamespaceApiKeysFn(userId),
  })
}

/**
 * Revoke an API key for the specified namespace on the server.
 * @param apiKeyId - The ID of the API key to revoke.
 * @returns A promise that resolves to the API response(void).
 */
export const revokeApiKeyOnNamespaceFn = createClientOnlyFn((apiKeyId: string) => {
  return authFetcher.delete<void>(`api-keys/${apiKeyId}`);
});