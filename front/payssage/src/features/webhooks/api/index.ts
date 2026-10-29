import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { itemsOf } from "@trieoh/api-client";
import { authFetcher } from "#/shared/lib/api/fetch";
import type {
  WebhookDelivery,
  WebhookEndpoint,
  WebhookEndpointCreateRequest,
  WebhookEvent,
} from "../model";

const unwrapList = <T>(response: {
  success: boolean;
  data?: unknown;
  message?: string;
}) => {
  if (!response.success) throw response;
  return itemsOf<T>(response.data);
};

export const createWebhookEndpointFn = createClientOnlyFn(
  (walletId: string, payload: WebhookEndpointCreateRequest) =>
    authFetcher.post<WebhookEndpoint>(
      `/wallets/${walletId}/webhooks/endpoints`,
      payload,
    ),
);

export const webhookEndpointsQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId, "webhooks", "endpoints"],
    queryFn: async () =>
      unwrapList<WebhookEndpoint>(
        await authFetcher.get<unknown>(
          `/wallets/${walletId}/webhooks/endpoints`,
        ),
      ),
  });

export const getWebhookEndpointFn = createClientOnlyFn((endpointId: string) =>
  authFetcher.get<WebhookEndpoint>(`/webhooks/endpoints/${endpointId}`),
);

export const deleteWebhookEndpointFn = createClientOnlyFn(
  (endpointId: string) =>
    authFetcher.delete<void>(`/webhooks/endpoints/${endpointId}`),
);

export const webhookDeliveriesQueryOptions = (endpointId: string) =>
  queryOptions({
    queryKey: ["webhooks", "endpoints", endpointId, "deliveries"],
    queryFn: async () =>
      unwrapList<WebhookDelivery>(
        await authFetcher.get<unknown>(
          `/webhooks/endpoints/${endpointId}/deliveries`,
        ),
      ),
  });

export const getWebhookDeliveryFn = createClientOnlyFn((deliveryId: string) =>
  authFetcher.get<WebhookDelivery>(`/webhooks/deliveries/${deliveryId}`),
);

export const webhookEventsQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId, "webhooks", "events"],
    queryFn: async () =>
      unwrapList<WebhookEvent>(
        await authFetcher.get<unknown>(`/wallets/${walletId}/webhooks/events`),
      ),
  });

export const getWebhookEventFn = createClientOnlyFn((eventId: string) =>
  authFetcher.get<WebhookEvent>(`/webhooks/events/${eventId}`),
);
