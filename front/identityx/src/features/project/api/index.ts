import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { ProjectCreateI, ProjectI } from "../model";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";

/**
 * Creates a new project on the server.
 * @param projectData - The data for the new project.
 * @param orgId - The organization ID to which the project belongs (optional).
 * @returns A promise that resolves to the API response containing the newly created project.
 */
export const createProjectFn = createClientOnlyFn((projectData: ProjectCreateI, orgId?: string) => {
  if (orgId)
    return authFetcher.post<ProjectI>(`/organizations/${orgId}/projects`, projectData);
  return authFetcher.post<ProjectI>("/projects", projectData);
});

/**
 * Fetches all projects from the server.
 * @param orgId - The organization ID to filter projects by (optional).
 * @returns A promise that resolves to an array of ProjectI objects.
 */
export const getProjectsFn = createClientOnlyFn(async (orgId?: string) => {
  if (orgId)
    return await tanstackQueryFetcher<ProjectI[]>(`/organizations/${orgId}/projects`);
  return await tanstackQueryFetcher<ProjectI[]>("/projects");
});

/**
 * Query options for fetching projects, compatible with React Query's useQuery hook.
 * @param orgId - The organization ID to filter projects by (optional).
 * @returns An object containing the query key and query function for fetching projects.
 */
export const allProjectsQueryOptions = (orgId?: string) => {
  return queryOptions({
    queryKey: ["organizations", orgId, "projects"],
    queryFn: () => getProjectsFn(orgId),
  });
};