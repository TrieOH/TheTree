import { useEffect, useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Check, ChevronDown, Link2, Unlink } from "lucide-react"
import { toast } from "sonner"
import { collectorsQueryOptions } from "#/features/collectors/api"
import { bindCollectorToWalletFn, unbindCollectorFromWalletFn, walletByIdQueryOptions } from "../api"
import type { WalletI } from "../model"
import { Button } from "#/shared/ui/shadcn/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#/shared/ui/shadcn/dropdown-menu"

export function WalletCollector({ wallet }: { wallet: WalletI }) {
  const queryClient = useQueryClient()
  const { data: collectors = [] } = useQuery(collectorsQueryOptions(wallet.organization_id))
  const [collectorId, setCollectorId] = useState(wallet.collector_id ?? "")

  useEffect(() => setCollectorId(wallet.collector_id ?? ""), [wallet.collector_id])

  const refreshWallet = () => queryClient.invalidateQueries({ queryKey: walletByIdQueryOptions(wallet.id).queryKey })
  const { mutate: bind, isPending: isBinding } = useMutation({
    mutationFn: () => bindCollectorToWalletFn(wallet.id, { collector_id: collectorId }),
    onSuccess: (response) => {
      if (!response.success) return toast.error(response.message || "Failed to bind collector")
      void refreshWallet()
      toast.success(wallet.collector_id ? "Collector changed" : "Collector bound")
    },
    onError: () => toast.error("Failed to bind collector"),
  })
  const { mutate: unbind, isPending: isUnbinding } = useMutation({
    mutationFn: () => unbindCollectorFromWalletFn(wallet.id),
    onSuccess: (response) => {
      if (!response.success) return toast.error(response.message || "Failed to unbind collector")
      setCollectorId("")
      void refreshWallet()
      toast.success("Collector unbound")
    },
    onError: () => toast.error("Failed to unbind collector"),
  })

  const activeCollectors = collectors.filter((collector) => !collector.revoked_at)
  const selectedCollector = activeCollectors.find((collector) => collector.id === collectorId)

  return (
    <section className="space-y-3 rounded-sm border bg-card p-4">
      <div>
        <h2 className="text-sm font-semibold">Wallet collector</h2>
        <p className="text-xs text-muted-foreground">Bind one collector to process payments for this wallet.</p>
      </div>
      <div className="space-y-2">
        <span className="text-[10px] font-black uppercase tracking-[0.2em]">Collector</span>
        <div className="flex flex-col gap-3 sm:flex-row">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={(
                <Button variant="outline" className="min-h-10! flex-1 justify-between rounded-sm px-3 py-2 font-normal" />
              )}
            >
              {selectedCollector ? (
                <span className="flex min-w-0 items-center gap-2">
                  <span className="capitalize">{selectedCollector.provider.replaceAll("_", " ")}</span>
                  <span className="truncate font-mono text-xs text-muted-foreground">{selectedCollector.provider_user_id}</span>
                </span>
              ) : <span className="text-muted-foreground">Select a collector</span>}
              <ChevronDown className="size-4 text-muted-foreground" />
            </DropdownMenuTrigger>
            <DropdownMenuContent className="">
              <DropdownMenuGroup>
                <DropdownMenuLabel>Available collectors</DropdownMenuLabel>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              {activeCollectors.map((collector) => (
                <DropdownMenuItem key={collector.id} onClick={() => setCollectorId(collector.id)}>
                  <span className="flex min-w-0 flex-1 flex-col">
                    <span className="capitalize">{collector.provider.replaceAll("_", " ")}</span>
                    <span className="truncate font-mono text-xs text-muted-foreground">{collector.provider_user_id}</span>
                  </span>
                  {collector.id === collectorId && <Check className="text-primary" />}
                </DropdownMenuItem>
              ))}
              {!activeCollectors.length && (
                <DropdownMenuItem disabled>No active collectors available</DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
          <Button
            disabled={!collectorId || collectorId === wallet.collector_id || isBinding}
            onClick={() => bind()}
          >
            <Link2 /> {isBinding ? "Saving..." : wallet.collector_id ? "Change collector" : "Bind collector"}
          </Button>
          {wallet.collector_id && (
            <Button variant="destructive" disabled={isUnbinding} onClick={() => unbind()}>
              <Unlink /> {isUnbinding ? "Unbinding..." : "Unbind"}
            </Button>
          )}
        </div>
      </div>
      {!activeCollectors.length && (
        <p className="text-xs text-muted-foreground">No active collectors are available for this wallet.</p>
      )}
    </section>
  )
}
