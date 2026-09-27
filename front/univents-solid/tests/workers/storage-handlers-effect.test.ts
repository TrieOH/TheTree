import { env } from "cloudflare:workers";
import { describe, expect, it } from "vitest";
import { Cause, Effect, Exit } from "effect";
import {
  handleStorageUploadEffect,
  StorageConfigError,
  validateEnvEffect,
} from "../../src/features/storage/api/storage-handlers";

describe("storage-handlers Effect pipeline", () => {
  describe("validateEnvEffect", () => {
    it("succeeds when all required environment variables are present", async () => {
      const exit = await Effect.runPromiseExit(validateEnvEffect(env));
      expect(Exit.isSuccess(exit)).toBe(true);
    });

    it("fails with StorageConfigError when an env var is missing", async () => {
      const brokenEnv = { ...env, S3_BUCKET: "" };
      const exit = await Effect.runPromiseExit(validateEnvEffect(brokenEnv));
      expect(Exit.isFailure(exit)).toBe(true);
      if (Exit.isFailure(exit)) {
        const res = Cause.findError(exit.cause) as { _tag: string; success?: unknown };
        expect(res._tag).toBe("Success");
        expect(res.success).toBeInstanceOf(StorageConfigError);
        expect((res.success as StorageConfigError).message).toContain("S3_BUCKET");
      }
    });
  });

  describe("handleStorageUploadEffect", () => {
    it("returns status 400 when JSON body is invalid", async () => {
      const request = new Request("http://localhost/storage/upload", {
        method: "POST",
        body: "invalid-json{",
        headers: { "Content-Type": "application/json" },
      });

      const response = await Effect.runPromise(handleStorageUploadEffect(request, env));
      expect(response.status).toBe(400);
      const data = (await response.json()) as { error: string };
      expect(data.error).toBe("Invalid JSON payload");
    });
  });
});
