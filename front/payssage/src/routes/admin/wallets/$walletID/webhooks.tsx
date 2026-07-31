import { createFileRoute } from "@tanstack/react-router";
import { WebhooksDashboard } from "#/features/webhooks/ui/webhooks-dashboard";

export const Route = createFileRoute("/admin/wallets/$walletID/webhooks")({
  component: WebhooksPage,
});

function WebhooksPage() {
  const { walletID } = Route.useParams();
  return <WebhooksDashboard walletId={walletID} />;
}
