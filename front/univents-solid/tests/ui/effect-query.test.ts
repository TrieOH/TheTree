import { describe, expect, it } from "vitest";
import { Cause, Effect, Exit, Schedule } from "effect";
import {
  apiEffect,
  isTransientNetworkError,
  MutationBusinessError,
  MutationNetworkError,
  toMutationFn,
} from "../../src/shared/lib/effect-query";

describe("effect-query bridge", () => {
  describe("isTransientNetworkError", () => {
    it("identifies 5xx server errors as transient", () => {
      expect(isTransientNetworkError({ status: 500 })).toBe(true);
      expect(isTransientNetworkError({ status: 502 })).toBe(true);
      expect(isTransientNetworkError({ status: 503 })).toBe(true);
      expect(isTransientNetworkError({ status: 429 })).toBe(true);
    });

    it("identifies 4xx client errors as non-transient", () => {
      expect(isTransientNetworkError({ status: 400 })).toBe(false);
      expect(isTransientNetworkError({ status: 401 })).toBe(false);
      expect(isTransientNetworkError({ status: 404 })).toBe(false);
      expect(isTransientNetworkError({ status: 409 })).toBe(false);
      expect(isTransientNetworkError({ status: 422 })).toBe(false);
    });

    it("identifies fetch TypeError as transient network failure", () => {
      expect(
        isTransientNetworkError(new TypeError("Failed to fetch")),
      ).toBe(true);
    });
  });

  describe("apiEffect", () => {
    it("returns successful value when promise resolves", async () => {
      const effect = apiEffect(async () => ({ id: "123", ok: true }));
      const result = await Effect.runPromise(effect);
      expect(result).toEqual({ id: "123", ok: true });
    });

    it("fails with MutationNetworkError on 5xx or fetch failure", async () => {
      const effect = apiEffect(async () => {
        const err = new Error("Gateway Timeout") as unknown as { status: number; message: string };
        err.status = 504;
        throw err;
      });

      const exit = await Effect.runPromiseExit(effect);
      expect(Exit.isFailure(exit)).toBe(true);
      if (Exit.isFailure(exit)) {
        const res = Cause.findError(exit.cause) as { _tag: string; success?: unknown };
        expect(res._tag).toBe("Success");
        const error = res.success;
        expect(error).toBeInstanceOf(MutationNetworkError);
        expect((error as MutationNetworkError).message).toBe("Gateway Timeout");
      }
    });

    it("fails with MutationBusinessError on 4xx validation or business failure", async () => {
      const effect = apiEffect(async () => {
        const err = new Error("Ingresso esgotado") as unknown as { status: number; message: string };
        err.status = 409;
        throw err;
      });

      const exit = await Effect.runPromiseExit(effect);
      expect(Exit.isFailure(exit)).toBe(true);
      if (Exit.isFailure(exit)) {
        const res = Cause.findError(exit.cause) as { _tag: string; success?: unknown };
        expect(res._tag).toBe("Success");
        const error = res.success;
        expect(error).toBeInstanceOf(MutationBusinessError);
        expect((error as MutationBusinessError).status).toBe(409);
        expect((error as MutationBusinessError).message).toBe("Ingresso esgotado");
      }
    });
  });

  describe("toMutationFn with retry policy", () => {
    it("retries transient failures and resolves when subsequent attempt succeeds", async () => {
      let callCount = 0;
      const effectFn = () =>
        apiEffect(async () => {
          callCount++;
          if (callCount < 2) {
            const err = new Error("Temporary 503") as unknown as { status: number; message: string };
            err.status = 503;
            throw err;
          }
          return { success: true };
        });

      const mutationFn = toMutationFn(effectFn, {
        retryTransient: true,
        retrySchedule: Schedule.recurs(2),
      });

      const result = await mutationFn(undefined);
      expect(result).toEqual({ success: true });
      expect(callCount).toBe(2);
    });

    it("does not retry business failures (fails immediately)", async () => {
      let callCount = 0;
      const effectFn = () =>
        apiEffect(async () => {
          callCount++;
          const err = new Error("Cart contains invalid items") as unknown as { status: number; message: string };
          err.status = 422;
          throw err;
        });

      const mutationFn = toMutationFn(effectFn, {
        retryTransient: true,
        retrySchedule: Schedule.recurs(3),
      });

      await expect(mutationFn(undefined)).rejects.toThrow("Cart contains invalid items");
      expect(callCount).toBe(1);
    });
  });
});
