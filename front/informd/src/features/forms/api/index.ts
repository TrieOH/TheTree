import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { FormCreateI, FormI } from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Create a new Form on the server.
 * @param formData - The data for the new form.
 * @returns A promise that resolves to the API response containing the newly created Form.
 */
export const createFormFn = createClientOnlyFn((formData: FormCreateI) => {
  return authFetcher.post<FormI>('forms', formData);
});

/**
 * Fetches all personal Forms for the current user from the server.
 * @returns A promise that resolves to an array of Form objects.
 */
export const getAllUserFormsFn = createClientOnlyFn(async () => {
  return tanstackQueryFetcher<FormI[]>("/forms");
});

/**
 * Fetches all archived Forms for the current user from the server.
 * @returns A promise that resolves to an array of archived Form objects.
 */
export const getAllUserArchivedFormsFn = createClientOnlyFn(async () => {
  return tanstackQueryFetcher<FormI[]>("/forms/archived");
});

/**
 * Query options for fetching all personal Forms for the current user, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all personal Forms.
 */
export const allUserFormsQueryOptions = () => {
  return queryOptions({
    queryKey: ["forms"],
    queryFn: () => getAllUserFormsFn(),
  })
}

/**
 * Query options for fetching all archived Forms for the current user, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all archived Forms.
 */
export const allUserArchivedFormsQueryOptions = () => {
  return queryOptions({
    queryKey: ["forms", "archived"],
    queryFn: () => getAllUserArchivedFormsFn(),
  })
}
// Manage Form Status

/**
 * Opens a Form on the server.
 * @param form_id - The ID of the form to open.
 * @returns A promise that resolves to the API response containing the updated Form.
 */
export const openFormFn = createClientOnlyFn((form_id: string) => {
  return authFetcher.post<FormI>(`forms/${form_id}/open`);
});

/**
 * Closes a Form on the server (if it is open).
 * @param form_id - The ID of the form to close.
 * @returns A promise that resolves to the API response containing the updated Form.
 */
export const closeFormFn = createClientOnlyFn((form_id: string) => {
  return authFetcher.post<FormI>(`forms/${form_id}/close`);
});

/**
 * Archives a Form on the server (if it is closed).
 * @param form_id - The ID of the form to archive.
 * @returns A promise that resolves to the API response containing the updated Form.
 */
export const archiveFormFn = createClientOnlyFn((form_id: string) => {
  return authFetcher.post<FormI>(`forms/${form_id}/archive`);
});

/**
 * Redrafts a Form on the server (if it is open and have zero submissions/responses).
 * @param form_id - The ID of the form to redraft.
 * @returns A promise that resolves to the API response containing the updated Form.
 */
export const redraftFormFn = createClientOnlyFn((form_id: string) => {
  return authFetcher.post<FormI>(`forms/${form_id}/redraft`);
});

// Response Count

/**
 * Fetches the response count for a specific Form from the server.
 * @param form_id - The ID of the form to fetch the response count for.
 * @returns A promise that resolves to the number of responses for the specified Form.
 */
export const getFormResponseCountFn = createClientOnlyFn((form_id: string) => {
  return tanstackQueryFetcher<{ count: number }>(`forms/${form_id}/responses/count`);
});

/**
 * Query options for fetching the response count for a specific Form, using TanStack Query.
 * @param form_id - The ID of the form to fetch the response count for.
 * @returns An object containing the query key and query function for fetching the response count.
 */
export const formResponseCountQueryOptions = (form_id: string) => {
  return queryOptions({
    queryKey: ["forms", form_id, "responses", "count"],
    queryFn: () => getFormResponseCountFn(form_id),
  })
}