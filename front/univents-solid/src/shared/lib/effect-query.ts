import { useMutation } from "@trieoh/front-core-solid";
import type { MutationObserverOptions } from "@tanstack/query-core";
import { Data, Effect, Schedule } from "effect";

export class MutationNetworkError extends Data.TaggedError("MutationNetworkError")<{
  readonly cause: unknown;
  readonly message: string;
}> {}

export class MutationBusinessError extends Data.TaggedError("MutationBusinessError")<{
  readonly status: number;
  readonly message: string;
  readonly cause?: unknown;
}> {}

export type ApiMutationError = MutationNetworkError | MutationBusinessError;

export const isTransientNetworkError = (err: unknown): boolean => {
  if (!err) return false;
  if (typeof err === "object" && err !== null && "status" in err) {
    const status = (err as { status: unknown }).status;
    if (typeof status === "number") {
      return status >= 500 || status === 429;
    }
  }
  if (err instanceof TypeError && err.message.toLowerCase().includes("fetch")) {
    return true;
  }
  return false;
};

/**
 * Wraps an async API call into an Effect, classifying failures into
 * MutationNetworkError (retryable) or MutationBusinessError (terminal).
 */
export const apiEffect = <T>(
  fn: () => Promise<T>,
): Effect.Effect<T, MutationNetworkError | MutationBusinessError> =>
  Effect.tryPromise({
    try: fn,
    catch: (cause) => {
      const isTransient = isTransientNetworkError(cause);
      const message =
        cause instanceof Error
          ? cause.message
          : typeof cause === "object" && cause !== null && "message" in cause
            ? String((cause as { message: unknown }).message)
            : "Erro inesperado na operação.";

      if (isTransient) {
        return new MutationNetworkError({ cause, message });
      }
      const status =
        typeof cause === "object" && cause !== null && "status" in cause
          ? Number((cause as { status: unknown }).status)
          : 400;
      return new MutationBusinessError({ status, message, cause });
    },
  });

export const defaultMutationRetrySchedule = Schedule.exponential("200 millis").pipe(
  Schedule.jittered,
  Schedule.upTo({ times: 2 }),
);

export interface ToMutationFnOptions {
  readonly retryTransient?: boolean;
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
}

export const toMutationFn = <TVariables, TData, TError = unknown>(
  effectFn: (variables: TVariables) => Effect.Effect<TData, TError, never>,
  options?: ToMutationFnOptions,
): ((variables: TVariables) => Promise<TData>) => {
  return (variables: TVariables) => {
    let eff = effectFn(variables);
    if (options?.retryTransient) {
      const schedule = options.retrySchedule ?? defaultMutationRetrySchedule;
      eff = Effect.retry(eff, {
        schedule,
        while: (err) =>
          typeof err === "object" &&
          err !== null &&
          "_tag" in err &&
          (err as { _tag: string })._tag === "MutationNetworkError",
      }) as typeof eff;
    }
    return Effect.runPromise(eff);
  };
};

export interface UseEffectMutationOptions<TData, TError, TVariables>
  extends Omit<
    MutationObserverOptions<TData, TError, TVariables>,
    "mutationKey" | "mutationFn"
  > {
  readonly mutationEffect: (variables: TVariables) => Effect.Effect<TData, TError, never>;
  readonly retryTransient?: boolean;
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
}

export function useEffectMutation<TData, TError = Error, TVariables = void>(
  options: UseEffectMutationOptions<TData, TError, TVariables>,
) {
  const { mutationEffect, retryTransient, retrySchedule, ...rest } = options;
  return useMutation<TData, TError, TVariables>({
    ...rest,
    mutationFn: toMutationFn(mutationEffect, { retryTransient, retrySchedule }),
  });
}
