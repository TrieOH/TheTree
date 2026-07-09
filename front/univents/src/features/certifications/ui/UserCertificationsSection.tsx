import { useMemo } from 'react'
import { useQueries, useQuery } from '@tanstack/react-query'
import { Award, CalendarDays, Layers3, Search, Sparkles } from 'lucide-react'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { certificationsByUserQueryOptions } from '@/features/certifications/api'
import { allEditionsQueryOptions } from '@/features/editions/api'
import { allActivitiesQueryOptions } from '@/features/activities/api'

type UserCertificationsSectionProps = {
  userId: string
  title?: string
  subtitle?: string
  eventId?: string
  editionId?: string
  onlyCurrentEvent?: boolean
}

function formatCertifiedAt(value: string) {
  return new Date(value).toLocaleString('pt-BR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function UserCertificationsSection({
  userId,
  title = 'Certificados',
  subtitle = 'Veja os certificados emitidos para sua conta.',
  eventId,
  editionId,
  onlyCurrentEvent = false,
}: UserCertificationsSectionProps) {
  const { data: certifications = [] } = useQuery(certificationsByUserQueryOptions(userId))
  const { data: editions = [] } = useQuery({
    ...allEditionsQueryOptions(eventId ?? ''),
    enabled: !!eventId,
  })

  const activitiesQueries = useQueries({
    queries: (eventId ? editions : []).map((edition) => ({
      ...allActivitiesQueryOptions(eventId!, edition.id),
      enabled: !!eventId,
    })),
  })

  const filteredCertifications = useMemo(() => {
    if (!eventId || (!editionId && !onlyCurrentEvent)) {
      return certifications
    }

    const editionSet = new Set<string>(editionId ? [editionId] : editions.map((edition) => edition.id))
    const activityEditionMap = new Map<string, string>()
    activitiesQueries.forEach((query) => {
      for (const activity of (query.data ?? []) as Array<{ id: string; edition_id: string }>) {
        activityEditionMap.set(activity.id, activity.edition_id)
      }
    })

    return certifications.filter((cert) => {
      if (cert.target_type === 'edition') return editionSet.has(cert.target_id)
      if (cert.target_type === 'activity') {
        const activityEditionId = activityEditionMap.get(cert.target_id)
        return activityEditionId ? editionSet.has(activityEditionId) : false
      }
      return true
    })
  }, [activitiesQueries, certifications, editionId, editions, eventId, onlyCurrentEvent])

  return (
    <section className="space-y-4">
      <Card className="overflow-hidden border-border/60 bg-background/80 backdrop-blur">
        <CardHeader className="border-b border-border/60 bg-linear-to-r from-background via-background to-muted/20">
          <CardTitle className="flex items-center gap-2 text-base">
            <Award className="size-4 text-primary" />
            {title}
          </CardTitle>
          <CardDescription>{subtitle}</CardDescription>
        </CardHeader>

        <CardContent className="p-4">
          {filteredCertifications.length > 0 ? (
            <div className="grid gap-3 grid-cols-[repeat(auto-fit,minmax(280px,1fr))]">
              {filteredCertifications.map((cert) => {
                const targetLabel = cert.target_type === 'edition' ? 'Edição' : 'Atividade'

                return (
                  <div
                    key={cert.id}
                    className="flex h-full flex-col gap-4 rounded-2xl border border-border/70 bg-card p-4 shadow-sm transition-colors hover:border-border hover:bg-muted/20"
                  >
                    <div className="space-y-2">
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge variant="outline" className="gap-1.5">
                          {cert.target_type === 'edition' ? <Layers3 className="size-3.5" /> : <Sparkles className="size-3.5" />}
                          {targetLabel}
                        </Badge>
                        <Badge variant="secondary" className="font-mono text-[10px] uppercase tracking-wider">
                          {cert.target_id.slice(0, 8)}
                        </Badge>
                      </div>

                      <div>
                        <p className="text-sm font-semibold text-foreground">
                          Certificado emitido
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatCertifiedAt(cert.certified_at)}
                        </p>
                      </div>

                      <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                        <span className="inline-flex items-center gap-1">
                          <CalendarDays className="size-3.5" />
                          {cert.target_type === 'edition' ? 'Vinculado à edição' : 'Vinculado à atividade'}
                        </span>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-border/60 bg-muted/20 px-6 py-12 text-center">
              <div className="flex size-12 items-center justify-center rounded-full bg-background shadow-sm">
                <Search className="size-5 text-muted-foreground" />
              </div>
              <p className="mt-4 text-sm font-medium text-foreground">Nenhum certificado encontrado</p>
              <p className="mt-1 max-w-sm text-sm text-muted-foreground">
                Quando você receber um certificado, ele vai aparecer aqui.
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </section>
  )
}
