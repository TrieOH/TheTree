import { ArrowUpRight, Clock, MapPin, Mail, Phone } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import type { EditionI } from '@/features/editions/model'
import { Button } from '@/shared/ui/shadcn/button'

interface OverviewTabProps {
  edition: EditionI
  eventId: string
}

export function OverviewTab({ edition, eventId }: OverviewTabProps) {
  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between py-4 border-b border-border">
        <div>
          <p className="font-medium text-foreground">Produtos exclusivos</p>
          <p className="text-sm text-muted-foreground">Itens limitados disponíveis</p>
        </div>
        <Link
          to="/events/$eventId/editions/$editionId/products"
          params={{ eventId, editionId: edition.id }}
        >
          <Button variant="ghost" className="gap-1.5">
            Ver produtos
            <ArrowUpRight className="w-4 h-4" />
          </Button>
        </Link>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div className="space-y-3">
          <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <Clock className="w-4 h-4" />
            Quando
          </h3>
          <div className="space-y-2 text-sm">
            <p><span className="text-muted-foreground">Início:</span> {formatDate(edition.starts_at)}</p>
            <p><span className="text-muted-foreground">Término:</span> {formatDate(edition.ends_at)}</p>
            {edition.registration_opens_at && (
              <p><span className="text-muted-foreground">Inscrições abrem:</span> {formatDate(edition.registration_opens_at)}</p>
            )}
            {edition.registration_closes_at && (
              <p><span className="text-muted-foreground">Inscrições fecham:</span> {formatDate(edition.registration_closes_at)}</p>
            )}
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <MapPin className="w-4 h-4" />
            Onde
          </h3>
          <div className="space-y-2 text-sm">
            <p className="font-medium text-foreground">{edition.location_name}</p>
            <p className="text-muted-foreground">{edition.location_address}</p>
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
            Contato
          </h3>
          <div className="space-y-2 text-sm">
            {edition.contact_email && (
              <p className="flex items-center gap-2">
                <Mail className="w-4 h-4 text-muted-foreground" />
                <a href={`mailto:${edition.contact_email}`} className="hover:text-primary transition-colors">
                  {edition.contact_email}
                </a>
              </p>
            )}
            {edition.contact_phone && (
              <p className="flex items-center gap-2">
                <Phone className="w-4 h-4 text-muted-foreground" />
                <a href={`tel:${edition.contact_phone}`} className="hover:text-primary transition-colors">
                  {edition.contact_phone}
                </a>
              </p>
            )}
            {edition.organizer_name && (
              <p className="text-muted-foreground">Organizado por: {edition.organizer_name}</p>
            )}
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
            Detalhes
          </h3>
          <div className="space-y-2 text-sm">
            <p><span className="text-muted-foreground">Tipo:</span> <span className="capitalize">{edition.type}</span></p>
            <p><span className="text-muted-foreground">Modalidade:</span> <span className="capitalize">{edition.monetary_type}</span></p>
            <p><span className="text-muted-foreground">Status:</span> <span className="capitalize">{edition.status.replace('_', ' ')}</span></p>
          </div>
        </div>
      </div>
    </div>
  )
}