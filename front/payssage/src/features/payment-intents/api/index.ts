import { createClientOnlyFn } from '@tanstack/react-start'
import type { CreateIntentRequest, Intent } from '../model'
import { authFetcher } from '#/shared/lib/api/fetch'
import { queryOptions } from '@tanstack/react-query'

export const listAllIntentsQueryOptions = () =>
  queryOptions({
    queryKey: ["intents", "personal"],
    queryFn: async () => {
      const response = await authFetcher.get<Intent[]>("/intents")
      if (!response.success) throw response
      return Array.isArray(response.data) ? response.data : []
    },
  })

export const listAllByWalletIntentsQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["intents", walletId],
    queryFn: async () => {
      const response = await authFetcher.get<Intent[]>(`/wallets/${walletId}/intents`)
      if (!response.success) throw response
      return Array.isArray(response.data) ? response.data : []
    },
  })

export const listAllByOrgIntentsQueryOptions = (organizationId: string) =>
  queryOptions({
    queryKey: ["intents", organizationId],
    queryFn: async () => {
      const response = await authFetcher.get<Intent[]>(`/organizations/${organizationId}/intents`)
      if (!response.success) throw response
      return Array.isArray(response.data) ? response.data : []
    },
  })


export const createWalletIntentFn = createClientOnlyFn(
  (walletId: string, payload: CreateIntentRequest) =>
    authFetcher.post<Intent>(`/wallets/${walletId}/intents`, payload),
)
