import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { useLayoutHeader } from "@trieoh/ui-base";
import { Link2, Receipt, Store, WalletCards, Webhook } from "lucide-react";
import { useMemo } from "react";
import { walletByIdQueryOptions } from "#/features/wallets/api";
import { cn } from "#/shared/lib/utils";

export const Route = createFileRoute("/admin/wallets/$walletID")({
  component: WalletLayout,
});

function WalletLayout() {
  const { walletID } = Route.useParams();
  const { data: wallet } = useQuery(walletByIdQueryOptions(walletID));
  const header = useMemo(
    () => (
      <div>
        <h1 className="text-lg font-semibold">{wallet?.name ?? "Wallet"}</h1>
        <p className="text-sm text-muted-foreground">
          Manage this wallet's collector and seller accounts.
        </p>
      </div>
    ),
    [wallet?.name],
  );
  useLayoutHeader(header);

  const tabs = [
    {
      label: "Collector",
      to: "/admin/wallets/$walletID",
      icon: Link2,
      exact: true,
    },
    {
      label: "Transactions",
      to: "/admin/wallets/$walletID/transactions",
      icon: Receipt,
      exact: true,
    },
    {
      label: "Sellers",
      to: "/admin/wallets/$walletID/sellers",
      icon: WalletCards,
      exact: true,
    },
    {
      label: "Connect seller",
      to: "/admin/wallets/$walletID/connect-seller",
      icon: Store,
      exact: true,
    },
    {
      label: "Webhooks",
      to: "/admin/wallets/$walletID/webhooks",
      icon: Webhook,
      exact: true,
    },
  ] as const;

  return (
    <div className="space-y-6">
      <div className="-mx-6 -mt-6 border-b border-border/60 bg-background/50 px-6">
        <div className="flex h-12 items-center gap-8">
          {tabs.map((tab) => (
            <Link
              key={tab.label}
              to={tab.to}
              params={{ walletID }}
              activeOptions={{ exact: tab.exact }}
              className="group relative flex h-full items-center gap-2 text-[10px] font-bold uppercase tracking-widest"
            >
              {({ isActive }) => (
                <>
                  <tab.icon
                    className={cn(
                      "size-3.5",
                      isActive
                        ? "text-primary"
                        : "text-muted-foreground group-hover:text-foreground",
                    )}
                  />
                  <span
                    className={
                      isActive
                        ? "text-foreground"
                        : "text-muted-foreground group-hover:text-foreground"
                    }
                  >
                    {tab.label}
                  </span>
                  {isActive && (
                    <div className="absolute inset-x-0 bottom-0 h-0.5 bg-primary" />
                  )}
                </>
              )}
            </Link>
          ))}
        </div>
      </div>
      <Outlet />
    </div>
  );
}
