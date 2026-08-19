import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { listAllByWalletIntentsQueryOptions } from "#/features/payment-intents/api";
import { TransactionsDashboard } from "#/features/payment-intents/ui/transactions-dashboard";
import { walletByIdQueryOptions } from "#/features/wallets/api";

export const Route = createFileRoute("/admin/wallets/$walletID/transactions")({
  component: RouteComponent,
});

function RouteComponent() {
  const { walletID } = Route.useParams();
  const { data: wallet } = useQuery(walletByIdQueryOptions(walletID));
  const { data: intents = [] } = useQuery(
    listAllByWalletIntentsQueryOptions(walletID),
  );
  return (
    <TransactionsDashboard
      title="Wallet transactions"
      description="Payment activity for this wallet."
      intents={intents}
      feeBps={wallet?.fee_bps ?? 0}
    />
  );
}
