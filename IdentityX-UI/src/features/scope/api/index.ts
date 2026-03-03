import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { Scope, ScopeCRUD } from "../model/types";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new scope on the server.
 * @param scopeData - The data for the new scope.
 * @returns A promise that resolves to the API response containing the newly created scope.
 */
export const createScopeFn = createClientOnlyFn((scopeData: Omit<ScopeCRUD, "id">) => {
  const { project_id, ...dataToSend } = scopeData;
  return authFetcher<Scope>(`/projects/${project_id}/scopes`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify(dataToSend),
  });
});


/**
 * Fetches all scopes from the server.
 * @returns A promise that resolves to an array of Scope objects.
 */
export const getScopesFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["scopes", string];
}) => {
  const [, projectId] = queryKey;

  try {
    return await tanstackQueryFetcher<Scope[]>(
      `/projects/${projectId}/scopes`
    );
  } catch {
    return [] as Scope[];
  }
});

export const scopesQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['scopes', project_id],
    queryFn: getScopesFn,
    enabled: !!project_id
  })
}


/**
 * Update a scope's metadata on the server.
 * @param metaData - The metadata for the scope.
 */
export const patchScopeMetaFn = createClientOnlyFn((metaData: Partial<ScopeCRUD>) => {
  const { project_id, id, meta } = metaData;
  return authFetcher<void>(`/projects/${project_id}/scopes/${id}/meta`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ meta }),
  });
});


/**
 * Deletes a scope from the server.
 * @param id - The ID of the scope to delete.
 * @returns A promise that resolves to the API response.
 */
export const deleteScopeFn = createClientOnlyFn(({project_id, id}: {project_id: string, id: string}) => {
  return authFetcher<void>(`/projects/${project_id}/scopes/${id}`, {
    method: "DELETE",
  });
});