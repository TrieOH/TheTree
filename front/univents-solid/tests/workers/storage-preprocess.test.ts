import { env } from "cloudflare:workers";
import { afterEach, describe, expect, it, vi } from "vitest";
import {
  handleStorageImagePreprocess,
  isSafeModerationOutput,
} from "../../src/features/storage/api/storage-handlers";

function preprocessRequest(file: File, path = "profiles/images") {
  const form = new FormData();
  form.set("file", file);
  form.set("path", path);
  form.set("idempotencyKey", "test-key");
  return new Request("http://localhost/storage/image/preprocess", {
    method: "POST",
    body: form,
  });
}

afterEach(() => {
  vi.restoreAllMocks();
});

describe("isSafeModerationOutput", () => {
  it.each([
    ["safe", true],
    ["Safe", true],
    ["safe.", true],
    ["unsafe", false],
    ["The image is safe.", false],
    ["", false],
  ])("maps %j to %s", (description, expected) => {
    expect(isSafeModerationOutput({ description })).toBe(expected);
  });
});

describe("handleStorageImagePreprocess", () => {
  it("rejects a disallowed content type without calling Workers AI", async () => {
    const run = vi.spyOn(env.AI, "run");

    const response = await handleStorageImagePreprocess(
      preprocessRequest(new File(["notes"], "notes.txt", { type: "text/plain" })),
      env,
    );

    expect(response.status).toBe(400);
    expect(run).not.toHaveBeenCalled();
  });

  it("reports approved:false when the model rejects the image", async () => {
    vi.spyOn(env.AI, "run").mockResolvedValue({
      description: "This image is unsafe.",
    } as never);

    const response = await handleStorageImagePreprocess(
      preprocessRequest(new File(["pixels"], "pic.png", { type: "image/png" })),
      env,
    );

    expect(response.status).toBe(200);
    await expect(response.json()).resolves.toEqual({ approved: false });
  });

  it("uploads and returns the public URL when the model approves", async () => {
    vi.spyOn(env.AI, "run").mockResolvedValue({
      description: "safe",
    } as never);
    const upload = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(new Response(null, { status: 200 }));

    const response = await handleStorageImagePreprocess(
      preprocessRequest(new File(["pixels"], "pic.png", { type: "image/png" })),
      env,
    );

    expect(upload).toHaveBeenCalledOnce();
    const body = (await response.json()) as {
      approved: boolean;
      publicUrl: string;
    };
    expect(body.approved).toBe(true);
    expect(body.publicUrl).toContain("/profiles/images/");
  });
});
