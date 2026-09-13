import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import { Effect } from "effect";
import { eventKeys } from "@/features/events/api/query-keys";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import {
  type PaymentProviderI,
  connectEventSellerFn,
  disconnectEventSellerFn,
} from "./index";

export const connectEventSellerEffect = (
  eventId: string,
  provider: PaymentProviderI,
) =>
  apiEffect(() =>
    withSpan("action:event-seller-connect", () =>
      connectEventSellerFn(eventId, provider),
    ),
  );

export const disconnectEventSellerEffect = (eventId: string) =>
  apiEffect(() =>
    withSpan("action:event-seller-disconnect", () =>
      disconnectEventSellerFn(eventId),
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
