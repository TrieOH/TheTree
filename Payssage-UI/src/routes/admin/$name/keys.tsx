import { createFileRoute } from '@tanstack/react-router'
import { Key, Plus, Copy, Trash2, Check, ShieldOff } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import { ConfirmModal } from '#/widgets/modal/modal'
import { apiKeyCreateSchema } from '#/features/keys/model'
import * as React from 'react'
import { toast } from 'sonner'
import { cn } from '#/shared/lib/utils'
import FormModal from '#/widgets/modal/form-modal'
import type { ApiKeyCreateI, ApiKeyI } from '#/features/keys/model'

export const Route = createFileRoute('/admin/$name/keys')({
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
  const [keys, setKeys] = React.useState(MOCK_KEYS)
  const [copiedId, setCopiedId] = React.useState<string | null>(null)
  const [isCreateOpen, setIsCreateOpen] = React.useState(false)
  const [deleteKeyId, setDeleteKeyId] = React.useState<string | null>(null)
  const [revokeKeyId, setRevokeKeyId] = React.useState<string | null>(null)

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text)
    setCopiedId(id)
    toast.success('Copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  const onCreateSubmit = (data: ApiKeyCreateI) => {
    const newKey: ApiKeyI = {
      id: `key_${Math.random().toString(36).substr(2, 7)}`,
      name: data.name,
      prefix: 'tr_live_',
      created_at: new Date().toISOString(),
      revoked_at: null,
    }

    setKeys([newKey, ...keys])
    setIsCreateOpen(false)
    toast.success('API Key created successfully')
  }

  const handleDelete = () => {
    if (!deleteKeyId) return
    setKeys(keys.filter(k => k.id !== deleteKeyId))
    setDeleteKeyId(null)
    toast.success('API Key deleted successfully')
  }

  const handleRevoke = () => {
    if (!revokeKeyId) return
    setKeys(keys.map(k => k.id === revokeKeyId ? { ...k, revoked_at: new Date().toISOString() } : k))
    setRevokeKeyId(null)
    toast.success('API Key revoked successfully')
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">API Keys</h2>
          <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">Programmatic access for your workspace.</p>
        </div>

        <Button
          onClick={() => setIsCreateOpen(true)}
          className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all"
        >
          <Plus className="w-4 h-4" />
          New Key
        </Button>
      </div>

      {/* List of Keys */}
      <div className="grid gap-4">
        {keys.map((key) => (
          <Card key={key.id} className="rounded-none border-border bg-card transition-colors">
            <CardContent className="p-4 md:p-6 flex flex-col md:flex-row md:items-center justify-between gap-6">
              <div className="flex items-start gap-3 sm:gap-4 min-w-0 flex-1 w-full">
                <div className="w-10 h-10 sm:w-12 sm:h-12 bg-primary/10 text-primary flex items-center justify-center shrink-0 border border-primary/20">
                  <Key className="w-5 h-5 sm:w-6 sm:h-6" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex flex-wrap items-center gap-2 sm:gap-3 mb-2">
                    <Badge variant="outline" className="rounded-none text-[8px] sm:text-[9px] uppercase font-black tracking-widest px-1.5 h-4 border-primary/30 text-primary">API KEY</Badge>
                    <button
                      onClick={() => copyToClipboard(key.id, key.id)}
                      className="flex items-center gap-1 group/id min-w-0"
                    >
                      <span className="text-[8px] sm:text-[10px] text-muted-foreground font-mono uppercase tracking-widest group-hover/id:text-primary transition-colors truncate">ID: {key.id}</span>
                      {copiedId === key.id ? (
                        <Check className="w-3 h-3 text-emerald-500 shrink-0" />
                      ) : (
                        <Copy className="w-3 h-3 text-muted-foreground opacity-0 group-hover/id:opacity-100 transition-opacity shrink-0" />
                      )}
                    </button>
                    <span className="hidden xs:inline-block border-l border-border h-3" />
                    <span className="text-[9px] font-black uppercase tracking-tight text-muted-foreground/60 truncate">{key.name}</span>
                  </div>
                  <p className="hidden md:block font-mono text-[10px] sm:text-xs md:text-sm font-bold truncate text-foreground/80 tracking-tight">
                    {key.prefix}••••••••••••
                  </p>
                </div>
              </div>

              <div className="flex flex-col items-center md:items-end gap-4 shrink-0 w-full md:w-auto">
                <div className="flex flex-wrap items-center justify-center md:justify-end gap-3 sm:gap-4 text-[8px] sm:text-xs text-muted-foreground font-medium uppercase tracking-wider w-full">
                  <span className="flex items-center gap-1.5">
                    <span className={cn("w-1.5 h-1.5 rounded-full", key.revoked_at ? "bg-destructive" : "bg-emerald-500")} />
                    {key.revoked_at ? "Revoked" : "Active"}
                  </span>
                  <span className="hidden xs:inline-block border-l border-border h-4" />
                  <span className="shrink-0">
                    {new Date(key.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
                  </span>
                </div>

                <div className="flex items-center justify-center md:justify-end gap-2 w-full">
                  <Button
                    variant="outline"
                    size="sm"
                    className="rounded-none h-8 gap-2 text-[9px] sm:text-[10px] uppercase font-black tracking-[0.2em] transition-all px-3"
                    onClick={() => copyToClipboard(`${key.prefix}SECRET_KEY`, key.id)}
                  >
                    {copiedId === key.id ? <Check className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-emerald-500" /> : <Copy className="w-3 sm:w-3.5 h-3 sm:h-3.5" />}
                    <span className="hidden xs:inline">{copiedId === key.id ? 'Copied' : 'Copy Key'}</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="rounded-none h-8 gap-2 text-[9px] sm:text-[10px] uppercase font-black tracking-[0.2em] transition-all px-3 border-destructive/30 text-destructive hover:bg-destructive/5 disabled:opacity-30"
                    onClick={() => setRevokeKeyId(key.id)}
                    disabled={!!key.revoked_at}
                  >
                    <ShieldOff className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
                    <span className="hidden xs:inline">Revoke</span>
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-none h-8 w-8 text-muted-foreground/30 cursor-not-allowed"
                    disabled
                  >
                    <Trash2 className="w-3.5 sm:w-4 h-3.5 sm:h-4" />
                  </Button>
                </div>

                <div className="md:hidden w-full pt-2 border-t border-border/50 mt-1">
                  <p className="font-mono text-[10px] font-bold text-center text-foreground/60 tracking-widest">
                    {key.prefix}••••••••••••
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Create Modal */}
      <FormModal<ApiKeyCreateI>
        title="Create API Key"
        description="Give your key a name to identify it later."
        buttonTitle="Generate Key"
        schema={apiKeyCreateSchema}
        formId="create-key-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={onCreateSubmit}
        fields={[
          {
            name: "name",
            label: "e.g. Production Mobile App",
            type: "text",
          }
        ]}
      />

      {/* Revoke Confirmation Modal */}
      <ConfirmModal
        isOpen={!!revokeKeyId}
        onClose={() => setRevokeKeyId(null)}
        onConfirm={handleRevoke}
        title="Revoke API Key"
        description="Are you sure you want to revoke this API key? This action will immediately invalidate the key and cannot be undone."
        confirmText="Revoke Key"
        variant="destructive"
      />

      {/* Delete Confirmation Modal (Kept but unreachable as button and functionally now is disabled/unimplemented) */}
      <ConfirmModal
        isOpen={!!deleteKeyId}
        onClose={() => setDeleteKeyId(null)}
        onConfirm={handleDelete}
        title="Delete API Key"
        description="Are you sure you want to delete this API key? This action cannot be undone and will immediately revoke access."
        confirmText="Delete Key"
        variant="destructive"
      />
    </div>
  )
}