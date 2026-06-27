import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyI, CreateApiKeyResponseI } from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Rotates an API key for a given project.
 * @param project_id - The ID of the project.
 * @param apiKeyData - The data for creating the new API key.
 * @returns A promise resolving to the response containing the new API key.
 */
export const rotateApiKeyFn = createClientOnlyFn((project_id: string, apiKeyData: ApiKeyCreateI) => {
  const dataToSend = {
    ...apiKeyData,
    create_for_service_account: apiKeyData.create_for_service_account === 'true',
  };
  return authFetcher.post<CreateApiKeyResponseI>(`/projects/${project_id}/api_keys`, dataToSend);
});

/**
 * Revokes an API key for a given project.
 * @param project_id - The ID of the project.
 * @param key_id - The ID of the API key to revoke.
 * @returns A promise resolving to the API response.
 */
export const revokeApiKeyFn = createClientOnlyFn((project_id: string, key_id: string) => {
  return authFetcher.delete<null>(`/projects/${project_id}/api_keys/${key_id}`);
});

/**
 * Fetches all API keys for a given project.
 * @param project_id - The ID of the project.
 * @returns A promise resolving to an array of ApiKeyI objects.
 */
export const getAllApiKeysFn = createClientOnlyFn((project_id: string) => {
  return tanstackQueryFetcher<ApiKeyI[]>(`/projects/${project_id}/api_keys`);
});

/**
 * Query options for fetching all API keys, using TanStack Query.
 * @param project_id - The ID of the project to fetch API keys for.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allApiKeysQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['projects', project_id, 'api_keys'],
    queryFn: () => getAllApiKeysFn(project_id),
  });
};
