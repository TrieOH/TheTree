import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { itemsOf } from "@trieoh/api-client";
import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import type { StepCreateI, StepI, StepUpdateI } from "../model";

/**
 * Create a new Step on the server.
 * @param formData - The data for the new step.
 * @param form_id - The ID of the Form to which the Step belongs.
 * @returns A promise that resolves to the API response containing the newly created Step.
 */
export const createStepFn = createClientOnlyFn(
  (formData: StepCreateI, form_id: string) => {
    return authFetcher.post<StepI>(`/forms/${form_id}/steps`, formData);
  },
);

/**
 * Bulk edit Steps on the server.
 * @param formData - The data for the updated steps.
 * @param form_id - The ID of the Form to which the Steps belong.
 * @returns A promise that resolves to the API response containing the updated Steps.
 */
export const bulkEditStepsFn = createClientOnlyFn(
  (formData: StepUpdateI[], form_id: string) => {
    return authFetcher.put<StepI>(`/forms/${form_id}/steps`, formData);
  },
);

/**
 * Fetches all Steps for the current user from the server.
 * @param form_id - The ID of the Form for which to fetch steps.
 * @returns A promise that resolves to an array of Step objects.
 */
export const getAllUserStepsFn = createClientOnlyFn(async (form_id: string) => {
  return itemsOf<StepI>(
    await tanstackQueryFetcher<unknown>(`/forms/${form_id}/steps`),
  );
});

/**
 * Query options for fetching all steps of a specific Form.
 * @param form_id - The ID of the Form for which to fetch steps.
 * @returns An object containing the query key and query function for fetching all steps of the specified Form.
 */
export const allFormsStepsQueryOptions = (form_id: string) => {
  return queryOptions({
    queryKey: ["forms", form_id, "steps"],
    queryFn: () => getAllUserStepsFn(form_id),
  });
};
