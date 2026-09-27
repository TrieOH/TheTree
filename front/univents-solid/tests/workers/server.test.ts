import { env } from "cloudflare:workers";
import { describe, expect, it } from "vitest";
import worker from "../../src/server";

const request = (url: string, init?: RequestInit) =>
  new Request(`http://localhost${url}`, init);

describe("worker routing", () => {
  it("dispatches POST /storage/upload to the storage handler", async () => {
    const response = await worker.fetch(
      request("/storage/upload", {
        method: "POST",
        body: JSON.stringify({
          filename: "photo.gif",
          contentType: "image/gif",
          size: 10,
        }),
      }),
      env,
    );

    expect(response.status).toBe(400);
    await expect(response.json()).resolves.toMatchObject({
      error: expect.stringContaining("Only"),
    });
  });

  it("does not route GET requests to the POST-only handlers", async () => {
    const response = await worker.fetch(request("/storage/upload"), env);

    // Non-POST traffic falls through to the static assets binding, so it must
    // never come back as the storage handler's JSON error shape.
    expect(response.status).not.toBe(400);
  });
});
