import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { Role, RoleCRUD } from "../model/types";
import { queryOptions } from "@tanstack/react-query";
import type { Permission } from "@/features/permission/model/types";


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

  return authFetcher<null>(`/projects/${project_id}/roles/${id}/description`, {
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

// Permissions

export const getRolePermissionsFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["rolePermissions", string, string];
}) => {
  const [, projectId, roleId] = queryKey;
  try {
    return await tanstackQueryFetcher<Permission[]>(`/projects/${projectId}/roles/${roleId}/permissions`);
  } catch {
    return [] as Permission[];
  }
});

export const rolePermissionsQueryOptions = (project_id: string, role_id: string) => {
  return queryOptions({
    queryKey: ['rolePermissions', project_id, role_id],
    queryFn: getRolePermissionsFn,
    enabled: !!project_id && !!role_id
  })
}


export const givePermissionToRoleFn = createClientOnlyFn((roleData: Role, permission_id: string) => {
  const { project_id, id } = roleData;
  return authFetcher<null>(`/projects/${project_id}/roles/${id}/permissions/${permission_id}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
  });
});

export const removePermissionOfRoleFn = createClientOnlyFn((roleData: Role, permission_id: string) => {
  const { project_id, id } = roleData;
  return authFetcher<null>(`/projects/${project_id}/roles/${id}/permissions/${permission_id}`, {
    method: "DELETE",
  });
});