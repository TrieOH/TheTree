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
      <DialogContent className="rounded-none border-border sm:max-w-106.25">
        <DialogHeader>
          <DialogTitle className="text-xl font-black uppercase tracking-tighter">
            {title}
          </DialogTitle>
          {description && (
            <DialogDescription className="text-xs font-bold uppercase tracking-widest opacity-70">
              {description}
            </DialogDescription>
          )}
        </DialogHeader>
        <div className="py-4">
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
            className="rounded-none font-bold uppercase tracking-widest text-xs"
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button
            variant={variant}
            onClick={onConfirm}
            className="rounded-none font-black uppercase tracking-widest text-xs px-6"
            disabled={isLoading}
          >
            {isLoading ? 'Processing...' : confirmText}
          </Button>
        </div>
      }
    />
  )
}
