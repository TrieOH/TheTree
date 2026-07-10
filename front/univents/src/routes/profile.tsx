import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@trieoh/identityx-sdk-ts/react'
import { useTheme } from 'next-themes'
import { BadgeCheck, Monitor, MoonStar, PencilRuler, SunMedium } from 'lucide-react'
import { useEffect, useState } from 'react'
import { requireAuth } from '@/features/auths/lib/route-guard'
// import { UserCertificationsSection } from '@/features/certifications/ui/UserCertificationsSection'
import { Button } from '@/shared/ui/shadcn/button'
import {
  applyThemePreference,
  readInplaceEditPreference,
  readThemePreference,
  writeInplaceEditPreference,
  writeThemePreference,
} from '@/shared/lib/ui-preferences'
import { cn } from '@/shared/lib/utils'
import type { ThemeMode } from '@/shared/lib/ui-preferences'

export const Route = createFileRoute('/profile')({
  beforeLoad: requireAuth,
  component: ProfilePage,
})

function ProfilePage() {
  const { auth } = useAuth()
  const { setTheme } = useTheme()
  const profile = auth.profile()
  // const userId = profile?.id ?? ''
  const [themeMode, setThemeMode] = useState<ThemeMode>('auto')
  const [inplaceEditEnabled, setInplaceEditEnabled] = useState(true)

  useEffect(() => {
    const storedTheme = readThemePreference()
    setThemeMode(storedTheme)
    setInplaceEditEnabled(readInplaceEditPreference())
    setTheme(storedTheme === 'auto' ? 'system' : storedTheme)
    applyThemePreference(storedTheme)
  }, [])

  const handleThemeChange = (nextTheme: ThemeMode) => {
    setThemeMode(nextTheme)
    writeThemePreference(nextTheme)
    applyThemePreference(nextTheme)
    setTheme(nextTheme === 'auto' ? 'system' : nextTheme)
  }

  const handleInplaceEditChange = (enabled: boolean) => {
    setInplaceEditEnabled(enabled)
    writeInplaceEditPreference(enabled)
  }

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
        <section className="mb-6 rounded-2xl border border-border/60 bg-card p-5 md:p-6">
          <div className="mb-4 space-y-1">
            <h2 className="text-base font-semibold text-foreground">Preferências</h2>
            <p className="text-sm text-muted-foreground">
              Ajuste tema e edição embutida sem sair do perfil.
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-2">
            <div className="space-y-3">
              <div className="flex items-center gap-2 text-sm font-medium text-foreground">
                <Monitor className="size-4 text-muted-foreground" />
                Tema
              </div>

              <div className="grid grid-cols-3 gap-2">
                {[
                  { value: 'auto', label: 'Auto', icon: Monitor },
                  { value: 'light', label: 'Light', icon: SunMedium },
                  { value: 'dark', label: 'Dark', icon: MoonStar },
                ].map((option) => {
                  const Icon = option.icon
                  const active = themeMode === option.value

                  return (
                    <Button
                      key={option.value}
                      type="button"
                      variant={active ? 'default' : 'outline'}
                      className={cn(
                        'h-12 justify-start gap-2 rounded-xl px-4',
                        active ? 'shadow-sm' : 'bg-background/70',
                      )}
                      onClick={() => { handleThemeChange(option.value as ThemeMode) }}
                    >
                      <Icon className="size-4" />
                      <span>{option.label}</span>
                    </Button>
                  )
                })}
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex items-center gap-2 text-sm font-medium text-foreground">
                <PencilRuler className="size-4 text-muted-foreground" />
                Edição embutida
              </div>

              <div className="rounded-xl border border-border/60 bg-muted/20 p-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="space-y-1">
                    <p className="text-sm font-medium text-foreground">
                      Mostrar ações de edição nos cards
                    </p>
                    <p className="text-sm text-muted-foreground">
                      Quando desligado, o botão de engrenagem e os controles inline ficam ocultos.
                    </p>
                  </div>

                  <Button
                    type="button"
                    variant={inplaceEditEnabled ? 'default' : 'outline'}
                    className="rounded-full px-4"
                    onClick={() => { handleInplaceEditChange(!inplaceEditEnabled) }}
                  >
                    {inplaceEditEnabled ? 'Ativo' : 'Desligado'}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* <UserCertificationsSection
          userId={userId}
          title="Certificados da conta"
          subtitle="Todos os certificados emitidos para o seu usuário."
        /> */}
      </div>
    </main>
  )
}
