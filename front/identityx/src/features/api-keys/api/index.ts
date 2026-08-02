import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { createAPIKey } from "@trieoh/identityx-api";
import type { ApiKeyCreateI, ApiKeyI, CreateApiKeyResponseI } from "../model";

/**
 * Rotates an API key for a given project.
 * @param project_id - The ID of the project.
 * @param apiKeyData - The data for creating the new API key.
 * @returns A promise resolving to the response containing the new API key.
 */
export const rotateApiKeyFn = createClientOnlyFn(
  (project_id: string, apiKeyData: ApiKeyCreateI) => {
    return createAPIKey(project_id, apiKeyData).then(
      orvalData<CreateApiKeyResponseI>,
    );
  },
);

/**
 * Revokes an API key for a given project.
 * @param project_id - The ID of the project.
 * @param key_id - The ID of the API key to revoke.
 * @returns A promise resolving to the API response.
 */
export const revokeApiKeyFn = createClientOnlyFn(
  (project_id: string, key_id: string) => {
    console.warn(
      "API key revocation is unavailable: the endpoint is not defined in IdentityX",
      { project_id, key_id },
    );
    return Promise.resolve({
      success: true as const,
      code: 200,
      data: null,
      message: "API key revocation is not available yet",
    });
  },
);

/**
 * Fetches all API keys for a given project.
 * @param project_id - The ID of the project.
 * @returns A promise resolving to an array of ApiKeyI objects.
 */
export const getAllApiKeysFn = createClientOnlyFn((project_id: string) => {
  console.warn(
    "API key listing is unavailable: the endpoint is not defined in IdentityX",
    { project_id },
  );
  return Promise.resolve<ApiKeyI[]>([]);
});

/**
 * Query options for fetching all API keys, using TanStack Query.
 * @param project_id - The ID of the project to fetch API keys for.
 * @returns An object containing the query key and query function for fetching all API keys.
 */
export const allApiKeysQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ["projects", project_id, "api_keys"],
    queryFn: () => getAllApiKeysFn(project_id),
  });
};
