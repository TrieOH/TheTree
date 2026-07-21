import { preprocessImageUpload } from '@/features/storage/api'
import {
  getUploadAssociationHandler,
  subscribeUploadAssociationHandlers,
} from './association-registry'
import { classifyUploadError } from './errors'
import { getRetryDelay, uploadQueueConfig } from './config'
import { uploadQueueStore } from './store'
import type { UploadTask } from '../model/types'

const QUEUE_LOCK_NAME = 'univents-upload-queue-processor'

class UploadQueueProcessor {
  private started = false
  private running = false
  private generation = 0
  private timer?: ReturnType<typeof setTimeout>
  private unsubscribe?: () => void
  private unsubscribeHandlers?: () => void

  start = async () => {
    if (this.started) return
    this.started = true
    const generation = ++this.generation
    await uploadQueueStore.initialize()
    if (!this.started || generation !== this.generation) return
    this.unsubscribe = uploadQueueStore.subscribe(this.schedule)
    this.unsubscribeHandlers = subscribeUploadAssociationHandlers(this.schedule)
    window.addEventListener('online', this.schedule)
    this.schedule()
  }

  stop = () => {
    this.started = false
    this.generation += 1
    this.unsubscribe?.()
    this.unsubscribeHandlers?.()
    window.removeEventListener('online', this.schedule)
    if (this.timer) clearTimeout(this.timer)
  }

  wake = () => this.schedule()

  private schedule = () => {
    if (!this.started || this.running) return
    if (this.timer) clearTimeout(this.timer)

    const now = Date.now()
    const pending = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) =>
          task.status === 'queued' ||
          task.status === 'waiting_retry' ||
          (task.status === 'paused' &&
            task.error?.code === 'ASSOCIATION_HANDLER_UNAVAILABLE' &&
            task.association &&
            getUploadAssociationHandler(task.association.handlerKey)),
      )

    if (pending.length === 0) return
    const nextTime = Math.min(
      ...pending.map((task) => task.nextAttemptAt ?? now),
    )
    this.timer = setTimeout(() => void this.pump(), Math.max(nextTime - now, 0))
  }

  private pump = async () => {
    if (this.running || !navigator.onLine) return
    this.running = true
    let acquiredProcessor = false

    try {
      if ('locks' in navigator) {
        await navigator.locks.request(
          QUEUE_LOCK_NAME,
          { ifAvailable: true },
          async (lock) => {
            if (!lock) return
            acquiredProcessor = true
            await this.processDueTasks()
          },
        )
      } else {
        acquiredProcessor = true
        await this.processDueTasks()
      }
    } finally {
      this.running = false
      if (acquiredProcessor) {
        this.schedule()
      } else {
        this.timer = setTimeout(() => void this.pump(), 500)
      }
    }
  }

  private processDueTasks = async () => {
    await uploadQueueStore.reload()
    const interruptedTasks = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) => task.status === 'uploading' || task.status === 'associating',
      )

    await Promise.all(
      interruptedTasks.map((task) =>
        uploadQueueStore.update(task.id, (value) => ({
          ...value,
          status: 'queued',
          nextAttemptAt: undefined,
          updatedAt: Date.now(),
        })),
      ),
    )

    const now = Date.now()
    const dueTasks = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) =>
          (task.status === 'queued' ||
            task.status === 'waiting_retry' ||
            (task.status === 'paused' &&
              task.error?.code === 'ASSOCIATION_HANDLER_UNAVAILABLE' &&
              task.association &&
              getUploadAssociationHandler(task.association.handlerKey))) &&
          (task.nextAttemptAt ?? 0) <= now,
      )

    for (const task of dueTasks) {
      if (!navigator.onLine) break
      await this.processTask(task)
    }
  }

  private processTask = async (task: UploadTask) => {
    try {
      let current = task

      if (current.stage === 'upload') {
        const claimedTask = await uploadQueueStore.update(
          current.id,
          (value) => ({
            ...value,
            status: 'uploading',
            error: undefined,
            nextAttemptAt: undefined,
            updatedAt: Date.now(),
          }),
        )
        if (!claimedTask) return
        current = claimedTask

        const uploadedUrl = await preprocessImageUpload(
          new File([current.file], current.fileName, {
            type: current.contentType,
          }),
          current.storagePath,
          current.id,
        )

        const uploadedTask = await uploadQueueStore.update(
          current.id,
          (value) => ({
            ...value,
            uploadedUrl,
            stage: 'association',
            status: 'queued',
            retryCount: 0,
            updatedAt: Date.now(),
          }),
        )
        if (!uploadedTask) return
        current = uploadedTask
      }

      if (!current.association) {
        await this.complete(current.id)
        return
      }

      const handler = getUploadAssociationHandler(
        current.association.handlerKey,
      )
      if (!handler) {
        await uploadQueueStore.update(current.id, (value) => ({
          ...value,
          status: 'paused',
          error: {
            kind: 'configuration',
            code: 'ASSOCIATION_HANDLER_UNAVAILABLE',
            message:
              'A integração necessária para associar esta imagem não está disponível.',
            retryable: false,
            requiresReplacement: false,
            occurredAt: Date.now(),
          },
          updatedAt: Date.now(),
        }))
        return
      }

      const associatingTask = await uploadQueueStore.update(
        current.id,
        (value) => ({
          ...value,
          status: 'associating',
          error: undefined,
          nextAttemptAt: undefined,
          updatedAt: Date.now(),
        }),
      )
      if (!associatingTask) return
      current = associatingTask

      await handler(current, current.uploadedUrl!)
      await this.complete(current.id)
    } catch (error) {
      await this.handleFailure(task.id, error)
    }
  }

  private complete = async (taskId: string) => {
    await uploadQueueStore.update(taskId, (task) => ({
      ...task,
      file: new Blob(),
      status: 'completed',
      error: undefined,
      nextAttemptAt: undefined,
      completedAt: Date.now(),
      updatedAt: Date.now(),
    }))
  }

  private handleFailure = async (taskId: string, error: unknown) => {
    const classified = classifyUploadError(error)
    await uploadQueueStore.update(taskId, (task) => {
      if (!classified.retryable) {
        return {
          ...task,
          status: classified.requiresReplacement ? 'rejected' : 'failed',
          error: classified,
          nextAttemptAt: undefined,
          updatedAt: Date.now(),
        }
      }

      const retryCount = task.retryCount + 1
      if (retryCount > uploadQueueConfig.maxRetries) {
        return {
          ...task,
          status: 'failed',
          retryCount,
          error: {
            ...classified,
            message: `${classified.message} O limite de tentativas foi atingido.`,
          },
          nextAttemptAt: undefined,
          updatedAt: Date.now(),
        }
      }

      return {
        ...task,
        status: 'waiting_retry',
        retryCount,
        error: classified,
        nextAttemptAt: Date.now() + getRetryDelay(retryCount),
        updatedAt: Date.now(),
      }
    })
  }
}

export const uploadQueueProcessor = new UploadQueueProcessor()
