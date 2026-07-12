import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { QueryClient } from '@tanstack/react-query'
import type { EditionCreateOutputI, EditionI } from '../model'
import {
  connectPaymentAccountToEditionFn,
  createEditionFn,
  disconnectPaymentAccountToEditionFn,
  patchEditionFn,
  publishEditionFn,
} from './index'
import { editionKeys } from './query-keys'

type CreateEditionInput = {
  eventId: string
  data: EditionCreateOutputI
}

type PublishEditionInput = {
  eventId: string
  editionId: string
}

type UpdateEditionInput = {
  eventId: string
  editionId: string
  data: EditionCreateOutputI
}

type ConnectPaymentInput = {
  eventId: string
  editionId: string
  credentialId: string
  provider: string
  publicKey: string
}

type DisconnectPaymentInput = {
  eventId: string
  editionId: string
}

function upsertById(editions: EditionI[] | undefined, edition: EditionI) {
  const list = editions ?? []
  const index = list.findIndex((item) => item.id === edition.id)

  if (index === -1) return [...list, edition]

  const next = [...list]
  next[index] = edition
  return next
}

function syncEditionCaches(queryClient: QueryClient, edition: EditionI) {
  queryClient.setQueryData<EditionI[]>(
    editionKeys.adminListByEvent(edition.event_id),
    (old) => upsertById(old, edition),
  )

  queryClient.setQueryData<EditionI[]>(
    editionKeys.publicListByEvent(edition.event_id),
    (old) => upsertById(old, edition),
  )
}

function syncEditionStatusInCaches(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  status: EditionI['status'],
) {
  const patch = (editions: EditionI[] | undefined) =>
    (editions ?? []).map((edition) => (
      edition.id === editionId
        ? { ...edition, status }
        : edition
    ))

  queryClient.setQueryData<EditionI[]>(editionKeys.adminListByEvent(eventId), patch)
  queryClient.setQueryData<EditionI[]>(editionKeys.publicListByEvent(eventId), patch)
}

function syncEditionPaymentInCaches(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  payment: Pick<EditionI, 'trie_payments_credential_id' | 'trie_payments_provider' | 'trie_payments_provider_public_key'>,
) {
  const patch = (editions: EditionI[] | undefined) =>
    (editions ?? []).map((edition) => (
      edition.id === editionId
        ? { ...edition, ...payment }
        : edition
    ))

  queryClient.setQueryData<EditionI[]>(editionKeys.adminListByEvent(eventId), patch)
  queryClient.setQueryData<EditionI[]>(editionKeys.publicListByEvent(eventId), patch)
}

export function useCreateEditionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, data }: CreateEditionInput) => createEditionFn(data, eventId),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao criar edição')
        return
      }

      syncEditionCaches(queryClient, res.data)
      toast.success('Edição criada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function usePublishEditionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId }: PublishEditionInput) => publishEditionFn(eventId, editionId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar edição')
        return
      }

      syncEditionStatusInCaches(queryClient, variables.eventId, variables.editionId, 'announced')
      toast.success('Edição publicada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useUpdateEditionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, data }: UpdateEditionInput) =>
      patchEditionFn(eventId, editionId, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao atualizar edição')
        return
      }

      syncEditionCaches(queryClient, res.data)
      toast.success('Edição atualizada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useConnectPaymentAccountToEditionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, credentialId, provider, publicKey }: ConnectPaymentInput) =>
      connectPaymentAccountToEditionFn(eventId, editionId, credentialId, provider, publicKey),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao conectar pagamento')
        return
      }

      syncEditionPaymentInCaches(queryClient, variables.eventId, variables.editionId, {
        trie_payments_credential_id: variables.credentialId,
        trie_payments_provider: variables.provider,
        trie_payments_provider_public_key: variables.publicKey,
      })
      toast.success('Pagamento conectado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useDisconnectPaymentAccountToEditionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId }: DisconnectPaymentInput) =>
      disconnectPaymentAccountToEditionFn(eventId, editionId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao desconectar pagamento')
        return
      }

      syncEditionPaymentInCaches(queryClient, variables.eventId, variables.editionId, {
        trie_payments_credential_id: null,
        trie_payments_provider: null,
        trie_payments_provider_public_key: null,
      })
      toast.success('Pagamento desconectado com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}
