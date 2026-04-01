import { createFileRoute, Link } from '@tanstack/react-router'
import { PencilLine, Share2, Save, X, ArrowRight, Calendar, Loader2 } from 'lucide-react'
import { useState, useEffect, useMemo } from 'react'
import { toast } from 'sonner'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { EventI } from '@/features/events/model'
import { Button } from '@/shared/ui/shadcn/button'
import { Badge } from '@/shared/ui/shadcn/badge'
import { cn } from '@/shared/lib/utils'
import { parseDatetimeLocal } from '@/shared/lib/date'
import InlineEdit from '@/shared/ui/inline-edit/InlineEdit'
import InlineImageEdit from '@/shared/ui/inline-edit/InlineImageEdit'
import { eventQueryOptions, ownEventQueryOptions, patchEventFn } from '@/features/events/api'
import { uploadAndModerateFile } from '@/features/storage/api'


export const Route = createFileRoute('/events/$eventId/')({
  component: RouteComponent,
  validateSearch: (search) => ({
    edit: search.edit === 'true' || search.edit === true,
  }),
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

  const { data: event, isLoading, error } = useQuery(
    edit ? ownEventQueryOptions(eventId) : eventQueryOptions(eventId)
  )

  const [eventData, setEventData] = useState<EventI | null>(null)

  useEffect(() => {
    if (event) {
      setEventData(event)
    }
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
    onError: () => {
      toast.error('Erro ao salvar alterações')
    }
  })

  const toggleEditMode = () => {
    void navigate({ search: (prev) => ({ ...prev, edit: !edit }) })
    setEditingField(null)
  }

  const updateField = (field: keyof EventI, value: string) => {
    if (!eventData) return
    setEventData((prev) => prev ? ({ ...prev, [field]: value }) : null)
  }

  const handleSave = () => {
    if (!eventData) return
    const { name, acronym, tagline, description, logo_url, banner_url } = eventData
    mutation.mutate({ name, acronym, tagline, description, logo_url, banner_url })
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
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
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

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .slice(0, 2)
      .toUpperCase()
  }

  const logoFallbackText = eventData.acronym
    ? eventData.acronym.slice(0, 4)
    : getInitials(eventData.name)

  const logoFontSize =
    logoFallbackText.length <= 2 ? 'text-lg' :
      logoFallbackText.length === 3 ? 'text-sm' :
        'text-xs'

  return (
    <div className="min-h-screen bg-background">
      {/* Banner */}
      <InlineImageEdit
        value={eventData.banner_url}
        onChange={(url) => { updateField('banner_url', url); }}
        isEditEnabled={edit}
        isEditing={editingField === 'banner_url'}
        onStartEdit={() => { setEditingField('banner_url'); }}
        onFinishEdit={() => { setEditingField(null); }}
        onUpload={(file) => uploadAndModerateFile(file, `events/${eventId}`)}
        renderDisplay={(url) => (
          <div className="relative w-full h-48 md:h-64 overflow-hidden">
            {url
              ? <img src={url} alt={eventData.name} className="w-full h-full object-cover" />
              : <div className="w-full h-full bg-linear-to-br from-primary via-accent to-secondary" />
            }
            <div className="absolute inset-0 bg-linear-to-b from-foreground/10 to-transparent" />
          </div>
        )}
      />

      <div className="mx-auto max-w-7xl">
        <main className="relative px-4 pb-8">
          <div className="relative -mt-20 md:-mt-16 bg-card rounded-2xl shadow-xl border border-border p-5 md:p-6">
            <div className="flex items-start justify-between mb-4">

              {/* Logo */}
              <InlineImageEdit
                value={eventData.logo_url}
                onChange={(url) => { updateField('logo_url', url); }}
                isEditEnabled={edit}
                isEditing={editingField === 'logo_url'}
                onStartEdit={() => { setEditingField('logo_url'); }}
                onFinishEdit={() => { setEditingField(null); }}
                className="h-16 w-16 shrink-0 -mt-12 md:-mt-14"
                onUpload={(file) => uploadAndModerateFile(file, `events/${eventId}`)}
                renderDisplay={(url) => (
                  <div className="h-full w-full rounded-xl bg-primary shadow-lg flex items-center justify-center overflow-hidden border-4 border-card">
                    {url
                      ? <img src={url} alt={eventData.acronym ?? eventData.name} className="h-full w-full object-cover" />
                      : <span className={cn("text-primary-foreground font-bold leading-none tracking-tight px-1 text-center break-all", logoFontSize)}>
                        {logoFallbackText}
                      </span>
                    }
                  </div>
                )}
              />

              {/* Status, Edition and Actions */}
              <div className="flex flex-col items-end gap-1">
                <div className="flex items-center gap-1">
                  {!edit ? (
                    <>
                      <Button
                        className="h-8 w-8 hover:text-foreground duration-200 transition-colors text-muted-foreground"
                        variant="ghost"
                        size="icon"
                        onClick={() => { void handleShare() }}
                        title="Compartilhar"
                      >
                        <Share2 className="h-4 w-4" />
                      </Button>
                      <Button
                        className="h-8 w-8 hover:text-foreground duration-200 transition-colors text-muted-foreground"
                        variant="ghost"
                        size="icon"
                        onClick={toggleEditMode}
                        title="Editar evento"
                      >
                        <PencilLine className="h-4 w-4" />
                      </Button>
                    </>
                  ) : (
                    <Button
                      className="h-8 w-8"
                      variant="ghost"
                      size="icon"
                      onClick={toggleEditMode}
                      title="Cancelar edição"
                    >
                      <X className="h-4 w-4 text-destructive" />
                    </Button>
                  )}
                </div>
                <Badge
                  variant="outline"
                  className={`text-xs font-medium ${statusColors[eventData.status]}`}
                >
                  {statusLabels[eventData.status]}
                </Badge>
                {eventData.is_series && eventData.editions_count > 0 && (
                  <span className="text-xs text-muted-foreground uppercase tracking-wider font-medium">
                    {eventData.editions_count}ª EDIÇÃO
                  </span>
                )}
              </div>
            </div>

            {/* Main Info */}
            <div className="space-y-3">
              {/* Event Name */}
              <div>
                <InlineEdit
                  value={eventData.name}
                  onChange={(val) => { updateField('name', val); }}
                  isEditEnabled={edit}
                  isEditing={editingField === 'name'}
                  onStartEdit={() => { setEditingField('name'); }}
                  onFinishEdit={() => { setEditingField(null); }}
                  className="text-2xl md:text-3xl font-bold tracking-tight text-foreground"
                />
              </div>

              {/* Acronym and Tagline */}
              <div className="flex flex-wrap items-center gap-2 text-sm md:text-base">
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
                {eventData.tagline && (
                  <>
                    <span className="text-muted-foreground">•</span>
                    <span className="text-muted-foreground italic">
                      <InlineEdit
                        value={eventData.tagline}
                        onChange={(val) => { updateField('tagline', val); }}
                        isEditEnabled={edit}
                        isEditing={editingField === 'tagline'}
                        onStartEdit={() => { setEditingField('tagline'); }}
                        onFinishEdit={() => { setEditingField(null); }}
                        className="text-muted-foreground italic"
                        placeholder="Tagline do evento..."
                      />
                    </span>
                  </>
                )}
              </div>

              {/* Description */}
              <div className="pt-2">
                <InlineEdit
                  value={eventData.description}
                  onChange={(val) => { updateField('description', val); }}
                  isEditEnabled={edit}
                  isEditing={editingField === 'description'}
                  onStartEdit={() => { setEditingField('description'); }}
                  onFinishEdit={() => { setEditingField(null); }}
                  multiline
                  className="text-base leading-relaxed text-muted-foreground block w-full"
                  placeholder="Adicione uma descrição do evento..."
                />
              </div>

              <div className="flex items-center gap-4 text-sm text-muted-foreground pt-2">
                <div className="flex items-center gap-1.5">
                  <Calendar className="h-4 w-4" />
                  <span>{parseDatetimeLocal(eventData.created_at).toLocaleDateString('pt-BR')}</span>
                </div>
              </div>
            </div>

            {/* Go To Editions*/}
            <div className="mt-6">
              <Link
                to="/events/$eventId/editions"
                params={{ eventId: eventData.id }}
                className={cn(
                  "flex items-center justify-center",
                  "w-full bg-primary/80 hover:bg-primary text-primary-foreground",
                  "rounded-sm h-12 font-medium transition-colors duration-300"
                )}
              >
                Ver Edições
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </div>
          </div>
        </main>
      </div>

      {/* Float Button Save */}
      {edit && isDirty && (
        <div className="fixed bottom-6 right-6 z-50 animate-in slide-in-from-bottom-2 fade-in duration-200">
          <Button
            size="lg"
            className="h-14 px-6 rounded-full shadow-xl hover:scale-105 active:scale-95 transition-all bg-primary text-primary-foreground hover:bg-primary/90 flex items-center gap-2"
            onClick={handleSave}
            disabled={mutation.isPending}
          >
            {mutation.isPending ? (
              <Loader2 className="h-5 w-5 animate-spin" />
            ) : (
              <Save className="h-5 w-5" />
            )}
            <span className="font-semibold">Salvar Alterações</span>
          </Button>
        </div>
      )}
    </div>
  )
}