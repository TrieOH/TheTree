import { createFileRoute } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { eventsQueryOptions } from '@/features/events/api'
import { EventCard } from '@/features/events/ui/EventCard'
import { Search, ChevronDown, X } from 'lucide-react'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import { cn } from '@/shared/lib/utils'

export const Route = createFileRoute('/events/')({
  component: EventsPage,
  loader: ({ context }) => context.queryClient.ensureQueryData(eventsQueryOptions()),
})

const filterOptions = [
  { value: 'all', label: 'Todos os eventos' },
  { value: 'series', label: 'Apenas séries' },
] as const

function EventsPage() {
  const { data: events } = useSuspenseQuery(eventsQueryOptions())
  const [filter, setFilter] = useState<'all' | 'series'>('all')
  const [isFilterOpen, setIsFilterOpen] = useState(false)

  const filteredEvents = filter === 'series'
    ? events.filter(e => e.is_series)
    : events

  const currentLabel = filterOptions.find(o => o.value === filter)?.label

  const handleFilterSelect = (value: 'all' | 'series') => {
    setFilter(value)
    setIsFilterOpen(false)
  }

  return (
    <div className="min-h-screen bg-background relative">
      {/* Overlay escuro quando filtro aberto (mobile) */}
      <AnimatePresence>
        {isFilterOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={() => setIsFilterOpen(false)}
            className="fixed inset-0 bg-background/60 backdrop-blur-sm z-40 sm:hidden"
          />
        )}
      </AnimatePresence>

      {/* Header */}
      <header className="sticky top-0 z-40 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-4 h-14">
            {/* Título */}
            <h1 className="text-lg md:text-xl font-semibold text-foreground shrink-0">
              Eventos
              <span className="ml-2 text-sm font-normal text-muted-foreground">
                ({filteredEvents.length})
              </span>
            </h1>

            {/* Desktop: Tabs */}
            <nav className="hidden sm:flex items-center bg-muted rounded-lg p-1 ml-auto">
              {filterOptions.map((option) => (
                <button
                  key={option.value}
                  onClick={() => setFilter(option.value)}
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

            {/* Mobile: Botão que abre bottom sheet */}
            <div className="sm:hidden ml-auto">
              <button
                onClick={() => setIsFilterOpen(true)}
                className="flex items-center gap-2 bg-muted hover:bg-muted/80 active:bg-muted/60 transition-colors rounded-lg px-3 py-2"
                aria-expanded={isFilterOpen}
                aria-haspopup="listbox"
              >
                <span className="text-sm font-medium text-foreground">
                  {currentLabel === 'Todos os eventos' ? 'Todos' : 'Séries'}
                </span>
                <ChevronDown
                  className={cn(
                    "w-4 h-4 text-muted-foreground shrink-0 transition-transform duration-200",
                    isFilterOpen && "rotate-180"
                  )}
                />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Mobile: Bottom Sheet de filtros */}
      <AnimatePresence>
        {isFilterOpen && (
          <>
            {/* Handle indicator */}
            <motion.div
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "100%" }}
              transition={{ type: "spring", damping: 25, stiffness: 300 }}
              className="fixed bottom-0 left-0 right-0 z-50 sm:hidden bg-card border-t border-border rounded-t-2xl shadow-2xl"
            >
              {/* Drag handle */}
              <div className="flex justify-center pt-3 pb-2">
                <div className="w-12 h-1.5 bg-muted-foreground/20 rounded-full" />
              </div>

              {/* Header do sheet */}
              <div className="flex items-center justify-between px-4 pb-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground">Filtrar eventos</h2>
                <button
                  onClick={() => setIsFilterOpen(false)}
                  className="p-2 -mr-2 hover:bg-muted rounded-full transition-colors"
                  aria-label="Fechar filtros"
                >
                  <X className="w-5 h-5 text-muted-foreground" />
                </button>
              </div>

              {/* Opções */}
              <div className="p-2 pb-8" role="listbox">
                {filterOptions.map((option) => (
                  <button
                    key={option.value}
                    onClick={() => handleFilterSelect(option.value)}
                    className={cn(
                      "w-full flex items-center justify-between px-4 py-3.5 rounded-xl text-sm transition-colors",
                      filter === option.value
                        ? "bg-primary/10 text-primary font-medium"
                        : "text-foreground hover:bg-muted"
                    )}
                    role="option"
                    aria-selected={filter === option.value}
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
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Conteúdo principal */}
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