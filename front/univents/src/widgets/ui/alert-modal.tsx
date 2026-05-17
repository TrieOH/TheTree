import { Loader2 } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/shadcn/dialog'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

interface AlertModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  onConfirm: () => void | Promise<void>
  variant?: 'default' | 'destructive' | 'success'
  loading?: boolean
}

export function AlertModal({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel = 'Confirmar',
  cancelLabel = 'Cancelar',
  onConfirm,
  variant = 'default',
  loading = false,
}: AlertModalProps) {
  const variantStyles = {
    default: 'bg-primary text-primary-foreground hover:bg-primary/90',
    destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
    success: 'bg-green-600 text-white hover:bg-green-700',
  }

  const handleConfirm = async () => {
    await onConfirm()
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md z-50">
        <DialogHeader className="pb-4 border-b border-border">
          <DialogTitle className="text-base font-semibold text-left">{title}</DialogTitle>
          {description && <p className="text-sm text-muted-foreground text-left mt-1">{description}</p>}
        </DialogHeader>

        <div className="p-4 flex gap-2">
          <Button
            variant="secondary"
            onClick={() => { onOpenChange(false); }}
            disabled={loading}
            className="flex-1 rounded-xl"
          >
            {cancelLabel}
          </Button>
          <Button
            onClick={() => { void handleConfirm() }}
            disabled={loading}
            className={cn("flex-1 rounded-xl", variantStyles[variant])}
          >
            {loading ? (
              <span className="flex items-center gap-2">
                <Loader2 className="w-4 h-4 animate-spin" />
                Processando...
              </span>
            ) : (
              confirmLabel
            )}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}