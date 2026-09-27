import { createSignal, onCleanup } from "solid-js";
import { uploadQueueProcessor } from "../lib/processor";
import { uploadQueueStore } from "../lib/store";
import type { UploadQueueSnapshot } from "../model/types";

export function useUploadQueue() {
  const [snapshot, setSnapshot] = createSignal<UploadQueueSnapshot>(
    uploadQueueStore.getSnapshot(),
    { ownedWrite: true },
  );

  const unsubscribe = uploadQueueStore.subscribe(() => {
    setSnapshot(uploadQueueStore.getSnapshot());
  });
  onCleanup(unsubscribe);

  return {
    get tasks() {
      return snapshot().tasks;
    },
    get initialized() {
      return snapshot().initialized;
    },
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
