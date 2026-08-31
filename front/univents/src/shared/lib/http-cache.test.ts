import { beforeEach, describe, expect, it, vi } from "vitest";

const { setResponseHeader } = vi.hoisted(() => ({
  setResponseHeader: vi.fn(),
}));

vi.mock("@tanstack/react-start/server", () => ({ setResponseHeader }));

import {
  CACHE_PRIVATE_NO_STORE,
  CACHE_PUBLIC_STATIC,
  preventResponseCaching,
  privateJsonResponse,
} from "./http-cache";

describe("HTTP cache policies", () => {
  beforeEach(() => vi.clearAllMocks());

  it("prevents shared and browser caches for private responses", () => {
    preventResponseCaching();

    expect(setResponseHeader).toHaveBeenCalledWith(
      "Cache-Control",
      CACHE_PRIVATE_NO_STORE,
    );
  });

  it("keeps stale-while-revalidate compatible with Cloudflare", () => {
    expect(CACHE_PUBLIC_STATIC).toContain("max-age=300");
    expect(CACHE_PUBLIC_STATIC).toContain("stale-while-revalidate=");
    expect(CACHE_PUBLIC_STATIC).not.toContain("s-maxage");
  });

  it("creates private JSON responses without discarding other headers", async () => {
    const response = privateJsonResponse(
      { uploadUrl: "signed" },
      { status: 201, headers: { "X-Test": "kept" } },
    );

    expect(response.status).toBe(201);
    expect(response.headers.get("Cache-Control")).toBe(CACHE_PRIVATE_NO_STORE);
    expect(response.headers.get("X-Test")).toBe("kept");
    await expect(response.json()).resolves.toEqual({ uploadUrl: "signed" });
  });
});
