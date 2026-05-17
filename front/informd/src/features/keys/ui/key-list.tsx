import { Key, Copy, Check, ShieldOff } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import { cn } from '#/shared/lib/utils'
import type { ApiKeyI } from '../model'
import { useState } from 'react'
import { toast } from 'sonner'

interface KeyListProps {
  keys: ApiKeyI[]
  isLoading: boolean
  onRevoke: (id: string) => void
}

export function KeyList({ keys, isLoading, onRevoke }: KeyListProps) {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text)
    setCopiedId(id)
    toast.success('Copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  if (isLoading) return <KeyListSkeleton />
  if (keys.length === 0) return <KeyListEmpty />

  return (
    <div className="grid gap-6">
      {keys.map((key) => (
        <Card key={key.id} className="rounded-none border-2 border-border bg-card hover:border-primary hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-all hover:-translate-x-1 hover:-translate-y-1">
          <CardContent className="p-6 flex flex-col md:flex-row md:items-center justify-between gap-8">
            <div className="flex items-start gap-4 min-w-0 flex-1 w-full">
              <div className="w-12 h-12 bg-primary text-primary-foreground flex items-center justify-center shrink-0 border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                <Key className="w-6 h-6" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex flex-wrap items-center gap-3 mb-3">
                  <Badge variant="outline" className="rounded-none text-[9px] uppercase font-black tracking-widest px-2 h-5 border-2 border-primary/30 text-primary bg-primary/5">API KEY</Badge>
                  <button
                    onClick={() => copyToClipboard(key.id, key.id)}
                    className="flex items-center gap-1.5 group/id min-w-0 bg-muted/30 px-2 py-0.5"
                  >
                    <span className="text-[10px] text-muted-foreground font-mono font-bold uppercase tracking-widest group-hover/id:text-primary transition-colors truncate">ID: {key.id}</span>
                    {copiedId === key.id ? (
                      <Check className="w-3.5 h-3.5 text-emerald-500 shrink-0" />
                    ) : (
                      <Copy className="w-3.5 h-3.5 text-muted-foreground opacity-50 group-hover/id:opacity-100 transition-opacity shrink-0" />
                    )}
                  </button>
                  <span className="text-[10px] font-black uppercase tracking-tight text-foreground/80 truncate">{key.name}</span>
                </div>
                <p className="hidden md:block font-mono text-sm font-bold truncate text-foreground/70 tracking-tight bg-muted/20 p-2 border-l-4 border-primary">
                  {key.prefix}••••••••••••••••••••••••••••••••
                </p>
              </div>
            </div>

            <div className="flex flex-col items-center md:items-end gap-5 shrink-0 w-full md:w-auto">
              <div className="flex flex-wrap items-center justify-center md:justify-end gap-4 text-[10px] text-muted-foreground font-black uppercase tracking-widest w-full">
                <span className="flex items-center gap-2">
                  <span className={cn("w-2 h-2 rounded-none border border-black/20 shadow-[2px_2px_0px_0px_rgba(0,0,0,0.1)]", key.revoked_at ? "bg-destructive" : "bg-emerald-500")} />
                  {key.revoked_at ? "Revoked" : "Active"}
                </span>
                <span className="hidden xs:inline-block border-l-2 border-border h-4" />
                <span className="shrink-0">
                  Created: {new Date(key.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
                </span>
              </div>

              <div className="flex items-center justify-center md:justify-end gap-3 w-full">
                <Button
                  variant="outline"
                  size="sm"
                  className="rounded-none h-10 gap-2 text-[10px] uppercase font-black tracking-widest transition-all px-4 border-2 hover:bg-primary hover:text-primary-foreground group"
                  onClick={() => copyToClipboard(`${key.prefix}SECRET_KEY`, key.id)}
                  disabled={!!key.revoked_at}
                >
                  {copiedId === key.id ? <Check className="w-4 h-4 text-emerald-500" /> : <Copy className="w-4 h-4 group-hover:scale-110 transition-transform" />}
                  <span>{copiedId === key.id ? 'Copied' : 'Copy Key'}</span>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  className="rounded-none h-10 gap-2 text-[10px] uppercase font-black tracking-widest transition-all px-4 border-2 border-destructive/30 text-destructive hover:bg-destructive hover:text-destructive-foreground disabled:opacity-30 group"
                  onClick={() => onRevoke(key.id)}
                  disabled={!!key.revoked_at}
                >
                  <ShieldOff className="w-4 h-4 group-hover:rotate-12 transition-transform" />
                  <span>Revoke</span>
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function KeyListSkeleton() {
  return (
    <div className="grid gap-6">
      {[1, 2, 3].map((i) => (
        <Card key={i} className="rounded-none border-2 border-border/50 bg-card animate-pulse">
          <CardContent className="p-6 flex flex-col md:flex-row md:items-center justify-between gap-8">
            <div className="flex items-start gap-4 flex-1 w-full">
              <div className="w-12 h-12 bg-muted shrink-0 border-2 border-border/20" />
              <div className="flex-1 space-y-3">
                <div className="h-4 bg-muted w-24" />
                <div className="h-4 bg-muted w-full max-w-[300px]" />
              </div>
            </div>
            <div className="flex flex-col items-end gap-5 w-full md:w-auto">
              <div className="h-4 bg-muted w-32" />
              <div className="h-10 bg-muted w-48" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function KeyListEmpty() {
  return (
    <Card className="rounded-none border-4 border-dashed border-border flex flex-col items-center justify-center py-20 px-6 text-center bg-muted/20 relative overflow-hidden">
      <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 -mr-16 -mt-16 rotate-12" />
      
      <div className="w-20 h-20 bg-primary text-primary-foreground flex items-center justify-center border-4 border-primary shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] mb-8">
        <Key className="w-10 h-10" />
      </div>
      <div className="max-w-md">
        <h3 className="text-3xl font-black uppercase tracking-tighter mb-4">No API Keys found</h3>
        <p className="text-muted-foreground text-sm uppercase tracking-widest font-bold opacity-60">
          Generate your first API key to start integrating our form services into your own applications.
        </p>
      </div>
    </Card>
  )
}
