import { motion } from 'motion/react'
import { MoreVertical, Tag, Eye, EyeOff, Trash2 } from 'lucide-react'
import { useState } from 'react'
import type { ProductI } from '@/features/products/model'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/shadcn/dropdown-menu'
import { cn } from '@/shared/lib/utils'

interface AdminProductCardProps {
  product: ProductI
  index: number
  onPublish: (product: ProductI) => void
  onSoftDelete: (product: ProductI) => void
  onRestore: (product: ProductI) => void
}

export function AdminProductCard({
  product,
  index,
  onPublish,
  onSoftDelete,
  onRestore,
}: AdminProductCardProps) {
  const [isActionsOpen, setIsActionsOpen] = useState(false)

  const handleAction = (type: 'publish' | 'delete' | 'restore') => {
    if (type === 'publish') onPublish(product)
    if (type === 'delete') onSoftDelete(product)
    if (type === 'restore') onRestore(product)
  }

  const isPublished = product.status === 'available'
  const isDeleted = product.deleted_at !== null

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05 }}
      className={cn(
        'relative flex flex-col justify-between p-4 rounded-xl border border-border bg-card',
        isDeleted && 'opacity-60 grayscale'
      )}
    >
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <Tag className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-foreground">{product.name}</h3>
            <p className="text-sm text-muted-foreground">
              {product.type} - R$ {(product.price_cents / 100).toFixed(2)}
            </p>
          </div>
        </div>

        <div className="hidden sm:flex">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button className="flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted">
                  <MoreVertical className="w-5 h-5 text-foreground" />
                </button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem onClick={() => { handleAction('publish'); }} disabled={isPublished}>
                <Eye className="mr-2 h-4 w-4" />
                <span>Publicar</span>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => { handleAction('delete'); }} disabled={isDeleted}>
                <Trash2 className="mr-2 h-4 w-4" />
                <span>Excluir</span>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => { handleAction('restore'); }} disabled={!isDeleted}>
                <EyeOff className="mr-2 h-4 w-4" />
                <span>Restaurar</span>
              </DropdownMenuItem>
              {/* Add edit link here once the edit page is available */}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="sm:hidden">
          <Drawer open={isActionsOpen} onOpenChange={setIsActionsOpen}>
            <DrawerTrigger asChild>
              <button className="flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted">
                <MoreVertical className="w-5 h-5 text-foreground" />
              </button>
            </DrawerTrigger>
            <DrawerContent className="z-60 rounded-t-2xl">
              <DrawerHeader className="pb-4 border-b">
                <DrawerTitle className="text-base font-semibold">Ações para {product.name}</DrawerTitle>
              </DrawerHeader>
              <div className="p-2 pb-8 space-y-1">
                <button
                  onClick={() => { handleAction('publish'); setIsActionsOpen(false); }}
                  className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                  disabled={isPublished}
                >
                  <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                    <Eye className="w-4 h-4 text-primary" />
                  </div>
                  <span className="font-medium">Publicar</span>
                </button>
                <button
                  onClick={() => { handleAction('delete'); setIsActionsOpen(false); }}
                  className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                  disabled={isDeleted}
                >
                  <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                    <Trash2 className="w-4 h-4 text-primary" />
                  </div>
                  <span className="font-medium">Excluir</span>
                </button>
                <button
                  onClick={() => { handleAction('restore'); setIsActionsOpen(false); }}
                  className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                  disabled={!isDeleted}
                >
                  <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                    <EyeOff className="w-4 h-4 text-primary" />
                  </div>
                  <span className="font-medium">Restaurar</span>
                </button>
                {/* Add edit link here once the edit page is available */}
              </div>
            </DrawerContent>
          </Drawer>
        </div>
      </div>
    </motion.div>
  )
}
