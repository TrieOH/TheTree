import { authFetcher } from "@/shared/lib/api/fetch";
import type { Project, ProjectCRUD } from "../model/types";
import { createClientOnlyFn } from "@tanstack/react-start";

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
})
