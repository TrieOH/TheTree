import { env } from "cloudflare:workers";
import { describe, expect, it, vi } from "vitest";
import worker from "../../src/server";

const request = (url: string, init?: RequestInit) =>
  new Request(`http://localhost${url}`, init);

describe("worker routing", () => {
  it("answers /api/health from D1", async () => {
    const response = await worker.fetch(request("/api/health"), env);

    expect(response.status).toBe(200);
    expect(response.headers.get("Cache-Control")).toBe("no-store");
    await expect(response.json()).resolves.toEqual({ ok: true });
  });

  it("rejects writes to the read API", async () => {
    const response = await worker.fetch(request("/api/health", { method: "POST" }), env);

    expect(response.status).toBe(405);
    expect(response.headers.get("Allow")).toBe("GET");
  });

  it("404s unknown API and webhook paths instead of serving the SPA", async () => {
    const assets = vi.spyOn(env.ASSETS, "fetch");

    expect((await worker.fetch(request("/api/nope"), env)).status).toBe(404);
    expect((await worker.fetch(request("/webhooks/github", { method: "POST" }), env)).status).toBe(404);
    expect(assets).not.toHaveBeenCalled();
  });

  it("hands every other path to the SPA assets", async () => {
    const assets = vi.spyOn(env.ASSETS, "fetch").mockResolvedValue(new Response("spa"));

    const response = await worker.fetch(request("/docs/decisions/adr-0019"), env);

    expect(assets).toHaveBeenCalledOnce();
    await expect(response.text()).resolves.toBe("spa");
  });
});
