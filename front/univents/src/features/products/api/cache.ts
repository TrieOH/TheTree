import type { QueryClient } from "@tanstack/react-query";
import { upsertById } from "../../../shared/lib/query-cache";
import type { ProductI, VariantI } from "../model";
import { productKeys } from "./query-keys";

function removeById<T extends { id: string }>(
  list: T[] | undefined,
  id: string,
) {
  return list?.filter((item) => item.id !== id);
}

export function syncProductCaches(queryClient: QueryClient, product: ProductI) {
  queryClient.setQueryData<ProductI[]>(
    productKeys.editionList(product.edition_id),
    (old) => upsertById(old, product),
  );
  queryClient.setQueryData<ProductI>(productKeys.detail(product.id), (old) =>
    old ? product : old,
  );
  void queryClient.invalidateQueries({
    queryKey: productKeys.vendorCodes(product.edition_id),
  });
}

export function removeProductCaches(
  queryClient: QueryClient,
  editionId: string,
  productId: string,
) {
  queryClient.setQueryData<ProductI[]>(
    productKeys.editionList(editionId),
    (old) => removeById(old, productId),
  );
  queryClient.removeQueries({ queryKey: productKeys.detail(productId) });
  queryClient.removeQueries({
    queryKey: productKeys.variants.byProduct(productId),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.vendorCodes(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.variants.vendorCodes(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(editionId),
  });
}

export function syncVariantCaches(queryClient: QueryClient, variant: VariantI) {
  queryClient.setQueryData<VariantI[]>(
    productKeys.variants.byProduct(variant.product_id),
    (old) => upsertById(old, variant),
  );
  void queryClient.invalidateQueries({
    queryKey: productKeys.variants.vendorCodes(variant.edition_id),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(variant.edition_id),
  });
}

export function removeVariantCaches(
  queryClient: QueryClient,
  editionId: string,
  productId: string,
  variantId: string,
) {
  queryClient.setQueryData<VariantI[]>(
    productKeys.variants.byProduct(productId),
    (old) => removeById(old, variantId),
  );
  void queryClient.invalidateQueries({
    queryKey: productKeys.variants.vendorCodes(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(editionId),
  });
}
