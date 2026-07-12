import type React from 'react'
import { useNavigate } from '@tanstack/react-router'
import { motion } from 'motion/react'
import {
  ArrowUpRight,
  Award,
  CalendarDays,
  CreditCard,
  FileText,
  MapPin,
  MoreVertical,
  Package,
  Pencil,
  ShieldCheck,
} from 'lucide-react'
import type { EditionI } from '../model'
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
interface EditionCardProps {
  edition: EditionI
  eventId: string
  index: number
  onEdit: (edition: EditionI) => void
  onPublish?: () => void
  onConnect?: () => void
  onDisconnect?: () => void
}

const statusConfig: Record<EditionI['status'], { label: string; dot: string; pill: string }> = {
  draft: {
    label: 'Rascunho',
    dot: 'bg-amber-500',
    pill: 'bg-amber-500/10 text-amber-700 border-amber-500/20',
  },
  announced: {
    label: 'Anunciada',
    dot: 'bg-sky-500',
    pill: 'bg-sky-500/10 text-sky-700 border-sky-500/20',
  },
  open: {
    label: 'Aberta',
    dot: 'bg-emerald-500',
    pill: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
  },
  ongoing: {
    label: 'Em andamento',
    dot: 'bg-cyan-500',
    pill: 'bg-cyan-500/10 text-cyan-700 border-cyan-500/20',
  },
  finished: {
    label: 'Finalizada',
    dot: 'bg-slate-500',
    pill: 'bg-slate-500/10 text-slate-700 border-slate-500/20',
  },
  completed: {
    label: 'Concluída',
    dot: 'bg-emerald-500',
    pill: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
  },
  cancelled: {
    label: 'Cancelada',
    dot: 'bg-rose-500',
    pill: 'bg-rose-500/10 text-rose-700 border-rose-500/20',
  },
  postponed: {
    label: 'Adiada',
    dot: 'bg-orange-500',
    pill: 'bg-orange-500/10 text-orange-700 border-orange-500/20',
  },
}

function MenuItems({
  isContext = false,
  edition,
  onEdit,
  eventId,
  onPublish,
  onConnect,
  onDisconnect,
}: {
  isContext?: boolean
  edition: EditionI
  onEdit: () => void
  eventId: string
  onPublish?: () => void
  onConnect?: () => void
  onDisconnect?: () => void
}) {
  const navigate = useNavigate()
  const Item = isContext ? ContextMenuItem : DropdownMenuItem
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator
  const go = (to: string) => () => {
    void navigate({
      to,
      params: { eventId, editionId: edition.id },
    })
  }
  const openEdition = go('/admin/events/$eventId/editions/$editionId')
  const openCertifications = go('/admin/events/$eventId/editions/$editionId/certifications')
  const openSignatures = go('/admin/events/$eventId/editions/$editionId/signatures')
  const openProducts = go('/admin/events/$eventId/editions/$editionId/products')
  const openCheckpoints = go('/admin/events/$eventId/editions/$editionId/checkpoints')
  const stop = (action?: () => void) => (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    action?.()
  }

  return (
    <>
      <Item onClick={stop(openEdition)}>
        <ArrowUpRight className="size-4" />
        <span>Abrir edição</span>
      </Item>
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar</span>
      </Item>
      <Separator />
      {edition.status === 'draft' && onPublish && (
        <Item onClick={stop(onPublish)}>
          <ArrowUpRight className="size-4" />
          <span>Publicar</span>
        </Item>
      )}
      <Item onClick={stop(openCertifications)}>
        <Award className="size-4" />
        <span>Certificações</span>
      </Item>
      <Item onClick={stop(openSignatures)}>
        <FileText className="size-4" />
        <span>Assinaturas</span>
      </Item>
      <Item onClick={stop(openProducts)}>
        <Package className="size-4" />
        <span>Produtos</span>
      </Item>
      <Item onClick={stop(openCheckpoints)}>
        <ShieldCheck className="size-4" />
        <span>Checkpoints</span>
      </Item>
      {(onConnect || onDisconnect) && <Separator />}
      {onConnect && (
        <Item onClick={stop(onConnect)}>
          <CreditCard className="size-4" />
          <span>Conectar pagamento</span>
        </Item>
      )}
      {onDisconnect && (
        <Item onClick={stop(onDisconnect)}>
          <CreditCard className="size-4" />
          <span>Desconectar pagamento</span>
        </Item>
      )}
    </>
  )
}

export function AdminEditionCard({
  edition,
  eventId,
  index,
  onEdit,
  onPublish,
  onConnect,
  onDisconnect,
}: EditionCardProps) {
  const navigate = useNavigate()
  const hasVisual = Boolean(edition.banner_url ?? edition.logo_url)
  const monetaryLabel =
    edition.monetary_type === 'free'
      ? 'Gratuita'
      : edition.monetary_type === 'paid'
        ? 'Paga'
        : 'Mista'
  const status = statusConfig[edition.status]

  const openEdition = () => {
    void navigate({
      to: '/admin/events/$eventId/editions/$editionId',
      params: { eventId, editionId: edition.id },
    })
  }
  const handleEdit = () => onEdit(edition)

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
            onClick={handleEdit}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                handleEdit()
              }
            }}
          >
            <div className="relative aspect-video overflow-hidden bg-muted">
              {hasVisual ? (
                <img
                  src={edition.banner_url ?? edition.logo_url ?? ''}
                  alt={edition.edition_name}
                  className={cn(
                    'h-full w-full object-cover transition-transform duration-700 ease-out',
                    'group-hover:scale-105',
                  )}
                  loading={index < 4 ? 'eager' : 'lazy'}
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
                  <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
                    <CalendarDays className="size-7 text-muted-foreground/40" />
                  </div>
                </div>
              )}

              <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/35 to-transparent" />

              <div className="absolute left-3 top-3 flex flex-wrap items-center gap-1.5">
                <span className={cn(
                  'inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm',
                  status.pill,
                )}>
                  <span className={cn('size-1.5 rounded-full', status.dot)} />
                  <span className="max-w-28 truncate">{status.label}</span>
                </span>
                <span className="inline-flex items-center gap-1 rounded-full border border-border/60 bg-background/75 px-2 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-sm">
                  <span className="max-w-28 truncate">{monetaryLabel}</span>
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
                        aria-label={`Abrir ações de ${edition.edition_name}`}
                      >
                        <MoreVertical className="size-4" />
                      </button>
                    }
                  />
                  <DropdownMenuContent align="end" className="w-56">
                    <MenuItems
                      edition={edition}
                      onEdit={handleEdit}
                      eventId={eventId}
                      onPublish={onPublish}
                      onConnect={onConnect}
                      onDisconnect={onDisconnect}
                    />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>

              <div className="absolute inset-x-0 bottom-0 p-3 sm:p-4">
                <div className="min-w-0 max-w-full space-y-0.5">
                  <h3 className="truncate text-base font-semibold leading-tight text-foreground transition-colors duration-300 group-hover:text-primary sm:text-lg">
                    {edition.edition_name}
                  </h3>
                  {edition.tagline && (
                    <p className="truncate text-xs leading-5 text-muted-foreground">
                      {edition.tagline}
                    </p>
                  )}
                </div>
              </div>
            </div>

            <div className="flex min-w-0 items-center justify-between gap-3 p-3 pt-2.5 sm:p-4 sm:pt-3">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex min-w-0 flex-wrap items-center gap-1.5 text-[11px] text-muted-foreground">
                  <span className="inline-flex min-w-0 max-w-full items-center gap-1">
                    <CalendarDays className="size-3 shrink-0" />
                    <span className="block min-w-0 truncate">{formatDateRange(edition.starts_at, edition.ends_at)}</span>
                  </span>
                  <span className="inline-flex min-w-0 max-w-full items-center gap-1">
                    <MapPin className="size-3 shrink-0" />
                    <span className="block min-w-0 truncate">{edition.location_name}</span>
                  </span>
                </div>

                <p className="block min-w-0 truncate text-[11px] text-muted-foreground">
                  {edition.location_address}
                </p>
              </div>

              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation()
                  openEdition()
                }}
                className={cn(
                  'inline-flex shrink-0 items-center gap-1 rounded-full px-3 py-1 text-[11px] font-medium',
                  'bg-secondary/60 text-secondary-foreground transition-colors hover:bg-secondary',
                )}
              >
                Painel
                <ArrowUpRight className="size-3.5" />
              </button>
            </div>
          </motion.article>
        }
      >
      </ContextMenuTrigger>

      <ContextMenuContent align="end" className="w-56">
        <MenuItems
          isContext
          edition={edition}
          onEdit={handleEdit}
          eventId={eventId}
          onPublish={onPublish}
          onConnect={onConnect}
          onDisconnect={onDisconnect}
        />
      </ContextMenuContent>
    </ContextMenu>
  )
}
