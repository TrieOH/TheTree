import { createLazyFileRoute } from '@tanstack/react-router'
import { Link } from '@tanstack/react-router'
import { motion } from 'motion/react'
import { useState } from 'react'
import type { ReactNode } from 'react'
import {
  CalendarDays,
  CheckCircle2,
  CircleAlert,
  Eye,
  LayoutGrid,
  Sparkles,
  ChevronRight,
  CalendarClock,
  CalendarX,
} from 'lucide-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/shadcn/card'
import { EmptyState } from '@trieoh/ui-base'
import { ManageEventModal } from '@/features/events/ui/ManageEventModal'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { ownEventsQueryOptions, patchEventFn, publishEventFn } from '@/features/events/api'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import type { EventCreateSubmitI, EventI } from '@/features/events/model'
import { cn } from '@/shared/lib/utils'

function QuickAction({
  children,
  disabled = false,
  to,
  params,
  onClick,
}: {
  children: ReactNode
  disabled?: boolean
} & (
    | { to: '/events/$eventId' | '/admin/events/$eventId/editions'; params: { eventId: string }; onClick?: never }
    | { to?: never; params?: never; onClick: () => void }
  )) {
  const baseClassName = cn(
    'flex items-center justify-between rounded-2xl border border-dashed border-border/70 bg-muted/15 px-4 py-4 text-left',
    'transition-colors hover:border-border hover:bg-muted/30',
    disabled && 'cursor-not-allowed opacity-60 hover:border-border/70 hover:bg-muted/15',
  )

  if (to) {
    return (
      <Link to={to} params={params} className={baseClassName}>
        {children}
      </Link>
    )
  }

  return (
    <button type="button" className={baseClassName} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  )
}

const statusConfig = {
  draft: {
    label: 'Rascunho',
    className: 'bg-amber-500/10 text-amber-700 border-amber-500/20',
    dot: 'bg-amber-500',
  },
  active: {
    label: 'Ativo',
    className: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
    dot: 'bg-emerald-500',
  },
  archived: {
    label: 'Arquivado',
    className: 'bg-slate-500/10 text-slate-700 border-slate-500/20',
    dot: 'bg-slate-500',
  },
  discontinued: {
    label: 'Descontinuado',
    className: 'bg-rose-500/10 text-rose-700 border-rose-500/20',
    dot: 'bg-rose-500',
  },
} as const

const editionStatusHint: Record<string, string> = {
  draft: 'Ainda em edição',
  announced: 'Pronta para divulgação',
  published: 'Visível ao público',
  canceled: 'Cancelada',
}

const editionTypeLabel: Record<string, string> = {
  year: 'Anual',
  season: 'Por temporada',
  number: 'Numerada',
  ordinal: 'Ordinal',
  custom: 'Personalizada',
}

const editionMonetaryLabel: Record<string, string> = {
  free: 'Gratuita',
  paid: 'Paga',
  mixed: 'Inscrições gratuitas e pagas',
}

export const Route = createLazyFileRoute('/admin/events/$eventId/')({
  component: EventOverviewRoute,
})

function EventOverviewRoute() {
  const queryClient = useQueryClient()
  const { eventId } = Route.useParams()
  const { data: events = [] } = useQuery(ownEventsQueryOptions())
  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false)
  const event = events.find((item) => item.id === eventId) ?? null
  const latestEdition = editions[0] ?? null
  const isPublished = event?.status === 'active'
  const status = event ? statusConfig[event.status] : statusConfig.draft

  const copyLink = () => {
    if (!event) return
    void navigator.clipboard.writeText(`${window.location.origin}/events/${event.id}`)
    toast.success('Link copiado')
  }

  const handlePublishEvent = () => {
    if (!event || isPublished) return
    publishEventMutation.mutate()
  }

  const publishEventMutation = useMutation({
    mutationFn: () => publishEventFn(eventId),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar evento')
        return
      }

      queryClient.setQueryData(
        ownEventsQueryOptions().queryKey,
        (old: typeof events = []) => old.map((item) => (
          item.id === eventId ? { ...item, status: 'active' as const } : item
        )),
      )
      toast.success('Evento publicado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const patchMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<EventCreateSubmitI> }) =>
      patchEventFn(id, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao atualizar evento')
        return
      }

      queryClient.setQueryData<EventI[]>(
        ownEventsQueryOptions().queryKey,
        (old = []) => old.map((item) => (item.id === res.data.id ? res.data : item)),
      )
      toast.success('Evento atualizado com sucesso!')
      setEditModalOpen(false)
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const metrics = [
    {
      label: 'Criado em',
      value: event ? format(new Date(event.created_at), 'dd MMM yyyy', { locale: ptBR }) : '—',
      hint: 'Data de criação do evento',
    },
    {
      label: 'Atualizado em',
      value: event ? format(new Date(event.updated_at), 'dd MMM yyyy', { locale: ptBR }) : '—',
      hint: 'Última alteração registrada',
    },
    {
      label: 'Contato',
      value: event?.contact_email ?? '—',
      hint: 'E-mail principal do evento',
    },
  ]
  const heroDescription = event?.tagline ?? event?.description ?? 'Sem descrição cadastrada para este evento.'
  const recentEditionDate = latestEdition
    ? format(new Date(latestEdition.starts_at), 'dd MMM yyyy', { locale: ptBR })
    : 'Sem edição recente'
  const checklist = [
    {
      label: 'Banner e logo cadastrados',
      done: Boolean(event?.banner_url || event?.logo_url),
    },
    {
      label: 'Descrição ou tagline preenchida',
      done: Boolean(event?.tagline || event?.description),
    },
    {
      label: 'Slug público disponível',
      done: Boolean(event?.id),
    },
    {
      label: 'Última edição cadastrada',
      done: Boolean(latestEdition),
    },
  ]

  return (
    <div className="relative overflow-hidden p-2">

      <div className="relative space-y-6">
        <motion.section
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35 }}
          className="overflow-hidden rounded-[1.75rem] border border-border/60 bg-card shadow-[0_1px_0_0_rgba(255,255,255,0.03),0_20px_40px_-24px_rgba(15,23,42,0.24)]"
        >
          <div className="relative flex flex-col gap-6 p-6">
            <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-primary/20 to-transparent" />
            <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-primary/5 via-transparent to-transparent" />

            <div className="space-y-5">
              <div className="flex items-center gap-2">
                <Badge variant="secondary" className="rounded-full px-3">
                  <LayoutGrid className="size-3.5" />
                  Overview
                </Badge>
                {event?.is_series && (
                  <Badge variant="outline" className="rounded-full px-3">
                    <Sparkles className="size-3.5" />
                    Série
                  </Badge>
                )}
              </div>

              <div className="space-y-3">
                <h1 className="max-w-3xl text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
                  {event?.name ?? 'Evento'}
                </h1>
                <p className="max-w-2xl text-sm leading-6 text-muted-foreground md:text-base">
                  {heroDescription}
                </p>
              </div>

              <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                  <CalendarDays className="size-3.5" />
                  {event?.slug ?? 'slug-do-evento'}
                </span>
                <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                  <Eye className="size-3.5" />
                  Link público
                </span>
                <span className={cn('inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1', status.className)}>
                  <span className={cn('size-1.5 rounded-full', status.dot)} />
                  {status.label}
                </span>
              </div>

              <div className="flex flex-wrap items-center gap-2 pt-1">
                <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
                  <LayoutGrid className="size-3.5" />
                  {event?.editions_count ?? '—'} edições
                </span>
                <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
                  {event?.is_series ? 'Série' : 'Único'}
                </span>
              </div>
            </div>
          </div>
        </motion.section>

        <section className="grid gap-4 md:grid-cols-3">
          {metrics.map((metric) => (
            <Card key={metric.label} className="border-border/60 bg-card/95 shadow-sm transition-shadow hover:shadow-md">
              <CardHeader className="pb-2">
                <CardDescription className="flex items-center gap-2 text-xs uppercase tracking-[0.22em]">
                  <span className="size-1.5 rounded-full bg-primary/60" />
                  {metric.label}
                </CardDescription>
                <CardTitle className="text-2xl font-semibold tracking-tight">{metric.value}</CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <p className="text-xs text-muted-foreground">{metric.hint}</p>
              </CardContent>
            </Card>
          ))}
        </section>

        <section className="grid gap-4 lg:grid-cols-[1.3fr_0.95fr]">
          <Card className="border-border/60 bg-card/95 shadow-sm">
            <CardHeader className="border-b border-border/60">
              <CardTitle>Checklist do evento</CardTitle>
              <CardDescription>
                Itens derivados dos dados já cadastrados no evento.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-3 py-5">
              {checklist.map((item) => (
                <div
                  key={item.label}
                  className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/15 px-4 py-3.5"
                >
                  <div className="flex items-center gap-3">
                    <div className={cn('size-2 rounded-full', item.done ? 'bg-emerald-500' : 'bg-amber-500')} />
                    <span className="text-sm text-foreground">{item.label}</span>
                  </div>
                  {item.done ? (
                    <CheckCircle2 className="size-4 text-emerald-500/70" />
                  ) : (
                    <CircleAlert className="size-4 text-amber-500/70" />
                  )}
                </div>
              ))}
            </CardContent>
          </Card>

          <Card className="border-border/60 bg-card/95 shadow-sm">
            <CardHeader className="border-b border-border/60 pb-3">
              <CardTitle>Edição recente</CardTitle>
              <CardDescription>
                Última edição cadastrada no admin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 py-4">
              {latestEdition ? (
                <>
                  <div className="rounded-2xl border border-border/60 bg-muted/10 px-4 py-4">
                    <div className="flex items-start justify-between gap-4">
                      <div className="min-w-0 space-y-3">
                        <div className="flex flex-wrap items-center gap-2">
                          <Badge variant="outline" className="rounded-full px-2.5">
                            <CalendarClock className="size-3.5" />
                            Última edição
                          </Badge>
                          <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs font-medium text-muted-foreground">
                            <span className="size-1.5 rounded-full bg-primary/70" />
                            {editionMonetaryLabel[latestEdition.monetary_type]}
                          </span>
                          <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
                            {editionStatusHint[latestEdition.status] ?? 'Status interno da edição'}
                          </span>
                        </div>

                        <div className="space-y-1.5">
                          <h3 className="truncate text-lg font-semibold tracking-tight text-foreground">
                            {latestEdition.edition_name}
                          </h3>
                          <p className="line-clamp-2 text-sm leading-6 text-muted-foreground">
                            {latestEdition.tagline ?? latestEdition.description ?? latestEdition.location_name}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="grid gap-2.5 sm:grid-cols-3">
                    <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
                      <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                        Data
                      </p>
                      <p className="mt-1.5 text-sm font-semibold text-foreground">
                        {recentEditionDate}
                      </p>
                    </div>
                    <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
                      <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                        Local
                      </p>
                      <p className="mt-1.5 truncate text-sm font-semibold text-foreground">
                        {latestEdition.location_name}
                      </p>
                    </div>
                    <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
                      <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                        Frequência
                      </p>
                      <p className="mt-1.5 text-sm font-semibold text-foreground">
                        {editionTypeLabel[latestEdition.type] ?? latestEdition.type}
                      </p>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-2">
                    <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
                      <span className="size-1.5 rounded-full bg-primary/70" />
                      {latestEdition.timezone}
                    </span>
                    {latestEdition.organizer_name || latestEdition.contact_email ? (
                      <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-background/70 px-2.5 py-1 text-xs text-muted-foreground">
                        <span className="size-1.5 rounded-full bg-primary/70" />
                        {latestEdition.organizer_name
                          ? `Organizador: ${latestEdition.organizer_name}`
                          : `Contato: ${latestEdition.contact_email}`}
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1.5 rounded-full border border-amber-500/20 bg-amber-500/10 px-2.5 py-1 text-xs text-amber-700">
                        <span className="size-1.5 rounded-full bg-amber-500" />
                        Organizador e contato não definidos
                      </span>
                    )}
                  </div>
                </>
              ) : (
                <EmptyState
                  icon={CalendarX}
                  eyebrow="Edição recente"
                  title="Nenhuma edição cadastrada ainda"
                  description="Quando a primeira edição for criada, ela aparece aqui com os detalhes principais e o link direto."
                  className="border-0 bg-transparent px-0 py-2 shadow-none"
                />
              )}
            </CardContent>
          </Card>
        </section>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Ações rápidas</CardTitle>
            <CardDescription>
              Atalhos para as operações mais comuns do evento.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-3 py-5 sm:grid-cols-2 xl:grid-cols-3">
            <QuickAction onClick={() => setEditModalOpen(true)}>
              <span className="text-sm font-medium text-foreground">Editar evento</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </QuickAction>

            {!isPublished && (
              <Button
                type="button"
                variant="outline"
                className="h-auto justify-between rounded-2xl border-dashed bg-muted/15 px-4 py-4 text-left hover:bg-muted/30"
                onClick={() => setPublishConfirmOpen(true)}
                disabled={publishEventMutation.isPending || !event}
              >
                <span className="text-sm font-medium text-foreground">Publicar evento</span>
                <ChevronRight className="size-4 text-muted-foreground" />
              </Button>
            )}

            <QuickAction to="/admin/events/$eventId/editions" params={{ eventId }}>
              <span className="text-sm font-medium text-foreground">Criar edição</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </QuickAction>

            {isPublished && (
              <QuickAction onClick={copyLink} disabled={!event}>
                <span className="text-sm font-medium text-foreground">Copiar link</span>
                <ChevronRight className="size-4 text-muted-foreground" />
              </QuickAction>
            )}

            {isPublished && (
              <QuickAction to="/events/$eventId" params={{ eventId }}>
                <span className="text-sm font-medium text-foreground">Abrir painel público</span>
                <ChevronRight className="size-4 text-muted-foreground" />
              </QuickAction>
            )}
          </CardContent>
        </Card>
      </div>

      <ManageEventModal
        key={event?.id ?? 'event-overview-edit'}
        open={editModalOpen}
        onOpenChange={setEditModalOpen}
        event={event ?? undefined}
        onUpdate={async (id, values) => {
          const res = await patchMutation.mutateAsync({ id, data: values })
          return res.success ? res.data : null
        }}
      />

      <AlertModal
        open={publishConfirmOpen}
        onOpenChange={setPublishConfirmOpen}
        title="Publicar evento?"
        description="Depois de publicar, o painel público ficará disponível para o evento."
        confirmLabel="Publicar evento"
        variant="default"
        loading={publishEventMutation.isPending}
        onConfirm={async () => {
          handlePublishEvent()
          setPublishConfirmOpen(false)
        }}
      />
    </div>
  )
}
