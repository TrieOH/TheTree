import { Check, Copy, AlertTriangle } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '#/shared/ui/shadcn/dialog'
import { useState } from 'react'
import { toast } from 'sonner'
import type { WebhookCreateResponseI } from '../model'

interface WebhookCreatedModalProps {
  webhook: WebhookCreateResponseI | null
  isOpen: boolean
  onClose: () => void
}

export function WebhookCreatedModal({ webhook, isOpen, onClose }: WebhookCreatedModalProps) {
  const [copied, setCopied] = useState(false)

  const copyToClipboard = () => {
    if (!webhook) return
    navigator.clipboard.writeText(webhook.secret)
    setCopied(true)
    toast.success('Webhook Secret copied to clipboard')
    setTimeout(() => setCopied(false), 2000)
  }

  if (!webhook) return null

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md rounded-none border-2 border-primary/20">
        <DialogHeader>
          <DialogTitle className="text-xl font-black uppercase tracking-tighter flex items-center gap-2">
            <Check className="w-6 h-6 text-emerald-500" />
            Webhook Created
          </DialogTitle>
          <DialogDescription className="text-xs uppercase tracking-widest font-bold opacity-70">
            Please copy your webhook secret now. You won't be able to see it again.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 py-4">
          <div className="relative group">
            <div className="absolute -inset-1 bg-linear-to-r from-primary/20 to-primary/10 blur opacity-25 group-hover:opacity-50 transition duration-1000"></div>
            <div className="relative flex items-center gap-2 bg-muted/50 p-4 border border-border font-mono text-sm break-all">
              <code className="flex-1 text-primary font-bold">{webhook.secret}</code>
              <Button
                size="icon"
                variant="ghost"
                onClick={copyToClipboard}
                className="shrink-0 hover:bg-primary/10 rounded-none"
              >
                {copied ? <Check className="w-4 h-4 text-emerald-500" /> : <Copy className="w-4 h-4" />}
              </Button>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 bg-amber-500/10 border border-amber-500/20 rounded-none">
            <AlertTriangle className="w-5 h-5 text-amber-500 shrink-0 mt-0.5" />
            <div className="space-y-1">
              <p className="text-[10px] font-black uppercase tracking-widest text-amber-600">Security Warning</p>
              <p className="text-[10px] font-bold text-amber-800/80 leading-relaxed uppercase tracking-tight">
                Store this secret securely. It is used to verify that the webhooks are indeed sent from our platform.
              </p>
            </div>
          </div>
        </div>

        <div className="flex justify-end">
          <Button
            onClick={onClose}
            className="rounded-none font-black uppercase tracking-widest w-full sm:w-auto"
          >
            I've saved the secret
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
