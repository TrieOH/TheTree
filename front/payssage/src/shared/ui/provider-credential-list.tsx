import { useState } from "react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import type { QueryKey } from "@tanstack/react-query"
import { PaginatedContainer } from "@trieoh/ui-base"
import { toast } from "sonner"
import type { Collector, Seller } from "#/features/oauth/model"
import { revokeProviderFn } from "#/features/oauth/api"
import { Badge } from "#/shared/ui/shadcn/badge"
import { Button } from "#/shared/ui/shadcn/button"

type ProviderCredential = Collector | Seller

function ProviderCredentialCard({ item, onRevoke, isRevoking }: {
  item: ProviderCredential
  onRevoke: (item: ProviderCredential) => void
  isRevoking: boolean
}) {
  return (
    <div className="rounded-sm border bg-card p-4">
      <div className="flex items-center justify-between gap-3">
        <span className="font-semibold capitalize">{item.provider.replaceAll("_", " ")}</span>
        <Badge variant={item.revoked_at ? "secondary" : "default"}>{item.revoked_at ? "Revoked" : "Active"}</Badge>
      </div>
      <p className="mt-3 text-xs text-muted-foreground">Provider user</p>
      <p className="font-mono text-sm">{item.provider_user_id}</p>
      <div className="mt-4 border-t pt-3">
        <Button
          variant="destructive"
          size="sm"
          disabled={Boolean(item.revoked_at) || isRevoking}
          onClick={() => onRevoke(item)}
        >
          {item.revoked_at ? "Revoked" : isRevoking ? "Revoking..." : "Revoke"}
        </Button>
      </div>
    </div>
  )
}

export function ProviderCredentialList({ items, flow, queryKey }: {
  items: ProviderCredential[]
  flow: "collector" | "seller"
  queryKey: QueryKey
}) {
  const [filter, setFilter] = useState("")
  const queryClient = useQueryClient()
  const { mutate: revoke, isPending: isRevoking, variables: revokingItem } = useMutation({
    mutationFn: (item: ProviderCredential) => revokeProviderFn(item.provider, { flow, id: item.id }),
    onSuccess: (response) => {
      if (!response.success) {
        toast.error(response.message || "Failed to revoke provider")
        return
      }
      void queryClient.invalidateQueries({ queryKey })
      toast.success("Provider revoked")
    },
    onError: () => toast.error("Failed to revoke provider"),
  })
  const search = filter.trim().toLowerCase()
  const filteredItems = search
    ? items.filter((item) =>
        item.provider.toLowerCase().includes(search) ||
        item.provider_user_id.toLowerCase().includes(search) ||
        (item.revoked_at ? "revoked" : "active").includes(search),
      )
    : items

  return (
    <PaginatedContainer<ProviderCredential>
      items={filteredItems}
      layout="grid"
      minItemWidth="16rem"
      gap="4"
      pageSize={9}
      sortFields={[
        { key: "provider", label: "Provider" },
        { key: "provider_user_id", label: "Provider User" },
        { key: "created_at", label: "Connected At" },
      ]}
      filterValue={filter}
      onFilterChange={setFilter}
      filterPlaceholder="Filter by provider, user or status..."
      itemLabel="connected accounts"
      renderItems={(slice) => slice.map((item) => (
        <ProviderCredentialCard
          key={item.id}
          item={item}
          onRevoke={revoke}
          isRevoking={isRevoking && revokingItem.id === item.id}
        />
      ))}
    />
  )
}
