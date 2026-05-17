import { motion } from 'motion/react'
import { Ticket, ArrowRight } from 'lucide-react'
import type { TicketI } from '@/features/tickets/model'
import { Button } from '@/shared/ui/shadcn/button'

interface TicketsTabProps {
  tickets: TicketI[]
  eventId: string
  editionId: string
}

export function TicketsTab({ tickets }: TicketsTabProps) {
  return (
    <div className="space-y-6">
      {tickets.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          Nenhum ticket disponível.
        </div>
      ) : (
        <div className="grid gap-4">
          {tickets.map((ticket, idx) => (
            <motion.div
              key={ticket.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: idx * 0.05 }}
              className="bg-card border border-border rounded-xl p-4 md:p-5 flex flex-col md:flex-row md:items-center justify-between gap-4 hover:border-primary/20 transition-colors"
            >
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center shrink-0">
                  <Ticket className="w-6 h-6 text-primary" />
                </div>
                <div className="space-y-1">
                  <h4 className="text-lg font-semibold text-foreground">{ticket.name}</h4>
                  {ticket.description && (
                    <p className="text-sm text-muted-foreground">{ticket.description}</p>
                  )}
                </div>
              </div>

              <Button variant="outline" className="rounded-xl shrink-0">
                Ver Detalhes
                <ArrowRight className="w-4 h-4 ml-2" />
              </Button>
            </motion.div>
          ))}
        </div>
      )}
    </div>
  )
}