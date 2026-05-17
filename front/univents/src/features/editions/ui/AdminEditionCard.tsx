import { motion } from "motion/react"
import { MapPin } from "lucide-react"
import { Link } from "@tanstack/react-router"
import type { EditionI } from "../model"
import { StatusBadge } from "@/shared/ui"
import { cn } from "@/shared/lib/utils"

interface EditionCardProps {
  edition: EditionI
  eventId: string
  index: number
  onPublish: () => void
  onConnect: () => void
  onDisconnect: () => void
}

export function AdminEditionCard({
  edition,
  eventId,
  index,
  onPublish,
  onConnect,
  onDisconnect
}: EditionCardProps) {
  const links = [
    { label: 'Produtos', to: '/events/$eventId/editions/$editionId/products' as const },
    { label: 'Checkpoints', to: '/admin/events/$eventId/editions/$editionId/checkpoints' as const },
  ]

  const formatDateRange = (start: string, end: string) => {
    const startDate = new Date(start)
    const endDate = new Date(end)
    const sameMonth = startDate.getMonth() === endDate.getMonth() && startDate.getFullYear() === endDate.getFullYear()

    if (sameMonth) {
      return `${startDate.getDate()}–${endDate.getDate()} ${startDate.toLocaleDateString('pt-BR', { month: 'short' })}`
    }
    return `${startDate.toLocaleDateString('pt-BR', { day: 'numeric', month: 'short' })} – ${endDate.toLocaleDateString('pt-BR', { day: 'numeric', month: 'short' })}`
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05 }}
      className={cn(
        "group relative bg-card border border-border rounded-xl p-4",
        "hover:border-primary/20 hover:shadow-sm transition-all"
      )}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap mb-1">
            <span className="text-sm font-semibold text-foreground">{edition.edition_name}</span>
            <StatusBadge status={edition.status} />
            <span className="text-xs text-muted-foreground capitalize">{edition.type}</span>
            <span className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">{edition.monetary_type}</span>
          </div>
          <div className="flex items-center gap-3 text-xs text-muted-foreground">
            <span>{formatDateRange(edition.starts_at, edition.ends_at)}</span>
            <span className="flex items-center gap-1">
              <MapPin className="w-3 h-3" />
              {edition.location_name}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-1 shrink-0">
          {edition.status === 'draft' && (
            <button
              onClick={onPublish}
              className={cn(
                "px-3 py-1.5 rounded-lg text-xs font-medium",
                "bg-primary/10 text-primary hover:bg-primary/20",
                "transition-colors"
              )}
            >
              Publicar
            </button>
          )}
          {edition.trie_payments_credential_id ? (
            <button
              onClick={onDisconnect}
              className={cn(
                "px-3 py-1.5 rounded-lg text-xs font-medium",
                "bg-destructive/10 text-destructive hover:bg-destructive/20",
                "transition-colors"
              )}
            >
              Desconectar
            </button>
          ) : (
            <button
              onClick={onConnect}
              className={cn(
                "px-3 py-1.5 rounded-lg text-xs font-medium",
                "bg-muted text-muted-foreground hover:bg-muted/80",
                "transition-colors"
              )}
            >
              Conectar pagamento
            </button>
          )}
        </div>
      </div>

      <div className="flex gap-2 mt-4 pt-4 border-t border-border/50 flex-wrap">
        {links.map(({ label, to }) => (
          <Link
            key={label}
            to={to}
            params={{ eventId, editionId: edition.id }}
            className={cn(
              "text-xs text-muted-foreground bg-muted hover:bg-muted/80",
              "px-3 py-1.5 rounded-md transition-colors font-medium"
            )}
          >
            {label}
          </Link>
        ))}
      </div>
    </motion.div>
  )
}