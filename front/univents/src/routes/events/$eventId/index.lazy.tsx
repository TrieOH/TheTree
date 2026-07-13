import { createLazyFileRoute, Link } from '@tanstack/react-router'
import {
  Share2, ArrowRight,
  Calendar, Mail, Link2, Globe, X as XIcon, Camera
} from 'lucide-react'
import { toast } from 'sonner'
import { useQuery } from '@tanstack/react-query'
import { Button } from '@/shared/ui/shadcn/button'
import { Badge } from '@/shared/ui/shadcn/badge'
import { cn } from '@/shared/lib/utils'
import { parseDatetimeLocal } from '@/shared/lib/date'
import { publicEventQueryOptions } from '@/features/events/api'
import { InfoRow, SectionCard, SocialChip } from '@/features/events/ui/EventDetailComponents'

export const Route = createLazyFileRoute('/events/$eventId/')({
  component: RouteComponent,
})

const statusColors = {
  active: 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20',
  inactive: 'bg-gray-500/10 text-gray-600 border-gray-500/20',
  draft: 'bg-amber-500/10 text-amber-600 border-amber-500/20',
  archived: 'bg-red-500/10 text-red-600 border-red-500/20',
  discontinued: 'bg-zinc-500/10 text-zinc-600 border-zinc-500/20',
}

const statusLabels = {
  active: 'Ativo',
  inactive: 'Inativo',
  draft: 'Rascunho',
  archived: 'Arquivado',
  discontinued: 'Descontinuado',
}

function RouteComponent() {
  const { eventId } = Route.useParams()

  const { data: event, isLoading, error } = useQuery(publicEventQueryOptions(eventId))

  const handleShare = async () => {
    if (!event) return
    const url = window.location.href
    try {
      if (typeof navigator.share === 'function') {
        await navigator.share({ title: event.name, url })
        return
      }
      await navigator.clipboard.writeText(url)
      toast.success('Link copiado!')
    } catch {
      toast.error('Erro ao compartilhar')
    }
  }

  if (isLoading) {
    return (
      <div className="flex h-[80vh] w-full items-center justify-center">
        <div className="text-muted-foreground">Carregando...</div>
      </div>
    )
  }

  if (error || !event) {
    return (
      <div className="flex h-[80vh] w-full flex-col items-center justify-center gap-4">
        <p className="text-muted-foreground">Erro ao carregar evento</p>
        <Link to="/events" className="text-primary hover:underline">Voltar para eventos</Link>
      </div>
    )
  }

  const getInitials = (name: string) =>
    name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()

  const logoFallbackText = event.acronym
    ?? (event.name ? getInitials(event.name) : '')

  const logoFontSize =
    logoFallbackText.length <= 2 ? 'text-lg' :
      logoFallbackText.length === 3 ? 'text-sm' : 'text-xs'

  const hasSocialLinks = !!(
    event.social_links?.website ??
    event.social_links?.twitter ??
    event.social_links?.instagram ??
    event.social_links?.linkedin
  )

  return (
    <div className="min-h-screen bg-background pb-24">
      <div className="relative">
        <div className="relative w-full h-40 sm:h-52 md:h-64">
          {event.banner_url ? (
            <img
              src={event.banner_url}
              alt={event.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-linear-to-br from-muted via-primary/25 to-secondary/25" />
          )}
          <div className="absolute inset-x-0 top-0 h-32 bg-linear-to-b from-muted/20 via-primary/10 to-transparent" />
          <div className="absolute inset-x-0 bottom-0 h-2/3 bg-linear-to-t from-background/80 via-secondary/20 to-transparent" />
          <div className="absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-background to-transparent" />
        </div>
      </div>

      <div className="mx-auto max-w-2xl px-4">
        <main className="space-y-2">
          <div className="bg-card rounded-xl shadow-xl border border-border/50">
            <div className="px-4 relative pt-0 flex items-end justify-between -mt-8">
              <div className="h-16 w-16 sm:h-20 sm:w-20 shrink-0 rounded-xl bg-primary shadow-lg flex items-center justify-center overflow-hidden ring-4 ring-card">
                {event.logo_url ? (
                  <img
                    src={event.logo_url}
                    alt={event.acronym ?? event.name}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <span className={cn(
                    'text-primary-foreground font-bold leading-none tracking-tight px-1 text-center break-all',
                    logoFontSize
                  )}>
                    {logoFallbackText}
                  </span>
                )}
              </div>

              <div className="flex flex-col items-end absolute right-4 top-10">
                <div className="flex items-center gap-2">
                  <Button
                    className="h-8 w-8 text-muted-foreground hover:text-foreground transition-colors"
                    variant="ghost"
                    size="icon"
                    onClick={() => void handleShare()}
                    title="Compartilhar"
                  >
                    <Share2 className="h-4 w-4" />
                  </Button>
                  <Badge
                    variant="outline"
                    className={cn('text-xs font-medium whitespace-nowrap', statusColors[event.status])}
                  >
                    {statusLabels[event.status]}
                  </Badge>
                </div>
                {event.is_series && event.editions_count > 0 && (
                  <span className="text-[11px] text-muted-foreground uppercase tracking-wider font-medium">
                    {event.editions_count}ª edição
                  </span>
                )}
              </div>
            </div>

            <div className="px-4 pt-2 pb-4 space-y-1">
              <h1 className="text-xl sm:text-2xl font-bold tracking-tight text-foreground leading-tight">
                {event.name}
              </h1>

              {/* Acronym + Tagline */}
              {(event.acronym || event.tagline) && (
                <div className="flex flex-wrap items-center gap-x-2 gap-y-0.5 text-sm">
                  {event.acronym && (
                    <span className="font-semibold text-primary">
                      {event.acronym}
                    </span>
                  )}
                  {event.acronym && event.tagline && (
                    <span className="text-muted-foreground/40 select-none">·</span>
                  )}
                  {event.tagline && (
                    <span className="text-muted-foreground italic">
                      {event.tagline}
                    </span>
                  )}
                </div>
              )}

              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <Calendar className="h-3.5 w-3.5 shrink-0" />
                <span>{parseDatetimeLocal(event.created_at).toLocaleDateString('pt-BR')}</span>
              </div>

              <div className="pt-2">
                <Link
                  to="/events/$eventId/editions"
                  params={{ eventId: event.id }}
                  className={cn(
                    'flex items-center justify-center gap-2',
                    'w-full bg-primary hover:bg-primary/90 active:scale-[.98]',
                    'text-primary-foreground rounded-xl h-10',
                    'text-sm font-semibold transition-all duration-150',
                  )}
                >
                  Ver Edições
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </div>
          </div>

          {/* Descrição */}
          {event.description && (
            <SectionCard label="Sobre">
              <div className="text-sm text-foreground/75 leading-relaxed w-full whitespace-pre-wrap">
                {event.description}
              </div>
            </SectionCard>
          )}

          <SectionCard label="Informações">
            <div className="grid grid-cols-2 gap-x-4 gap-y-3">
              <InfoRow label="Slug" value={`/${event.slug}`} mono />
              {event.is_series && (
                <InfoRow label="Edições" value={String(event.editions_count)} />
              )}
            </div>
          </SectionCard>

          {event.gallery_urls && event.gallery_urls.length > 0 && (
            <SectionCard label="Galeria">
              <div className="grid grid-cols-3 gap-1.5">
                {event.gallery_urls.map((url, i) => (
                  <a
                    key={i}
                    href={url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="aspect-square rounded-lg overflow-hidden bg-muted block"
                  >
                    <img
                      src={url}
                      alt={`Galeria ${i + 1}`}
                      className="h-full w-full object-cover hover:scale-105 transition-transform duration-200"
                    />
                  </a>
                ))}
              </div>
            </SectionCard>
          )}

          {(Boolean(event.contact_email) || hasSocialLinks) && (
            <SectionCard label="Contato">
              <div className="flex flex-wrap gap-2">
                {event.contact_email && (
                  <a
                    href={`mailto:${event.contact_email}`}
                    className={cn(
                      'flex items-center gap-2 px-3 py-2 rounded-lg',
                      'bg-background border border-border',
                      'text-sm text-foreground/80 hover:text-foreground',
                      'hover:bg-muted transition-colors min-w-0',
                    )}
                  >
                    <Mail className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                    <span className="truncate max-w-50">{event.contact_email}</span>
                  </a>
                )}
                {event.social_links?.website && (
                  <SocialChip href={event.social_links.website} label="Website" icon={<Globe className="h-3.5 w-3.5" />} />
                )}
                {event.social_links?.twitter && (
                  <SocialChip href={event.social_links.twitter} label="Twitter" icon={<XIcon className="h-3.5 w-3.5" />} />
                )}
                {event.social_links?.instagram && (
                  <SocialChip href={event.social_links.instagram} label="Instagram" icon={<Camera className="h-3.5 w-3.5" />} />
                )}
                {event.social_links?.linkedin && (
                  <SocialChip href={event.social_links.linkedin} label="LinkedIn" icon={<Link2 className="h-3.5 w-3.5" />} />
                )}
              </div>
            </SectionCard>
          )}
        </main>
      </div>
    </div>
  )
}
