import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { FormCreateI, FormI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Create a new Form for the specified namespace on the server.
 * @param formData - The data for the new form.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to the API response containing the newly created Form.
 */
export const createFormOnNamespaceFn = createClientOnlyFn((
  formData: FormCreateI,
  namespaceID: string
) => {
  return authFetcher.post<FormI>(`/namespaces/${namespaceID}/forms`, formData);
});


/**
 * Fetches all Forms for the specified namespace from the server.
 * @param namespaceID - The namespace ID
 * @returns A promise that resolves to an array of Form objects.
 */
export const getAllNamespaceFormsFn = createClientOnlyFn(async (namespaceID: string) => {
  try {
    return await tanstackQueryFetcher<FormI[]>(`/namespaces/${namespaceID}/forms`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all Form for a specific namespace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Form.
 */
export const allNamespaceFormsQueryOptions = (namespaceID: string) => {
  return queryOptions({
    queryKey: ['namespaces', namespaceID, "forms"],
    queryFn: () => getAllNamespaceFormsFn(namespaceID),
  })
}
