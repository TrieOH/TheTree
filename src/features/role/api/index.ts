import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { Role, RoleCRUD } from "../model/types";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new Role on the server.
 * @param roleData - The data for the new role.
 * @returns A promise that resolves to the API response containing the newly created role.
 */
export const createRoleFn = createClientOnlyFn((roleData: Omit<RoleCRUD, "id">) => {
  const { project_id, ...dataToSend } = roleData;
  return authFetcher<Role>(`/projects/${project_id}/roles`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify(dataToSend),
  });
});


/**
 * Updates an existing role on the server.
 * @param roleData - The data for the role to update, including its ID.
 * @returns A promise that resolves to the API response containing the updated role.
 */
export const patchRoleFn = createClientOnlyFn((roleData: RoleCRUD) => {
  const { id, project_id, ...dataToSend } = roleData;

  return authFetcher<null>(`/projects/${project_id}/roles/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(dataToSend),
  });
});


/**
 * Fetches all roles from the server.
 * @returns A promise that resolves to an array of Role objects.
 */
export const getRolesFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["roles", string];
}) => {
  const [, projectId] = queryKey;
  try {
    return await tanstackQueryFetcher<Role[]>(`/projects/${projectId}/roles`);
  } catch {
    return [] as Role[];
  }
});

export const roleQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['roles', project_id],
    queryFn: getRolesFn,
    enabled: !!project_id
  })
}