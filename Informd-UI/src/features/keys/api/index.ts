import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Create a new API key for the specified project on the server.
 * @param apiKeyData - The data for the new API key.
 * @param projectID - The project ID
 * @returns A promise that resolves to the API response containing the newly created API key.
 */
export const createApiKeyOnProjectFn = createClientOnlyFn((
  apiKeyData: ApiKeyCreateI,
  projectID: string
) => {
  return authFetcher.post<ApiKeyCreateResponseI>(
    `/projects/${projectID}/keys`,
    apiKeyData
  );
});


/**
 * Fetches all API keys for the specified project from the server.
 * @param projectID - The project ID
 * @returns A promise that resolves to an array of API key objects.
 */
export const getAllProjectApiKeysFn = createClientOnlyFn(async (projectID: string) => {
  try {
    return await tanstackQueryFetcher<ApiKeyI[]>(`/projects/${projectID}/keys`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all API keys for a specific project, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allProjectApiKeysQueryOptions = (projectID: string) => {
  return queryOptions({
    queryKey: ['projects', projectID, "keys"],
    queryFn: () => getAllProjectApiKeysFn(projectID),
  })
}


/**
 * Revoke an API key for the specified project on the server.
 * @param apiKeyId - The ID of the API key to revoke.
 * @param projectID - The project ID
 * @returns A promise that resolves to the API response(void).
 */
export const revokeApiKeyOnProjectFn = createClientOnlyFn((
  apiKeyId: string,
  projectID: string
) => {
  return authFetcher.delete<void>(
    `/projects/${projectID}/keys/${apiKeyId}`
  );
});