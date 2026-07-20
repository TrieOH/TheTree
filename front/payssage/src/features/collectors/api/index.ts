import { queryOptions } from "@tanstack/react-query"
import { authFetcher } from "#/shared/lib/api/fetch"
import type { Collector } from "#/features/oauth/model"

export const collectorsQueryOptions = (organizationId?: string) =>
  queryOptions({
    queryKey: ["collectors", organizationId ?? "personal"],
    queryFn: async () => {
      const path = organizationId
        ? `/organizations/${organizationId}/collectors`
        : "/collectors"
      const response = await authFetcher.get<Collector[]>(path)
      if (!response.success) throw response
      return Array.isArray(response.data) ? response.data : []
    },
  })
