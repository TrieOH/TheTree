import type { QueryKey } from "@tanstack/react-query";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { PaginatedContainer } from "@trieoh/ui-base";
import { CalendarDays, Check, Copy, ShieldOff } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { revokeProviderFn } from "#/features/oauth/api";
import type { Collector, Seller } from "#/features/oauth/model";
import { Badge } from "#/shared/ui/shadcn/badge";
import { Button } from "#/shared/ui/shadcn/button";

type ProviderCredential = Collector | Seller;
type ProviderFlow = "collector" | "seller";

const providerDetails: Partial<
  Record<string, { label: string; logo?: string }>
> = {
  mercadopago: {
    label: "Mercado Pago",
    logo: "/external-logos/MP_RGB_HANDSHAKE_color_vertical.svg",
  },
};

const formatDate = (value: string) =>
  new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
  }).format(new Date(value));

function ProviderCredentialCard({
  item,
  flow,
  onRevoke,
  isRevoking,
}: {
  item: ProviderCredential;
  flow: ProviderFlow;
  onRevoke: (item: ProviderCredential) => void;
  isRevoking: boolean;
}) {
  const [copied, setCopied] = useState(false);
  const provider = providerDetails[item.provider];
  const providerLabel = provider?.label ?? item.provider.replaceAll("_", " ");
  const isRevoked = Boolean(item.revoked_at);

  const copyProviderUserId = async () => {
    await navigator.clipboard.writeText(item.provider_user_id);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1500);
  };

  return (
    <article className="flex w-full min-w-0 max-w-full flex-col overflow-hidden rounded-md border bg-card">
      <div className="flex min-w-0 items-start justify-between gap-3 p-4">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-md bg-muted/60 p-2">
            {provider?.logo ? (
              <img
                src={provider.logo}
                alt=""
                className="size-full object-contain"
              />
            ) : (
              <span className="text-xs font-bold uppercase text-primary">
                {providerLabel.slice(0, 2)}
              </span>
            )}
          </div>
          <div className="min-w-0">
            <h3 className="truncate text-sm font-semibold capitalize">
              {providerLabel}
            </h3>
            <p className="mt-0.5 truncate text-xs capitalize text-muted-foreground">
              {flow} account
            </p>
          </div>
        </div>
        <Badge
          variant="secondary"
          className={
            isRevoked
              ? "gap-1.5"
              : "gap-1.5 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
          }
        >
          <span
            className={
              isRevoked
                ? "size-1.5 rounded-full bg-muted-foreground"
                : "size-1.5 rounded-full bg-emerald-500"
            }
          />
          {isRevoked ? "Revoked" : "Active"}
        </Badge>
      </div>

      <div className="mx-4 min-w-0 rounded-md bg-muted/40 px-3 py-2.5">
        <p className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
          Provider user
        </p>
        <button
          type="button"
          className="mt-1 flex w-full min-w-0 cursor-pointer items-center gap-2 text-left"
          onClick={() => void copyProviderUserId()}
          aria-label="Copy provider user ID"
        >
          <span className="min-w-0 flex-1 truncate font-mono text-xs text-foreground">
            {item.provider_user_id}
          </span>
          {copied ? (
            <Check className="size-3.5 shrink-0 text-emerald-500" />
          ) : (
            <Copy className="size-3.5 shrink-0 text-muted-foreground" />
          )}
        </button>
      </div>

      <div className="flex min-w-0 items-center justify-between gap-2 px-4 py-2.5 text-xs text-muted-foreground">
        <span className="flex min-w-0 flex-1 items-center gap-1.5">
          <CalendarDays className="size-3.5 shrink-0" />
          <span className="truncate">{formatDate(item.created_at)}</span>
        </span>
        <Button
          variant="ghost"
          size="sm"
          className="h-8 shrink-0 px-2.5 text-xs text-destructive hover:bg-destructive/10 hover:text-destructive"
          disabled={isRevoked || isRevoking}
          onClick={() => onRevoke(item)}
        >
          {isRevoked && <ShieldOff />}
          {isRevoked ? "Revoked" : isRevoking ? "Revoking..." : "Revoke"}
        </Button>
      </div>
    </article>
  );
}

export function ProviderCredentialList({
  items,
  flow,
  queryKey,
}: {
  items: ProviderCredential[];
  flow: ProviderFlow;
  queryKey: QueryKey;
}) {
  const [filter, setFilter] = useState("");
  const queryClient = useQueryClient();
  const {
    mutate: revoke,
    isPending: isRevoking,
    variables: revokingItem,
  } = useMutation({
    mutationFn: (item: ProviderCredential) =>
      revokeProviderFn(item.provider, { flow, id: item.id }),
    onSuccess: (response) => {
      if (!response.success) {
        toast.error(response.message || "Failed to revoke provider");
        return;
      }
      void queryClient.invalidateQueries({ queryKey });
      toast.success("Provider revoked");
    },
    onError: () => toast.error("Failed to revoke provider"),
  });
  const search = filter.trim().toLowerCase();
  const filteredItems = search
    ? items.filter(
        (item) =>
          item.provider.toLowerCase().includes(search) ||
          item.provider_user_id.toLowerCase().includes(search) ||
          (item.revoked_at ? "revoked" : "active").includes(search),
      )
    : items;

  return (
    <PaginatedContainer<ProviderCredential>
      items={filteredItems}
      layout="grid"
      minItemWidth="min(100%, 17rem)"
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
      renderItems={(slice) =>
        slice.map((item) => (
          <ProviderCredentialCard
            key={item.id}
            item={item}
            flow={flow}
            onRevoke={revoke}
            isRevoking={isRevoking && revokingItem.id === item.id}
          />
        ))
      }
    />
  );
}
