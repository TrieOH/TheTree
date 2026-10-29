import { queryOptions } from "@tanstack/react-query";
import { itemsOf } from "@trieoh/api-client";
import type { Seller } from "#/features/oauth/model";
import { authFetcher } from "#/shared/lib/api/fetch";

export const sellersQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId, "sellers"],
    queryFn: async () => {
      const response = await authFetcher.get<unknown>(
        `/wallets/${walletId}/sellers`,
      );
      if (!response.success) throw response;
      return itemsOf<Seller>(response.data);
    },
  });
