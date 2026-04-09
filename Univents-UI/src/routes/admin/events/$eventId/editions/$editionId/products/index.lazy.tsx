import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import {
  Plus,
  Tag,
  MoreVertical,
  ShieldCheck,
} from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { ProductCreateI, ProductI } from '@/features/products/model'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { cn } from '@/shared/lib/utils'
import { FormDrawer } from '@/widgets/form/ui/form-drawer'
import {
  allAdminProductsQueryOptions,
  allProductsQueryOptions,
  createProductFn,
  publishProductFn,
  restoreSoftDeletedProductFn,
  softDeleteProductFn,
} from '@/features/products/api'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { productCreateSchema } from '@/features/products/model'
import { getProductFields } from '@/features/products/model/field'
import { AdminProductCard } from '@/features/products/ui/AdminProductCard'

export const Route = createLazyFileRoute('/admin/events/$eventId/editions/$editionId/products/')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const { eventId, editionId } = Route.useParams()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [publishingProduct, setPublishingProduct] = useState<ProductI | null>(null)
  const [softDeletingProduct, setSoftDeletingProduct] = useState<ProductI | null>(null)
  const [restoringProduct, setRestoringProduct] = useState<ProductI | null>(null)
  const [isActionsOpen, setIsActionsOpen] = useState(false)

  const { data: products = [], isLoading: isProductsLoading } = useQuery(allAdminProductsQueryOptions(eventId, editionId))

  const createMutation = useMutation({
    mutationFn: (data: ProductCreateI) => createProductFn(data, eventId, editionId),
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData<ProductI[]>(
          allAdminProductsQueryOptions(eventId, editionId).queryKey,
          (old = []) => [...old, res.data]
        )
        setIsCreateOpen(false)
        toast.success('Produto criado com sucesso!')
      } else toast.error(res.message || 'Erro ao criar produto')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })

  const publishMutation = useMutation({
    mutationFn: ({ productId }: { productId: string }) =>
      publishProductFn(eventId, editionId, productId),
    onSuccess: async (res, variables) => {
      if (res.success) {
        queryClient.setQueryData<ProductI[]>(
          allAdminProductsQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((prod: ProductI) =>
            prod.id === variables.productId ? { ...prod, status: 'available' as const } : prod
          )
        )
        await queryClient.invalidateQueries(allProductsQueryOptions(eventId, editionId))
        setPublishingProduct(null)
        toast.success('Produto publicado com sucesso!')
      } else toast.error(res.message || 'Erro ao publicar produto')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })

  const softDeleteMutation = useMutation({
    mutationFn: ({ productId }: { productId: string }) =>
      softDeleteProductFn(eventId, editionId, productId),
    onSuccess: async (res, variables) => {
      if (res.success) {
        queryClient.setQueryData<ProductI[]>(
          allAdminProductsQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((prod: ProductI) =>
            prod.id === variables.productId ? { ...prod, deleted_at: new Date().toISOString() } : prod
          )
        )
        await queryClient.invalidateQueries(allProductsQueryOptions(eventId, editionId))
        setSoftDeletingProduct(null)
        toast.success('Produto excluído com sucesso!')
      } else toast.error(res.message || 'Erro ao excluir produto')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })

  const restoreMutation = useMutation({
    mutationFn: ({ productId }: { productId: string }) =>
      restoreSoftDeletedProductFn(eventId, editionId, productId),
    onSuccess: async (res, variables) => {
      if (res.success) {
        queryClient.setQueryData<ProductI[]>(
          allAdminProductsQueryOptions(eventId, editionId).queryKey,
          (old = []) => old.map((prod: ProductI) =>
            prod.id === variables.productId ? { ...prod, deleted_at: null } : prod
          )
        )
        await queryClient.invalidateQueries(allProductsQueryOptions(eventId, editionId))
        setRestoringProduct(null)
        toast.success('Produto restaurado com sucesso!')
      } else toast.error(res.message || 'Erro ao restaurar produto')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })


  const handleCreate = (data: ProductCreateI) => {
    createMutation.mutate(data)
  }

  const handlePublish = () => {
    if (!publishingProduct) return
    publishMutation.mutate({ productId: publishingProduct.id })
  }

  const handleSoftDelete = () => {
    if (!softDeletingProduct) return
    softDeleteMutation.mutate({ productId: softDeletingProduct.id })
  }

  const handleRestore = () => {
    if (!restoringProduct) return
    restoreMutation.mutate({ productId: restoringProduct.id })
  }

  const loading = createMutation.isPending || publishMutation.isPending || softDeleteMutation.isPending || restoreMutation.isPending

  return (
    <div className="min-h-screen bg-background relative pb-20 md:pb-0">
      <header className="sticky top-0 z-30 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <div className="flex items-center gap-2 shrink-0">
              <h1 className="text-lg md:text-xl font-semibold text-foreground">
                Produtos
                <span className="ml-2 text-sm font-normal text-muted-foreground">
                  ({products.length})
                </span>
              </h1>
            </div>

            <div className="hidden sm:flex items-center gap-2 ml-auto">
              <button
                onClick={() => { setIsCreateOpen(true) }}
                className={cn(
                  "flex items-center gap-2 px-4 py-2 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium"
                )}
              >
                <Plus className="w-4 h-4" />
                Novo produto
              </button>
            </div>

            <div className="sm:hidden flex items-center gap-1 ml-auto">
              <Drawer open={isActionsOpen} onOpenChange={setIsActionsOpen}>
                <DrawerTrigger asChild>
                  <button className={cn("flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted")}>
                    <MoreVertical className="w-5 h-5 text-foreground" />
                  </button>
                </DrawerTrigger>
                <DrawerContent className="z-60 rounded-t-2xl">
                  <DrawerHeader className="pb-4 border-b">
                    <DrawerTitle className="text-base font-semibold">Ações</DrawerTitle>
                  </DrawerHeader>
                  <div className="p-2 pb-8 space-y-1">
                    <button
                      onClick={() => { setIsActionsOpen(false); setIsCreateOpen(true) }}
                      className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                    >
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Plus className="w-4 h-4 text-primary" />
                      </div>
                      <span className="font-medium">Novo produto</span>
                    </button>
                  </div>
                </DrawerContent>
              </Drawer>
            </div>

            <Link
              to="/events/$eventId/editions/$editionId/products"
              params={{ eventId, editionId }}
              className={cn(
                "group relative flex items-center justify-center",
                "w-9 h-9 rounded-lg transition-all",
                "bg-primary text-primary-foreground",
                "hover:bg-primary/90",
                "shrink-0"
              )}
            >
              <ShieldCheck className="w-5 h-5" />
              <span
                className={cn(
                  "pointer-events-none absolute -bottom-9 right-0",
                  "whitespace-nowrap rounded-md px-2 py-1",
                  "bg-popover text-popover-foreground border text-xs shadow-md",
                  "opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0",
                  "transition-all"
                )}>
                Sair do admin de produtos
              </span>
            </Link>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
        <AnimatePresence mode="wait">
          {isProductsLoading ? (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="py-24 text-center text-muted-foreground"
            >
              Carregando produtos...
            </motion.div>
          ) : products.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="flex flex-col items-center justify-center py-24 space-y-6"
            >
              <div className="w-20 h-20 rounded-2xl bg-muted flex items-center justify-center">
                <Tag className="w-10 h-10 text-muted-foreground/30" />
              </div>
              <div className="text-center space-y-2">
                <h3 className="text-lg font-medium">Nenhum produto ainda</h3>
                <p className="text-sm text-muted-foreground max-w-xs">
                  Crie o primeiro produto para esta edição.
                </p>
              </div>
              <button
                onClick={() => { setIsCreateOpen(true) }}
                className={cn(
                  "mt-2 px-5 py-2.5 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium",
                  "active:scale-95 transition-all"
                )}
              >
                Criar produto
              </button>
            </motion.div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
              {products.map((product, idx) => (
                <AdminProductCard
                  key={product.id}
                  product={product}
                  index={idx}
                  onPublish={() => { setPublishingProduct(product) }}
                  onSoftDelete={() => { setSoftDeletingProduct(product) }}
                  onRestore={() => { setRestoringProduct(product) }}
                />
              ))}
            </div>
          )}
        </AnimatePresence>
      </main>

      <FormDrawer
        idPrefix="create-product-"
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        title="Novo produto"
        fields={getProductFields()}
        schema={productCreateSchema}
        onSubmit={handleCreate}
        submitLabel="Criar produto"
        loading={loading}
        closeOnSubmit={false}
        defaultValues={{ edition_scope_id: editionId }}
      />

      <AlertModal
        open={!!publishingProduct}
        onOpenChange={() => { setPublishingProduct(null) }}
        title="Publicar produto?"
        description={`Ao publicar "${publishingProduct?.name}", ele ficará visível para o público.`}
        confirmLabel="Publicar"
        onConfirm={handlePublish}
        variant="success"
        loading={loading}
      />

      <AlertModal
        open={!!softDeletingProduct}
        onOpenChange={() => { setSoftDeletingProduct(null) }}
        title="Excluir produto?"
        description={`Tem certeza que deseja excluir "${softDeletingProduct?.name}"? Ele será movido para a lixeira.`}
        confirmLabel="Excluir"
        onConfirm={handleSoftDelete}
        variant="destructive"
        loading={loading}
      />

      <AlertModal
        open={!!restoringProduct}
        onOpenChange={() => { setRestoringProduct(null) }}
        title="Restaurar produto?"
        description={`Tem certeza que deseja restaurar "${restoringProduct?.name}"? Ele voltará a ficar disponível.`}
        confirmLabel="Restaurar"
        onConfirm={handleRestore}
        variant="success"
        loading={loading}
      />
    </div>
  )
}
