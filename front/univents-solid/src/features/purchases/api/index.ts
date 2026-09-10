import type { QueryClient } from "@tanstack/query-core";
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
  Purchase,
  WsToken,
} from "@trieoh/univents-api/schemas";
import { resolvePurchaseCatalog } from "./purchase-catalog";
import { purchaseKeys } from "./query-keys";

export const myPurchasesQueryOptions = () => ({
  queryKey: purchaseKeys.mine(),
  queryFn: () => listMyPurchases().then(orvalData<MyPurchases>),
});

export const checkoutQueryOptions = (purchaseId: string) => ({
  queryKey: purchaseKeys.detail(purchaseId),
  queryFn: () => getCheckout(purchaseId).then(orvalData<Checkout>),
});

export const purchaseCatalogQueryOptions = (purchase: Purchase) => ({
  queryKey: purchaseKeys.catalog(purchase.purchase_id),
  queryFn: ({ client }: { client: QueryClient }) =>
    resolvePurchaseCatalog(
      client,
      purchase.items.map((item) => ({
        item_type: item.item_type,
        item_id: item.item_id,
        edition_id: purchase.edition_id,
      })),
    ),
});

export const createCheckout = (editionId: string, data: CreateCheckoutRequest) =>
  createEditionCheckout(editionId, data).then(orvalData<CheckoutResult>);

export const purchaseWsToken = (purchaseId: string) =>
  getWsToken({ purchase_id: purchaseId }).then(orvalData<WsToken>);
