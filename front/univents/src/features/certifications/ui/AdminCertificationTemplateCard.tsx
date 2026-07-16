import type React from 'react'
import { motion } from 'motion/react'
import {
  Check,
  FileText,
  MoreVertical,
  Pencil,
} from 'lucide-react'
import type { CertificationTemplateI } from '../model'
import { CertViewer } from './CertViewer'
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
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

interface AdminCertificationTemplateCardProps {
  template: CertificationTemplateI
  selected: boolean
  index?: number
  onSelect: (templateId: string) => void
  onEdit: () => void
  verifyUrl: string
  editionName: string
}

function MenuItems({
  isContext = false,
  onSelect,
  onEdit,
}: {
  isContext?: boolean
  onSelect: () => void
  onEdit: () => void
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator
  const stop = (action: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
    e.preventDefault()
    e.stopPropagation()
    action()
  }

  return (
    <>
      <Item onClick={stop(onSelect)}>
        <Check className="size-4" />
        <span>Selecionar template</span>
      </Item>
      <Separator />
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar template</span>
      </Item>
    </>
  )
}

export function AdminCertificationTemplateCard({
  template,
  selected,
  index = 0,
  onSelect,
  onEdit,
  verifyUrl,
  editionName,
}: AdminCertificationTemplateCardProps) {
  const handleSelect = () => onSelect(template.id)
  const handleEdit = () => onEdit()

  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <motion.article
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.05, duration: 0.35, ease: [0.25, 0.1, 0.25, 1] }}
            className={cn(
              'group relative flex w-full min-w-62.5 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left',
              'ring-1 ring-foreground/10 shadow-xs',
              'transform-gpu will-change-transform',
              'transition-all duration-300 ease-out',
              'hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm',
              'focus:outline-none focus-visible:outline-none focus-visible:ring-0',
            )}
            role="button"
            tabIndex={0}
            onClick={handleSelect}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                handleSelect()
              }
            }}
          >
            <div className="relative h-24 overflow-hidden bg-muted">
              {template.url ? (
                <img
                  src={template.url}
                  alt={template.title}
                  className="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-[1.03]"
                  loading={index < 4 ? 'eager' : 'lazy'}
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
                  <div className="flex size-12 items-center justify-center rounded-full border border-dashed border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
                    <FileText className="size-5 text-muted-foreground/40" />
                  </div>
                </div>
              )}

              <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/35 to-transparent" />

              <div className="absolute left-3 top-3 flex flex-wrap items-center gap-1.5">
                <span className="inline-flex items-center gap-1 rounded-full border border-border/60 bg-background/75 px-2 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-sm">
                  <FileText className="size-3.5" />
                  Template
                </span>
                {selected && (
                  <Badge className="bg-primary text-primary-foreground hover:bg-primary">
                    Selecionado
                  </Badge>
                )}
              </div>

              <div className="absolute right-2 top-2">
                <DropdownMenu>
                  <DropdownMenuTrigger
                    render={
                      <button
                        type="button"
                        onClick={(e) => e.stopPropagation()}
                        className={cn(
                          'inline-flex size-8 items-center justify-center rounded-full',
                          'bg-background/85 text-foreground shadow-sm backdrop-blur-sm',
                          'transition-colors hover:bg-background',
                        )}
                        aria-label={`Abrir ações de ${template.title}`}
                      >
                        <MoreVertical className="size-4" />
                      </button>
                    }
                  />
                  <DropdownMenuContent align="end" className="w-56">
                    <MenuItems
                      onSelect={handleSelect}
                      onEdit={handleEdit}
                    />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>

            <div className="flex items-center justify-between gap-3 p-3">
              <div className="min-w-0 flex-1 space-y-2">
                <div className="min-w-0 space-y-1">
                  <h3 className="line-clamp-2 text-sm font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary">
                    {template.title}
                  </h3>
                  <p className="text-[11px] text-muted-foreground">
                    {template.url ? 'Com fundo configurado' : 'Sem fundo'}
                  </p>
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant={selected ? 'default' : 'outline'}
                    size="sm"
                    className="h-8 flex-1 gap-2"
                    onClick={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      handleSelect()
                    }}
                  >
                    <Check className="size-3.5" />
                    {selected ? 'Selecionado' : 'Selecionar'}
                  </Button>
                  <CertViewer
                    template={template}
                    triggerLabel="Ver"
                    variables={{
                      activity_name: editionName,
                      certified_at: 'DD/MM/AAAA',
                      cert_hash: 'HASH-DE-EXEMPLO',
                      verify_url: verifyUrl,
                    }}
                  />
                </div>
              </div>
            </div>
          </motion.article>
        }
      />
      <ContextMenuContent className="w-56">
        <MenuItems
          isContext
          onSelect={handleSelect}
          onEdit={handleEdit}
        />
      </ContextMenuContent>
    </ContextMenu>
  )
}
