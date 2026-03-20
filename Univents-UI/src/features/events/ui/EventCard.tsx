import { motion } from 'motion/react'
import type { EventI } from '../model'

interface EventCardProps {
  event: EventI
  index?: number
}

export function EventCard({ event, index = 0 }: EventCardProps) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' }).replace('.', '')
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.1, duration: 0.4 }}
      className="group cursor-pointer"
    >
      {/* Placeholder de imagem */}
      <div className="aspect-4/3 bg-muted rounded-xl md:rounded-2xl mb-3 md:mb-4 overflow-hidden relative">
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-12 h-12 md:w-16 md:h-16 rounded-full border-2 border-dashed border-border flex items-center justify-center">
            <svg
              className="w-5 h-5 md:w-6 md:h-6 text-muted-foreground"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
        <div className="absolute inset-0 bg-accent/0 group-hover:bg-accent/10 transition-colors duration-300" />
      </div>

      {/* Info */}
      <div className="space-y-1 md:space-y-2">
        <h3 className="text-sm md:text-base font-medium leading-tight text-foreground group-hover:text-muted-foreground transition-colors line-clamp-2">
          {event.name}
        </h3>
        <p className="text-xs md:text-sm text-muted-foreground">
          {formatDate(event.created_at)}
        </p>
        <p className="text-xs text-muted-foreground/70">Ingresso digital</p>
      </div>
    </motion.div>
  )
}