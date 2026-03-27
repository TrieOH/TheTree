import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from '@/shared/ui/shadcn/drawer'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

interface AlertDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  onConfirm: () => void | Promise<void>
  variant?: 'default' | 'destructive' | 'success'
}

export function AlertDrawer({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel = 'Confirmar',
  cancelLabel = 'Cancelar',
  onConfirm,
  variant = 'default',
}: AlertDrawerProps) {
  const variantStyles = {
    default: 'bg-primary text-primary-foreground hover:bg-primary/90',
    destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
    success: 'bg-green-600 text-white hover:bg-green-700',
  }

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent className="z-60 rounded-t-2xl border-t border-border bg-card">
        <DrawerHeader className="pb-4 border-b border-border">
          <DrawerTitle className="text-base font-semibold text-left">{title}</DrawerTitle>
          {description && <p className="text-sm text-muted-foreground text-left mt-1">{description}</p>}
        </DrawerHeader>

        <div className="p-4 flex gap-2">
          <Button
            variant="secondary"
            onClick={() => { onOpenChange(false); }}
            className="flex-1 rounded-xl"
          >
            {cancelLabel}
          </Button>
          <Button
            onClick={() => { void onConfirm() }}
            className={cn("flex-1 rounded-xl", variantStyles[variant])}
          >
            {confirmLabel}
          </Button>
        </div>
      </DrawerContent>
    </Drawer>
  )
}