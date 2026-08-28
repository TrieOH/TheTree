import { describe, expect, it } from "vitest";
import {
  AUTH_RETURN_TO_STORAGE_KEY,
  clearAuthReturnTo,
  isAuthOnlyPath,
  readAuthReturnTo,
  requiresVerifiedEmail,
  safeInternalReturnTo,
  storeAuthReturnTo,
  verificationReturnTo,
} from "./auth-path";

function memoryStorage() {
  const values = new Map<string, string>();
  return {
    getItem: (key: string) => values.get(key) ?? null,
    setItem: (key: string, value: string) => values.set(key, value),
    removeItem: (key: string) => values.delete(key),
  };
}

describe("isAuthOnlyPath", () => {
  it.each([
    "/admin/events",
    "/checkouts/purchase-id",
    "/events/evento/checkout",
    "/profile",
    "/profile/config",
    "/profile/edit",
    "/profile/setup",
  ])("marks %s as private", (pathname) => {
    expect(isAuthOnlyPath(pathname)).toBe(true);
  });

  it.each([
    "/",
    "/events/evento",
    "/events/evento/programs",
    "/events/evento/store",
    "/profile/alguem",
  ])("marks %s as public", (pathname) => {
    expect(isAuthOnlyPath(pathname)).toBe(false);
  });
});

describe("safeInternalReturnTo", () => {
  it("uses the first internal destination", () => {
    expect(safeInternalReturnTo(undefined, "/events/evento/programs")).toBe(
      "/events/evento/programs",
    );
  });

  it("rejects external and protocol-relative destinations", () => {
    expect(safeInternalReturnTo("https://example.com", "//example.com")).toBe(
      "/profile",
    );
  });

  it("preserves a checkout nested in the profile setup return", () => {
    const destination = "/profile/setup?returnTo=%2Fevents%2Fevento%2Fcheckout";
    expect(safeInternalReturnTo(destination)).toBe(destination);
  });
});

describe("requiresVerifiedEmail", () => {
  it.each([
    "/admin/events",
    "/events/evento",
    "/events/evento/programs",
    "/events/evento/store",
    "/events/evento/checkout",
    "/checkouts/purchase-id",
    "/profile",
    "/profile/alguem",
    "/profile/setup",
  ])("requires verification on %s", (pathname) => {
    expect(requiresVerifiedEmail(pathname)).toBe(true);
  });

  it.each(["/", "/events", "/auth", "/auth/verify-email", "/terms"])(
    "does not require verification on %s",
    (pathname) => {
      expect(requiresVerifiedEmail(pathname)).toBe(false);
    },
  );
});

describe("verificationReturnTo", () => {
  it("stores the final checkout instead of nesting the profile setup", () => {
    expect(
      verificationReturnTo(
        "/profile/setup",
        "/profile/setup?returnTo=%2Fevents%2Fevento%2Fcheckout",
        "/events/evento/checkout",
      ),
    ).toBe("/events/evento/checkout");
  });
});

describe("auth return-to storage", () => {
  it("shares a valid internal destination until it expires", () => {
    const storage = memoryStorage();
    storeAuthReturnTo(storage, "/events/evento/checkout", 1_000);

    expect(readAuthReturnTo(storage, 1_001)).toBe("/events/evento/checkout");
    expect(readAuthReturnTo(storage, 601_001)).toBeNull();
    expect(storage.getItem(AUTH_RETURN_TO_STORAGE_KEY)).toBeNull();
  });

  it("rejects invalid data and clears explicitly", () => {
    const storage = memoryStorage();
    storage.setItem(AUTH_RETURN_TO_STORAGE_KEY, "invalid");
    expect(readAuthReturnTo(storage)).toBeNull();

    storeAuthReturnTo(storage, "/profile");
    clearAuthReturnTo(storage);
    expect(readAuthReturnTo(storage)).toBeNull();
  });
});
