import { createFileRoute, Link } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { PenLine, Plus, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import { allSignaturesQueryOptions, removeSignatureFn } from '@/features/signatures/api'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Badge } from '@/shared/ui/shadcn/badge'

export const Route = createFileRoute('/admin/events/$eventId_/editions/$editionId/signatures/')({
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const queryClient = useQueryClient()

  const { data: editions = [] } = useQuery(allAdminEditionsQueryOptions(eventId))
  const edition = editions.find((item) => item.id === editionId) ?? null

  const { data: signatures = [] } = useQuery(allSignaturesQueryOptions(eventId, editionId))

  const removeMutation = useMutation({
    mutationFn: (sigId: string) => removeSignatureFn(eventId, editionId, sigId),
    onSuccess: (res) => {
      if (res.success) {
        void queryClient.invalidateQueries({ queryKey: allSignaturesQueryOptions(eventId, editionId).queryKey })
        toast.success('Assinatura removida')
      } else {
        toast.error(res.message || 'Erro ao remover assinatura')
      }
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  return (
    <div className="min-h-screen bg-background text-foreground">
      <div className="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-10">
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div className="space-y-1">
            <p className="text-xs font-medium uppercase tracking-[0.24em] text-muted-foreground">Admin</p>
            <h1 className="text-2xl font-semibold">Assinaturas</h1>
            <p className="text-sm text-muted-foreground">
              Biblioteca de assinaturas para {edition?.edition_name ?? 'esta edição'}.
            </p>
          </div>
          <Link
            to="/admin/events/$eventId/editions/$editionId/signatures/editor"
            params={{ eventId, editionId }}
            className="inline-flex h-10 w-full items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 md:w-auto"
          >
            <Plus className="size-4" />
            Nova assinatura
          </Link>
        </div>

        <div className="mt-6 grid gap-6">
          <Card>
            <CardHeader className="border-b pb-3">
              <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                <PenLine className="size-4 text-muted-foreground" />
                Assinaturas salvas
              </CardTitle>
              <CardDescription className="text-xs">
                Desenhos ou imagens prontos para serem usados nas certificações.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-4">
              {signatures.length === 0 ? (
                <div className="rounded-xl border border-dashed p-6 text-sm text-muted-foreground">
                  Nenhuma assinatura cadastrada ainda.
                </div>
              ) : (
                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                  {signatures.map((signature) => (
                    <div key={signature.id} className="rounded-2xl border bg-card p-4">
                      <div className="mb-3 flex items-start justify-between gap-3">
                        <div className="space-y-1">
                          <p className="font-medium">{signature.title}</p>
                          <Badge variant="outline">Assinatura</Badge>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => { void removeMutation.mutateAsync(signature.id) }}
                          disabled={removeMutation.isPending}
                        >
                          <Trash2 className="size-4 text-destructive" />
                        </Button>
                      </div>
                      <div className="rounded-xl border bg-muted/10 p-3">
                        <img src={signature.url} alt={signature.title} className="h-28 w-full object-contain" />
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
