import { createClientOnlyFn } from "@tanstack/react-start";
import type { PaymentIntentsI } from "../model";
import { tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { queryOptions } from "@tanstack/react-query";

/**
 * Fetches all payment intents for the specified workspace from the server.
 * @returns A promise that resolves to an array of PaymentIntent objects.
 */
export const getAllPaymentIntentsOnWorkspaceFn = createClientOnlyFn(async (
  name: string
) => {
  try {
    return await tanstackQueryFetcher<PaymentIntentsI[]>(`/workspaces/${name}/intents`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all payment intents for a specific workspace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all payment intents.
 */
export const allWorkspacePaymentIntentsQueryOptions = (name: string) => {
  return queryOptions({
    queryKey: ['workspaces', name, "intents"],
    queryFn: () => getAllPaymentIntentsOnWorkspaceFn(name),
  })
}
