import { createFileRoute, Link } from '@tanstack/react-router'
import { motion, AnimatePresence } from 'motion/react'
import {
  CalendarX,
  ChevronDown,
  ArrowLeft,
  SlidersHorizontal,
  X,
} from 'lucide-react'
import { useState } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query'
import { cn } from '@/shared/lib/utils'
import { EditionCard } from '@/features/editions/ui/EditionCard'
import { allEditionsQueryOptions } from '@/features/editions/api'

const statusFilters = [
  { value: 'all', label: 'Todos' },
  { value: 'open', label: 'Inscrições abertas' },
  { value: 'ongoing', label: 'Em andamento' },
  { value: 'announced', label: 'Anunciados' },
  { value: 'finished', label: 'Encerrados' },
] as const

const typeFilters = [
  { value: 'all', label: 'Todos os tipos' },
  { value: 'year', label: 'Anual' },
  { value: 'season', label: 'Temporada' },
  { value: 'number', label: 'Numerada' },
  { value: 'ordinal', label: 'Ordinal' },
  { value: 'custom', label: 'Personalizado' },
] as const

export const Route = createFileRoute('/events/$eventId/editions/')({
  component: EventEditionsPage,
  loader: async ({ context: ctx, params }) => {
    await ctx.queryClient.ensureQueryData(
      allEditionsQueryOptions(params.eventId)
    )
  },
})

function EventEditionsPage() {
  const { eventId } = Route.useParams()
  const { data: editions } = useSuspenseQuery(allEditionsQueryOptions(eventId))
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [typeFilter, setTypeFilter] = useState<string>('all')
  const [isFilterOpen, setIsFilterOpen] = useState(false)

  const filteredEditions = editions.filter(edition => {
    const matchesStatus = statusFilter === 'all' || edition.status === statusFilter
    const matchesType = typeFilter === 'all' || edition.type === typeFilter
    return matchesStatus && matchesType
  })

  const activeFiltersCount = (statusFilter !== 'all' ? 1 : 0) + (typeFilter !== 'all' ? 1 : 0)

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between gap-4 h-14">
            <div className="flex items-center gap-3">
              <Link
                to="/events"
                className="flex items-center justify-center w-9 h-9 rounded-lg hover:bg-muted transition-colors"
                aria-label="Voltar para eventos"
              >
                <ArrowLeft className="w-5 h-5 text-muted-foreground" />
              </Link>
              <h1 className="text-lg md:text-xl font-semibold text-foreground">
                Edições
                <span className="ml-2 text-sm font-normal text-muted-foreground">
                  ({filteredEditions.length})
                </span>
              </h1>
            </div>

            <button
              onClick={() => { setIsFilterOpen(true); }}
              className={cn(
                "flex items-center gap-2 rounded-lg px-3 py-2 transition-colors",
                activeFiltersCount > 0
                  ? "bg-primary/10 text-primary hover:bg-primary/20"
                  : "bg-muted hover:bg-muted/80 text-foreground"
              )}
            >
              <SlidersHorizontal className="w-4 h-4" />
              <span className="hidden sm:inline text-sm font-medium">
                Filtros
                {activeFiltersCount > 0 && ` (${activeFiltersCount})`}
              </span>
              <ChevronDown className={cn(
                "w-4 h-4 transition-transform",
                isFilterOpen && "rotate-180"
              )} />
            </button>
          </div>
        </div>
      </div>

      <AnimatePresence>
        {isFilterOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => { setIsFilterOpen(false); }}
              className="fixed inset-0 bg-background/60 backdrop-blur-sm z-40"
            />
            <motion.div
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "100%" }}
              transition={{ type: "spring", damping: 25, stiffness: 300 }}
              className="fixed bottom-0 left-0 right-0 z-50 bg-card border-t border-border rounded-t-2xl md:max-w-md md:left-auto md:right-4 md:bottom-4 md:rounded-2xl md:border md:shadow-2xl"
            >
              <div className="flex justify-center pt-3 pb-2 md:hidden">
                <div className="w-12 h-1.5 bg-muted-foreground/20 rounded-full" />
              </div>

              {/* Header */}
              <div className="flex items-center justify-between px-4 pb-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground">Filtrar edições</h2>
                <button
                  onClick={() => { setIsFilterOpen(false); }}
                  className="p-2 -mr-2 hover:bg-muted rounded-full transition-colors"
                >
                  <X className="w-5 h-5 text-muted-foreground" />
                </button>
              </div>

              <div className="p-4 space-y-6 max-h-[60vh] overflow-y-auto">
                {/* Status */}
                <div className="space-y-3">
                  <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
                    Status
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {statusFilters.map((filter) => (
                      <button
                        key={filter.value}
                        onClick={() => { setStatusFilter(filter.value); }}
                        className={cn(
                          "px-3 py-2 rounded-lg text-sm transition-colors",
                          statusFilter === filter.value
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-foreground hover:bg-muted/80"
                        )}
                      >
                        {filter.label}
                      </button>
                    ))}
                  </div>
                </div>

                <div className="space-y-3">
                  <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
                    Tipo
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {typeFilters.map((filter) => (
                      <button
                        key={filter.value}
                        onClick={() => { setTypeFilter(filter.value); }}
                        className={cn(
                          "px-3 py-2 rounded-lg text-sm transition-colors",
                          typeFilter === filter.value
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-foreground hover:bg-muted/80"
                        )}
                      >
                        {filter.label}
                      </button>
                    ))}
                  </div>
                </div>

                {activeFiltersCount > 0 && (
                  <button
                    onClick={() => {
                      setStatusFilter('all')
                      setTypeFilter('all')
                    }}
                    className="w-full py-3 text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    Limpar filtros
                  </button>
                )}
              </div>

              {/* Footer */}
              <div className="p-4 border-t border-border">
                <button
                  onClick={() => { setIsFilterOpen(false); }}
                  className="w-full py-3 bg-primary text-primary-foreground rounded-xl font-medium hover:bg-primary/90 transition-colors"
                >
                  Aplicar
                </button>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
        {filteredEditions.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-8">
            {filteredEditions.map((edition, idx) => (
              <EditionCard
                key={edition.id}
                edition={edition}
                eventId={eventId}
                index={idx}
              />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-24 md:py-32 space-y-6">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
              <CalendarX className="w-8 h-8 text-muted-foreground/40" />
            </div>
            <div className="text-center space-y-2">
              <h3 className="text-lg font-medium text-foreground">
                Nenhuma edição encontrada
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