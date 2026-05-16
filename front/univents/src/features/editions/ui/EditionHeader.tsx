import { ArrowLeft, Calendar, MapPin, Clock, Ticket } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import type { EditionI } from '@/features/editions/model'

interface EditionHeaderProps {
  edition: EditionI
  eventId: string
}

export function EditionHeader({ edition, eventId }: EditionHeaderProps) {
  const formatDateRange = () => {
    const start = new Date(edition.starts_at)
    const end = new Date(edition.ends_at)
    const sameDay = start.toDateString() === end.toDateString()

    const opts: Intl.DateTimeFormatOptions = { day: '2-digit', month: 'long', year: 'numeric' }

    if (sameDay) {
      return start.toLocaleDateString('pt-BR', opts)
    }

    return `${start.toLocaleDateString('pt-BR', opts)} - ${end.toLocaleDateString('pt-BR', opts)}`
  }

  const formatTime = () => {
    const start = new Date(edition.starts_at)
    return start.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2">
        <Link
          to="/events/$eventId/editions"
          params={{ eventId }}
          className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Voltar para edições</span>
        </Link>
      </div>

      <div className="space-y-3">
        <h1 className="text-3xl md:text-4xl lg:text-5xl font-bold tracking-tight text-foreground">
          {edition.edition_name}
        </h1>
        {edition.tagline && (
          <p className="text-xl text-muted-foreground">{edition.tagline}</p>
        )}
        {edition.description && (
          <p className="text-base text-muted-foreground/80 max-w-3xl leading-relaxed">
            {edition.description}
          </p>
        )}
      </div>

      {/* Meta info */}
      <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4" />
          <span>{formatDateRange()}</span>
        </div>
        <div className="flex items-center gap-2">
          <Clock className="w-4 h-4" />
          <span>Início às {formatTime()}</span>
        </div>
        <div className="flex items-center gap-2">
          <MapPin className="w-4 h-4" />
          <span>{edition.location_name}</span>
        </div>
        <div className="flex items-center gap-2">
          <Ticket className="w-4 h-4" />
          <span className="capitalize">{edition.monetary_type === 'free' ? 'Gratuito' : edition.monetary_type === 'paid' ? 'Pago' : 'Misto'}</span>
        </div>
      </div>
    </div>
  )
}