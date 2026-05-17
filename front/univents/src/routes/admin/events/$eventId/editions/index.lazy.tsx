import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import {
  Plus,
  Calendar,
  MoreVertical,
  ShieldCheck,
} from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { EditionCreateI, EditionI } from '@/features/editions/model'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { cn } from '@/shared/lib/utils'
import { FormDrawer } from '@/widgets/form/ui/form-drawer'
import {
  allAdminEditionsQueryOptions,
  createEditionFn,
  disconnectPaymentAccountToEditionFn,
  publishEditionFn,
} from '@/features/editions/api'
import { AlertModal } from '@/widgets/ui/alert-modal'
import { editionCreateSchema } from '@/features/editions/model'
import { getEditionFields } from '@/features/editions/model/field'
import { ownEventsQueryOptions } from '@/features/events/api'
import {
  connectEditionSellerToWorkspaceFn,
  disconnectEditionSellerToWorkspaceFn
} from '@/features/payments/api'
import { env } from '@/env'
import { AdminEditionCard } from '@/features/editions/ui/AdminEditionCard'

export const Route = createLazyFileRoute('/admin/events/$eventId/editions/')({
  component: RouteComponent,
})

const defaultTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone

function RouteComponent() {
  const queryClient = useQueryClient()
  const { eventId } = Route.useParams()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [publishingEdition, setPublishingEdition] = useState<EditionI | null>(null)
  const [isActionsOpen, setIsActionsOpen] = useState(false)

  const { data: events = [] } = useQuery(ownEventsQueryOptions())
  const event = events.find(e => e.id === eventId) ?? null

  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))

  const createMutation = useMutation({
    mutationFn: (data: EditionCreateI) => createEditionFn(data, eventId),
    onSuccess: (res) => {
      if (res.success) {
        queryClient.setQueryData<EditionI[]>(
          allAdminEditionsQueryOptions(eventId).queryKey,
          (old = []) => [...old, res.data]
        )
        setIsCreateOpen(false)
        toast.success('Edição criada com sucesso!')
      } else toast.error(res.message || 'Erro ao criar edição')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })

  const publishMutation = useMutation({
    mutationFn: ({ editionId }: { editionId: string }) =>
      publishEditionFn(eventId, editionId),
    onSuccess: (res, variables) => {
      if (res.success) {
        queryClient.setQueryData<EditionI[]>(
          allAdminEditionsQueryOptions(eventId).queryKey,
          (old = []) => old.map((ed: EditionI) =>
            ed.id === variables.editionId ? { ...ed, status: 'announced' as const } : ed
          )
        )
        setPublishingEdition(null)
        toast.success('Edição publicada com sucesso!')
      } else toast.error(res.message || 'Erro ao publicar edição')
    },
    onError: () => toast.error('Erro ao conectar com o servidor')
  })

  const handleCreate = (data: EditionCreateI) => {
    createMutation.mutate(data)
  }

  const handlePublish = () => {
    if (!publishingEdition) return
    publishMutation.mutate({ editionId: publishingEdition.id })
  }

  const handleConnect = async (editionId: string) => {
    const res = await connectEditionSellerToWorkspaceFn({
      data: {
        provider: 'mercadopago',
        workspace_name: 'Univents',
        final_redirect_url: `${window.location.origin}/admin/events/${eventId}/editions/${editionId}/callback/payment`,
        provider_redirect_url: env.VITE_MERCADO_PAGO_CALLBACK_URL,
      },
    })
    if (res.success) window.location.href = res.data.redirect_url
  }

  const handleDisconnect = async (editionId: string, credentialId: string) => {
    const ws = await disconnectEditionSellerToWorkspaceFn({
      data: { workspace_name: 'Univents', credential_id: credentialId }
    })
    if (!ws.success) return
    const res = await disconnectPaymentAccountToEditionFn(eventId, editionId)
    if (res.success) {
      queryClient.setQueryData<EditionI[]>(
        allAdminEditionsQueryOptions(eventId).queryKey,
        (old = []) => old.map(ed =>
          ed.id === editionId
            ? { ...ed, trie_payments_credential_id: null, trie_payments_provider: null, trie_payments_provider_public_key: null }
            : ed
        )
      )
      toast.success('Conta desconectada com sucesso!')
    }
  }

  const loading = createMutation.isPending || publishMutation.isPending

  return (
    <div className="min-h-screen bg-background relative pb-20 md:pb-0">
      <header className="sticky top-0 z-30 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <div className="flex items-center gap-2 shrink-0">
              <h1 className="text-lg md:text-xl font-semibold text-foreground">
                Edições
                <span className="ml-2 text-sm font-normal text-muted-foreground">
                  ({editions.length})
                </span>
              </h1>
            </div>

            <div className="hidden sm:flex items-center gap-2 ml-auto">
              <button
                onClick={() => { setIsCreateOpen(true) }}
                className={cn(
                  "flex items-center gap-2 px-4 py-2 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium"
                )}
              >
                <Plus className="w-4 h-4" />
                Nova edição
              </button>
            </div>

            <div className="sm:hidden flex items-center gap-1 ml-auto">
              <Drawer open={isActionsOpen} onOpenChange={setIsActionsOpen}>
                <DrawerTrigger asChild>
                  <button className={cn("flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted")}>
                    <MoreVertical className="w-5 h-5 text-foreground" />
                  </button>
                </DrawerTrigger>
                <DrawerContent className="z-60 rounded-t-2xl">
                  <DrawerHeader className="pb-4 border-b">
                    <DrawerTitle className="text-base font-semibold">Ações</DrawerTitle>
                  </DrawerHeader>
                  <div className="p-2 pb-8 space-y-1">
                    <button
                      onClick={() => { setIsActionsOpen(false); setIsCreateOpen(true) }}
                      className="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl hover:bg-muted"
                    >
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Plus className="w-4 h-4 text-primary" />
                      </div>
                      <span className="font-medium">Nova edição</span>
                    </button>
                  </div>
                </DrawerContent>
              </Drawer>
            </div>

            <Link
              to="/events/$eventId/editions"
              params={{ eventId }}
              className={cn(
                "group relative flex items-center justify-center",
                "w-9 h-9 rounded-lg transition-all",
                "bg-primary text-primary-foreground",
                "hover:bg-primary/90",
                "shrink-0"
              )}
            >
              <ShieldCheck className="w-5 h-5" />
              <span
                className={cn(
                  "pointer-events-none absolute -bottom-9 right-0",
                  "whitespace-nowrap rounded-md px-2 py-1",
                  "bg-popover text-popover-foreground border text-xs shadow-md",
                  "opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0",
                  "transition-all"
                )}>
                Sair do admin
              </span>
            </Link>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
        <AnimatePresence mode="wait">
          {editions.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="flex flex-col items-center justify-center py-24 space-y-6"
            >
              <div className="w-20 h-20 rounded-2xl bg-muted flex items-center justify-center">
                <Calendar className="w-10 h-10 text-muted-foreground/30" />
              </div>
              <div className="text-center space-y-2">
                <h3 className="text-lg font-medium">Nenhuma edição ainda</h3>
                <p className="text-sm text-muted-foreground max-w-xs">
                  Crie a primeira edição para {event?.name ? `"${event.name}"` : 'este evento'}.
                </p>
              </div>
              <button
                onClick={() => { setIsCreateOpen(true) }}
                className={cn(
                  "mt-2 px-5 py-2.5 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium",
                  "active:scale-95 transition-all"
                )}
              >
                Criar edição
              </button>
            </motion.div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
              {editions.map((edition, idx) => (
                <AdminEditionCard
                  key={edition.id}
                  edition={edition}
                  eventId={eventId}
                  index={idx}
                  onPublish={() => { setPublishingEdition(edition); }}
                  onConnect={() => { void handleConnect(edition.id) }}
                  onDisconnect={() => {
                    if (edition.trie_payments_credential_id) {
                      void handleDisconnect(edition.id, edition.trie_payments_credential_id)
                    }
                  }}
                />
              ))}
            </div>
          )}
        </AnimatePresence>
      </main>

      <FormDrawer
        idPrefix="create-"
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        title="Nova edição"
        fields={getEditionFields()}
        defaultValues={{ timezone: defaultTimezone }}
        schema={editionCreateSchema}
        onSubmit={handleCreate}
        submitLabel="Criar edição"
        loading={loading}
      />

      <AlertModal
        open={!!publishingEdition}
        onOpenChange={() => { setPublishingEdition(null) }}
        title="Publicar edição?"
        description={`Ao publicar "${publishingEdition?.edition_name}", ela ficará visível para o público.`}
        confirmLabel="Publicar"
        onConfirm={handlePublish}
        variant="success"
        loading={loading}
      />
    </div>
  )
}

