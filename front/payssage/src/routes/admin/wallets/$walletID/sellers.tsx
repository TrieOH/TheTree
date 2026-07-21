import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { sellersQueryOptions } from "#/features/sellers/api"
import { ProviderCredentialList } from "#/shared/ui/provider-credential-list"

export const Route = createFileRoute("/admin/wallets/$walletID/sellers")({ component: SellersPage })

function SellersPage() {
  const { walletID } = Route.useParams()
  const options = sellersQueryOptions(walletID)
  const { data = [] } = useQuery(options)
  return <ProviderCredentialList items={data} flow="seller" queryKey={options.queryKey} />
}
