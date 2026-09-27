import { Data, Effect } from "effect";

export class StorageUnavailableError extends Data.TaggedError("StorageUnavailableError")<{
  readonly storageType: "local" | "session";
  readonly message: string;
}> {}

export class StorageReadError extends Data.TaggedError("StorageReadError")<{
  readonly storageType: "local" | "session";
  readonly key: string;
  readonly cause: unknown;
  readonly message: string;
}> {}

export class StorageWriteError extends Data.TaggedError("StorageWriteError")<{
  readonly storageType: "local" | "session";
  readonly key: string;
  readonly cause: unknown;
  readonly message: string;
}> {}

export class StorageParseError extends Data.TaggedError("StorageParseError")<{
  readonly storageType: "local" | "session";
  readonly key: string;
  readonly rawValue: string;
  readonly cause: unknown;
  readonly message: string;
}> {}

export type StorageError =
  | StorageUnavailableError
  | StorageReadError
  | StorageWriteError
  | StorageParseError;

/**
 * Safely accesses window.localStorage or window.sessionStorage.
 * Probing with a temporary key is the W3C standard approach to detect disabled storage
 * (e.g. Safari private browsing, sandboxed iframes, blocked third-party storage).
 */
export const resolveStorage = (
  type: "local" | "session",
): Effect.Effect<Storage, StorageUnavailableError> =>
  Effect.sync(() => {
    if (typeof window === "undefined") {
      return null;
    }
    try {
      const storage = type === "local" ? window.localStorage : window.sessionStorage;
      if (!storage) return null;
      // Probe to verify write permissions are allowed in current context
      const testKey = `__univents_test__`;
      storage.setItem(testKey, testKey);
      storage.removeItem(testKey);
      return storage;
    } catch {
      // Browser blocked storage access (e.g. SecurityError in opaque iframe)
      return null;
    }
  }).pipe(
    Effect.flatMap((storage) =>
      storage
        ? Effect.succeed(storage)
        : Effect.fail(
            new StorageUnavailableError({
              storageType: type,
              message: `Armazenamento ${type === "local" ? "localStorage" : "sessionStorage"} não disponível.`,
            }),
          ),
    ),
  );

export const getItemEffect = (
  type: "local" | "session",
  key: string,
): Effect.Effect<string | null, StorageUnavailableError | StorageReadError> =>
  Effect.gen(function* () {
    const storage = yield* resolveStorage(type);
    return yield* Effect.try({
      try: () => storage.getItem(key),
      catch: (cause) =>
        new StorageReadError({
          storageType: type,
          key,
          cause,
          message: `Erro ao ler a chave "${key}" de ${type}Storage.`,
        }),
    });
  });

export const setItemEffect = (
  type: "local" | "session",
  key: string,
  value: string,
): Effect.Effect<void, StorageUnavailableError | StorageWriteError> =>
  Effect.gen(function* () {
    const storage = yield* resolveStorage(type);
    return yield* Effect.try({
      try: () => storage.setItem(key, value),
      catch: (cause) =>
        new StorageWriteError({
          storageType: type,
          key,
          cause,
          message: `Erro ao salvar a chave "${key}" em ${type}Storage.`,
        }),
    });
  });

export const removeItemEffect = (
  type: "local" | "session",
  key: string,
): Effect.Effect<void, StorageUnavailableError | StorageWriteError> =>
  Effect.gen(function* () {
    const storage = yield* resolveStorage(type);
    return yield* Effect.try({
      try: () => storage.removeItem(key),
      catch: (cause) =>
        new StorageWriteError({
          storageType: type,
          key,
          cause,
          message: `Erro ao remover a chave "${key}" de ${type}Storage.`,
        }),
    });
  });

/**
 * Consumes and deletes a single-use entry (e.g. one-time WS authentication tokens).
 */
export const takeItemEffect = (
  type: "local" | "session",
  key: string,
): Effect.Effect<string | null, StorageUnavailableError | StorageReadError | StorageWriteError> =>
  Effect.gen(function* () {
    const value = yield* getItemEffect(type, key);
    if (value !== null) {
      yield* removeItemEffect(type, key);
    }
    return value;
  });

export const getJsonEffect = <T>(
  type: "local" | "session",
  key: string,
  fallback: T,
): Effect.Effect<T> =>
  Effect.gen(function* () {
    const raw = yield* getItemEffect(type, key).pipe(
      Effect.orElseSucceed(() => null),
    );
    if (raw === null) return fallback;

    return yield* Effect.try({
      try: () => JSON.parse(raw) as T,
      catch: () => fallback,
    }).pipe(Effect.orElseSucceed(() => fallback));
  });

export const setJsonEffect = <T>(
  type: "local" | "session",
  key: string,
  value: T,
): Effect.Effect<void, StorageUnavailableError | StorageWriteError> =>
  Effect.gen(function* () {
    const serialized = JSON.stringify(value);
    yield* setItemEffect(type, key, serialized);
  });

/* =========================================================================
 * Synchronous Convenience Wrappers (Ideal for Solid.js signals & callbacks)
 * ========================================================================= */

export const getItemSync = (
  type: "local" | "session",
  key: string,
  fallback: string | null = null,
): string | null =>
  Effect.runSync(
    getItemEffect(type, key).pipe(
      Effect.map((val) => val ?? fallback),
      Effect.orElseSucceed(() => fallback),
    ),
  );

export const setItemSync = (
  type: "local" | "session",
  key: string,
  value: string,
): boolean =>
  Effect.runSync(
    setItemEffect(type, key, value).pipe(
      Effect.match({
        onFailure: () => false,
        onSuccess: () => true,
      }),
    ),
  );

export const removeItemSync = (
  type: "local" | "session",
  key: string,
): void => {
  Effect.runSync(
    removeItemEffect(type, key).pipe(
      Effect.orElseSucceed(() => undefined),
    ),
  );
};

export const takeItemSync = (
  type: "local" | "session",
  key: string,
  fallback: string | null = null,
): string | null =>
  Effect.runSync(
    takeItemEffect(type, key).pipe(
      Effect.map((val) => val ?? fallback),
      Effect.orElseSucceed(() => fallback),
    ),
  );

export const getJsonSync = <T>(
  type: "local" | "session",
  key: string,
  fallback: T,
): T =>
  Effect.runSync(getJsonEffect(type, key, fallback));

export const setJsonSync = <T>(
  type: "local" | "session",
  key: string,
  value: T,
): boolean =>
  Effect.runSync(
    setJsonEffect(type, key, value).pipe(
      Effect.match({
        onFailure: () => false,
        onSuccess: () => true,
      }),
    ),
  );

export const appLocalStorage = {
  get: (key: string, fallback: string | null = null) =>
    getItemSync("local", key, fallback),
  set: (key: string, value: string) => setItemSync("local", key, value),
  remove: (key: string) => removeItemSync("local", key),
  take: (key: string, fallback: string | null = null) =>
    takeItemSync("local", key, fallback),
  getJson: <T>(key: string, fallback: T) => getJsonSync<T>("local", key, fallback),
  setJson: <T>(key: string, value: T) => setJsonSync<T>("local", key, value),
  getEffect: (key: string) => getItemEffect("local", key),
  setEffect: (key: string, value: string) => setItemEffect("local", key, value),
  getJsonEffect: <T>(key: string, fallback: T) =>
    getJsonEffect<T>("local", key, fallback),
  setJsonEffect: <T>(key: string, value: T) =>
    setJsonEffect<T>("local", key, value),
};

export const appSessionStorage = {
  get: (key: string, fallback: string | null = null) =>
    getItemSync("session", key, fallback),
  set: (key: string, value: string) => setItemSync("session", key, value),
  remove: (key: string) => removeItemSync("session", key),
  take: (key: string, fallback: string | null = null) =>
    takeItemSync("session", key, fallback),
  getJson: <T>(key: string, fallback: T) =>
    getJsonSync<T>("session", key, fallback),
  setJson: <T>(key: string, value: T) =>
    setJsonSync<T>("session", key, value),
  getEffect: (key: string) => getItemEffect("session", key),
  setEffect: (key: string, value: string) => setItemEffect("session", key, value),
  takeEffect: (key: string) => takeItemEffect("session", key),
};
