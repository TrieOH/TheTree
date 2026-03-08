import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Create a new API key for the specified workspace on the server.
 * @param apiKeyData - The data for the new API key.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to the API response containing the newly created API key.
 */
export const createApiKeyOnWorkspaceFn = createClientOnlyFn((
  apiKeyData: ApiKeyCreateI,
  workspaceName: string
) => {
  return authFetcher.post<ApiKeyCreateResponseI>(
    `/workspaces/${workspaceName}/keys`,
    apiKeyData
  );
});


/**
 * Fetches all API keys for the specified workspace from the server.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to an array of API key objects.
 */
export const getAllWorkspaceApiKeysFn = createClientOnlyFn(async (workspaceName: string) => {
  try {
    return await tanstackQueryFetcher<ApiKeyI[]>(`/workspaces/${workspaceName}/keys`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all API keys for a specific workspace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allWorkspaceApiKeysQueryOptions = (workspaceName: string) => {
  return queryOptions({
    queryKey: ['workspaces', workspaceName, "keys"],
    queryFn: () => getAllWorkspaceApiKeysFn(workspaceName),
  })
}


/**
 * Delete an API key for the specified workspace on the server.
 * @param workspaceName - The workspace name
 * @param apiKeyId - The API key ID
 * @returns A promise that resolves to the API response(void).
 */
export const deleteApiKeyOnWorkspaceFn = createClientOnlyFn((
  workspaceName: string,
  apiKeyId: string
) => {
  return authFetcher.delete<void>(
    `/workspaces/${workspaceName}/keys/${apiKeyId}`
  );
});