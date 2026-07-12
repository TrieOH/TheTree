import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useMemo } from 'react'
import { CalendarDays, CreditCard, MapPin, Share2, Ticket, Users } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import { cn } from '@/shared/lib/utils'

const statusConfig: Record<string, { label: string; dot: string }> = {
  draft: { label: 'Rascunho', dot: 'bg-amber-500' },
  announced: { label: 'Anunciada', dot: 'bg-amber-500' },
  open: { label: 'Aberta', dot: 'bg-emerald-500' },
  ongoing: { label: 'Em andamento', dot: 'bg-emerald-500' },
  finished: { label: 'Encerrada', dot: 'bg-slate-500' },
  completed: { label: 'Concluída', dot: 'bg-slate-500' },
  cancelled: { label: 'Cancelada', dot: 'bg-rose-500' },
  postponed: { label: 'Adiada', dot: 'bg-orange-500' },
}

export const Route = createLazyFileRoute('/admin/events/$eventId/editions/$editionId/')({
  component: AdminEditionDetailRoute,
})

function AdminEditionDetailRoute() {
  const { eventId, editionId } = Route.useParams()
  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const edition = useMemo(
    () => editions.find((item) => item.id === editionId) ?? null,
    [editions, editionId],
  )

  if (!edition) {
    return (
      <div className="flex h-[70vh] items-center justify-center p-4">
        <p className="text-sm text-muted-foreground">Edição não encontrada.</p>
      </div>
    )
  }

  const status = statusConfig[edition.status] ?? statusConfig.draft

  return (
    <div className="space-y-6 p-4">
      <Card className="border-border/60 bg-card/95 shadow-sm">
        <CardHeader className="border-b border-border/60 pb-4">
          <div className="space-y-3">
            <div className="flex flex-wrap items-center gap-2">
              <Badge variant="outline" className="rounded-full px-2.5">
                <span className={cn('mr-1.5 size-1.5 rounded-full', status.dot)} />
                {status.label}
              </Badge>
              <Badge variant="secondary" className="rounded-full px-2.5">
                <CalendarDays className="mr-1.5 size-3.5" />
                {edition.type}
              </Badge>
              <Badge variant="secondary" className="rounded-full px-2.5">
                <Ticket className="mr-1.5 size-3.5" />
                {edition.monetary_type === 'free' ? 'Gratuita' : edition.monetary_type === 'paid' ? 'Paga' : 'Mista'}
              </Badge>
            </div>

            <div className="space-y-1">
              <CardTitle className="text-2xl font-semibold tracking-tight">
                {edition.edition_name}
              </CardTitle>
              {edition.tagline && (
                <p className="max-w-2xl text-sm text-muted-foreground">
                  {edition.tagline}
                </p>
              )}
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-4 py-4">
          <div className="grid gap-3 md:grid-cols-3">
            <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
              <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                Período
              </p>
              <p className="mt-1.5 text-sm font-semibold text-foreground">
                {new Date(edition.starts_at).toLocaleDateString('pt-BR')} - {new Date(edition.ends_at).toLocaleDateString('pt-BR')}
              </p>
            </div>
            <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
              <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                Local
              </p>
              <p className="mt-1.5 text-sm font-semibold text-foreground">
                {edition.location_name}
              </p>
            </div>
            <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
              <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                Endereço
              </p>
              <p className="mt-1.5 line-clamp-2 text-sm text-foreground">
                {edition.location_address}
              </p>
            </div>
          </div>

          <div className="flex flex-wrap gap-2">
            <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
              <Users className="size-3.5" />
              {edition.organizer_name ?? edition.contact_email ?? 'Organizador não definido'}
            </span>
            <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
              <CreditCard className="size-3.5" />
              {edition.trie_payments_credential_id ? 'Pagamento conectado' : 'Pagamento não conectado'}
            </span>
            <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
              <MapPin className="size-3.5" />
              {edition.timezone}
            </span>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-3">
        <Link
          to="/admin/events/$eventId/editions/$editionId/activities"
          params={{ eventId, editionId }}
          className="rounded-2xl border border-border/60 bg-card/95 p-4 text-left shadow-sm transition-colors hover:bg-muted/30"
        >
          <p className="text-sm font-medium text-foreground">Atividades</p>
          <p className="mt-1 text-sm text-muted-foreground">Abrir programação da edição.</p>
        </Link>
        <Link
          to="/admin/events/$eventId/editions/$editionId/products"
          params={{ eventId, editionId }}
          className="rounded-2xl border border-border/60 bg-card/95 p-4 text-left shadow-sm transition-colors hover:bg-muted/30"
        >
          <p className="text-sm font-medium text-foreground">Produtos</p>
          <p className="mt-1 text-sm text-muted-foreground">Gerenciar itens e ingressos.</p>
        </Link>
        <Link
          to="/admin/events/$eventId/editions/$editionId/checkpoints"
          params={{ eventId, editionId }}
          className="rounded-2xl border border-border/60 bg-card/95 p-4 text-left shadow-sm transition-colors hover:bg-muted/30"
        >
          <p className="text-sm font-medium text-foreground">Checkpoints</p>
          <p className="mt-1 text-sm text-muted-foreground">Abrir o painel operacional.</p>
        </Link>
      </div>

      <div className="flex flex-wrap gap-2">
        <Button type="button" variant="outline" className="gap-2">
          <Share2 className="size-4" />
          Compartilhar
        </Button>
      </div>
    </div>
  )
}
