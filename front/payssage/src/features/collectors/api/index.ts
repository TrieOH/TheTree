import { queryOptions } from "@tanstack/react-query";
import { itemsOf } from "@trieoh/api-client";
import type { Collector } from "#/features/oauth/model";
import { authFetcher } from "#/shared/lib/api/fetch";

export const collectorsQueryOptions = (organizationId?: string) =>
  queryOptions({
    queryKey: ["collectors", organizationId ?? "personal"],
    queryFn: async () => {
      const path = organizationId
        ? `/organizations/${organizationId}/collectors`
        : "/collectors";
      const response = await authFetcher.get<unknown>(path);
      if (!response.success) throw response;
      return itemsOf<Collector>(response.data);
    },
  });
