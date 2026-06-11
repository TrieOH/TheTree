import { authFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Create a new API key for the current user on the server.
 * @param apiKeyData - The data for the new API key.
 * @returns A promise that resolves to the API response containing the newly created API key.
 */
export const createApiKeyFn = createClientOnlyFn((apiKeyData: ApiKeyCreateI) => {
  return authFetcher.post<ApiKeyCreateResponseI>('/api-keys', apiKeyData);
});


/**
 * Fetches all API keys from the server.
 * @returns A promise that resolves to an array of API key objects.
 */
export const getAllApiKeysFn = createClientOnlyFn(async () => {
  const res = await authFetcher.post<ApiKeyI[]>('/api-keys/bulk');
  return res.success ? res.data : []
});

/**
 * Query options for fetching all API keys for the current user, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allApiKeysQueryOptions = () => {
  return queryOptions({
    queryKey: ["keys"],
    queryFn: () => getAllApiKeysFn(),
  })
}

/**
 * Revoke an API key for the current user on the server.
 * @param apiKeyId - The ID of the API key to revoke.
 * @returns A promise that resolves to the API response(void).
 */
export const revokeApiKeyFn = createClientOnlyFn((apiKeyId: string) => {
  return authFetcher.delete<void>(`api-keys/${apiKeyId}`);
});