import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { getCheckout, listMyPurchases } from "@trieoh/univents-api";
import type { Checkout, MyPurchases } from "@trieoh/univents-api/schemas";

export const getCheckoutFn = createClientOnlyFn((purchaseId: string) =>
  getCheckout(purchaseId).then(orvalData<Checkout>),
);

export const listMyPurchasesFn = createClientOnlyFn(() =>
  listMyPurchases().then(orvalData<MyPurchases>),
);

export const checkoutQueryOptions = (purchaseId: string) =>
  queryOptions({
    queryKey: ["purchases", purchaseId],
    queryFn: () => getCheckoutFn(purchaseId),
    // ponytail: polling is temporary; replace it with the split-6 WebSocket.
    refetchInterval: (query) =>
      query.state.data?.status === "pending" ? 5_000 : false,
  });

export const myPurchasesQueryOptions = () =>
  queryOptions({
    queryKey: ["purchases"],
    queryFn: listMyPurchasesFn,
  });
