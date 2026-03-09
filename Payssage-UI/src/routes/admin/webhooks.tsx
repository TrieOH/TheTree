import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute } from '@tanstack/react-router'
import { AdminLayout } from '#/features/admin/ui/admin-layout'
import { Webhook, Plus, Copy, Check } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import type { WebhookI } from '#/features/webhooks/model'
import * as React from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/webhooks' as any)({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
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
    <AdminLayout>
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
            <Card key={webhook.id} className="rounded-none border-border bg-card hover:border-primary/30 transition-colors">
              <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex items-start gap-4 flex-1">
                  <div className="w-10 h-10 bg-primary/10 text-primary flex items-center justify-center shrink-0">
                    <Webhook className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-1.5">
                      <Badge variant="outline" className="rounded-none text-[9px] uppercase font-black tracking-widest px-1.5 h-4 border-primary/30 text-primary">POST</Badge>
                      <button 
                        onClick={() => copyToClipboard(webhook.id)}
                        className="flex items-center gap-1 group/id"
                      >
                        <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-widest group-hover/id:text-primary transition-colors">ID: {webhook.id}</span>
                        {copiedId === webhook.id ? (
                          <Check className="w-3 h-3 text-emerald-500" />
                        ) : (
                          <Copy className="w-3 h-3 text-muted-foreground opacity-0 group-hover/id:opacity-100 transition-opacity" />
                        )}
                      </button>
                    </div>
                    <p className="font-mono text-xs md:text-sm font-bold truncate text-foreground/80 tracking-tight">
                      {webhook.url}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-4 text-xs text-muted-foreground font-medium uppercase tracking-wider shrink-0">
                  <span className="flex items-center gap-2">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                    Active
                  </span>
                  <span className="hidden sm:inline-block border-l border-border h-4" />
                  <span>
                    Created {new Date(webhook.created_at).toLocaleDateString('en-US')}
                  </span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </AdminLayout>
  )
}
