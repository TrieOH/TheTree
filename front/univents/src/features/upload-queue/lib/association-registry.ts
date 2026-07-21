import type { UploadTask } from '../model/types'

export type UploadAssociationHandler = (
  task: UploadTask,
  uploadedUrl: string,
) => Promise<void>

const handlers = new Map<string, UploadAssociationHandler>()
const listeners = new Set<() => void>()

export function registerUploadAssociationHandler(
  key: string,
  handler: UploadAssociationHandler,
) {
  handlers.set(key, handler)
  for (const listener of listeners) listener()
  return () => {
    handlers.delete(key)
    for (const listener of listeners) listener()
  }
}

export function getUploadAssociationHandler(key: string) {
  return handlers.get(key)
}

export function subscribeUploadAssociationHandlers(listener: () => void) {
  listeners.add(listener)
  return () => listeners.delete(listener)
}
