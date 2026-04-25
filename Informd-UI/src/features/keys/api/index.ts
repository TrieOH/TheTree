import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Create a new API key for the specified namespace on the server.
 * @param apiKeyData - The data for the new API key.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to the API response containing the newly created API key.
 */
export const createApiKeyOnNamespaceFn = createClientOnlyFn((
  apiKeyData: ApiKeyCreateI,
  namespaceID: string
) => {
  return authFetcher.post<ApiKeyCreateResponseI>(
    `/namespaces/${namespaceID}/keys`,
    apiKeyData
  );
});


/**
 * Fetches all API keys for the specified namespace from the server.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to an array of API key objects.
 */
export const getAllNamespaceApiKeysFn = createClientOnlyFn(async (namespaceID: string) => {
  try {
    return await tanstackQueryFetcher<ApiKeyI[]>(`/namespaces/${namespaceID}/keys`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all API keys for a specific namespace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allNamespaceApiKeysQueryOptions = (namespaceID: string) => {
  return queryOptions({
    queryKey: ['namespaces', namespaceID, "keys"],
    queryFn: () => getAllNamespaceApiKeysFn(namespaceID),
  })
}


/**
 * Revoke an API key for the specified namespace on the server.
 * @param apiKeyId - The ID of the API key to revoke.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to the API response(void).
 */
export const revokeApiKeyOnNamespaceFn = createClientOnlyFn((
  apiKeyId: string,
  namespaceID: string
) => {
  return authFetcher.delete<void>(
    `/namespaces/${namespaceID}/keys/${apiKeyId}`
  );
});