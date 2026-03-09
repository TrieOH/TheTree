import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute } from '@tanstack/react-router'
import { AdminLayout } from '#/features/admin/ui/admin-layout'
import { Key, Plus, Copy, Trash2, Check } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import type { ApiKeyI } from '#/features/keys/model'
import * as React from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/keys' as any)({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
  component: RouteComponent,
})

const MOCK_KEYS: ApiKeyI[] = [
  {
    id: 'key_01jhc83',
    name: 'Production Mobile App',
    prefix: 'tr_live_',
    created_at: '2025-01-15T10:00:00Z',
    revoked_at: null,
  },
  {
    id: 'key_01jhc84',
    name: 'Local Testing',
    prefix: 'tr_test_',
    created_at: '2025-02-20T14:30:00Z',
    revoked_at: '2025-03-01T09:15:00Z',
  }
]

function RouteComponent() {
  const [copiedId, setCopiedId] = React.useState<string | null>(null)

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text)
    setCopiedId(id)
    toast.success('Copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  return (
    <AdminLayout>
      <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div className="space-y-1">
            <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">API Keys</h2>
            <p className="text-muted-foreground text-sm">Manage programmatic access to your workspaces.</p>
          </div>
          <Button className="rounded-none gap-2 h-10">
            <Plus className="w-4 h-4" />
            New Key
          </Button>
        </div>

        <div className="grid gap-4">
          {MOCK_KEYS.map((key) => (
            <Card key={key.id} className="rounded-none border-border bg-card">
              <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex items-start gap-4">
                  <div className="w-10 h-10 bg-primary/10 text-primary flex items-center justify-center shrink-0">
                    <Key className="w-5 h-5" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-bold text-base uppercase tracking-tight">{key.name}</h3>
                      {key.revoked_at ? (
                        <Badge variant="destructive" className="rounded-none text-[10px] uppercase font-black tracking-widest px-1.5 h-4.5">Revoked</Badge>
                      ) : (
                        <Badge variant="outline" className="rounded-none text-[10px] uppercase font-black tracking-widest px-1.5 h-4.5 border-emerald-500/50 text-emerald-600">Active</Badge>
                      )}
                    </div>
                    <code className="text-xs bg-muted px-2 py-1 rounded-none border border-border/50 font-mono text-muted-foreground uppercase tracking-widest">
                      {key.prefix}••••••••••••
                    </code>
                  </div>
                </div>

                <div className="flex items-center gap-2 self-end md:self-center">
                  <Button 
                    variant="outline" 
                    size="sm" 
                    className="rounded-none gap-2 text-xs uppercase font-bold tracking-widest"
                    onClick={() => copyToClipboard(`${key.prefix}SECRET_KEY`, key.id)}
                  >
                    {copiedId === key.id ? <Check className="w-3.5 h-3.5 text-emerald-500" /> : <Copy className="w-3.5 h-3.5" />}
                    {copiedId === key.id ? 'Copied' : 'Copy'}
                  </Button>
                  <Button variant="ghost" size="icon" className="rounded-none h-8 w-8 text-destructive hover:bg-destructive/5">
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </AdminLayout>
  )
}
