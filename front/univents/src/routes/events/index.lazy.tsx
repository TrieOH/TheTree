import { createLazyFileRoute } from '@tanstack/react-router'
import { Search, Settings, SlidersHorizontal } from 'lucide-react'
import { useEffect, useState } from 'react'
import { motion } from 'motion/react'
import { useQuery } from '@tanstack/react-query'
import { EventCard } from '@/features/events/ui/EventCard'
import { CreateEventCard } from '@/features/events/ui/CreateEventCard'
import { allOwnEventsQueryOptions, allPublicEventsQueryOptions } from '@/features/events/api'
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import { FABMenu } from '@/widgets/ui/fab-menu'
import {
  UI_PREFERENCES_CHANGE_EVENT,
  readInplaceEditPreference,
} from '@/shared/lib/ui-preferences'
import { ManageEventModal } from '@/features/events/ui/ManageEventModal'
import type { EventI } from '@/features/events/model'
import { useCreateEventMutation, usePatchEventMutation, usePublishEventMutation } from '@/features/events/api/mutations'

export const Route = createLazyFileRoute('/events/')({
  component: EventsPage,
})

const filterOptions = [
  { value: 'all', label: 'Todos os eventos' },
  { value: 'series', label: 'Apenas séries' },
] as const

const editFilterOptions = [
  { value: 'active', label: 'Ativos' },
  { value: 'draft', label: 'Rascunhos' },
] as const

type FilterValue = (typeof filterOptions)[number]['value'] | (typeof editFilterOptions)[number]['value']

function EventsPage() {
  const [isEditMode, setIsEditMode] = useState(false)
  const [inplaceEditEnabled, setInplaceEditEnabled] = useState(readInplaceEditPreference)
  const [filter, setFilter] = useState<FilterValue>('all')
  const [isFilterOpen, setIsFilterOpen] = useState(false)
  const [modalState, setModalState] = useState<{ open: boolean; event?: EventI }>({ open: false });
  const createMutation = useCreateEventMutation()
  const patchMutation = usePatchEventMutation()
  const publishMutation = usePublishEventMutation()

  const { data: publicEvents = [] } = useQuery({
    ...allPublicEventsQueryOptions(),
    enabled: !isEditMode
  })

  const { data: ownEvents, isFetching: isFetchingOwnEvents } = useQuery({
    ...allOwnEventsQueryOptions(),
    enabled: isEditMode,
  })

  const events = isEditMode
    ? (isFetchingOwnEvents ? publicEvents : (ownEvents ?? publicEvents))
    : publicEvents

  const sortedEvents = [...events].sort((a, b) => (
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  ))

  const visibleFilters = isEditMode
    ? [...filterOptions, ...editFilterOptions]
    : filterOptions

  const filteredEvents = sortedEvents.filter((event) => {
    if (filter === 'series') return event.is_series
    if (filter === 'active') return event.status === 'active'
    if (filter === 'draft') return event.status === 'draft'
    return true
  })

  const handleFilterSelect = (value: FilterValue) => {
    setFilter(value)
    setIsFilterOpen(false)
  }

  const handleEditToggle = () => {
    if (!inplaceEditEnabled) return

    setIsEditMode((current) => {
      const next = !current

      if (!next) {
        setFilter((currentFilter) => (
          currentFilter === 'active' || currentFilter === 'draft' ? 'all' : currentFilter
        ))
      }

      return next
    })
  }

  useEffect(() => {
    const syncPreferences = () => {
      const enabled = readInplaceEditPreference()
      setInplaceEditEnabled(enabled)

      if (!enabled) setIsEditMode(false)
    }

    syncPreferences()
    window.addEventListener(UI_PREFERENCES_CHANGE_EVENT, syncPreferences)

    return () => {
      window.removeEventListener(UI_PREFERENCES_CHANGE_EVENT, syncPreferences)
    }
  }, [])

  return (
    <div className="min-h-screen bg-background relative pb-24">
      {/* Header */}
      <header className="sticky top-0 z-30 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-2 h-14">
            <h1 className="text-lg md:text-xl font-semibold text-foreground">
              Eventos
              <span className="ml-2 text-sm font-normal text-muted-foreground">
                ({filteredEvents.length})
              </span>
            </h1>

            {/* Desktop */}
            <nav className="hidden sm:flex items-center bg-muted rounded-lg p-1 ml-auto">
              {visibleFilters.map((option) => (
                <Button
                  key={option.value}
                  type="button"
                  onClick={() => { setFilter(option.value) }}
                  className={cn(
                    "px-3 py-1.5 text-sm rounded-md transition-all whitespace-nowrap",
                    filter === option.value
                      ? "bg-background text-foreground shadow-sm"
                      : "text-muted-foreground hover:text-foreground"
                  )}
                  variant="ghost"
                >
                  {option.label === 'Todos os eventos' ? 'Todos' : option.label === 'Apenas séries' ? 'Séries' : option.label}
                </Button>
              ))}
            </nav>

            {/* Mobile */}
            <div className="sm:hidden! flex items-center gap-1 ml-auto">
              <Drawer open={isFilterOpen} onOpenChange={setIsFilterOpen}>
                <DrawerTrigger
                  render={
                    <Button
                      type="button"
                      className={cn(
                        "flex items-center justify-center w-9 h-9 rounded-lg transition-colors",
                        "hover:bg-muted active:bg-muted/60",
                        isFilterOpen && "bg-muted"
                      )}
                      aria-label="Filtrar eventos"
                      variant="ghost"
                    >
                      <SlidersHorizontal className="w-5 h-5 text-foreground" />
                    </Button>
                  }
                />
                <DrawerContent className="z-60 rounded-t-2xl border-t border-border bg-card">
                  <DrawerHeader className="pb-4 border-b border-border">
                    <DrawerTitle className="text-base font-semibold text-left">
                      Filtrar eventos
                    </DrawerTitle>
                  </DrawerHeader>

                  <div className="p-2 pb-8 space-y-1">
                    {visibleFilters.map((option) => (
                      <Button
                        key={option.value}
                        type="button"
                        onClick={() => { handleFilterSelect(option.value) }}
                        className={cn(
                          "w-full flex items-center justify-between px-4 py-3.5 rounded-xl text-sm transition-colors",
                          filter === option.value
                            ? "bg-primary/10 text-primary font-medium"
                            : "text-foreground hover:bg-muted"
                        )}
                        variant="ghost"
                      >
                        <span>{option.label}</span>
                        {filter === option.value && (
                          <motion.span
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            className="flex items-center justify-center w-5 h-5 bg-primary text-primary-foreground rounded-full text-xs"
                          >
                            ✓
                          </motion.span>
                        )}
                      </Button>
                    ))}
                  </div>
                </DrawerContent>
              </Drawer>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-6 md:gap-8">
          {isEditMode && (
            <CreateEventCard onClick={() => setModalState({ open: true, event: undefined })} />
          )}
          {filteredEvents.map((event, idx) => (
            <EventCard
              key={event.id}
              event={event}
              index={idx}
              showEditAction={isEditMode}
              onPublish={isEditMode ? (event) => publishMutation.mutate(event.id) : undefined}
              onEdit={() => { setModalState({ open: true, event: event }) }}
            />
          ))}
        </div>

        {!isEditMode && filteredEvents.length === 0 && (
          <div className="flex flex-col items-center justify-center py-12 md:py-16 space-y-6">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
              <Search className="w-8 h-8 text-muted-foreground/40" />
            </div>
            <div className="text-center space-y-1">
              <h3 className="text-lg font-medium text-foreground">
                Nenhum evento encontrado
              </h3>
              <p className="text-sm text-muted-foreground">
                Tente ajustar os filtros ou volte mais tarde.
              </p>
            </div>
          </div>
        )}
        <ManageEventModal
          key={modalState.event?.id ?? "event-create"}
          open={modalState.open}
          onOpenChange={(open) => setModalState((prev) => ({ ...prev, open }))}
          event={modalState.event}
          onCreate={(values): Promise<boolean> =>
            createMutation.mutateAsync(values).then(
              (res) => {
                if (!res.success) return false

                setModalState({ open: false, event: undefined })
                return true
              },
              () => false,
            )
          }
          onUpdate={(id, values): Promise<boolean> =>
            patchMutation.mutateAsync({ id, data: values }).then(
              (res) => {
                if (!res.success) return false

                setModalState({ open: false, event: undefined })
                return true
              },
              () => false,
            )
          }
        />
      </main>

      {inplaceEditEnabled && (
        <FABMenu
          mode="action"
          icon={Settings}
          onClick={handleEditToggle}
          active={isEditMode}
          ariaLabel={isEditMode ? 'Sair do modo de edição' : 'Ativar modo de edição'}
        />
      )}
    </div>
  )
}
