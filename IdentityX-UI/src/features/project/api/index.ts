import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { Project, ProjectCRUD } from "../model/types";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";

/**
 * Creates a new project on the server.
 * @param projectData - The data for the new project.
 * @returns A promise that resolves to the API response containing the newly created project.
 */
export const createProjectFn = createClientOnlyFn((projectData: Omit<ProjectCRUD, "id">) => {
    const dataToSend = {
    ...projectData,
    metadata: {}
  };

  return authFetcher<Project>("/projects", {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify(dataToSend),
  });
});

/**
 * Fetches all projects from the server.
 * @returns A promise that resolves to an array of Project objects.
 */
export const getProjectsFn = createClientOnlyFn(() => {
  return tanstackQueryFetcher<Project[]>("/projects").catch( _ => {
    return [] as Project[]
  })
});

export const projectsQueryOptions = queryOptions({
  queryKey: ['projects'],
  queryFn: getProjectsFn
})

/**
 * Updates an existing project on the server.
 * @param projectData - The data for the project to update, including its ID.
 * @returns A promise that resolves to the API response containing the updated project.
 */
export const patchProjectFn = createClientOnlyFn((projectData: ProjectCRUD) => {
  const { id, ...dataToSend } = projectData;

  return authFetcher<Project>(`/projects/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(dataToSend),
  });
});

/**
 * Deletes a project from the server.
 * @param id - The ID of the project to delete.
 * @returns A promise that resolves to the API response.
 */
export const deleteProjectFn = createClientOnlyFn((id: string) => {
  return authFetcher<void>(`/projects/${id}`, {
    method: "DELETE",
  });
});

