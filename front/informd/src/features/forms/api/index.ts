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
 * Query options for fetching all personal Forms for the current user, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all personal Forms.
 */
export const allUserFormsQueryOptions = () => {
  return queryOptions({
    queryKey: ["forms"],
    queryFn: () => getAllUserFormsFn(),
  })
}
