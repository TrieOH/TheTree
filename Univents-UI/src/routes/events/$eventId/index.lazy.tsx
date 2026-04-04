import { createLazyFileRoute, Link } from '@tanstack/react-router'
import {
  PencilLine, Share2, Save, X, ArrowRight,
  Calendar, Loader2, Mail, Link2, Globe, X as XIcon, Camera
} from 'lucide-react'
import { useState, useEffect, useMemo } from 'react'
import { toast } from 'sonner'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '@soramux/node-auth-sdk/react'
import type { EventI } from '@/features/events/model'
import { Button } from '@/shared/ui/shadcn/button'
import { Badge } from '@/shared/ui/shadcn/badge'
import { cn } from '@/shared/lib/utils'
import { parseDatetimeLocal } from '@/shared/lib/date'
import InlineEdit from '@/shared/ui/inline-edit/InlineEdit'
import InlineImageEdit from '@/shared/ui/inline-edit/InlineImageEdit'
import { eventQueryOptions, ownEventQueryOptions, patchEventFn } from '@/features/events/api'
import { uploadAndModerateFile } from '@/features/storage/api'
import { InfoRow, SectionCard, SocialChip } from '@/features/events/ui/EventDetailComponents'
import { usePermissions } from '@/features/auths/hooks/use-permissions'
import { canEditEvent } from '@/features/events/model/permissions'
import WaveSpinnerLoading from '@/shared/ui/loader/WaveSpinnerLoading'

export const Route = createLazyFileRoute('/events/$eventId/')({
  component: RouteComponent,
})

const statusColors = {
  active: 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20',
  inactive: 'bg-gray-500/10 text-gray-600 border-gray-500/20',
  draft: 'bg-amber-500/10 text-amber-600 border-amber-500/20',
  archived: 'bg-red-500/10 text-red-600 border-red-500/20',
  discontinued: 'bg-zinc-500/10 text-zinc-600 border-zinc-500/20',
}

const statusLabels = {
  active: 'Ativo',
  inactive: 'Inativo',
  draft: 'Rascunho',
  archived: 'Arquivado',
  discontinued: 'Descontinuado',
}

function RouteComponent() {
  const { edit } = Route.useSearch()
  const navigate = Route.useNavigate()
  const { eventId } = Route.useParams()
  const queryClient = useQueryClient()
  const { auth } = useAuth()
  const userProfile = auth.profile()

  const { data: event, isLoading, error } = useQuery(
    edit ? ownEventQueryOptions(eventId) : eventQueryOptions(eventId)
  )

  const { canEditEvent: hasEditPermission } = usePermissions(
    { canEditEvent },
    userProfile?.id
  )

  const [eventData, setEventData] = useState<EventI | null>(null)

  useEffect(() => {
    if (event) setEventData(event)
  }, [event])

  const [editingField, setEditingField] = useState<string | null>(null)

  const isDirty = useMemo(() => {
    if (!event || !eventData) return false
    return JSON.stringify(event) !== JSON.stringify(eventData)
  }, [event, eventData])

  const mutation = useMutation({
    mutationFn: (data: Partial<EventI>) => patchEventFn(eventId, data),
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData(ownEventQueryOptions(eventId).queryKey, res.data)
        toast.success('Alterações salvas!')
        void navigate({ search: (prev) => ({ ...prev, edit: false }) })
      } else toast.error('Erro ao salvar alterações')
    },
    onError: () => toast.error('Erro ao salvar alterações'),
  })

  const toggleEditMode = () => {
    void navigate({ search: (prev) => ({ ...prev, edit: !edit }) })
    setEditingField(null)
  }

  const updateField = (field: keyof EventI, value: string) => {
    if (!eventData) return
    setEventData((prev) => prev ? { ...prev, [field]: value } : null)
  }

  const handleSave = () => {
    if (!eventData) return
    const {
      name, acronym, tagline, description, logo_url, banner_url,
      slug, is_series, contact_email, gallery_urls
    } = eventData
    mutation.mutate({
      name, acronym, tagline, description, logo_url, banner_url,
      slug, is_series, contact_email, gallery_urls
    })
  }

  const handleShare = async () => {
    if (!eventData) return
    const url = window.location.href
    try {
      if (typeof navigator.share === 'function') {
        await navigator.share({ title: eventData.name, url })
        return
      }
      await navigator.clipboard.writeText(url)
      toast.success('Link copiado!')
    } catch {
      toast.error('Erro ao compartilhar')
    }
  }

  if (isLoading) {
    return (
      <div className="flex h-[80vh] w-full items-center justify-center">
        <WaveSpinnerLoading text='Carregando...' />
      </div>
    )
  }

  if (error || !eventData) {
    return (
      <div className="flex h-[80vh] w-full flex-col items-center justify-center gap-4">
        <p className="text-muted-foreground">Erro ao carregar evento</p>
        <Link to="/events" className="text-primary hover:underline">Voltar para eventos</Link>
      </div>
    )
  }

  const getInitials = (name: string) =>
    name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()

  const logoFallbackText = eventData.acronym
    ?? (eventData.name ? getInitials(eventData.name) : '')

  const logoFontSize =
    logoFallbackText.length <= 2 ? 'text-lg' :
      logoFallbackText.length === 3 ? 'text-sm' : 'text-xs'

  const hasSocialLinks = !!(
    eventData.social_links?.website ??
    eventData.social_links?.twitter ??
    eventData.social_links?.instagram ??
    eventData.social_links?.linkedin
  )

  return (
    <div className="min-h-screen bg-background pb-24">
      <div className="relative">
        <InlineImageEdit
          value={eventData.banner_url}
          onChange={(url) => { updateField('banner_url', url); }}
          isEditEnabled={edit}
          isEditing={editingField === 'banner_url'}
          onStartEdit={() => { setEditingField('banner_url'); }}
          onFinishEdit={() => { setEditingField(null); }}
          onUpload={(file) => uploadAndModerateFile(file, `events/${eventId}`)}
          renderDisplay={(url) => (
            <div className="relative w-full h-40 sm:h-52 md:h-64">
              {url
                ? <img src={url} alt={eventData.name} className="w-full h-full object-cover" />
                : <div className="w-full h-full bg-linear-to-br from-primary via-accent to-secondary" />
              }
              <div className="absolute inset-0 bg-linear-to-b from-transparent to-background/50" />
            </div>
          )}
        />
      </div>

      <div className="mx-auto max-w-2xl px-4">
        <main className="space-y-2">
          <div className="bg-card rounded-xl shadow-xl border">
            <div className="px-4 pt-0 flex items-end justify-between -mt-8 sm:-mt-10">
              <InlineImageEdit
                value={eventData.logo_url}
                onChange={(url) => { updateField('logo_url', url); }}
                isEditEnabled={edit}
                isEditing={editingField === 'logo_url'}
                onStartEdit={() => { setEditingField('logo_url'); }}
                onFinishEdit={() => { setEditingField(null); }}
                className="h-16 w-16 sm:h-20 sm:w-20 shrink-0"
                onUpload={(file) => uploadAndModerateFile(file, `events/${eventId}`)}
                renderDisplay={(url) => (
                  <div className="h-full w-full rounded-xl bg-primary shadow-lg flex items-center justify-center overflow-hidden ring-4 ring-card">
                    {url
                      ? <img src={url} alt={eventData.acronym ?? eventData.name} className="h-full w-full object-cover" />
                      : <span className={cn('text-primary-foreground font-bold leading-none tracking-tight px-1 text-center break-all', logoFontSize)}>
                        {logoFallbackText}
                      </span>
                    }
                  </div>
                )}
              />

              <div className="flex flex-col items-end gap-1.5 pb-1">
                <div className="flex items-center gap-2">
                  {!edit ? (
                    <div className="flex items-center gap-1 mr-1">
                      <Button
                        className="h-8 w-8 text-muted-foreground hover:text-foreground transition-colors"
                        variant="ghost" size="icon"
                        onClick={() => void handleShare()}
                        title="Compartilhar"
                      >
                        <Share2 className="h-4 w-4" />
                      </Button>
                      {hasEditPermission && (
                        <Button
                          className="h-8 w-8 text-muted-foreground hover:text-foreground transition-colors"
                          variant="ghost" size="icon"
                          onClick={toggleEditMode}
                          title="Editar evento"
                        >
                          <PencilLine className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  ) : (
                    <div className="flex items-center gap-1 mr-1">
                      {isDirty && (
                        <Button
                          className="h-8 w-8 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-500/10 transition-colors"
                          variant="ghost" size="icon"
                          onClick={handleSave}
                          disabled={mutation.isPending}
                          title="Salvar alterações"
                        >
                          {mutation.isPending ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                          ) : (
                            <Save className="h-4 w-4" />
                          )}
                        </Button>
                      )}
                      <Button
                        className="h-8 w-8 text-muted-foreground hover:text-destructive transition-colors"
                        variant="ghost" size="icon"
                        onClick={toggleEditMode}
                        title="Cancelar edição"
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    </div>
                  )}
                  <Badge
                    variant="outline"
                    className={cn('text-xs font-medium whitespace-nowrap', statusColors[eventData.status])}
                  >
                    {statusLabels[eventData.status]}
                  </Badge>
                </div>
                {eventData.is_series && eventData.editions_count > 0 && (
                  <span className="text-[11px] text-muted-foreground uppercase tracking-wider font-medium">
                    {eventData.editions_count}ª edição
                  </span>
                )}
              </div>
            </div>

            <div className="px-4 pt-2 pb-4 space-y-1">
              <InlineEdit
                value={eventData.name}
                onChange={(val) => { updateField('name', val); }}
                isEditEnabled={edit}
                isEditing={editingField === 'name'}
                onStartEdit={() => { setEditingField('name'); }}
                onFinishEdit={() => { setEditingField(null); }}
                className="text-xl sm:text-2xl font-bold tracking-tight text-foreground leading-tight"
              />

              {(eventData.acronym ?? eventData.tagline ?? edit) && (
                <div className="flex flex-wrap items-center gap-x-2 gap-y-0.5 text-sm">
                  {(eventData.acronym ?? edit) && (
                    <InlineEdit
                      value={eventData.acronym ?? ''}
                      onChange={(val) => { updateField('acronym', val); }}
                      isEditEnabled={edit}
                      isEditing={editingField === 'acronym'}
                      onStartEdit={() => { setEditingField('acronym'); }}
                      onFinishEdit={() => { setEditingField(null); }}
                      className="font-semibold text-primary"
                      placeholder="SIGLA"
                    />
                  )}
                  {eventData.acronym && (eventData.tagline ?? edit) && (
                    <span className="text-muted-foreground/40 select-none">·</span>
                  )}
                  {(eventData.tagline ?? edit) && (
                    <InlineEdit
                      value={eventData.tagline ?? ''}
                      onChange={(val) => { updateField('tagline', val); }}
                      isEditEnabled={edit}
                      isEditing={editingField === 'tagline'}
                      onStartEdit={() => { setEditingField('tagline'); }}
                      onFinishEdit={() => { setEditingField(null); }}
                      className="text-muted-foreground italic"
                      placeholder="Tagline do evento..."
                    />
                  )}
                </div>
              )}

              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <Calendar className="h-3.5 w-3.5 shrink-0" />
                <span>{parseDatetimeLocal(eventData.created_at).toLocaleDateString('pt-BR')}</span>
              </div>

              <div className="pt-2">
                <Link
                  to="/events/$eventId/editions"
                  params={{ eventId: eventData.id }}
                  className={cn(
                    'flex items-center justify-center gap-2',
                    'w-full bg-primary hover:bg-primary/90 active:scale-[.98]',
                    'text-primary-foreground rounded-xl h-10',
                    'text-sm font-semibold transition-all duration-150',
                  )}
                >
                  Ver Edições
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </div>
          </div>

          {eventData.description && (
            <SectionCard label="Sobre">
              <p className="text-sm text-foreground/75 leading-relaxed">
                {eventData.description}
              </p>
            </SectionCard>
          )}

          <SectionCard label="Informações">
            <div className="grid grid-cols-2 gap-x-4 gap-y-3">
              <InfoRow label="Slug" value={`/${eventData.slug}`} mono />
              {eventData.is_series && (
                <InfoRow label="Edições" value={String(eventData.editions_count)} />
              )}
            </div>
          </SectionCard>

          {eventData.gallery_urls.length > 0 && (
            <SectionCard label="Galeria">
              <div className="grid grid-cols-3 gap-1.5">
                {eventData.gallery_urls.map((url, i) => (
                  <a
                    key={i}
                    href={url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="aspect-square rounded-lg overflow-hidden bg-muted block"
                  >
                    <img
                      src={url}
                      alt={`Galeria ${i + 1}`}
                      className="h-full w-full object-cover hover:scale-105 transition-transform duration-200"
                    />
                  </a>
                ))}
              </div>
            </SectionCard>
          )}

          {(eventData.contact_email ?? hasSocialLinks) && (
            <SectionCard label="Contato">
              <div className="flex flex-wrap gap-2">
                {eventData.contact_email && (
                  <a
                    href={`mailto:${eventData.contact_email}`}
                    className={cn(
                      'flex items-center gap-2 px-3 py-2 rounded-lg',
                      'bg-background border border-border',
                      'text-sm text-foreground/80 hover:text-foreground',
                      'hover:bg-muted transition-colors min-w-0',
                    )}
                  >
                    <Mail className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                    <span className="truncate max-w-50">{eventData.contact_email}</span>
                  </a>
                )}
                {eventData.social_links?.website && (
                  <SocialChip href={eventData.social_links.website} label="Website" icon={<Globe className="h-3.5 w-3.5" />} />
                )}
                {eventData.social_links?.twitter && (
                  <SocialChip href={eventData.social_links.twitter} label="Twitter" icon={<XIcon className="h-3.5 w-3.5" />} />
                )}
                {eventData.social_links?.instagram && (
                  <SocialChip href={eventData.social_links.instagram} label="Instagram" icon={<Camera className="h-3.5 w-3.5" />} />
                )}
                {eventData.social_links?.linkedin && (
                  <SocialChip href={eventData.social_links.linkedin} label="LinkedIn" icon={<Link2 className="h-3.5 w-3.5" />} />
                )}
              </div>
            </SectionCard>
          )}
        </main>
      </div>
    </div>
  )
}
