import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createInitialProduct,
  createProductVariant,
  deleteProduct,
  deleteProductVariant,
  getProduct,
  getProductByVendorCode,
  getVariantByVendorCode,
  listEditionProducts,
  listProductVariants,
  patchProduct,
  patchProductVariant,
} from "@trieoh/univents-api";
import type {
  CreateInitialProductRequest,
  CreateProductVariantRequest,
  PatchProductRequest,
  PatchProductVariantRequest,
} from "@trieoh/univents-api/schemas";
// import { authFetcher } from "@/shared/lib/api/fetch";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
  VariantCreateOutputI,
  VariantI,
} from "../model";
import { productKeys } from "./query-keys";

// ──────────── Products ────────────

/** GET /editions/{edition_id}/products */
export const getProductsByEditionFn = createClientOnlyFn(
  async (editionId: string) => {
    return listEditionProducts(editionId, { public: true }).then(
      orvalData<ProductI[]>,
    );
  },
);

export const productsByEditionQueryOptions = (editionId: string) => {
  return queryOptions({
    queryKey: productKeys.editionList(editionId),
    queryFn: () => getProductsByEditionFn(editionId),
  });
};

/** GET /editions/{edition_id}/products/{vendor_code}:by-code */
export const getProductByVendorCodeFn = createClientOnlyFn(
  async (editionId: string, vendorCode: string) => {
    return getProductByVendorCode(editionId, vendorCode, { public: true }).then(
      orvalData<ProductI>,
    );
  },
);

export const productByVendorCodeQueryOptions = (
  editionId: string,
  vendorCode: string,
) => {
  return queryOptions({
    queryKey: productKeys.byVendorCode(editionId, vendorCode),
    queryFn: () => getProductByVendorCodeFn(editionId, vendorCode),
  });
};

/** GET /products/{product_id} */
export const getProductByIdFn = createClientOnlyFn(
  async (productId: string) => {
    return getProduct(productId, { public: true }).then(orvalData<ProductI>);
  },
);

export const productByIdQueryOptions = (productId: string) => {
  return queryOptions({
    queryKey: productKeys.detail(productId),
    queryFn: () => getProductByIdFn(productId),
  });
};

/** GET /products/{product_id}/variants */
export const getProductVariantsFn = createClientOnlyFn(
  async (productId: string) => {
    return listProductVariants(productId, { public: true }).then(
      orvalData<VariantI[]>,
    );
  },
);

export const productVariantsQueryOptions = (productId: string) => {
  return queryOptions({
    queryKey: productKeys.variants.byProduct(productId),
    queryFn: () => getProductVariantsFn(productId),
  });
};

/** POST /editions/{edition_id}/products — Create product + first variant */
export const createInitialProductFn = createClientOnlyFn(
  (data: CreateInitialProductOutputI, editionId: string) => {
    return createInitialProduct(
      editionId,
      data as CreateInitialProductRequest,
    ).then(orvalData<ProductI>);
  },
);

/** PATCH /products/{product_id} */
export const patchProductFn = createClientOnlyFn(
  (productId: string, data: ProductPatchOutputI) => {
    return patchProduct(productId, data as PatchProductRequest).then(
      orvalData<ProductI>,
    );
  },
);

/** DELETE /products/{product_id} */
export const deleteProductFn = createClientOnlyFn((productId: string) => {
  return deleteProduct(productId).then(orvalData<null>);
});

// ──────────── Variants ────────────

/** GET /editions/{edition_id}/variants/{vendor_code}:by-code */
export const getVariantByVendorCodeFn = createClientOnlyFn(
  async (editionId: string, vendorCode: string) => {
    return getVariantByVendorCode(editionId, vendorCode, { public: true }).then(
      orvalData<VariantI>,
    );
  },
);

export const variantByVendorCodeQueryOptions = (
  editionId: string,
  vendorCode: string,
) => {
  return queryOptions({
    queryKey: productKeys.variants.byVendorCode(editionId, vendorCode),
    queryFn: () => getVariantByVendorCodeFn(editionId, vendorCode),
  });
};

/** POST /products/{product_id}/variants */
export const createVariantFn = createClientOnlyFn(
  (productId: string, data: VariantCreateOutputI) => {
    return createProductVariant(
      productId,
      data as CreateProductVariantRequest,
    ).then(orvalData<VariantI>);
  },
);

/** PATCH /variants/{variant_id} */
export const patchVariantFn = createClientOnlyFn(
  (variantId: string, data: VariantCreateOutputI) => {
    return patchProductVariant(
      variantId,
      data as PatchProductVariantRequest,
    ).then(orvalData<VariantI>);
  },
);

/** DELETE /variants/{variant_id} */
export const deleteVariantFn = createClientOnlyFn((variantId: string) => {
  return deleteProductVariant(variantId).then(orvalData<null>);
});

// ──────────── WebSocket ────────────

export const getWebsocketAuthToken = createClientOnlyFn(() => {
  console.warn("Not Implemeted");
  return { data: { token: "" }, success: false };
  // return authFetcher.get<{ token: string }>("/ws/token");
});
