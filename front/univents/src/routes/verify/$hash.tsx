import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { BadgeCheck, FileX2, Hash, Loader2, ShieldCheck } from 'lucide-react'
import { verifyCertificationHashFn } from '@/features/certifications/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Separator } from '@/shared/ui/shadcn/separator'

export const Route = createFileRoute('/verify/$hash')({
  component: VerifyCertificationPage,
})

function formatCertifiedAt(value: string) {
  return new Date(value).toLocaleString('pt-BR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function VerifyCertificationPage() {
  const { hash } = Route.useParams()
  const { data, isLoading, isError } = useQuery({
    queryKey: ['certification-verify', hash],
    queryFn: () => verifyCertificationHashFn(hash),
    retry: false,
  })

  const verified = data?.success && data.data?.is_verified
  const payload = data?.success ? data.data : null

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-3xl px-4 py-10 md:px-6 md:py-14">
          <div className="space-y-3">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <ShieldCheck className="size-4" />
              Verificação pública
            </div>
            <h1 className="text-3xl font-semibold tracking-tight">Certificado</h1>
            <p className="max-w-2xl text-sm text-muted-foreground">
              Validação pública do certificado usando o hash da URL.
            </p>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-3xl px-4 py-6 md:px-6 md:py-8">
        <Card className="overflow-hidden border-border/60 bg-card shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle className="flex items-center gap-2 text-base">
              <BadgeCheck className="size-4 text-primary" />
              Status da verificação
            </CardTitle>
            <CardDescription>Hash: <span className="font-mono">{hash}</span></CardDescription>
          </CardHeader>

          <CardContent className="space-y-5 p-5">
            {isLoading ? (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="size-4 animate-spin" />
                Verificando certificado...
              </div>
            ) : verified ? (
              <>
                <div className="flex items-center gap-2 rounded-2xl border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-emerald-700">
                  <BadgeCheck className="size-5" />
                  Certificado verificado com sucesso
                </div>

                <div className="grid gap-3 sm:grid-cols-2">
                  <div className="rounded-2xl border bg-muted/20 p-4">
                    <p className="text-[10px] uppercase tracking-wider text-muted-foreground">Usuário</p>
                    <p className="mt-1 font-mono text-sm break-all">{payload?.user_id}</p>
                  </div>
                  <div className="rounded-2xl border bg-muted/20 p-4">
                    <p className="text-[10px] uppercase tracking-wider text-muted-foreground">Tipo</p>
                    <p className="mt-1 text-sm font-medium capitalize">{payload?.target_type}</p>
                  </div>
                </div>

                <div className="rounded-2xl border bg-muted/20 p-4">
                  <p className="text-[10px] uppercase tracking-wider text-muted-foreground">Target</p>
                  <p className="mt-1 font-mono text-sm break-all">{payload?.target_id}</p>
                </div>

                <Separator />

                <div className="text-sm text-muted-foreground">
                  Emitido em <span className="font-medium text-foreground">{payload ? formatCertifiedAt(payload.certified_at) : '-'}</span>
                </div>
              </>
            ) : (
              <div className="space-y-4">
                <div className="flex items-center gap-2 rounded-2xl border border-destructive/20 bg-destructive/10 px-4 py-3 text-destructive">
                  <FileX2 className="size-5" />
                  Não foi possível validar este certificado
                </div>
                <p className="text-sm text-muted-foreground">
                  {isError
                    ? 'A verificação falhou ao consultar o certificado.'
                    : 'O hash informado não corresponde a um certificado válido ou já não está ativo.'}
                </p>
              </div>
            )}

            {payload?.target_type && (
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline" className="gap-1.5">
                  <Hash className="size-3.5" />
                  {payload.target_type}
                </Badge>
                <Badge variant="secondary" className="font-mono text-[10px] uppercase tracking-wider">
                  {payload.is_verified ? 'verified' : 'unverified'}
                </Badge>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
