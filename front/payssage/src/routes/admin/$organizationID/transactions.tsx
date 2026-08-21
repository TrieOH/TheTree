import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { listAllByOrgIntentsQueryOptions } from "#/features/payment-intents/api";
import { TransactionsDashboard } from "#/features/payment-intents/ui/transactions-dashboard";
import { allWalletsQueryOptions } from "#/features/wallets/api";

export const Route = createFileRoute("/admin/$organizationID/transactions")({
  component: RouteComponent,
});

function RouteComponent() {
  const { organizationID } = Route.useParams();

  const { data: intents = [] } = useQuery(
    listAllByOrgIntentsQueryOptions(organizationID),
  );
  const { data: wallets = [] } = useQuery(
    allWalletsQueryOptions(organizationID),
  );
  const walletFees = Object.fromEntries(
    wallets.map((wallet) => [wallet.id, wallet.fee_bps]),
  );
  return (
    <TransactionsDashboard
      title="Organization transactions"
      description="Payment activity for wallets in this organization."
      intents={intents}
      walletFees={walletFees}
    />
  );
}
