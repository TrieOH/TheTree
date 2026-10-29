import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  completeEventPayments,
  connectEventPayments,
  disconnectEventPayments,
} from "@trieoh/univents-api";
import type {
  ConnectEventPaymentsResult,
  Event,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { eventKeys } from "@/features/events/api/query-keys";
import { syncEventCaches } from "@/features/events/api/cache";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import type { PaymentProviderI } from "./index";

export const completeEventSellerEffect = (
  eventId: string,
  sellerId: string,
  publicKey: string,
) =>
  apiEffect(() =>
    withSpan("action:event-seller-complete", () =>
      completeEventPayments(eventId, {
        seller_id: sellerId,
        public_key: publicKey,
      }).then(orvalData<Event>),
    ),
  );

export const useCompleteEventSellerMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    mutationEffect: ({ eventId, sellerId, publicKey }: {
      eventId: string;
      sellerId: string;
      publicKey: string;
    }) =>
      completeEventSellerEffect(eventId, sellerId, publicKey).pipe(
        Effect.tap((event) => Effect.sync(() => syncEventCaches(queryClient, event))),
      ),
  });
};

export const connectEventSellerEffect = (
  eventId: string,
  provider: PaymentProviderI,
) =>
  apiEffect(() =>
    withSpan("action:event-seller-connect", () =>
      connectEventPayments(eventId, { provider }).then(
        orvalData<ConnectEventPaymentsResult>,
      ),
    ),
  );

export const disconnectEventSellerEffect = (eventId: string) =>
  apiEffect(() =>
    withSpan("action:event-seller-disconnect", () =>
      disconnectEventPayments(eventId).then(orvalData<null>),
    ),
  );

export const useConnectEventSellerMutation = () => {
  return useEffectMutation({
    mutationEffect: ({
      eventId,
      provider,
    }: {
      eventId: string;
      provider: PaymentProviderI;
    }) =>
      connectEventSellerEffect(eventId, provider).pipe(
        Effect.tap(({ auth_url }) =>
          Effect.sync(() => {
            window.location.assign(auth_url);
          }),
        ),
      ),
  });
};

export const useDisconnectEventSellerMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    mutationEffect: (eventId: string) =>
      disconnectEventSellerEffect(eventId).pipe(
        Effect.tap(() =>
          Effect.tryPromise({
            try: () =>
              Promise.all([
                queryClient.invalidateQueries({ queryKey: eventKeys.lists() }),
                queryClient.invalidateQueries({
                  queryKey: eventKeys.details(),
                }),
              ]),
            catch: () => undefined,
          }),
        ),
      ),
  });
};
