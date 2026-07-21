import { useEffect, useRef } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { useAuth } from '@trieoh/identityx-sdk-ts/react'
import { uploadQueueProcessor } from '../lib/processor'
import { uploadQueueStore } from '../lib/store'
import { useUploadQueue } from '../hooks/use-upload-queue'
import type { UploadTaskStatus } from '../model/types'

const notificationStatuses = new Set<UploadTaskStatus>([
  'completed',
  'failed',
  'rejected',
])

export function UploadQueueProvider({
  children,
}: {
  children: React.ReactNode
}) {
  const navigate = useNavigate()
  const { auth, isAuthenticated } = useAuth()
  const accountId = isAuthenticated ? auth.profile()?.id : undefined
  const { tasks, initialized, retry } = useUploadQueue()
  const previousStatuses = useRef(new Map<string, UploadTaskStatus>())
  const notificationAccountId = useRef<string | undefined>(undefined)

  useEffect(() => {
    uploadQueueStore.setActiveAccount(accountId)
    previousStatuses.current.clear()

    if (!accountId) {
      uploadQueueProcessor.stop()
      return
    }

    void uploadQueueProcessor.start()
    return () => {
      uploadQueueProcessor.stop()
      uploadQueueStore.setActiveAccount(undefined)
    }
  }, [accountId])

  useEffect(() => {
    if (!initialized || !accountId) return

    if (notificationAccountId.current !== accountId) {
      notificationAccountId.current = accountId
      previousStatuses.current = new Map(
        tasks.map((task) => [task.id, task.status]),
      )
      return
    }

    if (previousStatuses.current.size === 0) {
      previousStatuses.current = new Map(
        tasks.map((task) => [task.id, task.status]),
      )
      return
    }

    for (const task of tasks) {
      const previousStatus = previousStatuses.current.get(task.id)
      if (
        previousStatus === task.status ||
        !notificationStatuses.has(task.status)
      )
        continue

      if (task.status === 'completed') {
        toast.success(`${task.label} foi enviada com sucesso.`, {
          id: `upload-${task.id}`,
        })
        continue
      }

      if (task.status === 'rejected' || task.error?.requiresReplacement) {
        toast.error(
          task.error?.message ?? `É necessário substituir ${task.label}.`,
          {
            id: `upload-${task.id}`,
            action: {
              label: 'Corrigir imagem',
              onClick: () => {
                if (task.correctionPath) {
                  void navigate({ to: task.correctionPath })
                } else {
                  void navigate({
                    to: '/admin/uploads',
                    search: { task: task.id },
                  })
                }
              },
            },
          },
        )
        continue
      }

      toast.error(
        task.error?.message ?? `Não foi possível enviar ${task.label}.`,
        {
          id: `upload-${task.id}`,
          action: task.error?.retryable
            ? {
                label: 'Tentar novamente',
                onClick: () => void retry(task.id),
              }
            : {
                label: 'Ver detalhes',
                onClick: () =>
                  void navigate({
                    to: '/admin/uploads',
                    search: { task: task.id },
                  }),
              },
        },
      )
    }

    previousStatuses.current = new Map(
      tasks.map((task) => [task.id, task.status]),
    )
  }, [accountId, initialized, navigate, retry, tasks])

  return children
}
