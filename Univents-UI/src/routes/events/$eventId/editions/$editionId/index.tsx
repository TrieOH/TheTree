import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import { useQuery } from '@tanstack/react-query'
import { EditionHeader } from '@/features/editions/ui/EditionHeader'
import { EditionTabs } from '@/features/editions/ui/EditionTab'
import { OverviewTab } from '@/features/editions/ui/OverviewTab'
import { ActivitiesTab } from '@/features/activities/ui/ActivitiesTab'
import { TicketsTab } from '@/features/tickets/ui/TicketsTab'
import { editionQueryOptions } from '@/features/editions/api'
import { allActivitiesQueryOptions } from '@/features/activities/api'
import { allTicketsQueryOptions } from '@/features/tickets/api'
import { Skeleton } from '@/shared/ui/shadcn/skeleton'

type TabValue = 'overview' | 'activities' | 'tickets'

export const Route = createFileRoute('/events/$eventId/editions/$editionId/')({
  component: EditionDetailPage,
})

function EditionDetailPage() {
  const { eventId, editionId } = Route.useParams()
  const [activeTab, setActiveTab] = useState<TabValue>('overview')

  const { data: edition, isLoading: isLoadingEdition } = useQuery(
    editionQueryOptions(eventId, editionId)
  );

  const { data: activities = [] } = useQuery(
    allActivitiesQueryOptions(eventId, editionId)
  );

  const { data: tickets = [] } = useQuery(
    allTicketsQueryOptions(eventId, editionId)
  );

  if (isLoadingEdition) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8 space-y-8">
        <div className="space-y-4">
          <Skeleton className="h-12 w-3/4" />
          <Skeleton className="h-6 w-1/2" />
        </div>
        <div className="flex gap-4">
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-24" />
        </div>
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (!edition) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold">Edição não encontrada</h1>
          <p className="text-muted-foreground">O link pode estar quebrado ou a edição foi removida.</p>
        </div>
      </div>
    )
  }

  const tabCounts = {
    activities: activities.filter(a => a.status === 'published').length,
    tickets: tickets.length,
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12 space-y-8">
        <EditionHeader edition={edition} eventId={eventId} />

        <EditionTabs
          activeTab={activeTab}
          onTabChange={setActiveTab}
          counts={tabCounts}
        />

        <AnimatePresence mode="wait">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
          >
            {activeTab === 'overview' && (
              <OverviewTab edition={edition} eventId={eventId} />
            )}
            {activeTab === 'activities' && (
              <ActivitiesTab activities={activities} />
            )}
            {activeTab === 'tickets' && (
              <TicketsTab tickets={tickets} eventId={eventId} editionId={editionId} />
            )}
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  )
}