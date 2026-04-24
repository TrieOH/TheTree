import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { FormCreateI, FormI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Create a new Form for the specified project on the server.
 * @param formData - The data for the new form.
 * @param projectID - The project ID
 * @returns A promise that resolves to the API response containing the newly created Form.
 */
export const createFormOnProjectFn = createClientOnlyFn((
  formData: FormCreateI,
  projectID: string
) => {
  return authFetcher.post<FormI>(`/projects/${projectID}/forms`, formData);
});


/**
 * Fetches all Forms for the specified project from the server.
 * @param projectID - The project ID
 * @returns A promise that resolves to an array of Form objects.
 */
export const getAllProjectFormsFn = createClientOnlyFn(async (projectID: string) => {
  try {
    return await tanstackQueryFetcher<FormI[]>(`/projects/${projectID}/forms`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all Form for a specific project, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Form.
 */
export const allProjectFormsQueryOptions = (projectID: string) => {
  return queryOptions({
    queryKey: ['projects', projectID, "forms"],
    queryFn: () => getAllProjectFormsFn(projectID),
  })
}
