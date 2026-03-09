import { createFileRoute } from '@tanstack/react-router'
import { Webhook, Plus, Copy, Check } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import type { WebhookI } from '#/features/webhooks/model'
import * as React from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/$name/webhooks')({
  component: RouteComponent,
})

const MOCK_WEBHOOKS: WebhookI[] = [
  {
    id: 'wh_01jhc83',
    workspace_id: 'ws_prod_01jhc83',
    url: 'https://api.my-app.com/v1/webhooks/trie',
    created_at: '2025-01-15T10:00:00Z',
  },
  {
    id: 'wh_01jhc84',
    workspace_id: 'ws_prod_01jhc83',
    url: 'https://api.test-env.io/hooks',
    created_at: '2025-02-20T14:30:00Z',
  }
]

function RouteComponent() {
  const [copiedId, setCopiedId] = React.useState<string | null>(null)

  const copyToClipboard = (id: string) => {
    navigator.clipboard.writeText(id)
    setCopiedId(id)
    toast.success('ID copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Webhooks</h2>
          <p className="text-muted-foreground text-sm">Receive real-time notifications about payment events.</p>
        </div>
        <Button className="rounded-none gap-2 h-10">
          <Plus className="w-4 h-4" />
          Add Endpoint
        </Button>
      </div>

      <div className="grid gap-4">
        {MOCK_WEBHOOKS.map((webhook) => (
          <Card key={webhook.id} className="rounded-sm border-border bg-card hover:border-primary/30 transition-colors">
            <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-6">
              <div className="flex items-start gap-3 sm:gap-4 min-w-0 flex-1">
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

              <div className="flex flex-wrap items-center gap-3 sm:gap-4 text-[8px] sm:text-xs text-muted-foreground font-medium uppercase tracking-wider shrink-0 justify-end">
                <span className="flex items-center gap-1.5">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  Active
                </span>
                <span className="hidden xs:inline-block border-l border-border h-4" />
                <span className="shrink-0">
                  {new Date(webhook.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
                </span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
