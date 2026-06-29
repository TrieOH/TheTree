import { createFileRoute } from '@tanstack/react-router'
import { useAuth, ModernIntrospect } from '@trieoh/identityx-sdk-ts/react'
import { useQuery } from '@tanstack/react-query'
import {
  AlertCircle,
  RefreshCw,
  Copy,
  Check,
  ShieldCheck,
  Terminal,
} from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useState, useCallback } from 'react'

export const Route = createFileRoute('/admin/introspect')({
  component: RouteComponent,
})

function CopyButton({ json }: { json: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(json).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }, [json])

  return (
    <button
      type='button'
      onClick={handleCopy}
      className={cn(
        'flex items-center gap-1.5 text-[11px] font-medium px-2.5 py-1.5 rounded-lg border transition-all',
        copied
          ? 'border-green-400/50 bg-green-50 dark:bg-green-950/20 text-green-700 dark:text-green-400'
          : 'border-border bg-muted/50 text-muted-foreground hover:text-foreground hover:bg-muted',
      )}
    >
      {copied ? (
        <><Check className="w-3 h-3" /> Copied!</>
      ) : (
        <><Copy className="w-3 h-3" /> Copy JSON</>
      )}
    </button>
  )
}

function RouteComponent() {
  const { auth } = useAuth()
  const [showRaw, setShowRaw] = useState(false)

  const { data, isLoading, isError, error, refetch, isFetching } = useQuery({
    queryKey: ['introspect'],
    queryFn: () => auth.introspect(),
    retry: 1,
  })

  const introspectResponse = data?.success ? data.data : undefined

  const capKeys = introspectResponse?.sub.capabilities
    ? Object.keys(introspectResponse.sub.capabilities)
    : []

  return (
    <div className="min-w-75 w-full max-w-5xl mx-auto flex flex-col gap-6 px-4 py-6 sm:px-6 lg:px-10 lg:py-10">
      {/* Header */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
            <ShieldCheck className="w-4 h-4 text-primary" />
          </div>
          <div>
            <h1 className="text-base font-semibold text-foreground leading-tight">Introspect</h1>
            <p className="text-[12px] text-muted-foreground">Inspect your current session token</p>
          </div>
        </div>

        {introspectResponse && (
          <div className="flex items-center gap-2 flex-wrap">
            <CopyButton json={JSON.stringify(introspectResponse, null, 2)} />
            <button
              type='button'
              onClick={() => setShowRaw(!showRaw)}
              className={cn(
                'flex items-center gap-1.5 text-[11px] font-medium px-2.5 py-1.5 rounded-lg border transition-all',
                showRaw
                  ? 'border-primary/40 bg-primary/5 text-primary'
                  : 'border-border bg-muted/50 text-muted-foreground hover:text-foreground hover:bg-muted',
              )}
            >
              <Terminal className="w-3 h-3" /> Raw
            </button>
            <button
              type='button'
              onClick={() => refetch()}
              disabled={isFetching}
              className="flex items-center gap-1.5 text-[11px] font-medium px-2.5 py-1.5 rounded-lg border border-border bg-muted/50 text-muted-foreground hover:text-foreground hover:bg-muted transition-all disabled:opacity-50"
            >
              <RefreshCw className={cn('w-3 h-3', isFetching && 'animate-spin')} />
              Refresh
            </button>
          </div>
        )}
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="flex flex-col items-center justify-center py-24 gap-4">
          <div className="relative w-12 h-12">
            <div className="absolute inset-0 rounded-full border-[3px] border-muted" />
            <div className="absolute inset-0 rounded-full border-[3px] border-primary border-t-transparent animate-spin" />
          </div>
          <p className="text-sm text-muted-foreground">Loading session info…</p>
        </div>
      )}

      {/* Error */}
      {isError && !isLoading && !introspectResponse && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-8 flex flex-col items-center gap-4 text-center">
          <div className="w-11 h-11 rounded-full bg-destructive/10 flex items-center justify-center">
            <AlertCircle className="w-5 h-5 text-destructive" />
          </div>
          <div className="space-y-0.5">
            <p className="text-sm font-semibold text-foreground">Failed to introspect</p>
            <p className="text-xs text-muted-foreground">
              {(error as Error)?.message ?? 'An unexpected error occurred.'}
            </p>
          </div>
          <button
            type='button'
            onClick={() => refetch()}
            className="inline-flex items-center gap-1.5 text-xs font-medium px-4 py-2 rounded-lg border border-border bg-background text-foreground hover:bg-muted transition-all shadow-sm"
          >
            <RefreshCw className="w-3.5 h-3.5" /> Try again
          </button>
        </div>
      )}

      {/* Content */}
      {introspectResponse && !isLoading && (
        <div className="flex flex-col gap-5">

          {/* Card + panels row */}
          <div className="flex flex-col lg:flex-row lg:items-start gap-5">

            <div className="w-full flex justify-center lg:block lg:w-75 lg:shrink-0">
              <div className="min-w-75 w-full">
                <ModernIntrospect
                  data={{
                    cred: {
                      id: introspectResponse.cred.id,
                      type: introspectResponse.cred.type,
                    },
                    sub: {
                      id: introspectResponse.sub.id,
                      project_id: introspectResponse.sub.project_id,
                      email: introspectResponse.sub.email,
                      type: introspectResponse.sub.type,
                      capabilities: Array.isArray(introspectResponse.sub.capabilities)
                        ? introspectResponse.sub.capabilities
                        : Object.keys(introspectResponse.sub.capabilities ?? {}).map(Number),
                      metadata: Array.isArray(introspectResponse.sub.metadata)
                        ? introspectResponse.sub.metadata
                        : Object.keys(introspectResponse.sub.metadata ?? {}).map(Number),
                    },
                  }}
                />
              </div>
            </div>

            {/* Capabilities */}
            <div className="w-full rounded-xl border border-border/60 bg-card shadow-sm overflow-hidden">
              <div className="px-4 py-2.5 border-b border-border/60">
                <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                  Capabilities
                </span>
              </div>
              <div className="p-4">
                {capKeys.length > 0 ? (
                  <div className="flex flex-wrap gap-1.5">
                    {capKeys.map((key) => (
                      <span
                        key={key}
                        className="text-[11px] font-mono px-2 py-1 rounded-md bg-muted/50 border border-border/60 text-muted-foreground"
                      >
                        {key}
                      </span>
                    ))}
                  </div>
                ) : (
                  <p className="text-[11px] text-muted-foreground text-center py-2 italic">
                    No capabilities
                  </p>
                )}
              </div>
            </div>

          </div>

          {/* Raw JSON */}
          {showRaw && (
            <div className="rounded-xl border border-border/60 bg-card shadow-sm overflow-hidden">
              <div className="px-4 py-2.5 border-b border-border/60 flex items-center justify-between gap-3">
                <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                  Raw JSON
                </span>
                <CopyButton json={JSON.stringify(introspectResponse, null, 2)} />
              </div>
              <pre className="p-4 text-[11px] font-mono leading-relaxed text-foreground overflow-x-auto whitespace-pre-wrap break-all max-h-80 overflow-y-auto">
                {JSON.stringify(introspectResponse, null, 2)}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  )
}