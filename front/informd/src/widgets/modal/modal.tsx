import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '#/shared/ui/shadcn/dialog'
import { Button } from '#/shared/ui/shadcn/button'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  children?: React.ReactNode
  footer?: React.ReactNode
}

export function Modal({
  isOpen,
  onClose,
  title,
  description,
  children,
  footer
}: ModalProps) {
  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="rounded-sm border-border sm:max-w-106.25 gap-0 max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-lg font-bold">
            {title}
          </DialogTitle>
          {description && (
            <DialogDescription className="text-sm text-muted-foreground">
              {description}
            </DialogDescription>
          )}
        </DialogHeader>
        <div className="py-4 px-2">
          {children}
        </div>
        {footer && (
          <DialogFooter>
            {footer}
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  )
}

interface ConfirmModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  description: string
  confirmText?: string
  variant?: 'default' | 'destructive'
  isLoading?: boolean
}

export function ConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  title,
  description,
  confirmText = 'Confirm',
  variant = 'default',
  isLoading
}: ConfirmModalProps) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      description={description}
      footer={
        <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 w-full">
          <Button
            variant="ghost"
            onClick={onClose}
            className="rounded-sm font-medium text-xs"
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button
            variant={variant}
            onClick={onConfirm}
            className="rounded-sm font-bold text-xs px-6"
            disabled={isLoading}
          >
            {isLoading ? 'Processing...' : confirmText}
          </Button>
        </div>
      }
    />
  )
}
