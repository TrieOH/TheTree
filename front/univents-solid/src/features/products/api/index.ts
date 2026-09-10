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
} from "@trieoh/univents-api/schemas";
import { productKeys } from "./query-keys";

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
