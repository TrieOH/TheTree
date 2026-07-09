import { createFileRoute, Link } from '@tanstack/react-router'
import { useMutation, useQuery } from '@tanstack/react-query'
import { FileText, Link2, Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { allAdminActivitiesQueryOptions } from '@/features/activities/api'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import {
  allCertificationTemplatesQueryOptions,
  setActivityCertificationTemplateFn,
  setEditionCertificationTemplateFn,
} from '@/features/certifications/api'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Label } from '@/shared/ui/shadcn/label'
import { cn } from '@/shared/lib/utils'
import { CertificationTemplatePreview } from '@/features/certifications/ui/CertificationTemplatePreview'

export const Route = createFileRoute('/admin/events/$eventId/editions/$editionId/certifications/')({
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const [selectedTemplateId, setSelectedTemplateId] = useState('')
  const [selectedActivityId, setSelectedActivityId] = useState('')

  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const edition = editions.find((item) => item.id === editionId) ?? null

  const { data: templates = [] } = useQuery(allCertificationTemplatesQueryOptions(eventId, editionId))
  const { data: activities = [] } = useQuery(allAdminActivitiesQueryOptions(eventId, editionId))

  const selectedTemplate = useMemo(
    () => templates.find((template) => template.id === selectedTemplateId) ?? templates[0] ?? null,
    [selectedTemplateId, templates]
  )

  const editionTemplateMutation = useMutation({
    mutationFn: () => {
      if (!selectedTemplate) throw new Error('Selecione um template')
      return setEditionCertificationTemplateFn(eventId, editionId, {
        certification_template_id: selectedTemplate.id,
      })
    },
    onSuccess: (res) => {
      if (res.success) toast.success('Template definido para a edição')
      else toast.error(res.message || 'Erro ao definir template da edição')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const activityTemplateMutation = useMutation({
    mutationFn: () => {
      if (!selectedTemplate || !selectedActivityId) throw new Error('Selecione template e atividade')
      return setActivityCertificationTemplateFn(eventId, editionId, selectedActivityId, {
        certification_template_id: selectedTemplate.id,
      })
    },
    onSuccess: (res) => {
      if (res.success) toast.success('Template definido para a atividade')
      else toast.error(res.message || 'Erro ao definir template da atividade')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const isPending = editionTemplateMutation.isPending || activityTemplateMutation.isPending
  const activityOptions = activities.filter((activity) => activity.status !== 'canceled')

  return (
    <div className="min-h-screen bg-background text-foreground">
      <div className="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-10">
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div className="space-y-1">
            <p className="text-xs font-medium uppercase tracking-[0.24em] text-muted-foreground">Admin</p>
            <h1 className="text-2xl font-semibold">Certificações</h1>
            <p className="text-sm text-muted-foreground">
              Templates e vínculos de certificados para {edition?.edition_name ?? 'esta edição'}.
            </p>
          </div>
          <Link
            to="/admin/events/$eventId/editions/$editionId/certifications/editor"
            params={{ eventId, editionId }}
            className="inline-flex h-10 w-full items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 md:w-auto"
          >
            <Plus className="size-4" />
            Novo template
          </Link>
        </div>

        <div className="mt-6 grid gap-6 lg:grid-cols-[1.35fr_1fr]">
          <section className="space-y-6">
            <Card>
              <CardHeader className="border-b pb-3">
                <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                  <FileText className="size-4 text-muted-foreground" />
                  Templates
                </CardTitle>
                <CardDescription className="text-xs">
                  Selecione um template para vincular ou abrir no editor.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3 p-4">
                {templates.length === 0 ? (
                  <div className="rounded-xl border border-dashed p-6 text-sm text-muted-foreground">
                    Nenhum template criado ainda.
                  </div>
                ) : (
                  templates.map((template) => {
                    const isSelected = template.id === selectedTemplate?.id
                    return (
                      <div
                        key={template.id}
                        onClick={() => setSelectedTemplateId(template.id)}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault()
                            setSelectedTemplateId(template.id)
                          }
                        }}
                        role="button"
                        tabIndex={0}
                        className={cn(
                          'flex w-full flex-col gap-3 rounded-2xl border p-4 text-left transition-colors md:flex-row md:items-center md:justify-between',
                          isSelected ? 'border-primary/40 bg-primary/5' : 'hover:bg-muted/40'
                        )}
                      >
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <p className="font-medium">{template.title}</p>
                            <Badge variant="outline">Template</Badge>
                            {isSelected && <Badge>Selecionado</Badge>}
                          </div>
                          <p className="text-xs text-muted-foreground">
                            {template.url ? 'Com fundo configurado' : 'Sem fundo'}
                          </p>
                        </div>
                        <div className="flex flex-wrap gap-2">
                          <CertificationTemplatePreview
                            eventId={eventId}
                            editionId={editionId}
                            template={template}
                            triggerLabel="Ver certificado"
                          />
                        </div>
                      </div>
                    )
                  })
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="border-b pb-3">
                <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                  <Link2 className="size-4 text-muted-foreground" />
                  Vínculo do template
                </CardTitle>
                <CardDescription className="text-xs">
                  Defina o template da edição ou sobrescreva para uma atividade.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4 p-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label className="text-xs text-muted-foreground">Template</Label>
                    <select
                      className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                      value={selectedTemplate?.id ?? ''}
                      onChange={(e) => setSelectedTemplateId(e.target.value)}
                    >
                      <option value="">Selecione um template</option>
                      {templates.map((template) => (
                        <option key={template.id} value={template.id}>
                          {template.title}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div className="space-y-2">
                    <Label className="text-xs text-muted-foreground">Atividade</Label>
                    <select
                      className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                      value={selectedActivityId}
                      onChange={(e) => setSelectedActivityId(e.target.value)}
                    >
                      <option value="">Selecione uma atividade</option>
                      {activityOptions.map((activity) => (
                        <option key={activity.id} value={activity.id}>
                          {activity.title}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row">
                  <Button
                    type="button"
                    className="flex-1"
                    disabled={!selectedTemplate || isPending}
                    onClick={() => { void editionTemplateMutation.mutateAsync() }}
                  >
                    Definir para a edição
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    className="flex-1"
                    disabled={!selectedTemplate || !selectedActivityId || isPending}
                    onClick={() => { void activityTemplateMutation.mutateAsync() }}
                  >
                    Definir para a atividade
                  </Button>
                </div>
              </CardContent>
            </Card>
          </section>
        </div>
      </div>
    </div>
  )
}
