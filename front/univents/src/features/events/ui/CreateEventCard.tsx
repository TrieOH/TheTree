import { PlusCircle } from 'lucide-react'
import { motion } from 'motion/react'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

interface CreateEventCardProps {
  onClick: () => void
}

export function CreateEventCard({ onClick }: CreateEventCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: [0.25, 0.1, 0.25, 1] }}
      className={cn(
        'group relative min-w-0 overflow-hidden rounded-2xl bg-card',
        'border border-transparent transition-all duration-300 ease-out',
        'hover:-translate-y-1 hover:border-primary/30 hover:shadow-xl hover:shadow-primary/10',
      )}
    >
      <Button
        type="button"
        onClick={onClick}
        variant="ghost"
        size="none"
        className={cn(
          'flex h-full min-h-full w-full flex-col items-stretch justify-start overflow-hidden rounded-2xl border-0 bg-card p-0 text-left font-normal shadow-none whitespace-normal',
          'focus:outline-none focus-visible:outline-none focus-visible:ring-0',
          'hover:bg-card active:translate-y-0',
        )}
      >
        <div className="relative aspect-4/3 overflow-hidden bg-muted">
          <div className="absolute inset-0 bg-linear-to-br from-primary/10 via-background/10 to-muted" />

          <div className="absolute inset-0 flex items-center justify-center">
            <div className="flex size-16 items-center justify-center rounded-full border border-border/60 bg-background/85 text-primary shadow-lg backdrop-blur-sm transition-transform duration-300 group-hover:scale-105">
              <PlusCircle className="size-7" />
            </div>
          </div>
        </div>

        <div className="min-w-0 space-y-2 p-4 text-left md:p-5">
          <h3 className="text-balance text-base font-medium leading-snug text-foreground md:text-lg">
            Criar evento
          </h3>
          <p className="text-pretty text-sm text-muted-foreground">
            Clique para abrir o fluxo de criação de um novo evento.
          </p>
        </div>
      </Button>
    </motion.div>
  )
}
