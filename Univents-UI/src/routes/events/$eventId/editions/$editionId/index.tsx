import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { motion, AnimatePresence } from 'motion/react'
import { mockEdition, mockActivities, mockCheckpoints, mockTickets } from '@/features/mock'
import { EditionHeader } from '@/features/editions/ui/EditionHeader'
import { EditionTabs } from '@/features/editions/ui/EditionTab'
import { OverviewTab } from '@/features/editions/ui/OverviewTab'
import { ActivitiesTab } from '@/features/activities/ui/ActivitiesTab'
import { CheckpointsTab } from '@/features/checkpoints/ui/CheckpointsTab'
import { TicketsTab } from '@/features/tickets/ui/TicketsTab'

type TabValue = 'overview' | 'activities' | 'checkpoints' | 'tickets'

export const Route = createFileRoute('/events/$eventId/editions/$editionId/')({
  component: EditionDetailPage,
})

function EditionDetailPage() {
  const { eventId, editionId } = Route.useParams()
  const [activeTab, setActiveTab] = useState<TabValue>('overview')

  const edition = mockEdition
  const activities = mockActivities
  const checkpoints = mockCheckpoints
  const tickets = mockTickets

  const tabCounts = {
    activities: activities.filter(a => a.status === 'published').length,
    checkpoints: checkpoints.length,
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
            {activeTab === 'checkpoints' && (
              <CheckpointsTab checkpoints={checkpoints} />
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