import { setupOauthOnWorkspaceFn } from '#/features/oauth/api'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import { createFileRoute, useParams } from '@tanstack/react-router'
import { ArrowRightFromLine, CreditCard, CheckCircle2, Zap } from 'lucide-react'
import z from 'zod'
import FormModal from '#/widgets/modal/form-modal'
import { oauthSetupSchema, type OauthSetupI } from '#/features/oauth/model'
import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'

const queryParams = z.object({
  status: z.string().optional().nullable(),
  provider: z.string().optional().nullable(),
})

export const Route = createFileRoute('/admin/$name/providers/')({
  component: RouteComponent,
  validateSearch: (search) => queryParams.parse(search)
})

function RouteComponent() {
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null)
  const { name } = useParams({ from: '/admin/$name/providers/' })
  const { status } = Route.useSearch()

  const { mutate: setupOauthOnWorkspace, isPending: isPendingSetup } = useMutation({
    mutationFn: (res: { data: OauthSetupI, provider: string }) =>
      setupOauthOnWorkspaceFn(res.data, name, res.provider),
    onSuccess: (response) => {
      if (response.success) {
        setSelectedProvider(null)
        window.location.href = response.data.redirect_url
      }
    },
  })
  const isConnected = status === "success"

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="space-y-1">
        <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Providers</h2>
        <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">
          Connect your payment gateways to start processing transactions.
        </p>
      </div>

      <div className="grid gap-6">
        <Card className="rounded-none border-border group transition-all duration-300 hover:border-primary/50 relative overflow-hidden">
          {isConnected && (
            <div className="absolute top-0 right-0 p-2">
              <Badge variant="outline" className="rounded-none text-[8px] font-black uppercase tracking-widest bg-emerald-500/10 text-emerald-500 border-emerald-500/20 gap-1">
                <CheckCircle2 className="w-2.5 h-2.5" />
                Connected
              </Badge>
            </div>
          )}

          <CardHeader className="flex flex-col sm:flex-row sm:items-center gap-6 pb-6">
            <div className="flex items-center justify-center w-16 h-16 bg-background border border-border group-hover:border-primary/30 transition-colors shrink-0 p-2">
              <img
                src="/external-logos/MP_RGB_HANDSHAKE_color_vertical.svg"
                alt="Mercado Pago"
                className="w-full h-full object-contain"
              />
            </div>

            <div className="space-y-1.5 flex-1">
              <div className="flex items-center gap-2">
                <CardTitle className="text-xl font-black uppercase tracking-tight">Mercado Pago</CardTitle>
              </div>
              <CardDescription className="text-xs font-mono uppercase tracking-widest max-w-lg">
                The leading payment solution in Latin America. Connect your account to accept Pix, credit cards, and more.
              </CardDescription>
            </div>

            <div className="flex items-center gap-3 shrink-0">
              <Button
                // onClick={performSetupOauth}
                onClick={() => setSelectedProvider("mercadopago")}
                variant={isConnected ? "outline" : "default"}
                className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all px-8"
                disabled={isConnected}
              >
                {!isConnected && <ArrowRightFromLine className="w-4 h-4" />}
                {isConnected ? "Connected" : "Connect"}
              </Button>
            </div>
          </CardHeader>

          <CardContent className="border-t border-border/40 bg-muted/5 flex flex-wrap gap-4 py-3 px-6">
            <div className="flex items-center gap-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              <CreditCard className="w-3.5 h-3.5" />
              Credit Cards
            </div>
            <div className="flex items-center gap-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              <Zap className="w-3.5 h-3.5" />
              Pix Support
            </div>
          </CardContent>
        </Card>
      </div>
      {selectedProvider && <FormModal<OauthSetupI>
        title="Configure OAuth Fee"
        description="Set the fee percentage that will be applied to transactions processed through this provider."
        buttonTitle="Perform OAuth Setup"
        schema={oauthSetupSchema}
        formId="setup-oauth-form"
        isOpen={!!selectedProvider}
        onClose={() => setSelectedProvider(null)}
        onSubmit={(data) => setupOauthOnWorkspace({ data, provider: selectedProvider })}
        fields={[
          {
            name: "fee_percent",
            label: "Fee (%)",
            type: "percentage",
            placeholder: "Ex: 1.5"
          }
        ]}
        disabled={isPendingSetup}
      />}
    </div>
  )
}
