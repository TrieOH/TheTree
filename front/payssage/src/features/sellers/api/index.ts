import { queryOptions } from "@tanstack/react-query";
import type { Seller } from "#/features/oauth/model";
import { authFetcher } from "#/shared/lib/api/fetch";

export const sellersQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId, "sellers"],
    queryFn: async () => {
      const response = await authFetcher.get<Seller[]>(
        `/wallets/${walletId}/sellers`,
      );
      if (!response.success) throw response;
      return Array.isArray(response.data) ? response.data : [];
    },
  });
