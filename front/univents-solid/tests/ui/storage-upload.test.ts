import { afterEach, describe, expect, it, vi } from "vitest";
import { Effect, Exit } from "effect";
import {
  uploadProfileImage,
  uploadProfileImageEffect,
  uploadProfileImagesBatch,
  validateImageFile,
  validateImageFileSync,
} from "../../src/features/storage/api";

afterEach(() => {
  vi.restoreAllMocks();
});

describe("validateImageFile & validateImageFileSync", () => {
  it("accepts valid image types under 10MB", () => {
    const validFile = new File(["dummy"], "avatar.png", { type: "image/png" });
    const result = validateImageFileSync(validFile);

    expect(result.ok).toBe(true);
    if (result.ok) {
      expect(result.file.name).toBe("avatar.png");
    }
  });

  it("rejects disallowed file types with InvalidFileTypeError", () => {
    const textFile = new File(["notes"], "notes.txt", { type: "text/plain" });

    const syncResult = validateImageFileSync(textFile);
    expect(syncResult.ok).toBe(false);
    if (!syncResult.ok) {
      expect(syncResult.error).toBe("Use uma imagem PNG, JPG ou WebP.");
    }

    const exit = Effect.runSyncExit(validateImageFile(textFile));
    expect(Exit.isFailure(exit)).toBe(true);
  });

  it("rejects files exceeding 10MB with FileSizeExceededError", () => {
    const largeFile = new File(["dummy"], "big.jpg", { type: "image/jpeg" });
    Object.defineProperty(largeFile, "size", { value: 11 * 1024 * 1024 });

    const syncResult = validateImageFileSync(largeFile);
    expect(syncResult.ok).toBe(false);
    if (!syncResult.ok) {
      expect(syncResult.error).toBe("A imagem deve ter no máximo 10 MB.");
    }

    const exit = Effect.runSyncExit(validateImageFile(largeFile));
    expect(Exit.isFailure(exit)).toBe(true);
  });
});

describe("uploadProfileImageEffect", () => {
  it("fails fast on invalid file without fetching", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch");
    const badFile = new File(["bad"], "bad.pdf", { type: "application/pdf" });

    const exit = await Effect.runPromiseExit(
      uploadProfileImageEffect(badFile, "pfpUrl"),
    );

    expect(Exit.isFailure(exit)).toBe(true);
    expect(fetchSpy).not.toHaveBeenCalled();
  });

  it("returns publicUrl on successful upload", async () => {
    vi.spyOn(globalThis, "fetch").mockImplementation(async () =>
      new Response(
        JSON.stringify({
          approved: true,
          publicUrl: "https://storage.univents.test/profiles/images/avatar.png",
        }),
        { status: 200 },
      ),
    );

    const file = new File(["data"], "avatar.png", { type: "image/png" });
    const url = await Effect.runPromise(uploadProfileImageEffect(file, "pfpUrl"));

    expect(url).toBe("https://storage.univents.test/profiles/images/avatar.png");
  });

  it("handles moderation rejection or missing publicUrl", async () => {
    vi.spyOn(globalThis, "fetch").mockImplementation(async () =>
      new Response(
        JSON.stringify({
          approved: false,
          error: "Conteúdo não aprovado pela moderação.",
        }),
        { status: 200 },
      ),
    );

    const file = new File(["data"], "banner.webp", { type: "image/webp" });
    const exit = await Effect.runPromiseExit(
      uploadProfileImageEffect(file, "bannerUrl"),
    );

    expect(Exit.isFailure(exit)).toBe(true);
    await expect(uploadProfileImage(file, "bannerUrl")).rejects.toThrow(
      "Conteúdo não aprovado pela moderação.",
    );
  });

  it("handles network failure gracefully", async () => {
    vi.spyOn(globalThis, "fetch").mockRejectedValue(new Error("Connection reset"));

    const file = new File(["data"], "avatar.jpg", { type: "image/jpeg" });
    const exit = await Effect.runPromiseExit(
      uploadProfileImageEffect(file, "pfpUrl"),
    );

    expect(Exit.isFailure(exit)).toBe(true);
  });
});

describe("uploadProfileImagesBatch & uploadProfileImagesBatchEffect", () => {
  it("uploads multiple images concurrently and partitions successes and failures", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockImplementation(
      async (_url, init) => {
        const formData = init?.body as FormData;
        const field = formData.get("field");
        if (field === "pfpUrl") {
          return new Response(
            JSON.stringify({
              approved: true,
              publicUrl: "https://storage.univents.test/pfp.png",
            }),
            { status: 200 },
          );
        }
        return new Response(
          JSON.stringify({
            approved: false,
            error: "Inapropriado",
          }),
          { status: 400 },
        );
      },
    );

    const avatarFile = new File(["a"], "avatar.png", { type: "image/png" });
    const bannerFile = new File(["b"], "banner.jpg", { type: "image/jpeg" });

    const batch = await uploadProfileImagesBatch([
      ["pfpUrl", avatarFile],
      ["bannerUrl", bannerFile],
    ]);

    expect(fetchSpy).toHaveBeenCalledTimes(2);
    expect(batch.uploaded.pfpUrl).toBe("https://storage.univents.test/pfp.png");
    expect(batch.uploaded.bannerUrl).toBeUndefined();
    expect(batch.failed).toEqual(["banner"]);
  });
});
