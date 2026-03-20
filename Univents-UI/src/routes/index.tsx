import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'
import { motion, AnimatePresence } from 'motion/react'
import { cn } from '@/shared/lib/utils'
import { ModeSelector } from '@/widgets/landing/ui/ModeSelector'
import { ParticipantView } from '@/widgets/landing/ui/ParticipantView'
import { OrganizerView } from '@/widgets/landing/ui/OrganizerView'
import { Footer } from '@/widgets/landing/ui/Footer'

const searchSchema = z.object({
  as: z.enum(['guest', 'host']).optional().default('guest'),
})

export const Route = createFileRoute('/')({
  component: Index,
  validateSearch: searchSchema,
})

export type Mode = 'guest' | 'host'

function Index() {
  const { as } = Route.useSearch()
  const navigate = Route.useNavigate()

  const setMode = (mode: Mode) => {
    void navigate({
      search: (prev) => ({ ...prev, as: mode }),
      replace: true,
    })
  }

  return (
    <div
      className={cn(
        "min-h-screen antialiased selection:bg-muted selection:text-foreground",
        "bg-background text-foreground"
      )}
    >
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="pt-24 pb-4 md:pt-12 md:pb-16">
          <div className="max-w-5xl mx-auto">
            <ModeSelector current={as} onChange={setMode} />
          </div>
        </div>

        <main className="pb-24 md:pb-32">
          <AnimatePresence mode="wait">
            <motion.div
              key={as}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3, ease: [0.25, 0.1, 0.25, 1] }}
            >
              {as === 'guest' ? <ParticipantView /> : <OrganizerView />}
            </motion.div>
          </AnimatePresence>
        </main>
      </div>

      <Footer />
    </div>
  )
}