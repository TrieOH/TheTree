import { useMutation, useQueryClient, type QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { ProductCreateOutputI, ProductI } from '../model'
import {
  createProductFn,
  patchProductFn,
  publishProductFn,
  restoreSoftDeletedProductFn,
  softDeleteProductFn,
} from './index'
import { productKeys } from './query-keys'

type CreateProductInput = {
  eventId: string
  editionId: string
  data: ProductCreateOutputI
}

type UpdateProductInput = {
  eventId: string
  editionId: string
  productId: string
  data: ProductCreateOutputI
}

type ProductActionInput = {
  eventId: string
  editionId: string
  productId: string
}

function upsertById(products: ProductI[] | undefined, product: ProductI) {
  const list = products ?? []
  const index = list.findIndex((item) => item.id === product.id)

  if (index === -1) return [...list, product]

  const next = [...list]
  next[index] = product
  return next
}

function syncProductCaches(queryClient: QueryClient, eventId: string, editionId: string, product: ProductI) {
  queryClient.setQueryData<ProductI[]>(
    productKeys.adminListByEdition(eventId, editionId),
    (old) => upsertById(old, product),
  )

  queryClient.setQueryData<ProductI[]>(
    productKeys.publicListByEdition(eventId, editionId),
    (old) => upsertById(old, product),
  )
}

function syncProductStatus(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  productId: string,
  status: ProductI['status'],
) {
  const patch = (products: ProductI[] | undefined) =>
    (products ?? []).map((product) => (
      product.id === productId
        ? { ...product, status }
        : product
    ))

  queryClient.setQueryData<ProductI[]>(productKeys.adminListByEdition(eventId, editionId), patch)
  queryClient.setQueryData<ProductI[]>(productKeys.publicListByEdition(eventId, editionId), patch)
}

function syncSoftDeletedProduct(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  productId: string,
  deletedAt: string | null,
) {
  const patch = (products: ProductI[] | undefined) =>
    (products ?? []).map((product) => (
      product.id === productId
        ? { ...product, deleted_at: deletedAt }
        : product
    ))

  queryClient.setQueryData<ProductI[]>(productKeys.adminListByEdition(eventId, editionId), patch)
  queryClient.setQueryData<ProductI[]>(productKeys.publicListByEdition(eventId, editionId), patch)
}

export function useCreateProductMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, data }: CreateProductInput) =>
      createProductFn(data, eventId, editionId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao criar produto')
        return
      }

      syncProductCaches(queryClient, variables.eventId, variables.editionId, res.data)
      toast.success('Produto criado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useUpdateProductMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, productId, data }: UpdateProductInput) =>
      patchProductFn(eventId, editionId, productId, data),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao atualizar produto')
        return
      }

      syncProductCaches(queryClient, variables.eventId, variables.editionId, res.data)
      toast.success('Produto atualizado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function usePublishProductMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, productId }: ProductActionInput) =>
      publishProductFn(eventId, editionId, productId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar produto')
        return
      }

      syncProductStatus(queryClient, variables.eventId, variables.editionId, variables.productId, 'available')
      toast.success('Produto publicado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useSoftDeleteProductMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, productId }: ProductActionInput) =>
      softDeleteProductFn(eventId, editionId, productId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao excluir produto')
        return
      }

      syncSoftDeletedProduct(
        queryClient,
        variables.eventId,
        variables.editionId,
        variables.productId,
        new Date().toISOString(),
      )
      toast.success('Produto excluído com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useRestoreSoftDeletedProductMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, productId }: ProductActionInput) =>
      restoreSoftDeletedProductFn(eventId, editionId, productId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao restaurar produto')
        return
      }

      syncSoftDeletedProduct(
        queryClient,
        variables.eventId,
        variables.editionId,
        variables.productId,
        null,
      )
      toast.success('Produto restaurado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}
