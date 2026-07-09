import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@trieoh/identityx-sdk-ts/react'
import { CalendarDays } from 'lucide-react'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { UserCertificationsSection } from '@/features/certifications/ui/UserCertificationsSection'

export const Route = createFileRoute('/events/$eventId/profile')({
  beforeLoad: requireAuth,
  component: EventProfilePage,
})

function EventProfilePage() {
  const { auth } = useAuth()
  const profile = auth.profile()
  const { eventId } = Route.useParams()

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-gradient-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-6xl px-4 py-10 md:px-6 md:py-14">
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <CalendarDays className="size-4" />
              Evento
            </div>
            <div className="space-y-2">
              <h1 className="text-3xl font-semibold tracking-tight">Certificados do evento</h1>
              <p className="max-w-2xl text-sm text-muted-foreground">
                {profile?.email ? `Usuário ${profile.email}` : 'Certificados vinculados ao evento atual.'}
              </p>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-6xl px-4 py-6 md:px-6 md:py-8">
        <UserCertificationsSection
          userId={profile?.id ?? ''}
          eventId={eventId}
          onlyCurrentEvent
          title="Certificados do evento"
          subtitle="Mostra apenas certificados emitidos para edições e atividades deste evento."
        />
      </div>
    </main>
  )
}
