import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type {
  CreateFieldRequestI,
  FieldI,
  FieldSelectConfigI,
  FieldUpdateI,
  CreateFieldSelectConfigRequestI,
} from "../model";
import { queryOptions } from "@tanstack/react-query";

/**
 * Create a new Field on the server.
 * @param formData - The data for the new field.
 * @param form_id - The ID of the Form to which the Field belongs.
 * @param step_id - The ID of the Step to which the Field belongs.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves to the API response containing the newly created Field.
 */
export const createFieldFn = createClientOnlyFn((
  formData: CreateFieldRequestI,
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return authFetcher.post<FieldI>(`namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields`, formData);
  return authFetcher.post<FieldI>(`/forms/${form_id}/steps/${step_id}/fields`, formData);
});

/**
 * Bulk edit Fields on the server.
 * @param formData - The data for the updated fields.
 * @param form_id - The ID of the Form to which the Fields belong.
 * @param step_id - The ID of the Step to which the Fields belong.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves to the API response containing the updated Fields.
 */
export const bulkEditFieldsFn = createClientOnlyFn((
  formData: FieldUpdateI[],
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return authFetcher.put<FieldI>(`namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields`, formData);
  return authFetcher.put<FieldI>(`/forms/${form_id}/steps/${step_id}/fields`, formData);
});

/**
 * Delete a Field on the server.
 * @param field_id - The ID of the Field to delete.
 * @param form_id - The ID of the Form to which the Field belongs.
 * @param step_id - The ID of the Step to which the Field belongs.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves when the Field is deleted.
 */
export const deleteFieldFn = createClientOnlyFn((
  field_id: string,
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return authFetcher.delete<void>(`namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields/${field_id}`);
  return authFetcher.delete<void>(`/forms/${form_id}/steps/${step_id}/fields/${field_id}`);
});

/**
 * Fetches all Fields for a specific Step from the server.
 * @param form_id - The ID of the Form for which to fetch fields.
 * @param step_id - The ID of the Step for which to fetch fields.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves to an array of Field objects.
 */
export const getAllUserFieldsFn = createClientOnlyFn(async (
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return tanstackQueryFetcher<FieldI[]>(`/namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields`);
  return tanstackQueryFetcher<FieldI[]>(`/forms/${form_id}/steps/${step_id}/fields`);
});

/**
 * Query options for fetching all fields of a specific Step.
 * @param form_id - The ID of the Form for which to fetch fields.
 * @param step_id - The ID of the Step for which to fetch fields.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns An object containing the query key and query function for fetching all fields of the specified Step.
 */
export const allStepsFieldsQueryOptions = (
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  return queryOptions({
    queryKey: ["forms", form_id, "steps", step_id, "fields", namespace_id],
    queryFn: () => getAllUserFieldsFn(form_id, step_id, namespace_id),
  });
};

// Field Select Config

/**
 * Edit the select configuration for a specific Field.
 * @param formData - The updated select configuration data.
 * @param field_id - The ID of the Field.
 * @param form_id - The ID of the Form to which the Field belongs.
 * @param step_id - The ID of the Step to which the Field belongs.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves to the updated select configuration.
 */
export const editFieldSelectConfigFn = createClientOnlyFn((
  formData: CreateFieldSelectConfigRequestI,
  field_id: string,
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return authFetcher.put<FieldSelectConfigI>(`namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields/${field_id}/select`, formData);
  return authFetcher.put<FieldSelectConfigI>(`/forms/${form_id}/steps/${step_id}/fields/${field_id}/select`, formData);
});

/**
 * Fetch the select configuration for a specific Field.
 * @param field_id - The ID of the Field.
 * @param form_id - The ID of the Form to which the Field belongs.
 * @param step_id - The ID of the Step to which the Field belongs.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns A promise that resolves to the select configuration.
 */
export const getFieldSelectConfigFn = createClientOnlyFn(async (
  field_id: string,
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  if (namespace_id)
    return tanstackQueryFetcher<FieldSelectConfigI>(`/namespaces/${namespace_id}/forms/${form_id}/steps/${step_id}/fields/${field_id}/select`);
  return tanstackQueryFetcher<FieldSelectConfigI>(`/forms/${form_id}/steps/${step_id}/fields/${field_id}/select`);
});

/**
 * Query options for fetching the select configuration for a specific Field.
 * @param field_id - The ID of the Field.
 * @param form_id - The ID of the Form to which the Field belongs.
 * @param step_id - The ID of the Step to which the Field belongs.
 * @param namespace_id - (Optional) The ID of the Namespace that the Form belongs to.
 * @returns An object containing the query key and query function for fetching the select configuration of the specified Field.
 */
export const allSelectConfigsQueryOptions = (
  field_id: string,
  form_id: string,
  step_id: string,
  namespace_id?: string
) => {
  return queryOptions({
    queryKey: ["forms", form_id, "steps", step_id, "fields", field_id, "select_config", namespace_id],
    queryFn: () => getFieldSelectConfigFn(field_id, form_id, step_id, namespace_id),
  });
};