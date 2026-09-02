import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { CreateCheckoutRequest } from "@trieoh/univents-api/schemas";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import { createCheckoutFn, refundPurchaseFn } from ".";
import { invalidatePurchaseCaches, syncCreatedCheckoutCache } from "./cache";

export function useCreateCheckoutMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      editionId,
      data,
    }: {
      editionId: string;
      data: CreateCheckoutRequest;
    }) => createCheckoutFn(editionId, data),
    onSuccess: (checkout) => syncCreatedCheckoutCache(queryClient, checkout),
  });
}

export function useRefundPurchaseMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ purchaseId }: { editionId: string; purchaseId: string }) =>
      refundPurchaseFn(purchaseId),
    onSuccess: (purchase) => {
      invalidatePurchaseCaches(
        queryClient,
        purchase.purchase_id,
        purchase.edition_id,
      );
      toast.success("Reembolso solicitado");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível solicitar o reembolso"),
      ),
  });
}
