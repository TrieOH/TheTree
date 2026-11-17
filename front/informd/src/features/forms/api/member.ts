import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { itemsOf } from "@trieoh/api-client";
import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import type { FormMemberI, MemberAddToFormI } from "../model/member";

/**
 * Adds a new member to a namespace on the server.
 * @param memberData - The data for the new member.
 * @param form_id - The ID of the form to add the member to.
 * @returns A promise that resolves to the API response containing the newly created member.
 */
export const addMemberToFormFn = createClientOnlyFn(
  (memberData: MemberAddToFormI, form_id: string) => {
    return authFetcher.post(`/forms/${form_id}/members`, memberData);
  },
);

/**
 * Removes a member from a namespace on the server.
 * @param user_id - The ID of the user to remove from the namespace.
 * @param form_id - The ID of the form to add the member to.
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromFormFn = createClientOnlyFn(
  (user_id: string, form_id: string) => {
    return authFetcher.delete(`/forms/${form_id}/members`, { user_id });
  },
);

/**
 * Fetches all members of a specific Form from the server.
 * @param form_id - The ID of the Form for which to fetch members.
 * @returns A promise that resolves to an array of FormMemberI objects.
 */
export const getAllFormsMembersFn = createClientOnlyFn(
  async (form_id: string) => {
    return itemsOf<FormMemberI>(
      await tanstackQueryFetcher<unknown>(`/forms/${form_id}/members`),
    );
  },
);

/**
 * Query options for fetching all members of a specific Form.
 * @param form_id - The ID of the Form for which to fetch members.
 * @returns An object containing the query key and query function for fetching all members of the specified Form.
 */
export const allFormsMembersQueryOptions = (form_id: string) => {
  return queryOptions({
    queryKey: ["forms", form_id, "members"],
    queryFn: () => getAllFormsMembersFn(form_id),
  });
};
