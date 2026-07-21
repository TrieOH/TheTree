import type { UploadTask } from '../model/types'

const DATABASE_NAME = 'univents-upload-queue'
const DATABASE_VERSION = 2
const TASK_STORE = 'tasks'

function openDatabase(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DATABASE_NAME, DATABASE_VERSION)

    request.onupgradeneeded = (event) => {
      const database = request.result
      if (!database.objectStoreNames.contains(TASK_STORE)) {
        database.createObjectStore(TASK_STORE, { keyPath: 'id' })
      }

      if (event.oldVersion < 2) {
        const store = request.transaction?.objectStore(TASK_STORE)
        const cursorRequest = store?.openCursor()
        if (cursorRequest) {
          cursorRequest.onsuccess = () => {
            const cursor = cursorRequest.result
            if (!cursor) return
            if (String(cursor.primaryKey).startsWith('upload-demo-')) {
              cursor.delete()
            }
            cursor.continue()
          }
        }
      }
    }

    request.onsuccess = () => resolve(request.result)
    request.onerror = () =>
      reject(
        request.error ?? new Error('Não foi possível abrir a fila de uploads.'),
      )
  })
}

function requestResult<T>(request: IDBRequest<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result)
    request.onerror = () =>
      reject(
        request.error ??
          new Error('Não foi possível acessar a fila de uploads.'),
      )
  })
}

export async function readUploadTasks(): Promise<UploadTask[]> {
  const database = await openDatabase()
  try {
    const transaction = database.transaction(TASK_STORE, 'readonly')
    return await requestResult(transaction.objectStore(TASK_STORE).getAll())
  } finally {
    database.close()
  }
}

export async function writeUploadTask(task: UploadTask): Promise<void> {
  const database = await openDatabase()
  try {
    const transaction = database.transaction(TASK_STORE, 'readwrite')
    await requestResult(transaction.objectStore(TASK_STORE).put(task))
  } finally {
    database.close()
  }
}

export async function deleteUploadTask(taskId: string): Promise<void> {
  const database = await openDatabase()
  try {
    const transaction = database.transaction(TASK_STORE, 'readwrite')
    await requestResult(transaction.objectStore(TASK_STORE).delete(taskId))
  } finally {
    database.close()
  }
}
