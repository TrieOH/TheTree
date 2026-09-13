import type { UploadTask } from "../model/types";

export type UploadAssociationHandler = (
  task: UploadTask,
  uploadedUrl: string,
) => Promise<void>;

const associationHandlers = new Map<string, UploadAssociationHandler>();
const listeners = new Set<() => void>();

function notifyListeners() {
  for (const listener of listeners) {
    listener();
  }
}

export function registerUploadAssociationHandler(
  handlerKey: string,
  handler: UploadAssociationHandler,
) {
  associationHandlers.set(handlerKey, handler);
  notifyListeners();
  return () => {
    if (associationHandlers.get(handlerKey) === handler) {
      associationHandlers.delete(handlerKey);
      notifyListeners();
    }
  };
}

export function getUploadAssociationHandler(handlerKey: string) {
  return associationHandlers.get(handlerKey);
}

export function subscribeUploadAssociationHandlers(listener: () => void) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}
