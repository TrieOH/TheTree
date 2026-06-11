import { Check, Copy, AlertTriangle } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle
} from '#/shared/ui/shadcn/dialog'
import { useState } from 'react'
import { toast } from 'sonner'
import type { ApiKeyCreateResponseI } from '../model'

interface ApiKeyCreatedModalProps {
  apiKey: ApiKeyCreateResponseI | null
  isOpen: boolean
  onClose: () => void
}

export function ApiKeyCreatedModal({ apiKey, isOpen, onClose }: ApiKeyCreatedModalProps) {
  const [copied, setCopied] = useState(false)

  const copyToClipboard = () => {
    if (!apiKey) return
    navigator.clipboard.writeText(apiKey.key)
    setCopied(true)
    toast.success('API Key copied to clipboard')
    setTimeout(() => setCopied(false), 2000)
  }

  if (!apiKey) return null

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md rounded-sm border border-foreground/10 shadow-xs bg-card">
        <DialogHeader>
          <DialogTitle className="text-base font-semibold flex items-center gap-2">
            <Check className="size-4 text-emerald-500" />
            API Key Generated
          </DialogTitle>
          <DialogDescription className="text-xs text-muted-foreground">
            Please copy your API key now. You won't be able to see it again.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-3 py-2">
          {/* Key display */}
          <div className="flex items-center gap-2 bg-muted px-3 py-2.5 rounded-sm ring-1 ring-foreground/10">
            <code className="flex-1 text-xs font-mono text-foreground break-all">{apiKey.key}</code>
            <Button
              size="icon"
              variant="ghost"
              onClick={copyToClipboard}
              className="shrink-0 size-7 hover:bg-foreground/10 text-muted-foreground hover:text-foreground duration-150 cursor-pointer outline-0"
            >
              {copied
                ? <Check className="size-3.5 text-emerald-500" />
                : <Copy className="size-3.5" />
              }
            </Button>
          </div>

          {/* Warning */}
          <div className="flex items-start gap-2.5 px-3 py-2.5 bg-amber-500/10 ring-1 ring-amber-500/20 rounded-sm">
            <AlertTriangle className="size-3.5 text-amber-500 shrink-0 mt-0.5" />
            <p className="text-xs text-amber-700 dark:text-amber-400/80 leading-relaxed">
              Store this key securely. Anyone with access can perform actions on behalf of your project.
            </p>
          </div>
        </div>

        <div className="flex justify-end pt-1">
          <Button
            onClick={onClose}
            size="sm"
            className="font-medium cursor-pointer"
          >
            I've saved the key
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}