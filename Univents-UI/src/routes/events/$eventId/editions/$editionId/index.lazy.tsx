import { createLazyFileRoute, Link } from '@tanstack/react-router'
import {
  ArrowRight,
  CalendarDays,
  ChevronRight,
  Clock,
  Mail,
  MapPin,
  Phone,
  Share2,
  Ticket,
  User,
} from 'lucide-react'
import { toast } from 'sonner'
import type { EditionI } from '@/features/editions/model'
import { formatDateRange } from '@/shared/lib/date'
import { cn } from '@/shared/lib/utils'
import EditionInfoCard from '@/features/editions/ui/EditionInfoCard'
import { MultiLocationMap } from '@/widgets/ui/map-embed'

const statusConfig: Record<
  EditionI['status'],
  { label: string; dot: string }
> = {
  ongoing: { label: 'Em andamento', dot: 'bg-emerald-500' },
  open: { label: 'Inscrições abertas', dot: 'bg-blue-500' },
  announced: { label: 'Anunciado', dot: 'bg-amber-500' },
  draft: { label: 'Rascunho', dot: 'bg-violet-500' },
  finished: { label: 'Finalizado', dot: 'bg-slate-400' },
  completed: { label: 'Concluído', dot: 'bg-slate-400' },
  cancelled: { label: 'Cancelado', dot: 'bg-red-500' },
  postponed: { label: 'Adiado', dot: 'bg-orange-500' },
}

const monetaryConfig: Record<
  EditionI['monetary_type'],
  { label: string; cls: string }
> = {
  free: { label: 'Gratuito', cls: 'bg-emerald-500/10 text-emerald-700 border-emerald-500/20' },
  paid: { label: 'Pago', cls: 'bg-amber-500/10 text-amber-700 border-amber-500/20' },
  mixed: { label: 'Misto', cls: 'bg-violet-500/10 text-violet-700 border-violet-500/20' },
}

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export const Route = createLazyFileRoute(
  '/events/$eventId/editions/$editionId/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const edition = Route.useLoaderData()

  if (!edition) {
    return (
      <div className="flex h-[80vh] w-full flex-col items-center justify-center gap-4">
        <p className="text-muted-foreground">Erro ao carregar edição</p>
        <Link
          to="/events/$eventId/editions"
          params={{ eventId }}
          className="text-primary hover:underline"
        >
          Voltar para edições
        </Link>
      </div>
    )
  }

  const status = statusConfig[edition.status]
  const monetary = monetaryConfig[edition.monetary_type]

  const handleShare = async () => {
    const url = window.location.href
    try {
      if (typeof navigator.share === 'function') {
        await navigator.share({ title: edition.edition_name, url })
        return
      }
      await navigator.clipboard.writeText(url)
      toast.success('Link copiado!')
    } catch {
      toast.error('Erro ao compartilhar')
    }
  }

  return (
    <div className="min-h-screen bg-background pb-32">

      {/* Hero banner*/}
      <div className="relative h-[40vh] sm:h-[52vh] min-h-75 max-h-150 overflow-hidden bg-muted">
        {edition.banner_url ? (
          <img
            src={edition.banner_url}
            alt={edition.edition_name}
            className="absolute inset-0 w-full h-full object-cover"
          />
        ) : (
          <div className="absolute inset-0 bg-primary/20" />
        )}
        <div className="absolute inset-x-0 top-0 h-32 bg-linear-to-b from-black/60 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-2/3 bg-linear-to-t from-black/80 via-black/20 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-background to-transparent" />

        {/* header actions */}
        <div className="absolute top-4 sm:top-5 inset-x-4 flex justify-between items-center">
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] sm:text-[11px] font-bold tracking-wider border bg-black/20 backdrop-blur-md border-white/20 text-white">
            <span
              className={cn(
                'w-1.5 h-1.5 rounded-full shrink-0',
                status.dot,
                edition.status === 'ongoing' && 'animate-pulse',
              )}
            />
            {status.label.toUpperCase()}
          </span>
          <button
            onClick={() => void handleShare()}
            className="w-8 h-8 sm:w-9 sm:h-9 rounded-full bg-black/20 backdrop-blur-md border border-white/20 flex items-center justify-center text-white active:scale-90 transition-transform"
          >
            <Share2 className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          </button>
        </div>

        {/* title */}
        <div className="absolute bottom-6 left-0 right-0 px-4 sm:px-6">
          <p className="text-[9px] sm:text-[10px] font-bold tracking-[0.2em] uppercase text-white/60 mb-1.5 sm:mb-2">
            {edition.type === 'year'
              ? `Edição ${new Date(edition.starts_at).getFullYear()}`
              : edition.edition_name}
          </p>
          <h1 className="text-3xl sm:text-4xl md:text-5xl font-black tracking-tight text-white leading-[1.1] mb-2 sm:mb-3">
            {edition.edition_name}
          </h1>
          {edition.tagline && (
            <p className="text-xs sm:text-sm font-medium text-white/80 leading-relaxed max-w-[95%] sm:max-w-[85%]">
              {edition.tagline}
            </p>
          )}
        </div>
      </div>

      {/* Activities button */}
      <div className="px-4 -mt-6 relative z-10 max-w-2xl mx-auto">
        <Link
          to="/events/$eventId/editions/$editionId"
          params={{ eventId, editionId }}
          className={cn(
            'flex items-center justify-between w-full',
            'bg-primary text-primary-foreground border border-primary/20',
            'rounded-2xl px-4 py-4 sm:px-5',
            'shadow-xl shadow-primary/20',
            'transition-all active:scale-[0.98] hover:brightness-110',
          )}
        >
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-primary-foreground/15 flex items-center justify-center shrink-0">
              <CalendarDays className="w-4.5 h-4.5" />
            </div>
            <div>
              <p className="text-[9px] font-bold tracking-widest uppercase text-primary-foreground/70">
                Programação
              </p>
              <p className="text-base font-bold tracking-tight">Ver Atividades</p>
            </div>
          </div>
          <div className="w-7 h-7 rounded-full bg-primary-foreground/10 flex items-center justify-center shrink-0">
            <ArrowRight className="w-3.5 h-3.5" />
          </div>
        </Link>
      </div>

      {/* Body */}
      <div className="px-4 pt-5 space-y-3 max-w-2xl mx-auto">

        {/* Date */}
        <EditionInfoCard
          icon={<CalendarDays className="w-4 h-4" />}
          label="Data"
          value={formatDateRange(edition.starts_at, edition.ends_at)}
          sub={`${formatTime(edition.starts_at)} – ${formatTime(edition.ends_at)} · ${edition.timezone}`}
          iconClass="bg-blue-500/10 text-blue-600"
        />

        {/* Location */}
        <EditionInfoCard
          icon={<MapPin className="w-4 h-4" />}
          label="Local"
          value={edition.location_name}
          sub={edition.location_address}
          iconClass="bg-amber-500/10 text-amber-600"
          footer={
            <MultiLocationMap
              locations={[
                { name: edition.location_name, address: edition.location_address },
              ]}
              height="300px"
              className='z-10'
            />
          }
        />

        {/* Registration */}
        <div className="bg-card border border-border rounded-2xl p-4">
          <p className="text-[10px] font-semibold tracking-widest uppercase text-muted-foreground mb-3">
            Inscrições
          </p>

          <div className="flex items-center justify-between gap-3 flex-wrap">
            {edition.registration_opens_at && edition.registration_closes_at ? (
              <div className="flex items-center gap-1.5">
                <Clock className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                <p className="text-sm font-medium text-foreground/80">
                  {new Date(edition.registration_opens_at).toLocaleDateString('pt-BR', {
                    day: '2-digit',
                    month: 'short',
                  })}
                  {' – '}
                  {new Date(edition.registration_closes_at).toLocaleDateString('pt-BR', {
                    day: '2-digit',
                    month: 'short',
                    year: 'numeric',
                  })}
                </p>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Período não definido</p>
            )}

            <span
              className={cn(
                'inline-flex items-center gap-1.5 text-[11px] font-semibold',
                'px-2.5 py-1 rounded-full border shrink-0',
                monetary.cls,
              )}
            >
              <Ticket className="w-3 h-3" />
              {monetary.label}
            </span>
          </div>

          {edition.organizer_name && (
            <div className="flex items-center gap-1.5 mt-3 pt-3 border-t border-border">
              <User className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
              <p className="text-xs text-muted-foreground">
                Organizado por{' '}
                <span className="font-medium text-foreground/80">
                  {edition.organizer_name}
                </span>
              </p>
            </div>
          )}
        </div>

        {/* Description */}
        {edition.description && (
          <div className="bg-card border border-border rounded-2xl p-4 sm:p-5">
            <p className="text-[10px] font-semibold tracking-widest uppercase text-muted-foreground mb-3">
              Sobre a edição
            </p>
            <p className="text-sm text-foreground/70 leading-relaxed whitespace-pre-wrap">
              {edition.description}
            </p>
          </div>
        )}

        {/* Contact */}
        {(edition.contact_email ?? edition.contact_phone) && (
          <div className="bg-card border border-border rounded-2xl p-4 sm:p-5">
            <p className="text-[10px] font-semibold tracking-widest uppercase text-muted-foreground mb-1">
              Contato
            </p>

            {edition.contact_email && (
              <a
                href={`mailto:${edition.contact_email}`}
                className="flex items-center gap-3 py-3 group"
              >
                <div className="w-8 h-8 rounded-lg bg-muted border border-border/60 flex items-center justify-center shrink-0 group-hover:bg-background transition-colors">
                  <Mail className="w-3.5 h-3.5 text-muted-foreground" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground mb-0.5">
                    E-mail
                  </p>
                  <p className="text-sm text-foreground/80 group-hover:text-foreground transition-colors truncate">
                    {edition.contact_email}
                  </p>
                </div>
                <ChevronRight className="w-4 h-4 text-muted-foreground/40 shrink-0" />
              </a>
            )}

            {edition.contact_email && edition.contact_phone && (
              <div className="border-t border-border" />
            )}

            {edition.contact_phone && (
              <a
                href={`tel:${edition.contact_phone}`}
                className="flex items-center gap-3 py-3 group"
              >
                <div className="w-8 h-8 rounded-lg bg-muted border border-border/60 flex items-center justify-center shrink-0 group-hover:bg-background transition-colors">
                  <Phone className="w-3.5 h-3.5 text-muted-foreground" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground mb-0.5">
                    Telefone
                  </p>
                  <p className="text-sm text-foreground/80 group-hover:text-foreground transition-colors">
                    {edition.contact_phone}
                  </p>
                </div>
                <ChevronRight className="w-4 h-4 text-muted-foreground/40 shrink-0" />
              </a>
            )}
          </div>
        )}
      </div>
    </div>
  )
}