import { ApiError, authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { User } from "../model/types";
import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { Permission } from "@/features/permission/model/types";
import type { Role } from "@/features/role/model/types";

/**
 * Fetches all users from the server.
 * @returns A promise that resolves to an array of User objects.
 */
export const getUsersFn = (async ({ queryKey,}: { queryKey: ["users", string] }) => {
  const [, projectId] = queryKey;
  try {
    return await tanstackQueryFetcher<User[]>(
      `/projects/${projectId}/users`
    );
  } catch {
    return [] as User[];
  }
});

export const usersQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['users', project_id],
    queryFn: getUsersFn,
    enabled: !!project_id
  })
}

// Permissions

export const getUserPermissionsFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["userPermissions", string, string, string | null];
}) => {
  const [, projectId, id, scope_id] = queryKey;
  try {
    let url = `/projects/${projectId}/identities/${id}/permissions`;
    if (scope_id) url += `?scope_id=${scope_id}`;
    return await tanstackQueryFetcher<Permission[]>(url);
  } catch {
    return [] as Permission[];
  }
});

export const userPermissionsQueryOptions = (project_id: string, id: string, scope_id: string | null) => {
  return queryOptions({
    queryKey: ['userPermissions', project_id, id, scope_id],
    queryFn: getUserPermissionsFn,
    enabled: !!project_id && !!id && scope_id !== undefined,
    refetchOnMount: true,
  })
}

export const givePermissionToUserFn = createClientOnlyFn(async (userData: User, permission_id: string, scope_id: string | null) => {
  const { project_id, id } = userData;
  const response = await authFetcher<void>(`/projects/${project_id}/identities/${id}/permissions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ permission_id, scope_id })
  });
  if (!response.success) throw new ApiError(response);
  
  return response.data;
});

export const removePermissionOfUserFn = createClientOnlyFn((userData: User, permission_id: string,  scope_id: string | null) => {
  const { project_id, id } = userData;
  return authFetcher<null>(`/projects/${project_id}/identities/${id}/permissions`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ permission_id, scope_id })
  });
});

// Roles

export const getUserRolesFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["userRoles", string, string];
}) => {
  const [, projectId, id] = queryKey;
  try {
    return await tanstackQueryFetcher<Role[]>(`/projects/${projectId}/identities/${id}/roles`);
  } catch {
    return [] as Role[];
  }
});

export const userRolesQueryOptions = (project_id: string, id: string) => {
  return queryOptions({
    queryKey: ['userRoles', project_id, id],
    queryFn: getUserRolesFn,
    enabled: !!project_id && !!id,
    refetchOnMount: true,
  })
}

export const giveRoleToUserFn = createClientOnlyFn(async (userData: User, role_id: string, scope_id: string | null) => {
  const { project_id, id } = userData;
  const response = await authFetcher<void>(`/projects/${project_id}/identities/${id}/roles`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ role_id, scope_id })
  });
  if (!response.success) throw new ApiError(response);
  return response.data;
});

export const removeRoleOfUserFn = createClientOnlyFn((userData: User, role_id: string,  scope_id: string | null) => {
  const { project_id, id } = userData;
  return authFetcher<null>(`/projects/${project_id}/identities/${id}/roles`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ role_id, scope_id })
  });
});
