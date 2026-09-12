import { orvalData } from "@trieoh/api-client";
import {
  listEditionProducts,
  listEditionStoreStock,
  listProductVariants,
} from "@trieoh/univents-api";
import type {
  Product,
  ProductVariant,
  StoreStockItem,
  TicketType,
  MyTicket,
} from "@trieoh/univents-api/schemas";
import { productKeys } from "./query-keys";
import type { QueryClient } from "@tanstack/query-core";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import type { EditionI } from "@/features/editions/model";
import { myTicketQueryOptions, ticketsQueryOptions } from "@/features/tickets/api";

export type StorePageData = {
  event: EventI;
  edition: EditionI | null;
  tickets?: TicketType[];
  products?: { product: Product; variants: ProductVariant[] }[];
  stock?: StoreStockItem[];
  heldTicket?: MyTicket | null;
};

export const productsQueryOptions = (editionId: string) => ({
  queryKey: productKeys.byEdition(editionId),
  queryFn: () =>
    listEditionProducts(editionId, { public: true }).then(orvalData<Product[]>),
});

export const productVariantsQueryOptions = (productId: string) => ({
  queryKey: productKeys.variants(productId),
  queryFn: () =>
    listProductVariants(productId, { public: true }).then(
      orvalData<ProductVariant[]>,
    ),
});

export const storeStockQueryOptions = (editionId: string) => ({
  queryKey: productKeys.stock(editionId),
  queryFn: () =>
    listEditionStoreStock(editionId, { public: true }).then(
      orvalData<StoreStockItem[]>,
    ),
});

export const storePageQueryOptions = (slug: string, authenticated: boolean, client: QueryClient) => ({
  queryKey: productKeys.storePage(slug, authenticated),
  queryFn: async (): Promise<StorePageData | null> => {
    const event = await publicEventBySlugQueryOptions(slug).queryFn();
    if (!event) return null;
    const editions = await allPublicEditionsQueryOptions(event.id).queryFn();
    const edition = editions.find((item) => new Date(item.ends_at).getTime() >= Date.now()) ?? editions.at(-1);
    if (!edition) return { event, edition: null };
    const [tickets, products, stock, heldTicket] = await Promise.all([
      client.fetchQuery(ticketsQueryOptions(edition.id)),
      client.fetchQuery(productsQueryOptions(edition.id)),
      client.fetchQuery(storeStockQueryOptions(edition.id)),
      authenticated ? client.fetchQuery(myTicketQueryOptions(edition.id)) : null,
    ]);
    const productItems = await Promise.all(products.map(async (product) => ({
      product,
      variants: await client.fetchQuery(productVariantsQueryOptions(product.id)),
    })));
    return { event, edition, tickets, products: productItems, stock, heldTicket };
  },
});
