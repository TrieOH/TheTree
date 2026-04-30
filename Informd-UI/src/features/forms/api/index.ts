import { authFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn, createServerFn } from "@tanstack/react-start";
import type { FormCreateI, FormI } from "../model";
import { queryOptions } from "@tanstack/react-query";
import { lookupResources } from "@soramux/node-perm-sdk";
import { serverPerm } from "#/shared/lib/api/server-auth";

const getFormsIds = createServerFn({ method: 'GET' })
  .inputValidator((data:{ userId: string, namespaceId: string}) => data)
  .handler(async ({ data }) => {
        const request = lookupResources().subject("user", data.userId)
      .permission("view").resourceType("form").build()
    //  const request = lookupResources().subject("namespace", data.namespaceId)
    //   .permission("view").resourceType("form")
    //   .context({"user": data.userId})
    //   .build()
    const userIds = [];
    const stream = serverPerm.lookupResources(request);
    for await (const response of stream) {
      console.log(response)
      if(response.result) userIds.push(response.result.resourceObjectId)
    }
  console.log(userIds)
    return userIds
  })

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
  return authFetcher.post<FormI>(`/forms?namespace_id=${namespaceID}`, formData);
});

/**
 * Fetches all Forms for the specified namespace from the server.
 * @param userId if the user that want to see the namespace forms
 * @returns A promise that resolves to an array of Form objects.
 */
export const getAllNamespaceFormsFn = createClientOnlyFn(async (
  namespaceId: string, userId: string
) => {
  const ids = await getFormsIds({ data: {namespaceId, userId} })
  const res = await authFetcher.post<FormI[]>('/forms/bulk', { ids });
  return res.success ? res.data : []
});

/**
 * Query options for fetching all Form for a specific namespace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all Form.
 */
export const allNamespaceFormsQueryOptions = (namespaceId: string, userId: string) => {
  return queryOptions({
    queryKey: ['namespaces', namespaceId, "forms"],
    queryFn: () => getAllNamespaceFormsFn(namespaceId, userId),
  })
}
