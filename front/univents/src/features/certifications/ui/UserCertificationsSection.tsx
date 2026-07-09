import { useMemo } from 'react'
import { useQueries, useQuery } from '@tanstack/react-query'
import { Award, CalendarDays, Layers3, Search, Sparkles } from 'lucide-react'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { certificationsByUserQueryOptions, certificationTemplateQueryOptions } from '@/features/certifications/api'
import { CertificationTemplatePreview } from '@/features/certifications/ui/CertificationTemplatePreview'
import type { CertificationTemplateI } from '@/features/certifications/model'
import { eventsQueryOptions } from '@/features/events/api'
import { allEditionsQueryOptions } from '@/features/editions/api'
import { allActivitiesQueryOptions } from '@/features/activities/api'
import type { ActivityI } from '@/features/activities/model'
import type { EditionI } from '@/features/editions/model'

type UserCertificationsSectionProps = {
  userId: string
  title?: string
  subtitle?: string
}

type RenderedCertification = {
  id: string
  target_type: 'edition' | 'activity'
  target_id: string
  certified_at: string
  hash: string
  event_id: string | null
  edition_id: string | null
  template: CertificationTemplateI | null
}

type TemplateQueryKey = {
  key: string
  eventId: string
  editionId: string
  templateId: string
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

function getOrigin() {
  if (typeof window !== 'undefined' && window.location?.origin)
    return window.location.origin

  return 'http://localhost:3002'
}

function substituteText(text: string, values: Record<string, string>) {
  return text.replace(/\{\{(\w+)\}\}/g, (_, key: string) => values[key] ?? `{{${key}}}`)
}

function applyTemplateVariables(
  template: CertificationTemplateI,
  values: Record<string, string>
): CertificationTemplateI {
  return {
    ...template,
    data: {
      ...template.data,
      elements: template.data.elements.map((element) => {
        if (element.type !== 'text') return element
        return {
          ...element,
          content: substituteText(element.content, values),
        }
      }),
    },
  }
}

export function UserCertificationsSection({
  userId,
  title = 'Certificados',
  subtitle = 'Veja os certificados emitidos para sua conta.',
}: UserCertificationsSectionProps) {
  const { data: certifications = [] } = useQuery(certificationsByUserQueryOptions(userId))
  const { data: events = [] } = useQuery(eventsQueryOptions())

  const editionQueries = useQueries({
    queries: events.flatMap((event) => ([
      {
        ...allEditionsQueryOptions(event.id),
        enabled: !!event.id,
      },
    ])),
  })

  const editionLookup = useMemo(() => {
    const editions = new Map<string, EditionI>()
    editionQueries.forEach((query) => {
      for (const edition of (query.data ?? []) as EditionI[]) {
        editions.set(edition.id, edition)
      }
    })
    return editions
  }, [editionQueries])

  const activityQueries = useQueries({
    queries: [...editionLookup.values()].flatMap((edition) => ([
      {
        ...allActivitiesQueryOptions(edition.event_id, edition.id),
        enabled: !!edition.event_id,
      },
    ])),
  })

  const activityLookup = useMemo(() => {
    const activities = new Map<string, ActivityI>()
    activityQueries.forEach((query) => {
      for (const activity of (query.data ?? []) as ActivityI[]) {
        activities.set(activity.id, activity)
      }
    })
    return activities
  }, [activityQueries])

  const templateQueryKeys = useMemo<TemplateQueryKey[]>(() => {
    const unique = new Map<string, TemplateQueryKey>()

    certifications.forEach((cert) => {
      const activity = cert.target_type === 'activity' ? activityLookup.get(cert.target_id) ?? null : null
      const editionId = cert.target_type === 'edition'
        ? cert.target_id
        : activity?.edition_id ?? null
      const edition = editionId ? editionLookup.get(editionId) ?? null : null
      const templateId = activity?.certification_template_id ?? edition?.certification_template_id ?? null

      if (!templateId || !edition?.event_id || !edition?.id) return

      const key = `${edition.event_id}:${edition.id}:${templateId}`
      if (!unique.has(key)) {
        unique.set(key, {
          key,
          eventId: edition.event_id,
          editionId: edition.id,
          templateId,
        })
      }
    })

    return [...unique.values()]
  }, [activityLookup, certifications, editionLookup])

  const templateQueries = useQueries({
    queries: templateQueryKeys.map((item) => ({
      ...certificationTemplateQueryOptions(item.eventId, item.editionId, item.templateId),
      enabled: true,
    })),
  })

  const templateByKey = useMemo(() => {
    const map = new Map<string, CertificationTemplateI>()
    templateQueries.forEach((query, index) => {
      const item = templateQueryKeys[index]
      if (item && query.data) map.set(item.key, query.data)
    })
    return map
  }, [templateQueries, templateQueryKeys])

  const finalCerts = useMemo<RenderedCertification[]>(() => {
    return certifications.map((cert) => {
      const activity = cert.target_type === 'activity' ? activityLookup.get(cert.target_id) ?? null : null
      const editionId = cert.target_type === 'edition'
        ? cert.target_id
        : activity?.edition_id ?? null
      const edition = editionId ? editionLookup.get(editionId) ?? null : null
      const templateId = activity?.certification_template_id ?? edition?.certification_template_id ?? null
      const templateKey =
        templateId && edition?.event_id && edition?.id
          ? `${edition.event_id}:${edition.id}:${templateId}`
          : null
      const template = templateKey ? templateByKey.get(templateKey) ?? null : null

      const filledTemplate = template
        ? applyTemplateVariables(template, {
          activity_name:
            cert.target_type === 'edition'
              ? edition?.edition_name ?? ''
              : activity?.title ?? edition?.edition_name ?? '',
          certified_at: formatCertifiedAt(cert.certified_at),
          cert_hash: cert.hash,
          verify_url: `${getOrigin()}/verify/${cert.hash}`,
        })
        : null

      return {
        ...cert,
        event_id: edition?.event_id ?? null,
        edition_id: edition?.id ?? null,
        template: filledTemplate,
      }
    })
  }, [activityLookup, certifications, editionLookup, templateByKey])

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
          {finalCerts.length > 0 ? (
            <div className="grid gap-3 grid-cols-[repeat(auto-fit,minmax(280px,1fr))]">
              {finalCerts.map((cert) => {
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
                          {cert.hash}
                        </Badge>
                      </div>

                      <div>
                        <p className="text-sm font-semibold text-foreground">Certificado emitido</p>
                        <p className="text-xs text-muted-foreground">{formatCertifiedAt(cert.certified_at)}</p>
                      </div>

                      <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                        <span className="inline-flex items-center gap-1">
                          <CalendarDays className="size-3.5" />
                          {cert.target_type === 'edition' ? 'Vinculado à edição' : 'Vinculado à atividade'}
                        </span>
                      </div>

                      {cert.template && cert.event_id && cert.edition_id ? (
                        <CertificationTemplatePreview
                          eventId={cert.event_id}
                          editionId={cert.edition_id}
                          template={cert.template}
                          triggerLabel="Abrir certificado"
                        />
                      ) : (
                        <Button type="button" variant="outline" size="sm" disabled className="w-fit">
                          Carregando certificado
                        </Button>
                      )}
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
