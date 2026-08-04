import {
  type QueryClient,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
  VariantCreateOutputI,
  VariantI,
} from "../model";
import {
  createInitialProductFn,
  createVariantFn,
  deleteProductFn,
  deleteVariantFn,
  patchProductFn,
  patchVariantFn,
} from "./index";
import { productKeys } from "./query-keys";

type CreateInitialProductInput = {
  editionId: string;
  data: CreateInitialProductOutputI;
};

type UpdateProductInput = {
  productId: string;
  data: ProductPatchOutputI;
};

type DeleteProductInput = {
  productId: string;
  editionId: string;
};

function upsertById<T extends { id: string }>(list: T[] | undefined, item: T) {
  const items = list ?? [];
  const index = items.findIndex((el) => el.id === item.id);

  if (index === -1) return [...items, item];

  const next = [...items];
  next[index] = item;
  return next;
}

function removeById<T extends { id: string }>(
  list: T[] | undefined,
  id: string,
) {
  const items = list ?? [];
  return items.filter((el) => el.id !== id);
}

function syncProductCaches(
  queryClient: QueryClient,
  editionId: string,
  product: ProductI,
) {
  queryClient.setQueryData<ProductI[]>(
    productKeys.editionList(editionId),
    (old) => upsertById(old, product),
  );
}

function removeProductFromCache(
  queryClient: QueryClient,
  editionId: string,
  productId: string,
) {
  queryClient.setQueryData<ProductI[]>(
    productKeys.editionList(editionId),
    (old) => removeById(old, productId),
  );
}

function syncVariantCaches(
  queryClient: QueryClient,
  productId: string,
  variant: VariantI,
) {
  queryClient.setQueryData<VariantI[]>(
    productKeys.variants.byProduct(productId),
    (old) => upsertById(old, variant),
  );
}

function removeVariantFromCache(
  queryClient: QueryClient,
  productId: string,
  variantId: string,
) {
  queryClient.setQueryData<VariantI[]>(
    productKeys.variants.byProduct(productId),
    (old) => removeById(old, variantId),
  );
}

export function useCreateInitialProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateInitialProductInput) =>
      createInitialProductFn(data, editionId),
    onSuccess: (product, variables) => {
      syncProductCaches(queryClient, variables.editionId, product);
      toast.success("Produto criado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar o produto")),
  });
}

export function useUpdateProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId, data }: UpdateProductInput) =>
      patchProductFn(productId, data),
    onSuccess: (product) => {
      syncProductCaches(queryClient, product.edition_id, product);
      toast.success("Produto atualizado com sucesso!");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar o produto"),
      ),
  });
}

export function useDeleteProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId }: DeleteProductInput) =>
      deleteProductFn(productId),
    onSuccess: (_res, variables) => {
      removeProductFromCache(
        queryClient,
        variables.editionId,
        variables.productId,
      );
      toast.success("Produto excluído com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível excluir o produto")),
  });
}

type CreateVariantInput = {
  productId: string;
  data: VariantCreateOutputI;
};

type UpdateVariantInput = {
  variantId: string;
  data: VariantCreateOutputI;
};

type DeleteVariantInput = {
  variantId: string;
  productId: string;
};

export function useCreateVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId, data }: CreateVariantInput) =>
      createVariantFn(productId, data),
    onSuccess: (variant, variables) => {
      syncVariantCaches(queryClient, variables.productId, variant);
      toast.success("Variação criada com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar a variação")),
  });
}

export function useUpdateVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ variantId, data }: UpdateVariantInput) =>
      patchVariantFn(variantId, data),
    onSuccess: (variant) => {
      syncVariantCaches(queryClient, variant.product_id, variant);
      toast.success("Variação atualizada com sucesso!");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar a variação"),
      ),
  });
}

export function useDeleteVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ variantId }: DeleteVariantInput) =>
      deleteVariantFn(variantId),
    onSuccess: (_res, variables) => {
      removeVariantFromCache(
        queryClient,
        variables.productId,
        variables.variantId,
      );
      toast.success("Variação excluída com sucesso!");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir a variação"),
      ),
  });
}
