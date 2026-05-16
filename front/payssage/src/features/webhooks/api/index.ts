import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import type { WebhookCreateI, WebhookCreateResponseI, WebhookI } from "../model";
import { queryOptions } from "@tanstack/react-query";


/**
 * Register a new Webhook for the specified workspace on the server.
 * @param webhookData - The data for the new webhook.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to the API response containing the newly created webhook endpoint.
 */
export const registerWebhookOnWorkspaceFn = createClientOnlyFn((
  webhookData: WebhookCreateI,
  workspaceName: string
) => {
  return authFetcher.post<WebhookCreateResponseI>(
    `/workspaces/${workspaceName}/webhooks`,
    webhookData
  );
});


/**
 * Fetches all webhooks for the specified workspace from the server.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to an array of Webhook objects.
 */
export const getAllWorkspaceWebhooksFn = createClientOnlyFn(async (workspaceName: string) => {
  try {
    return await tanstackQueryFetcher<WebhookI[]>(`/workspaces/${workspaceName}/webhooks`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all webhooks for a specific workspace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all webhooks.
 */
export const allWorkspaceWebhooksQueryOptions = (workspaceName: string) => {
  return queryOptions({
    queryKey: ['workspaces', workspaceName, "webhooks"],
    queryFn: () => getAllWorkspaceWebhooksFn(workspaceName),
  })
}


/**
 * Delete a webhook for the specified workspace on the server.
 * @param workspaceName - The workspace name
 * @param endpointId - The webhook endpoint ID
 * @returns A promise that resolves to the API response(void).
 */
export const deleteWebhookOnWorkspaceFn = createClientOnlyFn((
  workspaceName: string,
  endpointId: string
) => {
  return authFetcher.delete<void>(
    `/workspaces/${workspaceName}/webhooks/${endpointId}`
  );
});