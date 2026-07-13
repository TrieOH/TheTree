import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import {
  CalendarDays,
  CreditCard,
  CheckCircle2,
  CircleAlert,
  Eye,
  LayoutGrid,
  MapPin,
  Share2,
  Ticket,
  Users,
  ChevronRight,
  CalendarClock,
  CalendarX,
} from 'lucide-react'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { motion } from 'motion/react'
import { useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/shadcn/card'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import { usePublishEditionMutation, useUpdateEditionMutation } from '@/features/editions/api/mutations'
import { ManageEditionModal } from '@/features/editions/ui/ManageEditionModal'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { cn } from '@/shared/lib/utils'

const statusConfig: Record<string, { label: string; className: string; dot: string }> = {
  draft: {
    label: 'Rascunho',
    className: 'bg-amber-500/10 text-amber-700 border-amber-500/20',
    dot: 'bg-amber-500',
  },
  announced: {
    label: 'Anunciada',
    className: 'bg-amber-500/10 text-amber-700 border-amber-500/20',
    dot: 'bg-amber-500',
  },
  open: {
    label: 'Aberta',
    className: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
    dot: 'bg-emerald-500',
  },
  ongoing: {
    label: 'Em andamento',
    className: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20',
    dot: 'bg-emerald-500',
  },
  finished: {
    label: 'Encerrada',
    className: 'bg-slate-500/10 text-slate-700 border-slate-500/20',
    dot: 'bg-slate-500',
  },
  completed: {
    label: 'Concluída',
    className: 'bg-slate-500/10 text-slate-700 border-slate-500/20',
    dot: 'bg-slate-500',
  },
  cancelled: {
    label: 'Cancelada',
    className: 'bg-rose-500/10 text-rose-700 border-rose-500/20',
    dot: 'bg-rose-500',
  },
  postponed: {
    label: 'Adiada',
    className: 'bg-orange-500/10 text-orange-700 border-orange-500/20',
    dot: 'bg-orange-500',
  },
}

const editionStatusHint: Record<string, string> = {
  draft: 'Ainda em edição',
  announced: 'Pronta para divulgação',
  open: 'Aberta ao público',
  ongoing: 'Acontecendo agora',
  finished: 'Já encerrada',
  completed: 'Concluída',
  cancelled: 'Cancelada',
  postponed: 'Adiada',
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
  mixed: 'Mista',
}

function formatDate(value: string | null | undefined) {
  if (!value) return '—'
  const parsed = new Date(value)
  return Number.isNaN(parsed.getTime()) ? '—' : format(parsed, 'dd MMM yyyy', { locale: ptBR })
}

function textOr(value: string | null | undefined, fallback: string) {
  return value && value.trim() ? value : fallback
}

function hasText(value: unknown) {
  return typeof value === 'string' && value.trim().length > 0
}

function QuickAction({
  children,
  to,
  params,
}: {
  children: ReactNode
  to:
  | '/events/$eventId/editions/$editionId'
  | '/admin/events/$eventId/editions/$editionId/activities'
  | '/admin/events/$eventId/editions/$editionId/products'
  | '/admin/events/$eventId/editions/$editionId/checkpoints'
  | '/admin/events/$eventId/editions/$editionId/certifications'
  | '/admin/events/$eventId/editions/$editionId/signatures'
  params: { eventId: string; editionId: string }
}) {
  return (
    <Link
      to={to}
      params={params}
      className={cn(
        'group flex items-center justify-between rounded-2xl border border-dashed border-border/70 bg-muted/15 px-4 py-4 text-left',
        'transition-colors hover:border-border hover:bg-muted/30',
      )}
    >
      {children}
    </Link>
  )
}

export const Route = createLazyFileRoute('/admin/events/$eventId_/editions/$editionId/')({
  component: AdminEditionDetailRoute,
})

function AdminEditionDetailRoute() {
  const { eventId, editionId } = Route.useParams()
  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false)
  const publishEditionMutation = usePublishEditionMutation()
  const updateEditionMutation = useUpdateEditionMutation()
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
  const hasContact = hasText(edition.organizer_name) || hasText(edition.contact_email)
  const hasLocation = hasText(edition.location_name) && hasText(edition.location_address)
  const hasMedia = hasText(edition.logo_url) || hasText(edition.banner_url)
  const hasPayment = hasText(edition.trie_payments_credential_id)
  const isOperational = ['open', 'ongoing', 'finished', 'completed'].includes(edition.status)
  const metrics = [
    {
      label: 'Período',
      value: `${formatDate(edition.starts_at)} - ${formatDate(edition.ends_at)}`,
      hint: 'Data de início e fim da edição',
    },
    {
      label: 'Local',
      value: textOr(edition.location_name, 'Não definido'),
      hint: textOr(edition.location_address, 'Endereço não definido'),
    },
    {
      label: 'Contato',
      value: hasText(edition.organizer_name)
        ? edition.organizer_name
        : hasText(edition.contact_email)
          ? edition.contact_email
          : 'Não definido',
      hint: hasText(edition.contact_email) ? edition.contact_email : 'Organização da edição',
    },
  ]

  const checklist = [
    {
      label: 'Nome da edição preenchido',
      done: hasText(edition.edition_name),
    },
    {
      label: 'Período cadastrado',
      done: Boolean(edition.starts_at && edition.ends_at),
    },
    {
      label: 'Local e endereço preenchidos',
      done: hasLocation,
    },
    {
      label: 'Banner ou logo cadastrados',
      done: hasMedia,
    },
    {
      label: 'Contato ou responsável definido',
      done: hasContact,
    },
    {
      label: 'Pagamento conectado',
      done: hasPayment,
    },
    {
      label: 'Já saiu do rascunho',
      done: isOperational || edition.status === 'announced',
    },
  ]

  const actions = [
    {
      label: 'Atividades',
      description: 'Abrir programação da edição.',
      to: '/admin/events/$eventId/editions/$editionId/activities',
      icon: CalendarDays,
    },
    {
      label: 'Produtos',
      description: 'Gerenciar itens e ingressos.',
      to: '/admin/events/$eventId/editions/$editionId/products',
      icon: Ticket,
    },
    {
      label: 'Checkpoints',
      description: 'Abrir o painel operacional.',
      to: '/admin/events/$eventId/editions/$editionId/checkpoints',
      icon: LayoutGrid,
    },
    {
      label: 'Certificações',
      description: 'Ajustar templates e emissão.',
      to: '/admin/events/$eventId/editions/$editionId/certifications',
      icon: CalendarClock,
    },
    {
      label: 'Assinaturas',
      description: 'Gerenciar assinaturas e selos.',
      to: '/admin/events/$eventId/editions/$editionId/signatures',
      icon: Eye,
    },
  ]

  const copyLink = async () => {
    await navigator.clipboard.writeText(
      `${window.location.origin}/events/${eventId}/editions/${editionId}`,
    )
    toast.success('Link da edição copiado')
  }

  const handlePublish = () => {
    if (edition.status === 'announced') return
    publishEditionMutation.mutate({ eventId, editionId })
  }

  return (
    <div className="relative space-y-6">
      <motion.section
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.35 }}
        className="overflow-hidden rounded-md border border-border/60 bg-card shadow-[0_1px_0_0_rgba(255,255,255,0.03),0_20px_40px_-24px_rgba(15,23,42,0.24)]"
      >
        <div className="relative flex flex-col gap-6 p-6">
          <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-primary/20 to-transparent" />
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-primary/5 via-transparent to-transparent" />

          <div className="space-y-5">
            <div className="flex flex-wrap items-center gap-2">
              <Badge variant="outline" className="rounded-full px-3">
                <span className={cn('mr-1.5 size-1.5 rounded-full', status.dot)} />
                {status.label}
              </Badge>
              <Badge variant="secondary" className="rounded-full px-3">
                <CalendarDays className="mr-1.5 size-3.5" />
                {editionTypeLabel[edition.type] ?? edition.type}
              </Badge>
              <Badge variant="secondary" className="rounded-full px-3">
                <Ticket className="mr-1.5 size-3.5" />
                {editionMonetaryLabel[edition.monetary_type] ?? edition.monetary_type}
              </Badge>
            </div>

            <div className="space-y-3">
              <h1 className="max-w-3xl text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
                {edition.edition_name}
              </h1>
              <p className="max-w-2xl text-sm leading-6 text-muted-foreground md:text-base">
                {textOr(edition.tagline, textOr(edition.description, 'Detalhes e gestão da edição.'))}
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                <CalendarDays className="size-3.5" />
                {textOr(edition.timezone, 'Fuso não definido')}
              </span>
              <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                <Users className="size-3.5" />
                {hasText(edition.organizer_name)
                  ? edition.organizer_name
                  : hasText(edition.contact_email)
                    ? edition.contact_email
                    : 'Organizador não definido'}
              </span>
            </div>

            <div className="flex flex-wrap items-center gap-2 pt-1">
              <Button
                type="button"
                variant="outline"
                className="gap-2 rounded-full bg-background/80 backdrop-blur-md border-border/20 text-foreground"
                onClick={() => void copyLink()}
              >
                <Share2 className="size-4" />
                Compartilhar
              </Button>
              <Button
                type="button"
                variant="outline"
                className="gap-2 rounded-full"
                onClick={() => setEditModalOpen(true)}
              >
                <CalendarDays className="size-4" />
                Editar edição
              </Button>
              <span className="inline-flex items-center gap-2 rounded-full border border-border/60 bg-muted/60 px-3 py-2 text-xs text-muted-foreground">
                <CalendarClock className="size-4" />
                {editionStatusHint[edition.status] ?? 'Detalhes da edição'}
              </span>
              {edition.status === 'cancelled' && (
                <Button type="button" variant="outline" className="gap-2 rounded-full">
                  <CalendarX className="size-4" />
                  Cancelada
                </Button>
              )}
            </div>
          </div>
        </div>
      </motion.section>

      <section className="grid gap-4 md:grid-cols-3">
        {metrics.map((metric) => (
          <Card key={metric.label} className="border-border/60 bg-card/95 shadow-sm transition-shadow hover:shadow-md">
            <CardHeader className="pb-2">
              <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                {metric.label}
              </p>
              <CardTitle className="text-xl font-semibold tracking-tight">
                {metric.value}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
              <p className="text-sm text-muted-foreground">
                {metric.hint}
              </p>
            </CardContent>
          </Card>
        ))}
      </section>

      <section className="grid gap-4 md:grid-cols-3">
        {actions.map((action) => (
          <Link
            key={action.label}
            to={action.to}
            params={{ eventId, editionId }}
            className={cn(
              'group rounded-2xl border border-border/60 bg-card/95 p-4 text-left shadow-sm transition-all',
              'hover:-translate-y-0.5 hover:bg-muted/30 hover:shadow-md',
            )}
          >
            <div className="flex items-start justify-between gap-3">
              <div className="space-y-1.5">
                <div className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-muted">
                  <action.icon className="size-4 text-foreground" />
                </div>
                <p className="text-sm font-medium text-foreground">{action.label}</p>
                <p className="text-sm text-muted-foreground">{action.description}</p>
              </div>
              <ChevronRight className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </div>
          </Link>
        ))}
      </section>

      <section className="grid gap-4 lg:grid-cols-[0.95fr_1.05fr]">
        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60 pb-4">
            <CardTitle className="text-lg font-semibold tracking-tight">Resumo da edição</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 py-4">
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
                <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                  Período
                </p>
                <p className="mt-1.5 text-sm font-semibold text-foreground">
                  {formatDate(edition.starts_at)} - {formatDate(edition.ends_at)}
                </p>
              </div>
              <div className="rounded-2xl border border-border/60 bg-muted/10 px-3 py-3">
                <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">
                  Status
                </p>
                <p className="mt-1.5 text-sm font-semibold text-foreground">
                  {status.label}
                </p>
              </div>
            </div>

            <div className="grid gap-2">
              <div className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/10 px-3 py-2.5">
                <span className="text-xs font-medium text-muted-foreground">Local</span>
                <span className="text-sm font-medium text-foreground">{textOr(edition.location_name, 'Não definido')}</span>
              </div>
              <div className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/10 px-3 py-2.5">
                <span className="text-xs font-medium text-muted-foreground">Endereço</span>
                <span className="max-w-[65%] truncate text-sm font-medium text-foreground">
                  {textOr(edition.location_address, 'Não definido')}
                </span>
              </div>
              <div className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/10 px-3 py-2.5">
                <span className="text-xs font-medium text-muted-foreground">Contato</span>
                <span className="max-w-[65%] truncate text-sm font-medium text-foreground">
                  {hasText(edition.organizer_name)
                    ? edition.organizer_name
                    : hasText(edition.contact_email)
                      ? edition.contact_email
                      : 'Não definido'}
                </span>
              </div>
              <div className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/10 px-3 py-2.5">
                <span className="text-xs font-medium text-muted-foreground">Pagamento</span>
                <span className="max-w-[65%] truncate text-sm font-medium text-foreground">
                  {edition.trie_payments_credential_id ? 'Conectado' : 'Não conectado'}
                </span>
              </div>
              <div className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/10 px-3 py-2.5">
                <span className="text-xs font-medium text-muted-foreground">Banner/logo</span>
                <span className="max-w-[65%] truncate text-sm font-medium text-foreground">
                  {hasMedia ? 'Configurados' : 'Não configurados'}
                </span>
              </div>
            </div>

            <div className="flex flex-wrap gap-2 pt-1">
              <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/60 px-2.5 py-1 text-xs text-muted-foreground">
                <MapPin className="size-3.5" />
                {textOr(edition.timezone, 'Fuso não definido')}
              </span>
              <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/60 px-2.5 py-1 text-xs text-muted-foreground">
                <CreditCard className="size-3.5" />
                {editionMonetaryLabel[edition.monetary_type] ?? 'Sem categoria'}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Checklist da edição</CardTitle>
            <CardDescription>
              Itens derivados dos dados já cadastrados na edição.
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
                  <div className="space-y-0.5">
                    <p className="text-sm text-foreground">{item.label}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.done ? 'Preenchido' : 'Pendente'}
                    </p>
                  </div>
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
      </section>

      <Card className="border-border/60 bg-card/95 shadow-sm">
        <CardHeader className="border-b border-border/60">
          <CardTitle>Ações rápidas</CardTitle>
          <CardDescription>
            Atalhos para as operações mais comuns da edição.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-3 py-5 sm:grid-cols-2 xl:grid-cols-3">
          <Button
            type="button"
            variant="outline"
            className="h-auto justify-between rounded-2xl border-dashed bg-muted/15 px-4 py-4 text-left hover:bg-muted/30"
            onClick={() => setEditModalOpen(true)}
          >
            <span className="text-sm font-medium text-foreground">Editar edição</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </Button>

          {edition.status !== 'announced' && (
            <Button
              type="button"
              variant="secondary"
              className="h-auto justify-between rounded-2xl px-4 py-4 text-left"
              onClick={() => setPublishConfirmOpen(true)}
              disabled={publishEditionMutation.isPending}
            >
              <span className="text-sm font-medium text-foreground">Publicar edição</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </Button>
          )}

          {edition.status === 'announced' && (
            <Button
              type="button"
              variant="outline"
              className="h-auto justify-between rounded-2xl border-dashed bg-muted/15 px-4 py-4 text-left hover:bg-muted/30"
              onClick={() => void copyLink()}
            >
              <span className="text-sm font-medium text-foreground">Copiar link público</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </Button>
          )}


          <QuickAction to="/admin/events/$eventId/editions/$editionId/activities" params={{ eventId, editionId }}>
            <span className="text-sm font-medium text-foreground">Abrir atividades</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </QuickAction>

          <QuickAction to="/admin/events/$eventId/editions/$editionId/products" params={{ eventId, editionId }}>
            <span className="text-sm font-medium text-foreground">Abrir produtos</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </QuickAction>

          <QuickAction to="/admin/events/$eventId/editions/$editionId/checkpoints" params={{ eventId, editionId }}>
            <span className="text-sm font-medium text-foreground">Abrir checkpoints</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </QuickAction>

        </CardContent>
      </Card>

      <ManageEditionModal
        open={editModalOpen}
        onOpenChange={setEditModalOpen}
        edition={edition}
        onUpdate={async (id, values) => {
          const result = await updateEditionMutation.mutateAsync({ eventId, editionId: id, data: values })
          return result.success ? result.data : null
        }}
      />

      <AlertModal
        open={publishConfirmOpen}
        onOpenChange={setPublishConfirmOpen}
        title="Publicar edição?"
        description="Depois de publicar, a edição ficará disponível como anunciada."
        confirmLabel="Publicar edição"
        variant="default"
        loading={publishEditionMutation.isPending}
        onConfirm={async () => {
          handlePublish()
          setPublishConfirmOpen(false)
        }}
      />
    </div>
  )
}
