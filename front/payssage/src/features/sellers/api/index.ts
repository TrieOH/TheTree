import { queryOptions } from "@tanstack/react-query"
import { authFetcher } from "#/shared/lib/api/fetch"
import type { Seller } from "#/features/oauth/model"

export const sellersQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId, "sellers"],
    queryFn: async () => {
      const response = await authFetcher.get<Seller[]>(`/wallets/${walletId}/sellers`)
      if (!response.success) throw response
      return Array.isArray(response.data) ? response.data : []
    },
  })
