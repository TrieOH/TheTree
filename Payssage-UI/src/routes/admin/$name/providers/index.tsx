import { allWorkspaceMarketplaceConfigsQueryOptions, removeMarketplaceConfigFromWorkspaceFn, setupOauthOnWorkspaceFn, updateOauthOnWorkspaceFn } from '#/features/oauth/api'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import { createFileRoute, useParams } from '@tanstack/react-router'
import {
  ArrowRightFromLine,
  CreditCard,
  CheckCircle2,
  Zap,
  Trash2,
  RefreshCw,
  Copy,
  Eye,
  EyeOff
} from 'lucide-react'
import z from 'zod'
import FormModal from '#/widgets/modal/form-modal'
import { oauthSetupSchema } from '#/features/oauth/model'
import { useState, useMemo } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { bpsToPercentage } from '#/shared/lib/utils'
import type { OauthSetupI, OauthWorkspaceMarketplaceConfigI } from '#/features/oauth/model'
import { env } from '#/env'


const queryParams = z.object({
  status: z.string().optional().nullable(),
  provider: z.string().optional().nullable(),
})

export const Route = createFileRoute('/admin/$name/providers/')({
  component: RouteComponent,
  validateSearch: (search) => queryParams.parse(search)
})

const getProviderInfo = (provider: string) => {
  const p = provider.toLowerCase()
  if (p.includes('mercadopago')) {
    return {
      id: provider,
      name: 'Mercado Pago',
      logo: '/external-logos/MP_RGB_HANDSHAKE_color_vertical.svg',
      description: 'The leading payment solution in Latin America. Connect your account to accept Pix, credit cards, and more.',
      features: [
        { icon: <CreditCard className="w-3.5 h-3.5" />, label: 'Credit Cards' },
        { icon: <Zap className="w-3.5 h-3.5" />, label: 'Pix Support' }
      ]
    }
  }

  return {
    id: provider,
    name: provider.split(/[-_]/).map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' '),
    logo: '',
    description: `Connect your ${provider} account to start processing transactions.`,
    features: []
  }
}

function RouteComponent() {
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null)
  const [updatingConfig, setUpdatingConfig] = useState<OauthWorkspaceMarketplaceConfigI | null>(null)
  const [showCredential, setShowCredential] = useState<Record<string, boolean>>({})
  const { name } = useParams({ from: '/admin/$name/providers/' })
  const queryClient = useQueryClient()

  const supportedProviders = useMemo(() => {
    return env.VITE_SUPPORTED_PROVIDERS
  }, [])

  const { data: configs, isLoading: isLoadingConfigs } = useQuery(allWorkspaceMarketplaceConfigsQueryOptions(name))

  const { mutate: setupOauthOnWorkspace, isPending: isPendingSetup } = useMutation({
    mutationFn: (res: { data: OauthSetupI, provider: string }) =>
      setupOauthOnWorkspaceFn(
        res.data,
        name,
        res.provider,
        window.location.origin + `/admin/${name}/providers`
      ),
    onSuccess: (response) => {
      if (response.success) {
        setSelectedProvider(null)
        window.location.href = response.data.redirect_url
      }
    },
  })

  const { mutate: updateOauthOnWorkspace, isPending: isPendingUpdate } = useMutation({
    mutationFn: (data: OauthSetupI) => {
      if (!updatingConfig) throw new Error("No config selected for update")
      return updateOauthOnWorkspaceFn(data, name, updatingConfig.credential_id)
    },
    onSuccess: (response) => {
      if (response.success) {
        toast.success("Provider configuration updated.")
        queryClient.setQueryData(
          allWorkspaceMarketplaceConfigsQueryOptions(name).queryKey,
          (old: OauthWorkspaceMarketplaceConfigI[] = []) =>
            old.map(c => c.id === response.data.id ? response.data : c)
        )
        setUpdatingConfig(null)
      }
    },
    onError: () => {
      toast.error("Failed to update provider configuration.")
    }
  })

  const { mutate: removeMarketplaceConfig, isPending: isPendingRemove } = useMutation({
    mutationFn: (credential_id: string) => removeMarketplaceConfigFromWorkspaceFn(credential_id, name),
    onSuccess: (_, credential_id) => {
      toast.success("Provider disconnected successfully.")
      queryClient.setQueryData(
        allWorkspaceMarketplaceConfigsQueryOptions(name).queryKey,
        (old: OauthWorkspaceMarketplaceConfigI[] = []) =>
          old.filter(c => c.credential_id !== credential_id)
      )
    },
    onError: () => {
      toast.error("Failed to disconnect provider.")
    }
  })

  const toggleCredential = (id: string) => {
    setShowCredential(prev => ({ ...prev, [id]: !prev[id] }))
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success("Credential ID copied to clipboard")
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="space-y-1">
        <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Providers</h2>
        <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">
          Connect your payment gateways to start processing transactions.
        </p>
      </div>

      <div className="grid gap-6">
        {supportedProviders.map((providerId) => {
          const config = configs?.find(c => c.provider === providerId)
          const isConnected = !!config
          const info = getProviderInfo(providerId)

          return (
            <Card key={providerId} className="rounded-none border-border group transition-all duration-300 hover:border-primary/50 relative overflow-hidden">
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
                  {info.logo ? (
                    <img
                      src={info.logo}
                      alt={info.name}
                      className="w-full h-full object-contain"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-muted/20">
                      <CreditCard className="w-8 h-8 text-muted-foreground/40" />
                    </div>
                  )}
                </div>

                <div className="space-y-1.5 flex-1">
                  <div className="flex items-center gap-2">
                    <CardTitle className="text-xl font-black uppercase tracking-tight">{info.name}</CardTitle>
                  </div>
                  <CardDescription className="text-xs font-mono uppercase tracking-widest max-w-lg">
                    {info.description}
                  </CardDescription>
                  {isConnected && (
                    <div className="flex items-center gap-2 mt-2">
                      <span className="text-[10px] font-black uppercase tracking-widest text-muted-foreground/60">Credential ID:</span>
                      <div className="flex items-center gap-1 bg-muted/50 px-2 py-0.5 border border-border/50">
                        <span className="text-[10px] font-mono text-muted-foreground">
                          {showCredential[config.id] ? config.credential_id : "••••••••••••••••"}
                        </span>
                        <button
                          onClick={() => toggleCredential(config.id)}
                          className="p-1 hover:text-primary transition-colors"
                          title={showCredential[config.id] ? "Hide" : "Show"}
                        >
                          {showCredential[config.id] ? <EyeOff className="w-3 h-3" /> : <Eye className="w-3 h-3" />}
                        </button>
                        <button
                          onClick={() => copyToClipboard(config.credential_id)}
                          className="p-1 hover:text-primary transition-colors"
                          title="Copy"
                        >
                          <Copy className="w-3 h-3" />
                        </button>
                      </div>
                      <Badge variant="secondary" className="rounded-none text-[9px] font-mono uppercase">
                        Fee: {bpsToPercentage(config.fee_bps)}%
                      </Badge>
                    </div>
                  )}
                </div>

                <div className="flex items-center gap-3 shrink-0">
                  {isConnected ? (
                    <div className="flex gap-2">
                      <Button
                        onClick={() => setUpdatingConfig(config)}
                        variant="outline"
                        size="sm"
                        className="rounded-none gap-2 h-9 font-black uppercase tracking-widest transition-all"
                      >
                        <RefreshCw className="w-3.5 h-3.5" />
                        Update
                      </Button>
                      <Button
                        onClick={() => removeMarketplaceConfig(config.credential_id)}
                        variant="outline"
                        size="sm"
                        disabled={isPendingRemove}
                        className="rounded-none gap-2 h-9 font-black uppercase tracking-widest transition-all border-destructive/30 hover:bg-destructive/10 hover:text-destructive hover:border-destructive/50"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                        Disconnect
                      </Button>
                    </div>
                  ) : (
                    <Button
                      onClick={() => setSelectedProvider(providerId)}
                      variant="default"
                      className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all px-8"
                      disabled={isLoadingConfigs || isPendingSetup}
                    >
                      <ArrowRightFromLine className="w-4 h-4" />
                      Connect
                    </Button>
                  )}
                </div>
              </CardHeader>

              <CardContent className="border-t border-border/40 bg-muted/5 flex flex-wrap gap-4 py-3 px-6">
                {info.features.map((feature, idx) => (
                  <div key={idx} className="flex items-center gap-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                    {feature.icon}
                    {feature.label}
                  </div>
                ))}
              </CardContent>
            </Card>
          )
        })}
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

      {updatingConfig && <FormModal<OauthSetupI>
        title={`Update ${updatingConfig.provider} Configuration`}
        description="Adjust the fee percentage for this provider."
        buttonTitle="Update Configuration"
        schema={oauthSetupSchema}
        formId="update-oauth-form"
        isOpen={!!updatingConfig}
        onClose={() => setUpdatingConfig(null)}
        onSubmit={updateOauthOnWorkspace}
        defaultValues={{
          fee_percent: bpsToPercentage(updatingConfig.fee_bps)
        }}
        fields={[
          {
            name: "fee_percent",
            label: "Fee (%)",
            type: "percentage",
            placeholder: "Ex: 1.5"
          }
        ]}
        disabled={isPendingUpdate}
      />}
    </div>
  )
}

