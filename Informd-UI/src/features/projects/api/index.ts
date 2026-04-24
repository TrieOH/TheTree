import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { ProjectCreateI, ProjectI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new Project on the server.
 * @param projectData - The data for the new project.
 * @returns A promise that resolves to the API response containing the newly created project.
 */
export const createProjectFn = createClientOnlyFn((projectData: ProjectCreateI) => {
  return authFetcher.post<ProjectI>("/projects", projectData);
});


/**
 * Fetches all projects from the server.
 * @returns A promise that resolves to an array of Projects objects.
 */
export const getAllProjectsFn = createClientOnlyFn(async () => {
  try {
    return await tanstackQueryFetcher<ProjectI[]>("/projects");
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all projects, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all projects.
 */
export const allProjectsQueryOptions = () => {
  return queryOptions({
    queryKey: ['projects'],
    queryFn: getAllProjectsFn,
  })
}
