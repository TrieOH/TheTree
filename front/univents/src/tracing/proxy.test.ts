import { describe, expect, it, vi } from "vitest";

vi.mock("@trieoh/front-core/tracing/browser", () => ({
  addActiveSpanEvent: () => undefined,
  getTraceparent: () =>
    "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
  setActiveSpanAttributes: () => undefined,
  withSpan: async <T>(_name: string, fn: () => Promise<T> | T) => fn(),
}));

import {
  createTanStackIdentityXIntegration,
  createTanStackServerProxyFetchers,
} from "@trieoh/front-core/auth/tanstack/client";

describe("BFF trace propagation", () => {
  it("adds the browser traceparent to proxy requests", async () => {
    const proxy = vi.fn().mockResolvedValue({
      success: true,
      code: 200,
      data: { ok: true },
    });
    const { authFetcher } = createTanStackServerProxyFetchers(proxy);

    await authFetcher.get("/events");

    expect(proxy).toHaveBeenCalledWith({
      data: {
        path: "/events",
        method: "GET",
        headers: {
          "content-type": "application/json",
          traceparent:
            "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
        },
      },
    });
  });

  it("preserves existing headers while adding traceparent", async () => {
    const proxy = vi.fn().mockResolvedValue({
      success: true,
      code: 200,
      data: { ok: true },
    });
    const { authFetcher } = createTanStackServerProxyFetchers(proxy);

    await authFetcher.post("/events", {
      headers: { "x-request-id": "request-1" },
      body: JSON.stringify({ title: "Test" }),
    });

    expect(proxy).toHaveBeenCalledWith({
      data: {
        path: "/events",
        method: "POST",
        headers: {
          "content-type": "application/json",
          traceparent:
            "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
        },
        body: {
          body: '{"title":"Test"}',
          headers: { "x-request-id": "request-1" },
        },
      },
    });
  });
});

describe("BFF session synchronization", () => {
  it("marks the browser session unauthenticated when an API request returns 401", async () => {
    const proxy = vi.fn().mockResolvedValue({
      success: false,
      code: 401,
      error_id: "SESSION_MISSING",
      message: "Session expired",
    });
    const integration = createTanStackIdentityXIntegration({
      login: vi.fn(),
      logout: vi.fn(),
      refresh: vi.fn(),
      restore: vi.fn(),
      request: proxy,
    });
    const setAuthenticated = vi.fn();
    const setProfile = vi.fn();

    integration.authAdapter?.createAuth({
      callbacks: {},
      defaultAuth: {} as never,
      getProfile: () => null,
      setProfile,
      setAuthenticated,
    });

    await integration.authFetcher.get("/events/joined");

    expect(setProfile).toHaveBeenCalledWith(null);
    expect(setAuthenticated).toHaveBeenCalledWith(false);
  });
});
