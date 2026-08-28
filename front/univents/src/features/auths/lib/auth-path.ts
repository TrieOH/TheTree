export const AUTH_RETURN_TO_STORAGE_KEY = "univents:auth-return-to";

export function isAuthOnlyPath(pathname: string) {
  return (
    pathname.startsWith("/admin") ||
    pathname.startsWith("/checkouts/") ||
    /^\/events\/[^/]+\/checkout\/?$/.test(pathname) ||
    [
      "/profile",
      "/profile/",
      "/profile/config",
      "/profile/edit",
      "/profile/setup",
    ].includes(pathname)
  );
}

export function requiresVerifiedEmail(pathname: string) {
  if (pathname === "/auth/verify-email") return false;
  return (
    isAuthOnlyPath(pathname) ||
    /^\/events\/[^/]+(?:\/|$)/.test(pathname) ||
    pathname.startsWith("/profile/")
  );
}

export function safeInternalReturnTo(
  ...candidates: (string | null | undefined)[]
) {
  return (
    candidates.find(
      (candidate) => candidate?.startsWith("/") && !candidate.startsWith("//"),
    ) ?? "/profile"
  );
}

export function verificationReturnTo(
  pathname: string,
  href: string,
  setupReturnTo?: string,
) {
  return safeInternalReturnTo(
    pathname === "/profile/setup" ? setupReturnTo : undefined,
    href,
  );
}
