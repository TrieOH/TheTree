import { Link } from '@tanstack/react-router'
import { motion } from 'motion/react'
import {
  Pencil,
  MoreVertical,
  Eye,
  Users,
  ChevronRight,
  ArrowUpRight,
} from 'lucide-react'
import type { EventI, EventStatusI } from '@/features/events/model'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { cn } from '@/shared/lib/utils'

const statusConfig: Record<EventStatusI, {
  label: string;
  dot: string;
}> = {
  draft: {
    label: 'Rascunho',
    dot: 'bg-muted-foreground'
  },
  active: {
    label: 'Ativo',
    dot: 'bg-primary'
  },
  archived: {
    label: 'Arquivado',
    dot: 'bg-muted-foreground'
  },
  discontinued: {
    label: 'Descontinuado',
    dot: 'bg-destructive'
  },
}

interface AdminEventCardProps {
  event: EventI
  index: number
  onEdit: (event: EventI) => void
  // onDelete: (event: EventI) => void
  onPublish: (event: EventI) => void
}

export default function AdminEventCard({
  event,
  index,
  onEdit,
  onPublish
}: AdminEventCardProps) {
  const isDraft = event.status === 'draft'
  const status = statusConfig[event.status]

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.04, duration: 0.25 }}
      className={cn(
        "group bg-card rounded-xl border p-4",
        "hover:border-foreground/20 hover:shadow-sm",
        "transition-all duration-200"
      )}
    >
      <div className="flex gap-3 items-center">
        <div className={cn(
          "w-12 h-12 rounded-lg shrink-0 overflow-hidden",
          "bg-muted ring-1 ring-border shadow-sm",
          "flex items-center justify-center",
          event.logo_url ? "bg-white" : "bg-secondary/30"
        )}>
          {event.logo_url ? (
            <img
              src={event.logo_url}
              alt=""
              className="w-full h-full object-cover"
            />
          ) : (
            <span className="text-sm font-bold text-secondary-foreground/50">
              {event.acronym?.slice(0, 2).toUpperCase() ?? event.name.slice(0, 2).toUpperCase()}
            </span>
          )}
        </div>

        <div className="flex-1 min-w-0 flex flex-col justify-center h-12">
          <h3 className="font-semibold text-foreground text-sm leading-tight truncate">
            {event.name}
          </h3>

          <div className="flex items-center gap-2 mt-0.5">
            <span className={cn(
              "w-1.5 h-1.5 rounded-full",
              status.dot
            )} />
            <span className="text-[11px] text-muted-foreground">
              {status.label} · {event.editions_count} {event.editions_count === 1 ? 'edição' : 'edições'}
            </span>
            {event.is_series && (
              <span className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded bg-muted text-muted-foreground text-[10px] font-medium">
                <Users className="w-3 h-3" />
                Série
              </span>
            )}
          </div>
        </div>
      </div>

      <div className="mt-3">
        <code className="text-[11px] text-muted-foreground font-mono bg-muted/50 px-1.5 py-0.5 rounded">
          {event.slug}
        </code>
      </div>

      <div className="flex items-center justify-between gap-2 pt-3 mt-3 border-t border-border/50">
        <div className="hidden sm:flex items-center gap-1">
          {isDraft ? (
            <button
              onClick={() => { onPublish(event); }}
              className={cn(
                "flex items-center gap-1.5 px-2.5 py-1 rounded-md",
                "bg-primary text-primary-foreground hover:bg-primary/90",
                "text-xs font-medium",
                "active:scale-95 transition-all"
              )}
            >
              <Eye className="w-3 h-3" />
              Publicar
            </button>
          ) : (
            <span className="text-[11px] text-muted-foreground">
              {event.editions_count} {event.editions_count === 1 ? 'edição' : 'edições'}
            </span>
          )}
        </div>

        <div className="hidden sm:flex items-center gap-0.5 ml-auto">
          <button
            onClick={() => { onEdit(event); }}
            className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
            title="Editar"
          >
            <Pencil className="w-3.5 h-3.5" />
          </button>


          <Link
            to="/admin/events/$eventId/editions"
            params={{ eventId: event.id }}
            className={cn(
              "ml-1 flex items-center gap-1 px-2.5 py-1 rounded-md",
              "bg-secondary/50 text-secondary-foreground hover:bg-secondary",
              "text-xs font-medium",
              "transition-colors"
            )}
          >
            Edições
            <ChevronRight className="w-3 h-3" />
          </Link>
        </div>

        <div className="sm:hidden flex items-center justify-between w-full">
          {isDraft ? (
            <button
              onClick={() => { onPublish(event); }}
              className={cn(
                "flex items-center gap-1.5 px-2.5 py-1 rounded-md",
                "bg-primary text-primary-foreground",
                "text-xs font-medium"
              )}
            >
              <Eye className="w-3 h-3" />
              Publicar
            </button>
          ) : (
            <span className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span className={cn("w-1.5 h-1.5 rounded-full", status.dot)} />
              {status.label}
            </span>
          )}

          <Drawer>
            <DrawerTrigger asChild>
              <button className="p-1.5 rounded-md hover:bg-muted text-muted-foreground">
                <MoreVertical className="w-4 h-4" />
              </button>
            </DrawerTrigger>
            <DrawerContent className="z-50 rounded-t-2xl">
              <DrawerHeader className="border-b pb-4">
                <DrawerTitle className="text-base font-semibold line-clamp-1">
                  {event.name}
                </DrawerTitle>
              </DrawerHeader>
              <div className="p-3 space-y-1">
                {isDraft && (
                  <button
                    onClick={() => { onPublish(event); }}
                    className="w-full flex items-center gap-3 px-4 py-3 rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 active:scale-95 transition-all"
                  >
                    <Eye className="w-5 h-5" />
                    <span className="font-medium">Publicar evento</span>
                  </button>
                )}

                <Link
                  to="/admin/events/$eventId/editions"
                  params={{ eventId: event.id }}
                  className="w-full flex items-center gap-3 px-4 py-3 rounded-xl hover:bg-muted active:bg-muted/80 transition-colors"
                >
                  <ArrowUpRight className="w-5 h-5 text-muted-foreground" />
                  <span className="font-medium">Ver edições</span>
                </Link>

                <button
                  onClick={() => { onEdit(event); }}
                  className="w-full flex items-center gap-3 px-4 py-3 rounded-xl hover:bg-muted active:bg-muted/80 transition-colors"
                >
                  <Pencil className="w-5 h-5 text-muted-foreground" />
                  <span className="font-medium">Editar evento</span>
                </button>

                <div className="h-px bg-border my-2" />
              </div>
            </DrawerContent>
          </Drawer>
        </div>
      </div>
    </motion.div>
  )
}