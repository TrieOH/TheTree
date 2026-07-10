import { useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'
import { Settings } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

export interface FABMenuItem {
  id: string
  label: string
  icon: LucideIcon
  onClick: () => void
  danger?: boolean
}

interface FABMenuBaseProps {
  items: FABMenuItem[]
  className?: string
  ariaLabel?: string
}

interface FABMenuMenuProps extends FABMenuBaseProps {
  mode?: 'menu'
}

interface FABMenuActionProps {
  mode: 'action'
  className?: string
  ariaLabel?: string
  icon: LucideIcon
  onClick: () => void
  active?: boolean
}

type FABMenuProps = FABMenuMenuProps | FABMenuActionProps

export function FABMenu(props: FABMenuProps) {
  const [open, setOpen] = useState(false)

  if (props.mode === 'action') {
    const {
      className,
      ariaLabel = 'Ação flutuante',
      icon: Icon,
      onClick,
      active,
    } = props

    return (
      <div className={cn('fixed bottom-20 right-4 z-40 md:bottom-8 md:right-8', className)}>
        <Button
          type="button"
          onClick={onClick}
          aria-label={ariaLabel}
          className={cn(
            'size-14 rounded-full border shadow-xl backdrop-blur-xl',
            active
              ? 'bg-primary text-primary-foreground border-primary/40'
              : 'bg-background text-foreground border-border/60 hover:bg-muted',
          )}
          variant="default"
        >
          <Icon className="size-5" />
        </Button>
      </div>
    )
  }

  const {
    items,
    className,
    ariaLabel = 'Abrir menu flutuante',
  } = props

  return (
    <div className={cn('fixed bottom-20 right-4 z-40 md:bottom-8 md:right-8', className)}>
      <div className="relative flex flex-col items-end gap-2">
        <AnimatePresence>
          {open && items.map((item, idx) => {
            const Icon = item.icon

            return (
              <motion.div
                key={item.id}
                initial={{ opacity: 0, y: 12, scale: 0.9 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                exit={{ opacity: 0, y: 12, scale: 0.9 }}
                transition={{ delay: idx * 0.04, type: 'spring', stiffness: 260, damping: 22 }}
                className="flex items-center gap-2"
              >
                <span className="rounded-full border border-border/60 bg-background/90 px-3 py-1 text-xs text-muted-foreground shadow-lg backdrop-blur-xl">
                  {item.label}
                </span>
                <Button
                  type="button"
                  onClick={() => {
                    item.onClick()
                    setOpen(false)
                  }}
                  className={cn(
                    'size-11 rounded-full shadow-lg backdrop-blur-xl',
                    item.danger
                      ? 'bg-destructive text-destructive-foreground hover:bg-destructive/90'
                      : 'bg-background text-foreground hover:bg-muted',
                  )}
                  variant="default"
                >
                  <Icon className="size-5" />
                </Button>
              </motion.div>
            )
          })}
        </AnimatePresence>

        <Button
          type="button"
          onClick={() => { setOpen((v) => !v) }}
          aria-label={ariaLabel}
          className={cn(
            'size-14 rounded-full border shadow-xl backdrop-blur-xl',
            open
              ? 'bg-primary text-primary-foreground border-primary/40'
              : 'bg-background text-foreground border-border/60 hover:bg-muted',
          )}
          variant="default"
        >
          <motion.span
            animate={{ rotate: open ? 45 : 0 }}
            transition={{ type: 'spring', stiffness: 260, damping: 20 }}
            className="flex items-center justify-center"
          >
            <Settings className="size-5" />
          </motion.span>
        </Button>
      </div>
    </div>
  )
}
