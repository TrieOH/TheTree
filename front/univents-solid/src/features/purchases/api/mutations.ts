import { useMutation, useQueryClient } from "@trieoh/front-core-solid";
import type { CreateCheckoutRequest } from "@trieoh/univents-api/schemas";
import { createEditionCheckout } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { CheckoutResult } from "@trieoh/univents-api/schemas";
import { syncCreatedCheckoutCache } from "./cache";

export const useCreateCheckoutMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: async ({ editionId, data }: { editionId: string; data: CreateCheckoutRequest }) => {
    const checkout = await createEditionCheckout(editionId, data).then(orvalData<CheckoutResult>);
    syncCreatedCheckoutCache(queryClient, checkout);
    return checkout;
  }});
};
