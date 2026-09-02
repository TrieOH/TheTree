import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type {
  CreateInitialProductOutputI,
  ProductPatchOutputI,
  VariantCreateOutputI,
} from "../model";
import {
  removeProductCaches,
  removeVariantCaches,
  syncProductCaches,
  syncVariantCaches,
} from "./cache";
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

export function useCreateInitialProductMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateInitialProductInput) =>
      createInitialProductFn(data, editionId),
    onSuccess: (product) => {
      syncProductCaches(queryClient, product);
      void queryClient.invalidateQueries({
        queryKey: productKeys.storeStock(product.edition_id),
      });
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
      syncProductCaches(queryClient, product);
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
      removeProductCaches(
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
  editionId: string;
};

export function useCreateVariantMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId, data }: CreateVariantInput) =>
      createVariantFn(productId, data),
    onSuccess: (variant) => {
      syncVariantCaches(queryClient, variant);
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
      syncVariantCaches(queryClient, variant);
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
      removeVariantCaches(
        queryClient,
        variables.editionId,
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
