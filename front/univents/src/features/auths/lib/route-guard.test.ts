import { describe, expect, it } from "vitest";
import {
  isAuthOnlyPath,
  requiresVerifiedEmail,
  safeInternalReturnTo,
  verificationReturnTo,
} from "./auth-path";

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
