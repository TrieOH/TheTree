import { useQueryClient } from "@trieoh/front-core-solid";
import type { CreateCheckoutRequest, CheckoutResult } from "@trieoh/univents-api/schemas";
import { createEditionCheckout } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { syncCreatedCheckoutCache } from "./cache";

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
