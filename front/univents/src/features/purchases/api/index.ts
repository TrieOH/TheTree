import { queryOptions, useMutation } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createEditionCheckout,
  getCheckout,
  getWsToken,
  listMyPurchases,
} from "@trieoh/univents-api";
import type {
  Checkout,
  CheckoutResult,
  CreateCheckoutRequest,
  MyPurchases,
  WsToken,
} from "@trieoh/univents-api/schemas";

export const getCheckoutFn = createClientOnlyFn((purchaseId: string) =>
  getCheckout(purchaseId).then(orvalData<Checkout>),
);

export const listMyPurchasesFn = createClientOnlyFn(() =>
  listMyPurchases().then(orvalData<MyPurchases>),
);

export const getWsTokenFn = createClientOnlyFn((purchaseId: string) =>
  getWsToken({ purchase_id: purchaseId }).then(orvalData<WsToken>),
);

export const createCheckoutFn = createClientOnlyFn(
  (editionId: string, data: CreateCheckoutRequest) =>
    createEditionCheckout(editionId, data).then(orvalData<CheckoutResult>),
);

export function useCreateCheckoutMutation() {
  return useMutation({
    mutationFn: ({
      editionId,
      data,
    }: {
      editionId: string;
      data: CreateCheckoutRequest;
    }) => createCheckoutFn(editionId, data),
  });
}

export const checkoutQueryOptions = (purchaseId: string) =>
  queryOptions({
    queryKey: ["purchases", purchaseId],
    queryFn: () => getCheckoutFn(purchaseId),
  });

export const myPurchasesQueryOptions = () =>
  queryOptions({
    queryKey: ["purchases"],
    queryFn: listMyPurchasesFn,
  });
