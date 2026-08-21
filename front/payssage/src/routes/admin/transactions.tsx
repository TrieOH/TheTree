import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { listAllIntentsQueryOptions } from "#/features/payment-intents/api";
import { TransactionsDashboard } from "#/features/payment-intents/ui/transactions-dashboard";
import { allWalletsQueryOptions } from "#/features/wallets/api";

export const Route = createFileRoute("/admin/transactions")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: intents = [] } = useQuery(listAllIntentsQueryOptions());
  const { data: wallets = [] } = useQuery(allWalletsQueryOptions());
  const walletFees = Object.fromEntries(
    wallets.map((wallet) => [wallet.id, wallet.fee_bps]),
  );
  return (
    <div className="p-6">
      <TransactionsDashboard intents={intents} walletFees={walletFees} />
    </div>
  );
}
