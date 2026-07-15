import { createClientOnlyFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { SignatureCreateOutputI, SignatureI } from '@/features/signatures/model'
import { authFetcher, authQueryFetcher } from '@/shared/lib/api/fetch'
import { signatureKeys } from './query-keys'

export const createSignatureFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  payload: SignatureCreateOutputI
) => {
  return authFetcher.post<SignatureI>(
    `/events/${eventId}/editions/${editionId}/signatures`,
    payload
  )
})

export const getAllSignaturesFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  return authQueryFetcher<SignatureI[]>(`/events/${eventId}/editions/${editionId}/signatures`)
})

export const allSignaturesQueryOptions = (eventId: string, editionId: string) => queryOptions({
  queryKey: signatureKeys.byEdition(eventId, editionId),
  queryFn: () => getAllSignaturesFn(eventId, editionId),
})

export const getSignatureFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  sigId: string
) => {
  return authQueryFetcher<SignatureI>(`/events/${eventId}/editions/${editionId}/signatures/${sigId}`)
})

export const signatureQueryOptions = (eventId: string, editionId: string, sigId: string) => queryOptions({
  queryKey: signatureKeys.byId(eventId, editionId, sigId),
  queryFn: () => getSignatureFn(eventId, editionId, sigId),
})

export const removeSignatureFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  sigId: string
) => {
  return authFetcher.delete<null>(
    `/events/${eventId}/editions/${editionId}/signatures/${sigId}`
  )
})
