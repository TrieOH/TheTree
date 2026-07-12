import { Link, useNavigate } from '@tanstack/react-router'
import type React from 'react'
import { motion } from 'motion/react'
import {
  ArrowUpRight,
  CalendarDays,
  Copy,
  Eye,
  MoreVertical,
  Pencil,
  Sparkles,
  Users,
} from 'lucide-react'
import type { EventI, EventStatusI } from '@/features/events/model'
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
import { toast } from 'sonner'

const statusConfig: Record<EventStatusI, {
  label: string
  dot: string
  pill: string
}> = {
  draft: {
    label: 'Rascunho',
    dot: 'bg-amber-500',
    pill: 'bg-amber-500/10 text-amber-700 border-amber-500/20',
  },
  active: {
    label: 'Ativo',
    dot: 'bg-emerald-500',
    pill: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
  },
  archived: {
    label: 'Arquivado',
    dot: 'bg-slate-500',
    pill: 'bg-slate-500/10 text-slate-700 border-slate-500/20',
  },
  discontinued: {
    label: 'Descontinuado',
    dot: 'bg-rose-500',
    pill: 'bg-rose-500/10 text-rose-700 border-rose-500/20',
  },
}

interface AdminEventCardProps {
  event: EventI
  index?: number
  onEdit: (event: EventI) => void
  onPublish: (event: EventI) => void
}

function MenuItems({
  event,
  isContext = false,
  onEdit,
  onPublish,
  onOpenEditions,
  onOpenDashboard,
}: {
  event: EventI
  isContext?: boolean
  onEdit: () => void
  onPublish: () => void
  onOpenEditions: () => void
  onOpenDashboard: () => void
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator
  const stop = (action: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
    e.preventDefault()
    e.stopPropagation()
    action()
  }
  const copyLink = () => {
    const url = `${window.location.origin}/events/${event.id}`
    void navigator.clipboard.writeText(url)
    toast.success("Link copied to clipboard");
  }

  return (
    <>
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar</span>
      </Item>
      <Item onClick={stop(onOpenDashboard)}>
        <ArrowUpRight className="size-4" />
        <span>Ver dashboard</span>
      </Item>
      <Separator />
      {event.status === 'draft' && (
        <Item onClick={stop(onPublish)}>
          <Eye className="size-4" />
          <span>Publicar</span>
        </Item>
      )}
      <Item onClick={stop(copyLink)}>
        <Copy className="size-4" />
        <span>Copiar link</span>
      </Item>
      <Separator />
      <Item onClick={stop(onOpenEditions)}>
        <ArrowUpRight className="size-4" />
        <span>Ver edições</span>
      </Item>
    </>
  )
}

export default function AdminEventCard({
  event,
  index = 0,
  onEdit,
  onPublish,
}: AdminEventCardProps) {
  const navigate = useNavigate()
  const status = statusConfig[event.status]
  const hasVisual = Boolean(event.banner_url ?? event.logo_url)
  const handleEdit = () => onEdit(event)
  const handlePublish = () => onPublish(event)
  const handleOpenDashboard = () => {
    void navigate({
      to: '/admin/events/$eventId',
      params: { eventId: event.id },
    })
  }
  const handleOpenEditions = () => {
    void navigate({
      to: '/admin/events/$eventId/editions',
      params: { eventId: event.id },
    })
  }

  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <motion.article
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.05, duration: 0.35, ease: [0.25, 0.1, 0.25, 1] }}
            className={cn(
              'group relative flex min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left',
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
                  src={event.banner_url ?? event.logo_url ?? ''}
                  alt="Representação Visual do Evento"
                  className={cn(
                    "h-full w-full object-cover transition-transform",
                    "duration-700 ease-out group-hover:scale-105",
                  )}
                  loading={index < 4 ? 'eager' : 'lazy'}
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
                  <div className="flex size-20 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
                    <span className="text-2xl font-semibold text-muted-foreground/40">
                      {event.acronym ?? event.name.charAt(0)}
                    </span>
                  </div>
                </div>
              )}

              <div className="absolute left-4 top-4 flex flex-wrap items-center gap-2">
                <span className={cn(
                  'inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-[11px] font-medium backdrop-blur-sm',
                  status.pill,
                )}>
                  <span className={cn('size-1.5 rounded-full', status.dot)} />
                  {status.label}
                </span>

                {event.is_series && (
                  <span className="inline-flex items-center gap-1 rounded-full bg-background/85 px-2.5 py-1 text-[11px] font-medium text-foreground backdrop-blur-sm">
                    <Sparkles className="size-3.5 text-primary" />
                    Série
                  </span>
                )}
              </div>

              <div className="absolute right-4 top-4">
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
                        aria-label={`Abrir ações de ${event.name}`}
                      >
                        <MoreVertical className="size-4" />
                      </button>
                    }
                  />
                  <DropdownMenuContent align="end" className="w-56">
                    <MenuItems
                      event={event}
                      onEdit={handleEdit}
                      onPublish={handlePublish}
                      onOpenDashboard={handleOpenDashboard}
                      onOpenEditions={handleOpenEditions}
                    />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>

              <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
                <div className="min-w-0 space-y-1">
                  <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
                    {event.name}
                  </h3>
                  {event.tagline && (
                    <p className="line-clamp-2 max-w-2xl text-xs text-muted-foreground">
                      {event.tagline}
                    </p>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                  <span className="inline-flex items-center gap-1.5">
                    <CalendarDays className="size-3.5" />
                    {event.editions_count} {event.editions_count === 1 ? 'edição' : 'edições'}
                  </span>
                  <span className="inline-flex items-center gap-1.5 truncate">
                    <Users className="size-3.5" />
                    {event.contact_email}
                  </span>
                </div>
                <code className="block truncate text-[11px] font-mono text-muted-foreground/80">
                  {event.slug}
                </code>
              </div>

              <Link
                to="/admin/events/$eventId"
                params={{ eventId: event.id }}
                className={cn(
                  'inline-flex shrink-0 items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-medium',
                  'bg-secondary/60 text-secondary-foreground transition-colors hover:bg-secondary',
                )}
                onClick={(e) => e.stopPropagation()}
              >
                Painel
                <ArrowUpRight className="size-3.5" />
              </Link>
            </div>
          </motion.article>
        }
      >
      </ContextMenuTrigger>

      <ContextMenuContent align="end" className="w-56">
        <MenuItems
          event={event}
          isContext
          onEdit={handleEdit}
          onPublish={handlePublish}
          onOpenDashboard={handleOpenDashboard}
          onOpenEditions={handleOpenEditions}
        />
      </ContextMenuContent>

    </ContextMenu>
  )
}
