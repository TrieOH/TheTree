type FlushTask = () => Promise<void>;

const tasks = new Map<string, FlushTask>();

export function registerImageUploadTask(id: string, task: FlushTask) {
  tasks.set(id, task);

  return () => {
    tasks.delete(id);
  };
}

export async function flushImageUploadTasks() {
  for (const task of tasks.values()) await task();
}

export function clearImageUploadTasks() {
  tasks.clear();
}
