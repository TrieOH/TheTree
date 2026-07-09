import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@trieoh/identityx-sdk-ts/react'
import { BadgeCheck } from 'lucide-react'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { UserCertificationsSection } from '@/features/certifications/ui/UserCertificationsSection'

export const Route = createFileRoute('/profile')({
  beforeLoad: requireAuth,
  component: ProfilePage,
})

function ProfilePage() {
  const { auth } = useAuth()
  const profile = auth.profile()
  const userId = profile?.id ?? ''

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-6xl px-4 py-10 md:px-6 md:py-14">
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <BadgeCheck className="size-4" />
              Perfil
            </div>
            <div className="space-y-2">
              <h1 className="text-3xl font-semibold tracking-tight">Meus certificados</h1>
              <p className="max-w-2xl text-sm text-muted-foreground">
                {profile?.email ? `Conta autenticada: ${profile.email}` : 'Central de certificados da sua conta.'}
              </p>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-6xl px-4 py-6 md:px-6 md:py-8">
        <UserCertificationsSection
          userId={userId}
          title="Certificados da conta"
          subtitle="Todos os certificados emitidos para o seu usuário."
        />
      </div>
    </main>
  )
}
