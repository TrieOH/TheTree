import { createLazyFileRoute, useRouter } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { Package, Plus } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { SortState } from '@trieoh/ui-base'
import { Button } from '@/shared/ui/shadcn/button'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { productsByEditionQueryOptions } from '@/features/products/api'
import {
  useCreateInitialProductMutation,
  useDeleteProductMutation,
  useUpdateProductMutation,
} from '@/features/products/api/mutations'
import type { ProductI } from '@/features/products/model'
import { AdminProductCard } from '@/features/products/ui/AdminProductCard'
import { ManageProductModal } from '@/features/products/ui/ManageProductModal'
import { EditProductModal } from '@/features/products/ui/EditProductModal'

export const Route = createLazyFileRoute('/admin/events/$eventId_/editions/$editionId/products/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const router = useRouter()
  const { data: products = [] } = useQuery(productsByEditionQueryOptions(editionId))

  const createProductMutation = useCreateInitialProductMutation()
  const updateProductMutation = useUpdateProductMutation()
  const deleteProductMutation = useDeleteProductMutation()

  const [filter, setFilter] = useState('')
  const [sort, setSort] = useState<SortState<ProductI>>({
    field: 'vendor_code',
    direction: 'asc',
  })
  const [modalOpen, setModalOpen] = useState(false)
  const [productToDelete, setProductToDelete] = useState<ProductI | null>(null)
  const [productToEdit, setProductToEdit] = useState<ProductI | null>(null)

  const filteredProducts = useMemo(() => {
    const search = filter.trim().toLowerCase()

    return [...products]
      .filter((product) => {
        if (!search) return true

        return [
          product.vendor_code,
          String(product.requires_registration),
        ].some((value) => value.toLowerCase().includes(search))
      })
      .sort((a, b) => {
        const direction = sort.direction === 'asc' ? 1 : -1

        if (sort.field === 'requires_registration') {
          return (Number(a.requires_registration) - Number(b.requires_registration)) * direction
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
          { key: 'vendor_code', label: 'Código' },
          { key: 'requires_registration', label: 'Cadastro', comparator: (a, b) => Number(a.requires_registration) - Number(b.requires_registration) },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por código..."
        itemLabel="produtos"
        headerActions={
          <Button
            type="button"
            className="h-9 gap-2"
            onClick={() => setModalOpen(true)}
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
              onEdit={setProductToEdit}
              onDelete={setProductToDelete}
              onManageVariants={() => {
                router.navigate({
                  to: '/admin/events/$eventId/editions/$editionId/products/$productId/variants',
                  params: { eventId, editionId, productId: product.id },
                })
              }}
            />
          ))
        }
      />

      <ManageProductModal
        key={'product-create'}
        open={modalOpen}
        editionId={editionId}
        onOpenChange={setModalOpen}
        onCreate={async (values) => {
          const res = await createProductMutation.mutateAsync({
            editionId,
            data: values,
          })
          return res.success ? res.data : false
        }}
      />

      {productToEdit ? (
        <EditProductModal
          key={productToEdit.id}
          open
          product={productToEdit}
          onOpenChange={(open) => {
            if (!open) setProductToEdit(null)
          }}
          onUpdate={async (values) => {
            const res = await updateProductMutation.mutateAsync({
              productId: productToEdit.id,
              data: values,
            })
            return res.success ? res.data : false
          }}
        />
      ) : null}

      <AlertModal
        open={Boolean(productToDelete)}
        onOpenChange={() => setProductToDelete(null)}
        title="Excluir produto?"
        description={
          productToDelete
            ? `Ao excluir o produto "${productToDelete.vendor_code}".`
            : undefined
        }
        confirmLabel="Excluir produto"
        variant="destructive"
        loading={deleteProductMutation.isPending}
        onConfirm={async () => {
          if (!productToDelete) return
          await deleteProductMutation.mutateAsync({
            productId: productToDelete.id,
            editionId,
          })
          setProductToDelete(null)
        }}
      />
    </div>
  )
}