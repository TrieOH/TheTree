import { useQueryClient } from "@trieoh/front-core-solid";
import type { CreateCheckoutRequest, CheckoutResult, Purchase } from "@trieoh/univents-api/schemas";
import { createEditionCheckout, refundPurchase } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { syncCreatedCheckoutCache } from "./cache";
import { purchaseKeys } from "./query-keys";

export const createCheckoutEffect = (
  editionId: string,
  data: CreateCheckoutRequest,
) =>
  apiEffect(() =>
    createEditionCheckout(editionId, data).then(orvalData<CheckoutResult>),
  );

export const useCreateCheckoutMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      editionId,
      data,
    }: {
      editionId: string;
      data: CreateCheckoutRequest;
    }) =>
      createCheckoutEffect(editionId, data).pipe(
        Effect.tap((checkout) =>
          Effect.sync(() => {
            syncCreatedCheckoutCache(queryClient, checkout);
          }),
        ),
      ),
  });
};

export interface RefundPurchaseInput {
  purchaseId: string;
  editionId?: string;
}

export const refundPurchaseEffect = (purchaseId: string) =>
  apiEffect(() =>
    refundPurchase(purchaseId).then(orvalData<Purchase>),
  );

export const useRefundPurchaseMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ purchaseId, editionId }: RefundPurchaseInput) =>
      refundPurchaseEffect(purchaseId).pipe(
        Effect.tap((purchase) =>
          Effect.sync(() => {
            if (editionId) {
              void queryClient.invalidateQueries({
                queryKey: purchaseKeys.edition(editionId),
              });
            }
            void queryClient.invalidateQueries({
              queryKey: purchaseKeys.detail(purchase.purchase_id),
            });
            void queryClient.invalidateQueries({
              queryKey: purchaseKeys.mine(),
            });
          }),
        ),
      ),
  });
};
