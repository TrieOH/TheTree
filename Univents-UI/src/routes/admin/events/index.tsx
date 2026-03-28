import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import {
  Plus,
  ShieldCheck,
  MoreVertical,
  Calendar,
} from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import type { EventCreateI, EventI } from '@/features/events/model';
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { cn } from '@/shared/lib/utils'
import { eventCreateSchema } from '@/features/events/model'
import { FormDrawer } from '@/widgets/form/ui/form-drawer'
import { AlertDrawer } from '@/widgets/ui/alert-drawer'
import { getEventFields } from '@/features/events/model/field'
import { ownEventsQueryOptions } from '@/features/events/api'
import AdminEventCard from '@/features/events/ui/AdminEventCard'

export const Route = createFileRoute('/admin/events/')({
  component: AdminEventsPage,
})

function AdminEventsPage() {
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [editingEvent, setEditingEvent] = useState<EventI | null>(null)
  const [deletingEvent, setDeletingEvent] = useState<EventI | null>(null)
  const [publishingEvent, setPublishingEvent] = useState<EventI | null>(null)
  const [isActionsOpen, setIsActionsOpen] = useState(false)
  const [loading, setLoading] = useState(false)

  const { data: events = [] } = useQuery(ownEventsQueryOptions())

  const handleCreate = (data: EventCreateI) => {
    setLoading(true)
    console.log('create', data)
    setLoading(false)
  }

  const handleEdit = (data: EventCreateI) => {
    setLoading(true)
    console.log('edit', data)
    setLoading(false)
    setEditingEvent(null)
  }

  const handleDelete = () => {
    if (!deletingEvent) return
    console.log('delete', deletingEvent.id)
    setDeletingEvent(null)
  }

  const handlePublish = () => {
    if (!publishingEvent) return
    console.log('publish', publishingEvent.id)
    setPublishingEvent(null)
  }

  const getInitialData = (event: EventI | null): Partial<EventCreateI> => event ? {
    name: event.name,
    slug: event.slug,
    acronym: event.acronym,
    contact_email: event.contact_email ?? undefined,
    tagline: event.tagline,
    description: event.description,
    is_series: event.is_series,
    logo_url: event.logo_url,
    banner_url: event.banner_url,
  } : {}

  return (
    <div className="min-h-screen bg-background relative pb-20 md:pb-0">
      <header className="sticky top-0 z-30 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <h1 className="text-lg md:text-xl font-semibold text-foreground shrink-0 flex items-center gap-2">
              Eventos
              <span className="text-sm font-normal text-muted-foreground">({events.length})</span>
            </h1>

            <div className="hidden sm:flex items-center gap-2 ml-auto">
              <button
                onClick={() => { setIsCreateOpen(true); }}
                className={cn(
                  "flex items-center gap-2 px-4 py-2 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium"
                )}
              >
                <Plus className="w-4 h-4" />
                Novo evento
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
                      <span className="font-medium">Novo evento</span>
                    </button>
                  </div>
                </DrawerContent>
              </Drawer>
            </div>

            <Link
              to="/events"
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
          {events.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="flex flex-col items-center justify-center py-24 space-y-6"
            >
              <div className="w-20 h-20 rounded-2xl bg-muted flex items-center justify-center">
                <Calendar className="w-10 h-10 text-muted-foreground/30" />
              </div>
              <div className="text-center space-y-2">
                <h3 className="text-lg font-medium">Nenhum evento ainda</h3>
                <p className="text-sm text-muted-foreground max-w-xs">
                  Crie seu primeiro evento para começar a gerenciar inscrições e programação.
                </p>
              </div>
              <button
                onClick={() => { setIsCreateOpen(true); }}
                className={cn(
                  "mt-2 px-5 py-2.5 rounded-lg",
                  "bg-primary text-primary-foreground hover:bg-primary/90",
                  "text-sm font-medium",
                  "active:scale-95 transition-all"
                )}
              >
                Criar evento
              </button>
            </motion.div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
              {events.map((event, idx) => (
                <AdminEventCard
                  key={event.id}
                  event={event}
                  index={idx}
                  onEdit={setEditingEvent}
                  onDelete={setDeletingEvent}
                  onPublish={setPublishingEvent}
                />
              ))}
            </div>
          )}
        </AnimatePresence>
      </main>

      <FormDrawer
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        title="Novo evento"
        fields={getEventFields()}
        schema={eventCreateSchema}
        onSubmit={handleCreate}
        submitLabel="Criar evento"
        loading={loading}
      />

      <FormDrawer
        open={!!editingEvent}
        onOpenChange={(open) => {
          if (!open) setEditingEvent(null)
        }}
        title="Editar evento"
        fields={getEventFields(editingEvent?.id)}
        schema={eventCreateSchema}
        onSubmit={handleEdit}
        defaultValues={getInitialData(editingEvent)}
        submitLabel="Salvar alterações"
        loading={loading}
      />

      <AlertDrawer
        open={!!publishingEvent}
        onOpenChange={() => { setPublishingEvent(null); }}
        title="Publicar evento?"
        description={`Ao publicar "${publishingEvent?.name}", ele ficará visível para o público.`}
        confirmLabel="Publicar"
        onConfirm={handlePublish}
        variant="success"
      />

      <AlertDrawer
        open={!!deletingEvent}
        onOpenChange={() => { setDeletingEvent(null); }}
        title="Deletar evento?"
        description={`Tem certeza que deseja deletar "${deletingEvent?.name}"?`}
        confirmLabel="Deletar"
        onConfirm={handleDelete}
        variant="destructive"
      />
    </div>
  )
}
