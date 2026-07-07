import { Check, Copy, Pencil, Sparkles, User, Wallet } from "lucide-react"
import { useMemo, useState } from "react"
import { toast } from "sonner"
import { bpsToPercentage, cn } from "#/shared/lib/utils"
import type { WalletI, WalletSetSandboxI } from "../model"
import { Button } from "#/shared/ui/shadcn/button"
import { truncateString } from "@trieoh/shared-utils"

interface PropsI {
  data: WalletI
  onEditFee: (wallet: WalletI) => void
  onSetSandbox: (walletId: string, data: WalletSetSandboxI) => void
  isSettingSandbox?: boolean
}

export default function WalletCard({
  data,
  onEditFee,
  onSetSandbox,
  isSettingSandbox = false,
}: PropsI) {
  const [copied, setCopied] = useState(false)

  const owner = useMemo(() => {
    return data.organization_id ? "Organization" : "Personal"
  }, [data.organization_id])

  const copyId = async () => {
    await navigator.clipboard.writeText(data.id)
    setCopied(true)
    toast.success("Wallet ID copied")
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <>
      <div
        className={cn(
          "w-full overflow-hidden rounded-sm border border-border/70 bg-card transition-all",
          "hover:border-primary hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0.08)]",
        )}
      >
        <div className="px-4 pt-4">
          <div className="flex items-start justify-between gap-3">
            <div className="flex min-w-0 items-start gap-3">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-sm bg-primary/10 text-primary">
                <Wallet className="size-5" />
              </div>
              <div className="min-w-0">
                <div className="truncate text-sm font-semibold leading-tight text-foreground">
                  {data.name}
                </div>
                <div className="mt-0.5 flex items-center gap-2 text-[11px] font-mono text-muted-foreground">
                  {data.organization_id ? (
                    <span className="inline-flex min-w-0 items-center gap-1 truncate">
                      <User className="size-3.5 shrink-0" />
                      {owner}
                    </span>
                  ) : (
                    <span className="inline-flex min-w-0 items-center gap-1 truncate">
                      <User className="size-3.5 shrink-0" />
                      {owner}
                    </span>
                  )}
                </div>
              </div>
            </div>

            <Button
              type="button"
              size="icon-sm"
              variant="ghost"
              className={cn("shrink-0 rounded-sm", copied && "text-emerald-600")}
              onClick={copyId}
            >
              {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
            </Button>
          </div>
        </div>

        <div className="mx-4 mt-3 border-t border-border/60" />

        <div className="grid gap-3 px-4 py-3">
          <div className="flex items-center justify-between gap-3">
            <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
              Wallet ID
            </span>
            <span className="max-w-48 truncate font-mono text-[11px] text-foreground">
              {truncateString(data.id, 8, 4)}
            </span>
          </div>

          <div className="flex items-center justify-between gap-3">
            <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
              Fee
            </span>
            <div className="flex items-center gap-2">
              <span className="font-mono text-sm text-foreground">
                {bpsToPercentage(data.fee_bps).toFixed(2)}%
              </span>
              <Button
                type="button"
                variant="ghost"
                size="icon-xs"
                className="rounded-sm text-muted-foreground hover:text-foreground"
                onClick={() => onEditFee(data)}
              >
                <Pencil className="size-3.5" />
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-between gap-3">
            <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
              Environment
            </span>
            <span
              className={cn(
                "inline-flex items-center gap-1 rounded-sm px-2 py-0.5 text-[10px] font-bold uppercase tracking-widest",
                data.sandbox ? "bg-amber-500/10 text-amber-600" : "bg-emerald-500/10 text-emerald-600",
              )}
            >
              <Sparkles className="size-3" />
              {data.sandbox ? "Sandbox" : "Production"}
            </span>
          </div>
        </div>

        <div className="flex items-center justify-between gap-2 border-t border-border/60 px-4 py-3">
          <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            Sandbox Mode
          </span>
          <Button
            type="button"
            variant={data.sandbox ? "secondary" : "outline"}
            size="sm"
            className="rounded-sm"
            disabled={isSettingSandbox}
            onClick={() =>
              onSetSandbox(data.id, {
                sandbox: !data.sandbox,
                organization_id: data.organization_id,
              })
            }
          >
            {data.sandbox ? "Disable" : "Enable"}
          </Button>
        </div>
      </div>
    </>
  )
}
