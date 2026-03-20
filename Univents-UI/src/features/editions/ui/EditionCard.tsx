import { motion } from 'motion/react'
import { Link as TanstackLink } from '@tanstack/react-router'
import {
  Calendar,
  MapPin,
  ArrowUpRight,
  Sun,
  Trophy,
  Sparkles,
  CalendarDays,
  HashIcon,
} from 'lucide-react'
import type { EditionI } from '../model'
import { cn } from '@/shared/lib/utils'

const typeIcons = {
  year: CalendarDays,
  season: Sun,
  number: HashIcon,
  ordinal: Trophy,
  custom: Sparkles,
} as const

const statusConfig = {
  draft: { label: 'Rascunho', variant: 'muted' },
  announced: { label: 'Anunciado', variant: 'blue' },
  open: { label: 'Inscrições Abertas', variant: 'green' },
  ongoing: { label: 'Em Andamento', variant: 'amber' },
  finished: { label: 'Encerrado', variant: 'muted' },
  completed: { label: 'Concluído', variant: 'green' },
  cancelled: { label: 'Cancelado', variant: 'destructive' },
  postponed: { label: 'Adiado', variant: 'amber' },
} as const

interface EditionCardProps {
  edition: EditionI
  eventId: string
  index?: number
  className?: string
}

export function EditionCard({ edition, eventId, index = 0, className }: EditionCardProps) {
  const status = statusConfig[edition.status]
  const TypeIcon = typeIcons[edition.type]

  const formatDate = () => {
    return new Date(edition.starts_at).toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    })
  }

  const hasVisual = Boolean(edition.banner_url ?? edition.logo_url)

  return (
    <TanstackLink
      to="/events/$eventId/editions/$editionId"
      params={{ eventId, editionId: edition.id }}
      className="block"
    >
      <motion.article
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: index * 0.08, duration: 0.4, ease: [0.25, 0.1, 0.25, 1] }}
        className={cn(
          "group relative bg-card border border-border rounded-lg overflow-hidden min-w-72",
          "hover:border-primary/30 transition-all duration-500",
          "hover:shadow-xl hover:shadow-foreground/5",
          "hover:-translate-y-1",
          className
        )}
      >
        {/* Imagem */}
        <div className="aspect-4/3 bg-muted relative overflow-hidden">
          {hasVisual ? (
            <img
              src={edition.banner_url ?? edition.logo_url ?? ""}
              alt=""
              className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
              loading={index < 2 ? "eager" : "lazy"}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-muted to-muted/50">
              <div className="w-20 h-20 rounded-full border-2 border-dashed border-border/50 flex items-center justify-center">
                <TypeIcon className="w-8 h-8 text-muted-foreground/30" />
              </div>
            </div>
          )}

          {/* Status badge */}
          <div className="absolute top-3 left-3 md:top-4 md:left-4">
            <span className={cn(
              "px-2.5 py-1 rounded-full text-xs font-medium border backdrop-blur-sm",
              status.variant === 'green' && "bg-green-500/10 text-green-600 border-green-200/50",
              status.variant === 'blue' && "bg-blue-500/10 text-blue-600 border-blue-200/50",
              status.variant === 'amber' && "bg-amber-500/10 text-amber-600 border-amber-200/50",
              status.variant === 'destructive' && "bg-destructive/10 text-destructive border-destructive/20",
              status.variant === 'muted' && "bg-muted/80 text-muted-foreground border-border",
            )}>
              {status.label}
            </span>
          </div>

          {/* Arrow */}
          <div className="absolute top-3 right-3 md:top-4 md:right-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <div className="w-8 h-8 rounded-full bg-background/90 backdrop-blur-sm flex items-center justify-center">
              <ArrowUpRight className="w-4 h-4 text-foreground" />
            </div>
          </div>
        </div>

        {/* Conteúdo */}
        <div className="p-4 md:p-5 space-y-3">
          {/* Título */}
          <div className="space-y-1">
            <h3 className="text-lg font-semibold leading-tight text-foreground group-hover:text-primary transition-colors duration-300 line-clamp-2">
              {edition.edition_name}
            </h3>
            {edition.tagline && (
              <p className="text-sm text-muted-foreground line-clamp-1">{edition.tagline}</p>
            )}
          </div>

          {/* Meta - separado em linhas */}
          <div className="space-y-1.5 text-sm text-muted-foreground">
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4 shrink-0" />
              <span>{formatDate()}</span>
            </div>
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 shrink-0" />
              <span className="truncate">{edition.location_name}</span>
            </div>
          </div>
        </div>
      </motion.article>
    </TanstackLink>
  )
}