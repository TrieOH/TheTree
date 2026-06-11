import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { FullFormI, SubmitRequestI } from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Submit the form answer to the server.
 * @param form_id - The ID of the Form for which to fetch responses.
 * @param submitData - The answer to submit.
 * @returns A promise that resolves to the API response.
 */
export const submitFormFn = createClientOnlyFn((
  form_id: string,
  submitData: SubmitRequestI
) => {
  return authFetcher.post<void>(`/forms/${form_id}/responses`, submitData);
});

/**
 * Fetches all Form Responses for the current user from the server.
 * @param form_id - The ID of the Form for which to fetch responses.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to. If not provided, fetches responses without namespace context.
 * @returns A promise that resolves to an array of FullFormI objects.
 */
export const getFullFormResponseDetailsFn = createClientOnlyFn(async (
  form_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return tanstackQueryFetcher<FullFormI>(`/namespaces/${namespace_id}/forms/${form_id}/full`);
  return tanstackQueryFetcher<FullFormI>(`/forms/${form_id}/full`);
});

/**
 * Query options for fetching all Form Responses for a specific Form.
 * @param form_id - The ID of the Form for which to fetch responses.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to. If not provided, fetches responses without namespace context.
 * @returns An object containing the query key and query function for fetching all responses of the specified Form.
 */
export const allFormsResponsesQueryOptions = (
  form_id: string,
  namespace_id?: string
) => {
  return queryOptions({
    queryKey: ["forms", form_id, "full", namespace_id],
    queryFn: () => getFullFormResponseDetailsFn(form_id, namespace_id),
  })
}
