import { createFileRoute } from "@tanstack/react-router";
import { ProviderConnectSection } from "#/features/oauth/ui/provider-connect-button";

export const Route = createFileRoute("/admin/wallets/$walletID/connect-seller")(
  { component: ConnectSellerPage },
);

function ConnectSellerPage() {
  const { walletID } = Route.useParams();
  return <ProviderConnectSection flow="seller" walletId={walletID} />;
}
