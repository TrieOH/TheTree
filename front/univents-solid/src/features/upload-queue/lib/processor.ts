import { Effect } from "effect";
import { preprocessImageUpload } from "@/features/storage/api/index";
import type { UploadTask } from "../model/types";
import {
  getUploadAssociationHandler,
  subscribeUploadAssociationHandlers,
} from "./association-registry";
import { getRetryDelay, uploadQueueConfig } from "./config";
import { classifyUploadError } from "./errors";
import { uploadQueueStore } from "./store";

const QUEUE_LOCK_NAME = "univents-upload-queue-processor";

class UploadQueueProcessor {
  private started = false;
  private running = false;
  private generation = 0;
  private timer?: ReturnType<typeof setTimeout>;
  private unsubscribe?: () => void;
  private unsubscribeHandlers?: () => void;

  start = async () => {
    if (this.started) return;
    this.started = true;
    const generation = ++this.generation;
    await uploadQueueStore.initialize();
    if (!this.started || generation !== this.generation) return;
    this.unsubscribe = uploadQueueStore.subscribe(this.schedule);
    this.unsubscribeHandlers = subscribeUploadAssociationHandlers(
      this.schedule,
    );
    window.addEventListener("online", this.schedule);
    this.schedule();
  };

  stop = () => {
    this.started = false;
    this.generation += 1;
    this.unsubscribe?.();
    this.unsubscribeHandlers?.();
    window.removeEventListener("online", this.schedule);
    if (this.timer) clearTimeout(this.timer);
  };

  wake = () => this.schedule();

  private schedule = () => {
    if (!this.started || this.running) return;
    if (this.timer) clearTimeout(this.timer);

    const now = Date.now();
    const pending = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) =>
          task.status === "queued" ||
          task.status === "waiting_retry" ||
          (task.status === "paused" &&
            task.error?.code === "ASSOCIATION_HANDLER_UNAVAILABLE" &&
            task.association &&
            getUploadAssociationHandler(task.association.handlerKey)),
      );

    if (pending.length === 0) return;
    const nextTime = Math.min(
      ...pending.map((task) => task.nextAttemptAt ?? now),
    );
    this.timer = setTimeout(
      () => void this.pump(),
      Math.max(nextTime - now, 0),
    );
  };

  private pump = async () => {
    if (this.running || !navigator.onLine) return;
    this.running = true;
    let acquiredProcessor = false;

    try {
      if ("locks" in navigator) {
        await navigator.locks.request(
          QUEUE_LOCK_NAME,
          { ifAvailable: true },
          async (lock) => {
            if (!lock) return;
            acquiredProcessor = true;
            await this.processDueTasks();
          },
        );
      } else {
        acquiredProcessor = true;
        await this.processDueTasks();
      }
    } finally {
      this.running = false;
      if (acquiredProcessor) {
        this.schedule();
      } else {
        this.timer = setTimeout(() => void this.pump(), 500);
      }
    }
  };

  private processDueTasks = async () => {
    await uploadQueueStore.reload();
    const interruptedTasks = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) => task.status === "uploading" || task.status === "associating",
      );

    await Promise.all(
      interruptedTasks.map((task) =>
        uploadQueueStore.update(task.id, (value) => ({
          ...value,
          status: "queued",
          nextAttemptAt: undefined,
          updatedAt: Date.now(),
        })),
      ),
    );

    const now = Date.now();
    const dueTasks = uploadQueueStore
      .getSnapshot()
      .tasks.filter(
        (task) =>
          (task.status === "queued" ||
            task.status === "waiting_retry" ||
            (task.status === "paused" &&
              task.error?.code === "ASSOCIATION_HANDLER_UNAVAILABLE" &&
              task.association &&
              getUploadAssociationHandler(task.association.handlerKey))) &&
          (task.nextAttemptAt ?? 0) <= now,
      );

    for (const task of dueTasks) {
      if (!navigator.onLine) break;
      await this.processTask(task);
    }
  };

  private processTask = async (task: UploadTask) => {
    const complete = this.complete;
    const taskPipeline = Effect.gen(function* () {
      let current = task;

      if (current.stage === "upload") {
        const claimedTask = yield* Effect.tryPromise({
          try: () =>
            uploadQueueStore.update(current.id, (value) => ({
              ...value,
              status: "uploading",
              error: undefined,
              nextAttemptAt: undefined,
              updatedAt: Date.now(),
            })),
          catch: (err) => err,
        });
        if (!claimedTask) return;
        current = claimedTask;

        const uploadedUrl = yield* Effect.tryPromise({
          try: () =>
            preprocessImageUpload(
              new File([current.file], current.fileName, {
                type: current.contentType,
              }),
              current.storagePath,
              current.id,
            ),
          catch: (err) => err,
        });

        const uploadedTask = yield* Effect.tryPromise({
          try: () =>
            uploadQueueStore.update(current.id, (value) => ({
              ...value,
              uploadedUrl,
              stage: "association",
              status: "queued",
              retryCount: 0,
              updatedAt: Date.now(),
            })),
          catch: (err) => err,
        });
        if (!uploadedTask) return;
        current = uploadedTask;
      }

      if (!current.association) {
        yield* Effect.tryPromise({
          try: () => complete(current.id),
          catch: (err) => err,
        });
        return;
      }

      const handler = getUploadAssociationHandler(
        current.association.handlerKey,
      );
      if (!handler) {
        yield* Effect.tryPromise({
          try: () =>
            uploadQueueStore.update(current.id, (value) => ({
              ...value,
              status: "paused",
              error: {
                kind: "configuration",
                code: "ASSOCIATION_HANDLER_UNAVAILABLE",
                message:
                  "A integração necessária para associar esta imagem não está disponível.",
                retryable: true,
                requiresReplacement: false,
                occurredAt: Date.now(),
              },
              updatedAt: Date.now(),
            })),
          catch: (err) => err,
        });
        return;
      }

      const claimedAssociationTask = yield* Effect.tryPromise({
        try: () =>
          uploadQueueStore.update(current.id, (value) => ({
            ...value,
            status: "associating",
            error: undefined,
            nextAttemptAt: undefined,
            updatedAt: Date.now(),
          })),
        catch: (err) => err,
      });
      if (!claimedAssociationTask) return;
      current = claimedAssociationTask;

      yield* Effect.tryPromise({
        try: () => handler(current, current.uploadedUrl!),
        catch: (err) => err,
      });
      yield* Effect.tryPromise({
        try: () => complete(current.id),
        catch: (err) => err,
      });
    });

    const outcome = await Effect.runPromise(
      Effect.match(taskPipeline, {
        onFailure: (err) => ({ ok: false as const, err }),
        onSuccess: () => ({ ok: true as const }),
      }),
    );

    if (!outcome.ok) {
      await this.handleFailure(task.id, outcome.err);
    }
  };

  private complete = async (taskId: string) => {
    await uploadQueueStore.update(taskId, (task) => ({
      ...task,
      status: "completed",
      error: undefined,
      nextAttemptAt: undefined,
      updatedAt: Date.now(),
    }));
  };

  private handleFailure = async (taskId: string, error: unknown) => {
    const classified = classifyUploadError(error);
    const task = uploadQueueStore
      .getSnapshot()
      .tasks.find((item) => item.id === taskId);
    if (!task) return;

    if (classified.requiresReplacement) {
      await uploadQueueStore.update(taskId, (value) => ({
        ...value,
        status: "rejected",
        error: classified,
        nextAttemptAt: undefined,
        updatedAt: Date.now(),
      }));
      return;
    }

    const nextRetryCount = task.retryCount + 1;
    if (
      !classified.retryable ||
      nextRetryCount >= uploadQueueConfig.maxRetries
    ) {
      await uploadQueueStore.update(taskId, (value) => ({
        ...value,
        status: "failed",
        error: classified,
        nextAttemptAt: undefined,
        updatedAt: Date.now(),
      }));
      return;
    }

    const delay = getRetryDelay(nextRetryCount);
    await uploadQueueStore.update(taskId, (value) => ({
      ...value,
      status: "waiting_retry",
      error: classified,
      retryCount: nextRetryCount,
      nextAttemptAt: Date.now() + delay,
      updatedAt: Date.now(),
    }));
    this.schedule();
  };
}

export const uploadQueueProcessor = new UploadQueueProcessor();
