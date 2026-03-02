import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { Permission, PermissionCRUD } from "../model/types";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new permission on the server.
 * @param permData - The data for the new permission.
 * @returns A promise that resolves to the API response containing the newly created permission.
 */
export const createPermissionFn = createClientOnlyFn((permData: Omit<PermissionCRUD, "id">) => {
  const { project_id, ...dataToSend } = permData;
  return authFetcher<Permission>(`/projects/${project_id}/permissions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify(dataToSend),
  });
});


/**
 * Fetches all permissions from the server.
 * @returns A promise that resolves to an array of Permission objects.
 */
export const getPermissionsFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["permissions", string, string?, string?];
}) => {
  const [, projectId, object, action] = queryKey;
  const params = new URLSearchParams()
  if(object) params.append("object", object)
  if(action) params.append("action", action)

  try {
    return await tanstackQueryFetcher<Permission[]>(
      `/projects/${projectId}/permissions?${params.toString()}`
    );
  } catch {
    return [] as Permission[];
  }
});

export const permissionsQueryOptions = (project_id: string, object?: string, action?: string) => {
  return queryOptions({
    queryKey: ['permissions', project_id, object, action],
    queryFn: getPermissionsFn,
    enabled: !!project_id
  })
}

/**
 * Deletes a permission from the server.
 * @param id - The ID of the permission to delete.
 * @returns A promise that resolves to the API response.
 */
export const deletePermissionFn = createClientOnlyFn(({project_id, id}: {project_id: string, id: string}) => {
  return authFetcher<void>(`/projects/${project_id}/permissions/${id}`, {
    method: "DELETE",
  });
});


/**
 * Update a permission's metadata on the server.
 * @param metaData - The metadata for the permission.
 */
export const patchPermissionMetaFn = createClientOnlyFn((metaData: Partial<PermissionCRUD>) => {
  const { project_id, id, meta } = metaData;
  return authFetcher<void>(`/projects/${project_id}/permissions/${id}/meta`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ meta }),
  });
});