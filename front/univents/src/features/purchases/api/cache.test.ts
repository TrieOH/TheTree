import { QueryClient } from "@tanstack/react-query";
import type { CheckoutResult, MyPurchases } from "@trieoh/univents-api/schemas";
import { describe, expect, it } from "vitest";
import { badgeKeys } from "../../badges/api/query-keys";
import { productKeys } from "../../products/api/query-keys";
import { programKeys } from "../../programs/api/query-keys";
import { ticketKeys } from "../../tickets/api/query-keys";
import { invalidatePurchaseCaches, syncCreatedCheckoutCache } from "./cache";
import { purchaseQueryKeys } from "./query-keys";

const checkout = (overrides: Partial<CheckoutResult> = {}): CheckoutResult => ({
  purchase_id: "purchase-1",
  edition_id: "edition-1",
  status: "pending",
  total_cents: 1_000,
  currency: "BRL",
  expires_at: "2026-01-01T00:10:00Z",
  items: [],
  ws_token: "token",
  ...overrides,
});

describe("purchase cache synchronization", () => {
  it("seeds the complete checkout without creating a partial purchase list", () => {
    const queryClient = new QueryClient();
    const created = checkout();

    syncCreatedCheckoutCache(queryClient, created);

    expect(
      queryClient.getQueryData(purchaseQueryKeys.detail("purchase-1")),
    ).toEqual(created);
    expect(queryClient.getQueryData(purchaseQueryKeys.mine())).toBeUndefined();
  });

  it("prepends the checkout to an already loaded purchase list", () => {
    const queryClient = new QueryClient();
    const created = checkout();
    const previous = checkout({ purchase_id: "purchase-2" });
    queryClient.setQueryData<MyPurchases>(purchaseQueryKeys.mine(), {
      purchases: [previous, created],
    });

    syncCreatedCheckoutCache(queryClient, created);

    expect(
      queryClient.getQueryData<MyPurchases>(purchaseQueryKeys.mine()),
    ).toEqual({ purchases: [created, previous] });
  });

  it("invalidates stock and buyer resources affected by a checkout", () => {
    const queryClient = new QueryClient();
    const keys = [
      purchaseQueryKeys.edition("edition-1"),
      productKeys.storeStock("edition-1"),
      ticketKeys.myTicket("edition-1"),
      programKeys.myParticipations("edition-1"),
      badgeKeys.user("actor-1"),
    ];
    for (const key of keys) queryClient.setQueryData(key, {});

    syncCreatedCheckoutCache(queryClient, checkout());

    for (const key of keys) {
      expect(queryClient.getQueryState(key)?.isInvalidated).toBe(true);
    }
  });

  it("invalidates purchase reads and derived resources after realtime updates", () => {
    const queryClient = new QueryClient();
    const keys = [
      purchaseQueryKeys.detail("purchase-1"),
      purchaseQueryKeys.mine(),
      purchaseQueryKeys.edition("edition-1"),
      productKeys.storeStock("edition-1"),
      ticketKeys.myTicket("edition-1"),
      programKeys.myParticipations("edition-1"),
      badgeKeys.user("actor-1"),
    ];
    for (const key of keys) queryClient.setQueryData(key, {});

    invalidatePurchaseCaches(queryClient, "purchase-1", "edition-1");

    for (const key of keys) {
      expect(queryClient.getQueryState(key)?.isInvalidated).toBe(true);
    }
  });
});
