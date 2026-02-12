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
  console.log(scopeData)
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