import type { QueryClient } from "@tanstack/query-core";
import { orvalData } from "@trieoh/api-client";
import {
  createEditionCheckout,
  getCheckout,
  getWsToken,
  listMyPurchases,
} from "@trieoh/univents-api";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import type { EditionI } from "@/features/editions/model";
import { myTicketQueryOptions } from "@/features/tickets/api";
import type {
  Checkout,
  CheckoutResult,
  CreateCheckoutRequest,
  MyPurchases,
  Purchase,
  WsToken,
  MyTicket,
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

export type CheckoutPageData = {
  event: EventI;
  edition: EditionI;
  heldTicket: MyTicket | null;
};

export const checkoutPageQueryOptions = (slug: string) => ({
  queryKey: purchaseKeys.page(slug, true),
  queryFn: async (): Promise<CheckoutPageData | null> => {
    const event = await publicEventBySlugQueryOptions(slug).queryFn();
    if (!event) return null;
    const editions = await allPublicEditionsQueryOptions(event.id).queryFn();
    const now = Date.now();
    const sorted = [...editions].sort((a, b) => a.starts_at.localeCompare(b.starts_at));
    const edition = sorted.find((item) => new Date(item.starts_at).getTime() <= now && new Date(item.ends_at).getTime() >= now)
      ?? sorted.find((item) => new Date(item.starts_at).getTime() > now)
      ?? sorted.at(-1);
    if (!edition) return null;
    const heldTicket = await myTicketQueryOptions(edition.id).queryFn();
    return { event, edition, heldTicket };
  },
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
