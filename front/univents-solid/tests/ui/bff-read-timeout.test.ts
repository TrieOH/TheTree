import { afterEach, describe, expect, it, vi } from "vitest";
import { createIdentityXBff } from "@trieoh/front-core/auth/bff/core";

// A token that is valid for a long time, so the BFF never tries to refresh.
const session = {
  read: () => ({
    tokens: {
      access_token: "a.b.c",
      refresh_token: "r.r.r",
      access_expires_at: new Date(Date.now() + 3_600_000).toISOString(),
      refresh_expires_at: new Date(Date.now() + 7_200_000).toISOString(),
      domain: "app.test",
    },
  }),
  write: vi.fn(),
  clear: vi.fn(),
};

function createBff(readTimeoutMs?: number) {
  return createIdentityXBff({
    identityX: { baseURL: "https://identityx.test" },
    apiBaseURL: "https://api.test",
    readTimeoutMs,
    session: () => session,
    request: () => new Request("https://app.test/_bff", { method: "POST" }),
    observability: { log: () => undefined },
  });
}

/** A fetch that never answers on its own, but honours the abort signal. */
function hangingFetch() {
  return vi.fn(
    (_url: string, init?: RequestInit) =>
      new Promise<Response>((_resolve, reject) => {
        init?.signal?.addEventListener("abort", () => reject(init.signal?.reason));
      }),
  );
}

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

describe("BFF read timeout", () => {
  it("gives up on a hung upstream read and reports it as a 504, instead of hanging", async () => {
    vi.useFakeTimers();
    const fetchMock = hangingFetch();
    vi.stubGlobal("fetch", fetchMock);

    const pending = createBff(1_000).request({ target: "api", path: "/events", method: "GET" });
    await vi.advanceTimersByTimeAsync(1_000);

    await expect(pending).resolves.toMatchObject({
      success: false,
      code: 504,
      error_id: "UPSTREAM_TIMEOUT",
    });
  });

  it("reports an unreachable upstream read as a 503", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => { throw new TypeError("fetch failed"); }));

    await expect(
      createBff().request({ target: "api", path: "/events", method: "GET" }),
    ).resolves.toMatchObject({ success: false, code: 503, error_id: "UPSTREAM_UNAVAILABLE" });
  });

  it("never cuts a write short: it may already have landed upstream", async () => {
    const fetchMock = vi.fn(async (_url: string, init?: RequestInit) => {
      expect(init?.signal).toBeUndefined();
      return new Response(JSON.stringify({ code: 201, data: { id: "1" } }), { status: 201 });
    });
    vi.stubGlobal("fetch", fetchMock);

    await expect(
      createBff(1).request({ target: "api", path: "/events", method: "POST", body: { a: 1 } }),
    ).resolves.toMatchObject({ success: true, code: 201 });
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("lets a write's network failure surface as an error rather than inventing an outcome", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => { throw new TypeError("fetch failed"); }));

    await expect(
      createBff().request({ target: "api", path: "/events", method: "POST", body: {} }),
    ).rejects.toThrow("fetch failed");
  });

  it("can be turned off", async () => {
    const fetchMock = vi.fn(async (_url: string, init?: RequestInit) => {
      expect(init?.signal).toBeUndefined();
      return new Response(JSON.stringify({ code: 200, data: [] }), { status: 200 });
    });
    vi.stubGlobal("fetch", fetchMock);

    await createBff(0).request({ target: "api", path: "/events", method: "GET" });
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });
});
