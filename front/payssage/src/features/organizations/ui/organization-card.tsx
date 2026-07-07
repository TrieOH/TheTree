import { cn } from "@/shared/lib/utils"
import { Building2, Copy, Ellipsis, Check, Users } from "lucide-react"
import type { OrganizationI } from "../model"
import { useState } from "react"
import { toast } from "sonner"
import { Button } from "@/shared/ui/shadcn/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card"

interface PropsI {
  data: OrganizationI
}

export default function OrganizationCard({ data }: PropsI) {
  const [copied, setCopied] = useState(false)

  const copyId = async () => {
    await navigator.clipboard.writeText(data.id)
    setCopied(true)
    toast.success("Organization ID copied")
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <Card className="group overflow-hidden rounded-sm border border-border/70 bg-card transition-all hover:border-primary hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0.08)]">
      <CardHeader className="space-y-4 p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="flex min-w-0 items-start gap-3">
            <div className="flex size-10 shrink-0 items-center justify-center rounded-sm bg-primary/10 text-primary">
              <Building2 className="size-5" />
            </div>
            <div className="min-w-0">
              <CardTitle className="truncate text-base font-bold leading-none">
                {data.name}
              </CardTitle>
              <CardDescription className="mt-1 truncate font-mono text-xs">
                @{data.slug}
              </CardDescription>
            </div>
          </div>

          <Button
            type="button"
            size="icon"
            variant="ghost"
            className={cn("size-8 shrink-0 rounded-sm", copied && "text-emerald-600")}
            onClick={copyId}
          >
            {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
          </Button>
        </div>
      </CardHeader>

      <CardContent className="space-y-3 p-4 pt-0">
        <div className="flex items-center justify-between gap-3 text-xs font-medium text-muted-foreground">
          <span className="flex items-center gap-2">
            <Users className="size-3.5" />
            Owner
          </span>
          <span className="max-w-40 truncate font-mono text-foreground" title={data.owner_id}>
            {data.owner_id}
          </span>
        </div>

        <div className="flex items-center justify-between gap-3 text-xs font-medium text-muted-foreground">
          <span>Created</span>
          <span>{new Date(data.created_at).toLocaleDateString("en-US", { month: "short", year: "numeric" })}</span>
        </div>

        <div className="flex items-center justify-between gap-3 border-t border-border/60 pt-3 text-[10px] font-bold uppercase tracking-widest">
          <span className="inline-flex items-center gap-1.5 text-muted-foreground">
            <Ellipsis className="size-3.5" />
            Organization
          </span>
          <span className="text-primary">Active</span>
        </div>
      </CardContent>
    </Card>
  )
}
