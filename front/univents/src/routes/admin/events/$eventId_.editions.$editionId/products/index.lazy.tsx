import { createLazyFileRoute } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { Package, Plus } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import { Button } from '@/shared/ui/shadcn/button'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { allAdminProductsQueryOptions } from '@/features/products/api'
import {
  useCreateProductMutation,
  usePublishProductMutation,
  useRestoreSoftDeletedProductMutation,
  useSoftDeleteProductMutation,
  useUpdateProductMutation,
} from '@/features/products/api/mutations'
import type { ProductI } from '@/features/products/model'
import { AdminProductCard } from '@/features/products/ui/AdminProductCard'
import { ManageProductModal } from '@/features/products/ui/ManageProductModal'

const STATUS_SORT_ORDER: Record<ProductI['status'], number> = {
  draft: 0,
  available: 1,
  sold_out: 2,
  unavailable: 3,
}

export const Route = createLazyFileRoute('/admin/events/$eventId_/editions/$editionId/products/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const { data: products = [] } = useQuery(allAdminProductsQueryOptions(eventId, editionId))
  const createProductMutation = useCreateProductMutation()
  const updateProductMutation = useUpdateProductMutation()
  const publishProductMutation = usePublishProductMutation()
  const softDeleteProductMutation = useSoftDeleteProductMutation()
  const restoreSoftDeletedProductMutation = useRestoreSoftDeletedProductMutation()
  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<ProductI>>({
    field: 'name',
    direction: 'asc',
  })
  const [publishingProduct, setPublishingProduct] = useState<ProductI | null>(null)
  const [deletingProduct, setDeletingProduct] = useState<ProductI | null>(null)
  const [restoringProduct, setRestoringProduct] = useState<ProductI | null>(null)
  const [modalState, setModalState] = useState<{ open: boolean; product?: ProductI }>({
    open: false,
  })

  const filteredProducts = useMemo(() => {
    const search = filter.trim().toLowerCase()

    return [...products]
      .filter((product) => {
        if (!search) return true

        return [
          product.name,
          product.description ?? '',
          product.type,
          product.status,
        ].some((value) => value.toLowerCase().includes(search))
      })
      .sort((a, b) => {
        const direction = sort.direction === 'asc' ? 1 : -1

        if (sort.field === 'status') {
          return (STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status]) * direction
        }

        if (sort.field === 'price_cents') {
          return (a.price_cents - b.price_cents) * direction
        }

        if (sort.field === 'inventory_remaining') {
          return (a.inventory_remaining - b.inventory_remaining) * direction
        }

        return String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
      })
  }, [filter, products, sort])

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<ProductI>
        items={filteredProducts}
        layout="grid"
        minItemWidth="16rem"
        pageSize={4}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: 'name', label: 'Nome' },
          { key: 'price_cents', label: 'Preço', comparator: (a, b) => a.price_cents - b.price_cents },
          { key: 'status', label: 'Status', comparator: (a, b) => STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status] },
          { key: 'inventory_remaining', label: 'Estoque', comparator: (a, b) => a.inventory_remaining - b.inventory_remaining },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, descrição, tipo ou status..."
        itemLabel="produtos"
        headerActions={
          <Button
            type="button"
            className="h-9 gap-2"
            onClick={() => setModalState({ open: true, product: undefined })}
          >
            <Plus className="size-4" />
            Novo produto
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Package}
            eyebrow="Produtos"
            title="Nenhum produto encontrado"
            description="Crie o primeiro produto para começar a vender nessa edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((product, index) => (
            <AdminProductCard
              key={product.id}
              product={product}
              index={index}
              onEdit={(currentProduct) => setModalState({ open: true, product: currentProduct })}
              onPublish={() => { setPublishingProduct(product); }}
              onSoftDelete={() => { setDeletingProduct(product); }}
              onRestore={() => { setRestoringProduct(product); }}
            />
          ))
        }
      />

      <ManageProductModal
        key={modalState.product?.id ?? 'product-create'}
        open={modalState.open}
        editionId={editionId}
        product={modalState.product}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }))
            return
          }

          setModalState({ open: false, product: undefined })
        }}
        onCreate={async (values) => {
          const res = await createProductMutation.mutateAsync({
            eventId,
            editionId,
            data: values,
          })

          return res.success ? res.data : false
        }}
        onUpdate={async (productId, values) => {
          const res = await updateProductMutation.mutateAsync({
            eventId,
            editionId,
            productId,
            data: values,
          })

          return res.success ? res.data : false
        }}
      />

      <AlertModal
        open={Boolean(publishingProduct)}
        onOpenChange={() => setPublishingProduct(null)}
        title="Publicar produto?"
        description={
          publishingProduct
            ? `Ao publicar "${publishingProduct.name}", ele ficará disponível para os participantes.`
            : undefined
        }
        confirmLabel="Publicar produto"
        variant="default"
        loading={publishProductMutation.isPending}
        onConfirm={async () => {
          if (!publishingProduct) return
          await publishProductMutation.mutateAsync({
            eventId,
            editionId,
            productId: publishingProduct.id,
          })
          setPublishingProduct(null)
        }}
      />

      <AlertModal
        open={Boolean(deletingProduct)}
        onOpenChange={() => setDeletingProduct(null)}
        title="Excluir produto?"
        description={
          deletingProduct
            ? `Ao excluir "${deletingProduct.name}", o produto será removido da listagem admin.`
            : undefined
        }
        confirmLabel="Excluir produto"
        variant="destructive"
        loading={softDeleteProductMutation.isPending}
        onConfirm={async () => {
          if (!deletingProduct) return
          await softDeleteProductMutation.mutateAsync({
            eventId,
            editionId,
            productId: deletingProduct.id,
          })
          setDeletingProduct(null)
        }}
      />

      <AlertModal
        open={Boolean(restoringProduct)}
        onOpenChange={() => setRestoringProduct(null)}
        title="Restaurar produto?"
        description={
          restoringProduct
            ? `Ao restaurar "${restoringProduct.name}", ele volta para a listagem admin.`
            : undefined
        }
        confirmLabel="Restaurar produto"
        variant="default"
        loading={restoreSoftDeletedProductMutation.isPending}
        onConfirm={async () => {
          if (!restoringProduct) return
          await restoreSoftDeletedProductMutation.mutateAsync({
            eventId,
            editionId,
            productId: restoringProduct.id,
          })
          setRestoringProduct(null)
        }}
      />
    </div>
  )
}
