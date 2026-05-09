import { Webhook, Copy, Check, Trash2 } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import type { WebhookI } from '../model'
import { useState } from 'react'
import { toast } from 'sonner'

interface WebhookListProps {
  webhooks: WebhookI[]
  isLoading: boolean
  onDelete: (id: string) => void
}

export function WebhookList({ webhooks, isLoading, onDelete }: WebhookListProps) {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const copyToClipboard = (id: string) => {
    navigator.clipboard.writeText(id)
    setCopiedId(id)
    toast.success('ID copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  if (isLoading) {
    return <WebhookListSkeleton />
  }

  if (webhooks.length === 0) {
    return <WebhookListEmpty />
  }

  return (
    <div className="grid gap-4">
      {webhooks.map((webhook) => (
        <Card key={webhook.id} className="rounded-none border-border bg-card transition-colors">
          <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-6">
            <div className="flex items-start gap-3 sm:gap-4 min-w-0 flex-1 w-full">
              <div className="w-10 h-10 sm:w-12 sm:h-12 bg-primary/10 text-primary flex items-center justify-center shrink-0 border border-primary/20">
                <Webhook className="w-5 h-5 sm:w-6 sm:h-6" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex flex-wrap items-center gap-2 sm:gap-3 mb-2">
                  <Badge variant="outline" className="rounded-none text-[8px] sm:text-[9px] uppercase font-black tracking-widest px-1.5 h-4 border-primary/30 text-primary">POST</Badge>
                  <button
                    onClick={() => copyToClipboard(webhook.id)}
                    className="flex items-center gap-1 group/id min-w-0"
                  >
                    <span className="text-[8px] sm:text-[10px] text-muted-foreground font-mono uppercase tracking-widest group-hover/id:text-primary transition-colors truncate">ID: {webhook.id}</span>
                    {copiedId === webhook.id ? (
                      <Check className="w-3 h-3 text-emerald-500 shrink-0" />
                    ) : (
                      <Copy className="w-3 h-3 text-muted-foreground opacity-0 group-hover/id:opacity-100 transition-opacity shrink-0" />
                    )}
                  </button>
                </div>
                <p className="font-mono text-[10px] sm:text-xs md:text-sm font-bold truncate text-foreground/80 tracking-tight">
                  {webhook.url}
                </p>
              </div>
            </div>

            <div className="flex flex-col items-center md:items-end gap-4 shrink-0 w-full md:w-auto">
              <div className="flex flex-wrap items-center justify-center md:justify-end gap-3 sm:gap-4 text-[8px] sm:text-xs text-muted-foreground font-medium uppercase tracking-wider w-full">
                <span className="flex items-center gap-1.5">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  Active
                </span>
                <span className="hidden xs:inline-block border-l border-border h-4" />
                <span className="shrink-0">
                  {new Date(webhook.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
                </span>
              </div>

              <div className="flex items-center justify-center md:justify-end gap-2 w-full">
                <Button
                  variant="outline"
                  size="sm"
                  className="rounded-none h-8 gap-2 text-[9px] sm:text-[10px] uppercase font-black tracking-[0.2em] transition-all px-3 border-destructive/30 text-destructive hover:bg-destructive/5"
                  onClick={() => onDelete(webhook.id)}
                >
                  <Trash2 className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
                  <span className="hidden xs:inline">Delete</span>
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function WebhookListSkeleton() {
  return (
    <div className="grid gap-4">
      {[1, 2].map((i) => (
        <Card key={i} className="rounded-none border-border bg-card animate-pulse">
          <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-6">
            <div className="flex items-start gap-4 flex-1 w-full">
              <div className="w-12 h-12 bg-muted shrink-0" />
              <div className="flex-1 space-y-2">
                <div className="h-4 bg-muted w-24" />
                <div className="h-4 bg-muted w-64" />
              </div>
            </div>
            <div className="flex flex-col items-end gap-4 w-full md:w-auto">
              <div className="h-4 bg-muted w-32" />
              <div className="h-8 bg-muted w-32" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function WebhookListEmpty() {
  return (
    <div className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-border rounded-none bg-muted/5 space-y-4">
      <div className="w-16 h-16 bg-muted/10 flex items-center justify-center rounded-full border border-border">
        <Webhook className="w-8 h-8 text-muted-foreground/50" />
      </div>
      <div className="text-center space-y-1">
        <h3 className="text-lg font-black uppercase tracking-tighter">No Webhooks found</h3>
        <p className="text-muted-foreground text-xs uppercase tracking-widest font-bold max-w-70">
          Add your first webhook endpoint to receive real-time notifications about payment events.
        </p>
      </div>
    </div>
  )
}
