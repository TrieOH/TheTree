import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@trieoh/identityx-sdk-ts/react'
import { Layers3 } from 'lucide-react'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { UserCertificationsSection } from '@/features/certifications/ui/UserCertificationsSection'

export const Route = createFileRoute('/events/$eventId/editions/$editionId/profile')({
  beforeLoad: requireAuth,
  component: EditionProfilePage,
})

function EditionProfilePage() {
  const { auth } = useAuth()
  const profile = auth.profile()
  const { eventId, editionId } = Route.useParams()

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-6xl px-4 py-10 md:px-6 md:py-14">
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <Layers3 className="size-4" />
              Edição
            </div>
            <div className="space-y-2">
              <h1 className="text-3xl font-semibold tracking-tight">Certificados da edição</h1>
              <p className="max-w-2xl text-sm text-muted-foreground">
                {profile?.email ? `Usuário ${profile.email}` : 'Certificados vinculados à edição atual.'}
              </p>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-6xl px-4 py-6 md:px-6 md:py-8">
        <UserCertificationsSection
          userId={profile?.id ?? ''}
          eventId={eventId}
          editionId={editionId}
          onlyCurrentEvent
          title="Certificados da edição"
          subtitle="Mostra apenas certificados emitidos para esta edição e suas atividades."
        />
      </div>
    </main>
  )
}
