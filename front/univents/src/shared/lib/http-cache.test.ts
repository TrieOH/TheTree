import { beforeEach, describe, expect, it, vi } from "vitest";

const { setResponseHeader } = vi.hoisted(() => ({
  setResponseHeader: vi.fn(),
}));

vi.mock("@tanstack/react-start/server", () => ({ setResponseHeader }));

import {
  CACHE_PRIVATE_NO_STORE,
  CACHE_PUBLIC_STATIC,
  preventResponseCaching,
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
});
