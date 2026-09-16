import { orvalData } from "@trieoh/api-client";
import {
  createInitialProduct,
  createProductVariant,
  deleteProduct,
  deleteProductVariant,
  listEditionProducts,
  listEditionStoreStock,
  listProductVariants,
  patchProduct,
  patchProductVariant,
} from "@trieoh/univents-api";
import type {
  CreateInitialProductRequest,
  CreateProductVariantRequest,
  PatchProductRequest,
  PatchProductVariantRequest,
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
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
  VariantCreateOutputI,
  VariantI,
} from "../model";

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
    listEditionProducts(editionId).then(orvalData<ProductI[]>),
});

export const productsByEditionQueryOptions = productsQueryOptions;

export const productVariantsQueryOptions = (productId: string) => ({
  queryKey: productKeys.variants(productId),
  queryFn: () =>
    listProductVariants(productId).then(
      orvalData<VariantI[]>,
    ),
});

export const storeStockQueryOptions = (editionId: string) => ({
  queryKey: productKeys.stock(editionId),
  queryFn: () =>
    listEditionStoreStock(editionId, { public: true }).then(
      orvalData<StoreStockItem[]>,
    ),
});

export const storePageQueryOptions = (
  slug: string,
  authenticated: boolean,
  client: QueryClient,
) => ({
  queryKey: productKeys.storePage(slug, authenticated),
  queryFn: async (): Promise<StorePageData | null> => {
    const event = await publicEventBySlugQueryOptions(slug).queryFn();
    if (!event) return null;
    const editions = await allPublicEditionsQueryOptions(event.id).queryFn();
    const edition =
      editions.find((item) => new Date(item.ends_at).getTime() >= Date.now()) ??
      editions.at(-1);
    if (!edition) return { event, edition: null };
    const [tickets, products, stock, heldTicket] = await Promise.all([
      client.fetchQuery(ticketsQueryOptions(edition.id)),
      client.fetchQuery(productsQueryOptions(edition.id)),
      client.fetchQuery(storeStockQueryOptions(edition.id)),
      authenticated
        ? client.fetchQuery(myTicketQueryOptions(edition.id))
        : null,
    ]);
    const productItems = await Promise.all(
      products.map(async (product) => ({
        product,
        variants: await client.fetchQuery(
          productVariantsQueryOptions(product.id),
        ),
      })),
    );
    return {
      event,
      edition,
      tickets,
      products: productItems,
      stock,
      heldTicket,
    };
  },
});

export const createInitialProductFn = (
  data: CreateInitialProductOutputI,
  editionId: string,
) =>
  createInitialProduct(
    editionId,
    data as unknown as CreateInitialProductRequest,
  ).then(orvalData<ProductI>);

export const patchProductFn = (
  productId: string,
  data: ProductPatchOutputI,
) =>
  patchProduct(
    productId,
    data as unknown as PatchProductRequest,
  ).then(orvalData<ProductI>);

export const deleteProductFn = (productId: string) =>
  deleteProduct(productId).then(orvalData<{ success?: boolean }>);

export const createVariantFn = (
  productId: string,
  data: VariantCreateOutputI,
) =>
  createProductVariant(
    productId,
    data as unknown as CreateProductVariantRequest,
  ).then(orvalData<VariantI>);

export const patchVariantFn = (
  variantId: string,
  data: Partial<VariantCreateOutputI>,
) =>
  patchProductVariant(
    variantId,
    data as unknown as PatchProductVariantRequest,
  ).then(orvalData<VariantI>);

export const deleteVariantFn = (variantId: string) =>
  deleteProductVariant(variantId).then(orvalData<{ success?: boolean }>);
