import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { WorkspaceCreateI, WorkspaceI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new Workspace on the server.
 * @param workspaceData - The data for the new workspace.
 * @returns A promise that resolves to the API response containing the newly created workspace.
 */
export const createWorkspaceFn = createClientOnlyFn((workspaceData: WorkspaceCreateI) => {
  return authFetcher.post<WorkspaceI>("/workspaces", workspaceData);
});


/**
 * Enable Workspace sandbox mode on the server.
 * @param name - The workspace name
 * @returns A promise that resolves to the API response containing the updated workspace.
 */
export const enableWorkspaceSandboxModeFn = createClientOnlyFn((name: string) => {
  return authFetcher.post<WorkspaceI>(`/workspaces/${name}/sandbox/enable`);
});


/**
 * Disable Workspace sandbox mode on the server.
 * @param name - The workspace name
 * @returns A promise that resolves to the API response containing the newly updated workspace.
 */
export const disableWorkspaceSandboxModeFn = createClientOnlyFn((name: string) => {
  return authFetcher.post<WorkspaceI>(`/workspaces/${name}/sandbox/disable`);
});


/**
 * Fetches all workspaces from the server.
 * @returns A promise that resolves to an array of Workspaces objects.
 */
export const getAllWorkspacesFn = createClientOnlyFn(async () => {
  try {
    return await tanstackQueryFetcher<WorkspaceI[]>("/workspaces");
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all workspaces, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all workspaces.
 */
export const allWorkspacesQueryOptions = () => {
  return queryOptions({
    queryKey: ['workspaces'],
    queryFn: getAllWorkspacesFn,
  })
}
