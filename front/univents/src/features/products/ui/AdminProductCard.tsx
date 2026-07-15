import type React from 'react'
import { motion } from 'motion/react'
import { MoreVertical, Tag, Eye, EyeOff, Trash2, Pencil, Package, Calendar, Info } from 'lucide-react'
import type { ProductI } from '@/features/products/model'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/shared/ui/shadcn/context-menu'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/shadcn/dropdown-menu'
import { cn } from '@/shared/lib/utils'
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

function MenuItems({
  product,
  isContext = false,
  onEdit,
  onPublish,
  onSoftDelete,
  onRestore,
}: {
  product: ProductI
  isContext?: boolean
  onEdit: () => void
  onPublish: () => void
  onSoftDelete: () => void
  onRestore: () => void
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator
  const stop = (action: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
    e.preventDefault()
    e.stopPropagation()
    action()
  }

  const isDeleted = product.deleted_at !== null

  return (
    <>
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar</span>
      </Item>
      {product.status === 'draft' && !isDeleted && (
        <Item onClick={stop(onPublish)}>
          <Eye className="size-4" />
          <span>Publicar</span>
        </Item>
      )}
      <Separator />
      {!isDeleted ? (
        <Item onClick={stop(onSoftDelete)}>
          <Trash2 className="size-4" />
          <span>Excluir</span>
        </Item>
      ) : (
        <Item onClick={stop(onRestore)}>
          <EyeOff className="size-4" />
          <span>Restaurar</span>
        </Item>
      )}
    </>
  )
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

  const article = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.35, ease: [0.25, 0.1, 0.25, 1] }}
      className={cn(
        'group relative flex w-full min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left',
        'ring-1 ring-foreground/10 shadow-xs',
        'transform-gpu will-change-transform',
        'transition-all duration-300 ease-out',
        'hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm',
        'focus:outline-none focus-visible:outline-none focus-visible:ring-0',
        isDeleted && 'opacity-60 grayscale',
      )}
      role="button"
      tabIndex={0}
      onClick={() => handleAction('edit')}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          handleAction('edit')
        }
      }}
    >
      <div className="relative aspect-video overflow-hidden bg-muted">
        {product.thumbnail_url ? (
          <img
            src={product.thumbnail_url}
            alt={product.name}
            className={cn(
              'h-full w-full object-cover transition-transform duration-700 ease-out',
              'group-hover:scale-105',
            )}
            loading={index < 4 ? 'eager' : 'lazy'}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
            <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
              <Tag className="size-7 text-muted-foreground/40" />
            </div>
          </div>
        )}

        <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/35 to-transparent" />

        <div className="absolute left-3 top-3 flex flex-wrap items-center gap-1.5">
          <span className={cn(
            'inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm',
            isDeleted
              ? 'border-rose-500/20 bg-rose-500/10 text-rose-700'
              : isPublished
                ? 'border-emerald-500/20 bg-emerald-500/10 text-emerald-700'
                : 'border-amber-500/20 bg-amber-500/10 text-amber-700',
          )}>
            <span className={cn(
              'size-1.5 rounded-full',
              isDeleted
                ? 'bg-rose-500'
                : isPublished
                  ? 'bg-emerald-500'
                  : 'bg-amber-500',
            )} />
            <span className="max-w-28 truncate">
              {isDeleted ? 'Excluído' : isPublished ? 'Publicado' : 'Rascunho'}
            </span>
          </span>

          <span className="inline-flex items-center gap-1 rounded-full border border-border/60 bg-background/75 px-2 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-sm">
            <Package className="size-3.5" />
            <span className="max-w-28 truncate">{typeLabels[product.type] || product.type}</span>
          </span>
        </div>

        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(e) => e.stopPropagation()}
                  className={cn(
                    'inline-flex size-9 items-center justify-center rounded-full',
                    'bg-background/85 text-foreground shadow-sm backdrop-blur-sm',
                    'transition-colors hover:bg-background',
                  )}
                  aria-label={`Abrir ações de ${product.name}`}
                >
                  <MoreVertical className="size-4" />
                </button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems
                product={product}
                onEdit={() => handleAction('edit')}
                onPublish={() => handleAction('publish')}
                onSoftDelete={() => handleAction('delete')}
                onRestore={() => handleAction('restore')}
              />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
          <div className="min-w-0 space-y-1">
            <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
              {product.name}
            </h3>
            {product.description && (
              <p className="line-clamp-2 max-w-2xl text-xs text-muted-foreground">
                {product.description}
              </p>
            )}
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
        <div className="min-w-0 flex-1 space-y-1.5">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <Info className="size-3.5 shrink-0" />
            <span className="font-semibold text-foreground">R$ {(product.price_cents / 100).toFixed(2)}</span>
          </div>

          {product.has_inventory ? (
            <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <Package className="size-3.5 shrink-0" />
              <span className="truncate">
                {product.inventory_remaining} de {product.inventory_quantity} em estoque
              </span>
            </div>
          ) : (
            <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <Package className="size-3.5 shrink-0" />
              <span className="truncate">Estoque ilimitado</span>
            </div>
          )}

          {(product.available_from ?? product.available_until) && (
            <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <Calendar className="size-3.5 opacity-60" />
              <span className="truncate">
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
    </motion.article>
  )

  return (
    <ContextMenu>
      <ContextMenuTrigger render={article} />
      <ContextMenuContent className="w-56">
        <MenuItems
          isContext
          product={product}
          onEdit={() => handleAction('edit')}
          onPublish={() => handleAction('publish')}
          onSoftDelete={() => handleAction('delete')}
          onRestore={() => handleAction('restore')}
        />
      </ContextMenuContent>
    </ContextMenu>
  )
}
