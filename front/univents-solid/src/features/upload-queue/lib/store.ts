import type {
  EnqueueUploadInput,
  UploadQueueSnapshot,
  UploadTask,
} from "../model/types";
import { deleteUploadTask, readUploadTasks, writeUploadTask } from "./database";

type Listener = () => void;

class UploadQueueStore {
  private snapshot: UploadQueueSnapshot = { tasks: [], initialized: false };
  private allTasks: UploadTask[] = [];
  private activeAccountId?: string;
  private listeners = new Set<Listener>();
  private initialization?: Promise<void>;
  private channel?: BroadcastChannel;

  getSnapshot = () => this.snapshot;

  getActiveAccountId = () => this.activeAccountId;

  setActiveAccount = (accountId?: string) => {
    if (this.activeAccountId === accountId) return;
    this.activeAccountId = accountId;
    this.publishSnapshot(this.snapshot.initialized);
  };

  subscribe = (listener: Listener) => {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  };

  initialize = async () => {
    if (typeof indexedDB === "undefined") return;
    if (this.initialization) return this.initialization;

    this.initialization = (async () => {
      if (typeof BroadcastChannel !== "undefined") {
        this.channel = new BroadcastChannel("univents-upload-queue");
        this.channel.addEventListener("message", () => void this.reload(false));
      }

      this.allTasks = await readUploadTasks();
      this.publishSnapshot(true);
    })();

    return this.initialization;
  };

  reload = async (broadcast = false) => {
    if (typeof indexedDB === "undefined") return;
    this.allTasks = await readUploadTasks();
    this.publishSnapshot(true);
    if (broadcast) this.channel?.postMessage("changed");
  };

  enqueue = async (input: EnqueueUploadInput) => {
    await this.initialize();
    if (!this.activeAccountId) {
      throw new Error(
        "É necessário estar autenticado para adicionar um upload.",
      );
    }
    const now = Date.now();
    const task: UploadTask = {
      id: crypto.randomUUID(),
      accountId: this.activeAccountId,
      owner: input.owner,
      mediaType: input.mediaType,
      label: input.label ?? input.file.name,
      file: input.file,
      fileName: input.file.name,
      contentType: input.file.type,
      size: input.file.size,
      storagePath: input.storagePath,
      correctionPath: input.correctionPath,
      association: input.association,
      stage: "upload",
      status: "queued",
      retryCount: 0,
      createdAt: now,
      updatedAt: now,
    };

    await this.persist(task);
    return task;
  };

  update = async (
    taskId: string,
    updater: (task: UploadTask) => UploadTask,
  ) => {
    await this.initialize();
    const current = this.snapshot.tasks.find((task) => task.id === taskId);
    if (!current) return undefined;
    const next = updater(current);
    await this.persist(next);
    return next;
  };

  retry = async (taskId: string) => {
    return this.update(taskId, (task) => ({
      ...task,
      status: "queued",
      nextAttemptAt: undefined,
      error: undefined,
      retryCount: 0,
      updatedAt: Date.now(),
    }));
  };

  replaceFile = async (taskId: string, file: File) => {
    return this.update(taskId, (task) => ({
      ...task,
      file,
      fileName: file.name,
      contentType: file.type,
      size: file.size,
      label: file.name,
      stage: "upload",
      status: "queued",
      uploadedUrl: undefined,
      nextAttemptAt: undefined,
      error: undefined,
      retryCount: 0,
      updatedAt: Date.now(),
    }));
  };

  remove = async (taskId: string) => {
    await this.initialize();
    const task = this.snapshot.tasks.find((item) => item.id === taskId);
    if (!task) return;

    this.allTasks = this.allTasks.filter((item) => item.id !== taskId);
    await deleteUploadTask(taskId);
    this.publishSnapshot(true);
    this.channel?.postMessage("changed");
  };

  private persist = async (task: UploadTask) => {
    const index = this.allTasks.findIndex((item) => item.id === task.id);
    if (index >= 0) {
      this.allTasks[index] = task;
    } else {
      this.allTasks.unshift(task);
    }
    await writeUploadTask(task);
    this.publishSnapshot(true);
    this.channel?.postMessage("changed");
  };

  private publishSnapshot = (initialized: boolean) => {
    const tasks = this.activeAccountId
      ? this.allTasks.filter((task) => task.accountId === this.activeAccountId)
      : [];
    this.snapshot = { tasks, initialized };
    for (const listener of this.listeners) {
      listener();
    }
  };
}

export const uploadQueueStore = new UploadQueueStore();
