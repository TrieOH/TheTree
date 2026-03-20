import { createFileRoute, Link } from '@tanstack/react-router'
import { Calendar, MapPin, ChevronRight, Package } from 'lucide-react'
import { Button } from '@/shared/ui/shadcn/button'

export interface EditionI {
  id: string;
  event_id: string;
  goauth_scope_id: string;
  type: string;
  edition_name: string;
  tagline: string | null;
  description: string | null;
  status: "draft" | "announced" | "open" | "ongoing" | "finished" |
  "completed" | "cancelled" | "postponed";
  monetary_type: "free" | "paid" | "mixed";
  registration_opens_at: string | null;
  registration_closes_at: string | null;
  starts_at: string;
  ends_at: string;
  timezone: string;
  location_name: string;
  location_address: string;
  logo_url: string | null;
  banner_url: string | null;
  contact_email: string | null;
  contact_phone: string | null;
  organizer_name: string | null;
  trie_payments_credential_id: string | null;
  trie_payments_provider: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

export const Route = createFileRoute('/events/$eventId/editions/')({
  component: EventEditionsPage,
})

function EventEditionsPage() {
  const { eventId } = Route.useParams()

  const mockEditions: EditionI[] = [
    {
      id: '1',
      event_id: eventId,
      goauth_scope_id: 'scope-1',
      type: 'physical',
      edition_name: 'Summer Festival 2026',
      tagline: 'O maior evento do verão',
      description: 'Uma experiência única com os melhores produtos.',
      status: 'open',
      monetary_type: 'paid',
      registration_opens_at: '2026-01-01T10:00:00Z',
      registration_closes_at: '2026-05-19T23:59:59Z',
      starts_at: '2026-05-20T14:00:00Z',
      ends_at: '2026-05-22T22:00:00Z',
      timezone: 'America/Sao_Paulo',
      location_name: 'Arena Central',
      location_address: 'Av. das Nações, 1000',
      logo_url: null,
      banner_url: null,
      contact_email: 'contato@summerfest.com',
      contact_phone: '11999999999',
      organizer_name: 'Summer Corp',
      trie_payments_credential_id: null,
      trie_payments_provider: null,
      created_by: 'user-1',
      created_at: '2025-10-10T10:00:00Z',
      updated_at: '2025-10-10T10:00:00Z',
      deleted_at: null,
    },
    {
      id: '2',
      event_id: eventId,
      goauth_scope_id: 'scope-2',
      type: 'online',
      edition_name: 'Winter Edition 2026',
      tagline: 'Conectando o mundo',
      description: 'Acesse de qualquer lugar e garanta seus produtos digitais.',
      status: 'announced',
      monetary_type: 'mixed',
      registration_opens_at: '2026-06-01T10:00:00Z',
      registration_closes_at: '2026-07-14T23:59:59Z',
      starts_at: '2026-07-15T18:00:00Z',
      ends_at: '2026-07-15T21:00:00Z',
      timezone: 'America/Sao_Paulo',
      location_name: 'Plataforma Online',
      location_address: 'Link enviado por e-mail',
      logo_url: null,
      banner_url: null,
      contact_email: 'contato@winterfest.com',
      contact_phone: '11888888888',
      organizer_name: 'Winter Events',
      trie_payments_credential_id: null,
      trie_payments_provider: null,
      created_by: 'user-1',
      created_at: '2025-11-10T10:00:00Z',
      updated_at: '2025-11-10T10:00:00Z',
      deleted_at: null,
    }
  ]

  const statusMap: Record<string, { label: string, color: string }> = {
    open: { label: 'Inscrições Abertas', color: 'bg-green-500/10 text-green-600 border-green-200' },
    announced: { label: 'Anunciado', color: 'bg-blue-500/10 text-blue-600 border-blue-200' },
    ongoing: { label: 'Em andamento', color: 'bg-amber-500/10 text-amber-600 border-amber-200' },
    finished: { label: 'Encerrado', color: 'bg-muted text-muted-foreground border-border' },
  }

  return (
    <div className="max-w-[1200px] mx-auto px-6 py-12 md:py-20 space-y-12 md:space-y-16">
      <div className="space-y-4">
        <h1 className="text-3xl md:text-5xl font-bold tracking-tight">Edições disponíveis</h1>
        <p className="text-lg text-muted-foreground max-w-2xl">
          Selecione a edição do evento para explorar e adquirir produtos exclusivos.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12">
        {mockEditions.map((edition) => (
          <div
            key={edition.id}
            className="group relative flex flex-col bg-card border border-border rounded-[2.5rem] overflow-hidden hover:border-primary/30 transition-all duration-500 hover:shadow-2xl hover:shadow-primary/5"
          >
            <div className="aspect-[21/9] bg-muted relative overflow-hidden">
              <div className="absolute inset-0 flex items-center justify-center">
                <Package className="w-12 h-12 text-muted-foreground/20" />
              </div>
              <div className="absolute top-6 left-6">
                <span className={`px-4 py-1.5 rounded-full text-xs font-semibold border ${statusMap[edition.status]?.color || statusMap.finished.color}`}>
                  {statusMap[edition.status]?.label || edition.status}
                </span>
              </div>
            </div>

            <div className="p-8 md:p-10 flex-1 flex flex-col space-y-6">
              <div className="space-y-3">
                <h3 className="text-2xl md:text-3xl font-bold leading-tight group-hover:text-primary transition-colors">
                  {edition.edition_name}
                </h3>
                {edition.tagline && (
                  <p className="text-lg text-muted-foreground font-medium">{edition.tagline}</p>
                )}
              </div>

              <div className="grid grid-cols-1 gap-4 pt-2">
                <div className="flex items-center gap-3 text-muted-foreground">
                  <div className="w-10 h-10 rounded-xl bg-muted flex items-center justify-center shrink-0">
                    <Calendar className="w-5 h-5" />
                  </div>
                  <span className="text-base font-medium">
                    {new Date(edition.starts_at).toLocaleDateString('pt-BR', { day: '2-digit', month: 'long', year: 'numeric' })}
                  </span>
                </div>
                <div className="flex items-center gap-3 text-muted-foreground">
                  <div className="w-10 h-10 rounded-xl bg-muted flex items-center justify-center shrink-0">
                    <MapPin className="w-5 h-5" />
                  </div>
                  <span className="text-base font-medium line-clamp-1">{edition.location_name}</span>
                </div>
              </div>

              <div className="pt-6 mt-auto">
                <Link
                  to="/events/$eventId/editions/$editionId/products"
                  params={{ eventId, editionId: edition.id }}
                  className="block"
                >
                  <Button size="lg" className="w-full rounded-2xl h-14 text-lg font-semibold group-hover:scale-[1.02] transition-transform">
                    Explorar Produtos
                    <ChevronRight className="w-5 h-5 ml-2" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="pt-12 border-t border-border">
        <Link to="/events">
          <Button variant="ghost" size="lg" className="text-muted-foreground hover:text-foreground text-lg rounded-xl">
            ← Voltar para todos os eventos
          </Button>
        </Link>
      </div>
    </div>
  )
}
