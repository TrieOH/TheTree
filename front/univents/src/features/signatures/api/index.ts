import { createClientOnlyFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { SignatureCreateI, SignatureI } from '@/features/signatures/model'
import { authFetcher, tanstackQueryFetcher } from '@/shared/lib/api/fetch'

export const createSignatureFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  payload: SignatureCreateI
) => {
  return authFetcher.post<SignatureI>(
    `/events/${eventId}/editions/${editionId}/signatures`,
    payload
  )
})

export const getAllSignaturesFn = createClientOnlyFn(async (
  eventId: string,
  editionId: string
) => {
  try {
    return await tanstackQueryFetcher<SignatureI[]>(
      `/events/${eventId}/editions/${editionId}/signatures`
    )
  } catch {
    return []
  }
})

export const allSignaturesQueryOptions = (eventId: string, editionId: string) => queryOptions({
  queryKey: ['signatures', eventId, editionId],
  queryFn: () => getAllSignaturesFn(eventId, editionId),
})

export const getSignatureFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  sigId: string
) => {
  return tanstackQueryFetcher<SignatureI>(
    `/events/${eventId}/editions/${editionId}/signatures/${sigId}`
  )
})

export const signatureQueryOptions = (eventId: string, editionId: string, sigId: string) => queryOptions({
  queryKey: ['signatures', eventId, editionId, sigId],
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
