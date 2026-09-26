import { describe, expect, it } from "vitest";
import { resolveStorageUrl } from "../../src/shared/lib/storage-url";

describe("resolveStorageUrl", () => {
  it("returns empty string for null, undefined, or empty string", () => {
    expect(resolveStorageUrl(null)).toBe("");
    expect(resolveStorageUrl(undefined)).toBe("");
    expect(resolveStorageUrl("")).toBe("");
    expect(resolveStorageUrl("   ")).toBe("");
  });

  it("returns absolute URLs untouched", () => {
    expect(resolveStorageUrl("https://example.com/banner.png")).toBe(
      "https://example.com/banner.png",
    );
    expect(resolveStorageUrl("http://localhost:9000/univents/banner.png")).toBe(
      "http://localhost:9000/univents/banner.png",
    );
  });

  it("returns blob and data URLs untouched", () => {
    expect(resolveStorageUrl("blob:http://localhost:3000/123-abc")).toBe(
      "blob:http://localhost:3000/123-abc",
    );
    expect(resolveStorageUrl("data:image/png;base64,iVBORw0KGgoAAAANSU=")).toBe(
      "data:image/png;base64,iVBORw0KGgoAAAANSU=",
    );
  });

  it("resolves relative storage paths using base storage URL", () => {
    const resolved = resolveStorageUrl("events/event-1/banner.webp");
    expect(resolved).toBeTruthy();
    expect(resolved).toContain("events/event-1/banner.webp");
    // Ensure no double slashes except protocol
    expect(resolved.replace(/^https?:\/\//, "")).not.toContain("//");
  });

  it("handles leading slashes in relative paths gracefully", () => {
    const resolved = resolveStorageUrl("/events/event-1/logo.webp");
    expect(resolved).toBeTruthy();
    expect(resolved).toContain("events/event-1/logo.webp");
    expect(resolved.replace(/^https?:\/\//, "")).not.toContain("//");
  });
});
