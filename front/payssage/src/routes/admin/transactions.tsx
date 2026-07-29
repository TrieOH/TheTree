import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { listAllIntentsQueryOptions } from "#/features/payment-intents/api";
import { TransactionsDashboard } from "#/features/payment-intents/ui/transactions-dashboard";

export const Route = createFileRoute("/admin/transactions")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: intents = [] } = useQuery(listAllIntentsQueryOptions());
  return (
    <div className="p-6">
      <TransactionsDashboard intents={intents} />
    </div>
  );
}
