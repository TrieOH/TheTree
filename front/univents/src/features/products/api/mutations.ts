import {
  type QueryClient,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
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
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao criar produto");
        return;
      }

      syncProductCaches(queryClient, variables.editionId, res.data);
      toast.success("Produto criado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useUpdateProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId, data }: UpdateProductInput) =>
      patchProductFn(productId, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao atualizar produto");
        return;
      }

      const editionId = res.data.edition_id;
      syncProductCaches(queryClient, editionId, res.data);
      toast.success("Produto atualizado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useDeleteProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId }: DeleteProductInput) =>
      deleteProductFn(productId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao excluir produto");
        return;
      }

      removeProductFromCache(
        queryClient,
        variables.editionId,
        variables.productId,
      );
      toast.success("Produto excluído com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
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
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao criar variação");
        return;
      }

      syncVariantCaches(queryClient, variables.productId, res.data);
      toast.success("Variação criada com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useUpdateVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ variantId, data }: UpdateVariantInput) =>
      patchVariantFn(variantId, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao atualizar variação");
        return;
      }

      const productId = res.data.product_id;
      syncVariantCaches(queryClient, productId, res.data);
      toast.success("Variação atualizada com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useDeleteVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ variantId }: DeleteVariantInput) =>
      deleteVariantFn(variantId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao excluir variação");
        return;
      }

      removeVariantFromCache(
        queryClient,
        variables.productId,
        variables.variantId,
      );
      toast.success("Variação excluída com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
