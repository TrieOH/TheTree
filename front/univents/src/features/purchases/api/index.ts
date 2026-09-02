import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createEditionCheckout,
  getCheckout,
  getWsToken,
  listEditionPurchases,
  listMyPurchases,
  refundPurchase,
} from "@trieoh/univents-api";
import type {
  Checkout,
  CheckoutResult,
  CreateCheckoutRequest,
  EditionPurchase,
  MyPurchases,
  Purchase,
  WsToken,
} from "@trieoh/univents-api/schemas";
import { purchaseQueryKeys } from "./query-keys";

export const getCheckoutFn = createClientOnlyFn((purchaseId: string) =>
  getCheckout(purchaseId).then(orvalData<Checkout>),
);

export const listMyPurchasesFn = createClientOnlyFn(() =>
  listMyPurchases().then(orvalData<MyPurchases>),
);

export const listEditionPurchasesFn = createClientOnlyFn((editionId: string) =>
  listEditionPurchases(editionId).then(orvalData<EditionPurchase[]>),
);

export const refundPurchaseFn = createClientOnlyFn((purchaseId: string) =>
  refundPurchase(purchaseId).then(orvalData<Purchase>),
);

export const getWsTokenFn = createClientOnlyFn((purchaseId: string) =>
  getWsToken({ purchase_id: purchaseId }).then(orvalData<WsToken>),
);

export const createCheckoutFn = createClientOnlyFn(
  (editionId: string, data: CreateCheckoutRequest) =>
    createEditionCheckout(editionId, data).then(orvalData<CheckoutResult>),
);

export const checkoutQueryOptions = (purchaseId: string) =>
  queryOptions({
    queryKey: purchaseQueryKeys.detail(purchaseId),
    queryFn: () => getCheckoutFn(purchaseId),
  });

export const myPurchasesQueryOptions = () =>
  queryOptions({
    queryKey: purchaseQueryKeys.mine(),
    queryFn: listMyPurchasesFn,
  });

export const editionPurchasesQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: purchaseQueryKeys.edition(editionId),
    queryFn: () => listEditionPurchasesFn(editionId),
  });
