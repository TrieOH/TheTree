import { motion } from 'motion/react'
import { MoreVertical, Tag, Eye, EyeOff, Trash2, Pencil, Package, Calendar, Info } from 'lucide-react'
import type { ProductI } from '@/features/products/model'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/shadcn/dropdown-menu'
import { cn } from '@/shared/lib/utils'
import { StatusBadge } from '@/shared/ui'
import { formatDateRange } from '@/shared/lib/date'

interface AdminProductCardProps {
  product: ProductI
  index: number
  onEdit: (product: ProductI) => void
  onPublish: (product: ProductI) => void
  onSoftDelete: (product: ProductI) => void
  onRestore: (product: ProductI) => void
}

const typeLabels: Record<string, string> = {
  merchandise: 'Mercadoria',
  ticket: 'Ingresso',
  token: 'Token',
  bundle: 'Pacote',
}

export function AdminProductCard({
  product,
  index,
  onEdit,
  onPublish,
  onSoftDelete,
  onRestore,
}: AdminProductCardProps) {
  const handleAction = (type: 'edit' | 'publish' | 'delete' | 'restore') => {
    if (type === 'edit') onEdit(product)
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
        'group relative flex flex-col bg-card border border-border rounded-xl overflow-hidden transition-all duration-200',
        'hover:border-primary/20 hover:shadow-sm',
        isDeleted && 'opacity-60 grayscale'
      )}
    >
      <div className="flex flex-col sm:flex-row gap-4 p-4">
        {/* Thumbnail */}
        <div className="w-full sm:w-24 h-40 sm:h-24 rounded-lg bg-muted shrink-0 overflow-hidden border border-border/50 flex items-center justify-center">
          {product.thumbnail_url ? (
            <img src={product.thumbnail_url} alt="" className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
          ) : (
            <Tag className="w-8 h-8 text-muted-foreground/30" />
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0 space-y-2">
          <div className="flex items-start justify-between gap-2">
            <div className="space-y-1 min-w-0">
              <div className="flex items-center gap-2 flex-wrap">
                <h3 className="text-base font-semibold text-foreground truncate">{product.name}</h3>
                <StatusBadge status={product.status === 'available' ? 'published' : product.status} />
              </div>
              <p className="text-sm text-muted-foreground line-clamp-1">{product.description}</p>
            </div>

            <div className="shrink-0">
              <DropdownMenu>
                <DropdownMenuTrigger
                  render={
                    <button className="flex items-center justify-center w-8 h-8 rounded-lg hover:bg-muted transition-colors">
                      <MoreVertical className="w-4 h-4 text-muted-foreground" />
                    </button>
                  }
                />
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem onClick={() => { handleAction('edit'); }}>
                    <Pencil className="mr-2 h-4 w-4" />
                    <span>Editar</span>
                  </DropdownMenuItem>
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
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Package className="w-3.5 h-3.5" />
              <span className="bg-muted px-1.5 py-0.5 rounded font-medium">{typeLabels[product.type] || product.type}</span>
            </div>

            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Info className="w-3.5 h-3.5" />
              <span className="font-semibold text-foreground">R$ {(product.price_cents / 100).toFixed(2)}</span>
            </div>

            {product.has_inventory && (
              <div className="flex items-center gap-1.5 text-xs text-muted-foreground col-span-2 sm:col-span-1">
                <Package className="w-3.5 h-3.5 text-primary/60" />
                <span className={cn(
                  "font-medium px-1.5 py-0.5 rounded",
                  product.inventory_remaining === 0 ? "bg-destructive/10 text-destructive" : "bg-primary/10 text-primary"
                )}>
                  {product.inventory_remaining} / {product.inventory_quantity}
                </span>
              </div>
            )}
          </div>

          {(product.available_from ?? product.available_until) && (
            <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <Calendar className="w-3.5 h-3.5 opacity-60" />
              <span>
                {product.available_from && product.available_until
                  ? formatDateRange(product.available_from, product.available_until)
                  : product.available_from
                    ? `Desde ${new Date(product.available_from).toLocaleDateString('pt-BR')}`
                    : product.available_until
                      ? `Até ${new Date(product.available_until).toLocaleDateString('pt-BR')}`
                      : null}
              </span>
            </div>
          )}
        </div>
      </div>
    </motion.div>
  )
}
