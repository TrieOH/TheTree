import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { Search, SlidersHorizontal, ShieldCheck } from 'lucide-react'
import { useState } from 'react'
import { motion } from 'motion/react'
import { useAuth } from '@soramux/node-auth-sdk/react'
import { EventCard } from '@/features/events/ui/EventCard'
import { cn } from '@/shared/lib/utils'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/ui/shadcn/drawer'
import {
  canCreateEvent,
  canEditEvent,
  canPublishEvent
} from '@/features/events/model/permissions'
import { usePermissions } from '@/features/auths/hooks/use-permissions'

export const Route = createLazyFileRoute('/events/')({
  component: EventsPage,
})

const filterOptions = [
  { value: 'all', label: 'Todos os eventos' },
  { value: 'series', label: 'Apenas séries' },
] as const

function EventsPage() {
  const { auth } = useAuth();
  const userProfile = auth.profile()
  const events = Route.useLoaderData()
  const { some: somePerms } = usePermissions(
    { canEditEvent, canPublishEvent, canCreateEvent },
    userProfile?.id
  )
  const isAdmin = somePerms('canEditEvent', 'canPublishEvent', 'canCreateEvent')

  const [filter, setFilter] = useState<'all' | 'series'>('all')
  const [isFilterOpen, setIsFilterOpen] = useState(false)

  const filteredEvents = filter === 'series'
    ? events.filter(e => e.is_series)
    : events

  const handleFilterSelect = (value: 'all' | 'series') => {
    setFilter(value)
    setIsFilterOpen(false)
  }

  return (
    <div className="min-h-screen bg-background relative">
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
              {filterOptions.map((option) => (
                <button
                  key={option.value}
                  onClick={() => { setFilter(option.value) }}
                  className={cn(
                    "px-3 py-1.5 text-sm rounded-md transition-all whitespace-nowrap",
                    filter === option.value
                      ? "bg-background text-foreground shadow-sm"
                      : "text-muted-foreground hover:text-foreground"
                  )}
                >
                  {option.label === 'Todos os eventos' ? 'Todos' : 'Séries'}
                </button>
              ))}
            </nav>

            {/* Mobile */}
            <div className="sm:hidden flex items-center gap-1 ml-auto">
              <Drawer open={isFilterOpen} onOpenChange={setIsFilterOpen}>
                <DrawerTrigger asChild>
                  <button
                    className={cn(
                      "flex items-center justify-center w-9 h-9 rounded-lg transition-colors",
                      "hover:bg-muted active:bg-muted/60",
                      isFilterOpen && "bg-muted"
                    )}
                    aria-label="Filtrar eventos"
                  >
                    <SlidersHorizontal className="w-5 h-5 text-foreground" />
                  </button>
                </DrawerTrigger>

                <DrawerContent className="z-60 rounded-t-2xl border-t border-border bg-card">
                  <DrawerHeader className="pb-4 border-b border-border">
                    <DrawerTitle className="text-base font-semibold text-left">
                      Filtrar eventos
                    </DrawerTitle>
                  </DrawerHeader>

                  <div className="p-2 pb-8 space-y-1">
                    {filterOptions.map((option) => (
                      <button
                        key={option.value}
                        onClick={() => { handleFilterSelect(option.value) }}
                        className={cn(
                          "w-full flex items-center justify-between px-4 py-3.5 rounded-xl text-sm transition-colors",
                          filter === option.value
                            ? "bg-primary/10 text-primary font-medium"
                            : "text-foreground hover:bg-muted"
                        )}
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
                      </button>
                    ))}
                  </div>
                </DrawerContent>
              </Drawer>

              {isAdmin && (
                <Link
                  to="/admin/events"
                  className={cn(
                    "flex items-center justify-center",
                    "w-9 h-9 rounded-lg transition-all duration-200",
                    "text-muted-foreground hover:text-foreground",
                    "hover:bg-muted active:bg-muted/60"
                  )}
                  aria-label="Painel administrativo"
                >
                  <ShieldCheck className="w-5 h-5" />
                </Link>
              )}
            </div>

            {/* Desktop: Admin */}
            {isAdmin && (
              <Link
                to="/admin/events"
                className={cn(
                  "hidden sm:flex group relative items-center justify-center",
                  "w-9 h-9 rounded-lg transition-all duration-200",
                  "text-muted-foreground hover:text-foreground",
                  "hover:bg-muted active:bg-muted/60",
                  "shrink-0"
                )}
                aria-label="Painel administrativo"
              >
                <ShieldCheck className="w-5 h-5 transition-transform duration-200 group-hover:scale-110" />
                <span className={cn(
                  "pointer-events-none absolute -bottom-9 right-0",
                  "whitespace-nowrap rounded-md px-2 py-1",
                  "bg-popover text-popover-foreground border border-border",
                  "text-xs shadow-md",
                  "opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0",
                  "transition-all duration-150"
                )}>
                  Modo admin
                </span>
              </Link>
            )}
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
        {filteredEvents.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-6 md:gap-8">
            {filteredEvents.map((event, idx) => (
              <EventCard key={event.id} event={event} index={idx} />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-24 md:py-32 space-y-6">
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
      </main>
    </div>
  )
}
