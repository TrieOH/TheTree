import { describe, expect, it, vi } from "vitest";

vi.mock("@trieoh/front-core/tracing/browser", () => ({
  getTraceparent: () =>
    "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
}));

import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";

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
