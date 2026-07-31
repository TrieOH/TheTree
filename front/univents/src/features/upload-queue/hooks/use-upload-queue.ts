import { useSyncExternalStore } from "react";
import { uploadQueueProcessor } from "../lib/processor";
import { uploadQueueStore } from "../lib/store";
import type { UploadQueueSnapshot } from "../model/types";

const serverSnapshot: UploadQueueSnapshot = { tasks: [], initialized: false };

export function useUploadQueue() {
  const snapshot = useSyncExternalStore(
    uploadQueueStore.subscribe,
    uploadQueueStore.getSnapshot,
    () => serverSnapshot,
  );

  return {
    ...snapshot,
    enqueue: uploadQueueStore.enqueue,
    retry: async (taskId: string) => {
      const task = await uploadQueueStore.retry(taskId);
      uploadQueueProcessor.wake();
      return task;
    },
    replaceFile: async (taskId: string, file: File) => {
      const task = await uploadQueueStore.replaceFile(taskId, file);
      uploadQueueProcessor.wake();
      return task;
    },
    remove: uploadQueueStore.remove,
  };
}
