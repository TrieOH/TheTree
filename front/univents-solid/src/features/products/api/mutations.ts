import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createInitialProduct,
  createProductVariant,
  deleteProduct,
  deleteProductVariant,
  patchProduct,
  patchProductVariant,
} from "@trieoh/univents-api";
import type {
  CreateInitialProductRequest,
  CreateProductVariantRequest,
  PatchProductRequest,
  PatchProductVariantRequest,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
  VariantCreateOutputI,
  VariantI,
} from "../model";
import { productKeys } from "./query-keys";

export interface CreateInitialProductInput {
  editionId: string;
  data: CreateInitialProductOutputI;
}

export interface UpdateProductInput {
  productId: string;
  editionId: string;
  data: ProductPatchOutputI;
}

export interface DeleteProductInput {
  productId: string;
  editionId: string;
}

export interface CreateVariantInput {
  productId: string;
  data: VariantCreateOutputI;
}

export interface UpdateVariantInput {
  variantId: string;
  productId: string;
  data: Partial<VariantCreateOutputI>;
}

export interface DeleteVariantInput {
  variantId: string;
  productId: string;
}

export const createInitialProductEffect = (
  editionId: string,
  data: CreateInitialProductOutputI,
) =>
  apiEffect(() =>
    withSpan("action:product-create", () =>
      createInitialProduct(
        editionId,
        data as unknown as CreateInitialProductRequest,
      ).then(orvalData<ProductI>),
    ),
  );

export const updateProductEffect = (
  productId: string,
  data: ProductPatchOutputI,
) =>
  apiEffect(() =>
    withSpan("action:product-patch", () =>
      patchProduct(
        productId,
        data as unknown as PatchProductRequest,
      ).then(orvalData<ProductI>),
    ),
  );

export const deleteProductEffect = (productId: string) =>
  apiEffect(() =>
    withSpan("action:product-delete", () =>
      deleteProduct(productId).then(orvalData<{ success?: boolean }>),
    ),
  );

export const createVariantEffect = (
  productId: string,
  data: VariantCreateOutputI,
) =>
  apiEffect(() =>
    withSpan("action:variant-create", () =>
      createProductVariant(
        productId,
        data as unknown as CreateProductVariantRequest,
      ).then(orvalData<VariantI>),
    ),
  );

export const updateVariantEffect = (
  variantId: string,
  data: Partial<VariantCreateOutputI>,
) =>
  apiEffect(() =>
    withSpan("action:variant-patch", () =>
      patchProductVariant(
        variantId,
        data as unknown as PatchProductVariantRequest,
      ).then(orvalData<VariantI>),
    ),
  );

export const deleteVariantEffect = (variantId: string) =>
  apiEffect(() =>
    withSpan("action:variant-delete", () =>
      deleteProductVariant(variantId).then(orvalData<{ success?: boolean }>),
    ),
  );

export const useCreateInitialProductMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateInitialProductInput) =>
      createInitialProductEffect(editionId, data).pipe(
        Effect.tap((_product: ProductI) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: productKeys.stock(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: ["events", "catalog"],
            });
          }),
        ),
      ),
  });
};

export const useUpdateProductMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ productId, editionId, data }: UpdateProductInput) =>
      updateProductEffect(productId, data).pipe(
        Effect.tap((_product: ProductI) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: productKeys.detail(productId),
            });
          }),
        ),
      ),
  });
};

export const useDeleteProductMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ productId, editionId }: DeleteProductInput) =>
      deleteProductEffect(productId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: productKeys.stock(editionId),
            });
          }),
        ),
      ),
  });
};

export const useCreateVariantMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ productId, data }: CreateVariantInput) =>
      createVariantEffect(productId, data).pipe(
        Effect.tap((_variant: VariantI) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.variants(productId),
            });
          }),
        ),
      ),
  });
};

export const useUpdateVariantMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ variantId, productId, data }: UpdateVariantInput) =>
      updateVariantEffect(variantId, data).pipe(
        Effect.tap((_variant: VariantI) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.variants(productId),
            });
          }),
        ),
      ),
  });
};

export const useDeleteVariantMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ variantId, productId }: DeleteVariantInput) =>
      deleteVariantEffect(variantId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: productKeys.variants(productId),
            });
          }),
        ),
      ),
  });
};
