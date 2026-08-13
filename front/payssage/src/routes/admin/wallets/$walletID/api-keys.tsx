import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { WalletApiKey } from "#/features/keys/ui/wallet-api-key";
import { walletByIdQueryOptions } from "#/features/wallets/api";

export const Route = createFileRoute("/admin/wallets/$walletID/api-keys")({
  component: WalletApiKeysPage,
});

function WalletApiKeysPage() {
  const { walletID } = Route.useParams();
  const { data: wallet } = useQuery(walletByIdQueryOptions(walletID));

  return wallet ? (
    <div className="space-y-4">
      <div>
        <h2 className="text-base font-semibold">Wallet API keys</h2>
        <p className="text-sm text-muted-foreground">
          Create a key owned by this wallet owner. The secret is shown only
          once.
        </p>
      </div>
      <WalletApiKey wallet={wallet} />
    </div>
  ) : null;
}
