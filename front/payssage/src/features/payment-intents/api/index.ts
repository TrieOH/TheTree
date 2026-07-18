import { createClientOnlyFn } from '@tanstack/react-start'
import type { CreateIntentRequest, Intent } from '../model'
import { authFetcher } from '#/shared/lib/api/fetch'

export const createWalletIntentFn = createClientOnlyFn(
  (walletId: string, payload: CreateIntentRequest) =>
    authFetcher.post<Intent>(`/wallets/${walletId}/intents`, payload),
)
