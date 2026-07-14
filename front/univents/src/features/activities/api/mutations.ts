import { useMutation, useQueryClient, type QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { ActivityCreateOutputI, ActivityI } from '../model'
import { createActivityFn, publishActivityFn, updateActivityFn } from './index'
import { activityKeys } from './query-keys'

type CreateActivityInput = {
  eventId: string
  editionId: string
  data: ActivityCreateOutputI
}

type UpdateActivityInput = {
  eventId: string
  editionId: string
  activityId: string
  data: ActivityCreateOutputI
}

type PublishActivityInput = {
  eventId: string
  editionId: string
  activityId: string
}

function upsertById(activities: ActivityI[] | undefined, activity: ActivityI) {
  const list = activities ?? []
  const index = list.findIndex((item) => item.id === activity.id)

  if (index === -1) return [...list, activity]

  const next = [...list]
  next[index] = activity
  return next
}

function syncActivityCaches(queryClient: QueryClient, eventId: string, editionId: string, activity: ActivityI) {
  queryClient.setQueryData<ActivityI[]>(
    activityKeys.adminListByEdition(eventId, editionId),
    (old) => upsertById(old, activity),
  )

  queryClient.setQueryData<ActivityI[]>(
    activityKeys.publicListByEdition(eventId, editionId),
    (old) => upsertById(old, activity),
  )
}

function syncActivityStatusInCaches(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  activityId: string,
  status: ActivityI['status'],
) {
  const patch = (activities: ActivityI[] | undefined) =>
    (activities ?? []).map((activity) => (
      activity.id === activityId
        ? { ...activity, status }
        : activity
    ))

  queryClient.setQueryData<ActivityI[]>(activityKeys.adminListByEdition(eventId, editionId), patch)
  queryClient.setQueryData<ActivityI[]>(activityKeys.publicListByEdition(eventId, editionId), patch)
}

export function useCreateActivityMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, data }: CreateActivityInput) =>
      createActivityFn(data, eventId, editionId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao criar atividade')
        return
      }

      syncActivityCaches(queryClient, variables.eventId, variables.editionId, res.data)
      void queryClient.invalidateQueries({ queryKey: activityKeys.adminListByEdition(variables.eventId, variables.editionId) })
      toast.success('Atividade criada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function useUpdateActivityMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, activityId, data }: UpdateActivityInput) =>
      updateActivityFn(activityId, data, eventId, editionId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao atualizar atividade')
        return
      }

      syncActivityCaches(queryClient, variables.eventId, variables.editionId, res.data)
      void queryClient.invalidateQueries({ queryKey: activityKeys.adminListByEdition(variables.eventId, variables.editionId) })
      toast.success('Atividade atualizada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}

export function usePublishActivityMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ eventId, editionId, activityId }: PublishActivityInput) =>
      publishActivityFn(eventId, editionId, activityId),
    onSuccess: (res, variables) => {
      if (!res.success) {
        toast.error(res.message || 'Erro ao publicar atividade')
        return
      }

      syncActivityStatusInCaches(
        queryClient,
        variables.eventId,
        variables.editionId,
        variables.activityId,
        'published',
      )
      toast.success('Atividade publicada com sucesso!')
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })
}
