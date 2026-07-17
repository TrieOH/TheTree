import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { walletByIdQueryOptions } from "#/features/wallets/api"
import { WalletCollector } from "#/features/wallets/ui/wallet-collector"

export const Route = createFileRoute("/admin/wallets/$walletID/")({ component: CollectorPage })

function CollectorPage() {
  const { walletID } = Route.useParams()
  const { data: wallet } = useQuery(walletByIdQueryOptions(walletID))
  return wallet ? <WalletCollector wallet={wallet} /> : null
}
