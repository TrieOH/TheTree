import type { QueryClient } from "@tanstack/react-query";
import type {
  Checkout,
  CheckoutResult,
  MyPurchases,
} from "@trieoh/univents-api/schemas";
import { badgeKeys } from "../../badges/api/query-keys";
import { productKeys } from "../../products/api/query-keys";
import { programKeys } from "../../programs/api/query-keys";
import { ticketKeys } from "../../tickets/api/query-keys";
import { purchaseQueryKeys } from "./query-keys";

export function syncCreatedCheckoutCache(
  queryClient: QueryClient,
  checkout: CheckoutResult,
) {
  queryClient.setQueryData<Checkout>(
    purchaseQueryKeys.detail(checkout.purchase_id),
    checkout,
  );
  queryClient.setQueryData<MyPurchases>(purchaseQueryKeys.mine(), (old) =>
    old
      ? {
          ...old,
          purchases: [
            checkout,
            ...old.purchases.filter(
              (purchase) => purchase.purchase_id !== checkout.purchase_id,
            ),
          ],
        }
      : old,
  );
  invalidatePurchaseDerivedCaches(queryClient, checkout.edition_id);
}

export function invalidatePurchaseCaches(
  queryClient: QueryClient,
  purchaseId: string,
  editionId?: string,
) {
  void queryClient.invalidateQueries({
    queryKey: purchaseQueryKeys.detail(purchaseId),
  });
  void queryClient.invalidateQueries({ queryKey: purchaseQueryKeys.mine() });
  if (editionId) {
    invalidatePurchaseDerivedCaches(queryClient, editionId);
  }
}

function invalidatePurchaseDerivedCaches(
  queryClient: QueryClient,
  editionId: string,
) {
  void queryClient.invalidateQueries({
    queryKey: purchaseQueryKeys.edition(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: ticketKeys.myTicket(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: programKeys.myParticipations(editionId),
  });
  void queryClient.invalidateQueries({ queryKey: badgeKeys.users() });
}
