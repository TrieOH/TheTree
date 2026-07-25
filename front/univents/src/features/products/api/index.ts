import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { CreateInitialProductOutputI, ProductI, ProductPatchOutputI, VariantCreateOutputI, VariantI } from "../model";
import { authFetcher, publicQueryFetcher } from "@/shared/lib/api/fetch";
import { productKeys } from "./query-keys";

// ──────────── Products ────────────

/** GET /editions/{edition_id}/products */
export const getProductsByEditionFn = createClientOnlyFn(async (editionId: string) => {
  return publicQueryFetcher<ProductI[]>(`/editions/${editionId}/products`);
});

export const productsByEditionQueryOptions = (editionId: string) => {
  return queryOptions({
    queryKey: productKeys.editionList(editionId),
    queryFn: () => getProductsByEditionFn(editionId),
  })
};

/** GET /editions/{edition_id}/products/{vendor_code}:by-code */
export const getProductByVendorCodeFn = createClientOnlyFn(async (editionId: string, vendorCode: string) => {
  return publicQueryFetcher<ProductI>(`/editions/${editionId}/products/${vendorCode}:by-code`);
});

export const productByVendorCodeQueryOptions = (editionId: string, vendorCode: string) => {
  return queryOptions({
    queryKey: productKeys.byVendorCode(editionId, vendorCode),
    queryFn: () => getProductByVendorCodeFn(editionId, vendorCode),
  })
};

/** GET /products/{product_id} */
export const getProductByIdFn = createClientOnlyFn(async (productId: string) => {
  return publicQueryFetcher<ProductI>(`/products/${productId}`);
});

export const productByIdQueryOptions = (productId: string) => {
  return queryOptions({
    queryKey: productKeys.detail(productId),
    queryFn: () => getProductByIdFn(productId),
  })
};

/** GET /products/{product_id}/variants */
export const getProductVariantsFn = createClientOnlyFn(async (productId: string) => {
  return publicQueryFetcher<VariantI[]>(`/products/${productId}/variants`);
});

export const productVariantsQueryOptions = (productId: string) => {
  return queryOptions({
    queryKey: productKeys.variants.byProduct(productId),
    queryFn: () => getProductVariantsFn(productId),
  })
};

/** POST /editions/{edition_id}/products — Create product + first variant */
export const createInitialProductFn = createClientOnlyFn((
  data: CreateInitialProductOutputI,
  editionId: string
) => {
  return authFetcher.post<ProductI>(`/editions/${editionId}/products`, data);
});

/** PATCH /products/{product_id} */
export const patchProductFn = createClientOnlyFn((
  productId: string,
  data: ProductPatchOutputI
) => {
  return authFetcher.patch<ProductI>(`/products/${productId}`, data);
});

/** DELETE /products/{product_id} */
export const deleteProductFn = createClientOnlyFn((productId: string) => {
  return authFetcher.delete<null>(`/products/${productId}`);
});

// ──────────── Variants ────────────

/** GET /editions/{edition_id}/variants/{vendor_code}:by-code */
export const getVariantByVendorCodeFn = createClientOnlyFn(async (editionId: string, vendorCode: string) => {
  return publicQueryFetcher<VariantI>(`/editions/${editionId}/variants/${vendorCode}:by-code`);
});

export const variantByVendorCodeQueryOptions = (editionId: string, vendorCode: string) => {
  return queryOptions({
    queryKey: productKeys.variants.byVendorCode(editionId, vendorCode),
    queryFn: () => getVariantByVendorCodeFn(editionId, vendorCode),
  })
};

/** POST /products/{product_id}/variants */
export const createVariantFn = createClientOnlyFn((
  productId: string, data: VariantCreateOutputI
) => {
  return authFetcher.post<VariantI>(`/products/${productId}/variants`, data);
});

/** PATCH /variants/{variant_id} */
export const patchVariantFn = createClientOnlyFn((variantId: string, data: VariantCreateOutputI) => {
  return authFetcher.patch<VariantI>(`/variants/${variantId}`, data);
});

/** DELETE /variants/{variant_id} */
export const deleteVariantFn = createClientOnlyFn((variantId: string) => {
  return authFetcher.delete<null>(`/variants/${variantId}`);
});

// ──────────── WebSocket ────────────

export const getWebsocketAuthToken = createClientOnlyFn(() => {
  return authFetcher.get<{ token: string }>("/ws/token");
});