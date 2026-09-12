import type { QueryClient } from "@trieoh/front-core-solid";
import type { Checkout, CheckoutResult, MyPurchases } from "@trieoh/univents-api/schemas";
import { purchaseKeys } from "./query-keys";

export function syncCreatedCheckoutCache(queryClient: QueryClient, checkout: CheckoutResult) {
  queryClient.setQueryData<Checkout>(purchaseKeys.detail(checkout.purchase_id), checkout);
  queryClient.setQueryData<MyPurchases>(purchaseKeys.mine(), (old) =>
    old ? { ...old, purchases: [checkout, ...old.purchases.filter((item) => item.purchase_id !== checkout.purchase_id)] } : old,
  );
  void queryClient.invalidateQueries({ queryKey: ["events", "catalog"] });
  void queryClient.invalidateQueries({ queryKey: ["products", "store-page"] });
}
