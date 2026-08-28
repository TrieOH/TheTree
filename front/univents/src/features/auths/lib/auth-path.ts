export const AUTH_RETURN_TO_STORAGE_KEY = "univents:auth-return-to";
const AUTH_RETURN_TO_TTL_MS = 10 * 60 * 1000;

type ReturnToStorage = Pick<Storage, "getItem" | "setItem" | "removeItem">;

export function storeAuthReturnTo(
  storage: ReturnToStorage,
  path: string,
  now = Date.now(),
) {
  const safePath = safeInternalReturnTo(path);
  if (safePath !== path) return;
  storage.setItem(
    AUTH_RETURN_TO_STORAGE_KEY,
    JSON.stringify({ path, expiresAt: now + AUTH_RETURN_TO_TTL_MS }),
  );
}

export function readAuthReturnTo(storage: ReturnToStorage, now = Date.now()) {
  const raw = storage.getItem(AUTH_RETURN_TO_STORAGE_KEY);
  if (!raw) return null;
  try {
    const value = JSON.parse(raw) as { path?: unknown; expiresAt?: unknown };
    if (
      typeof value.path === "string" &&
      safeInternalReturnTo(value.path) === value.path &&
      typeof value.expiresAt === "number" &&
      value.expiresAt > now
    ) {
      return value.path;
    }
  } catch {
    // Invalid and legacy values are cleared below.
  }
  storage.removeItem(AUTH_RETURN_TO_STORAGE_KEY);
  return null;
}

export function clearAuthReturnTo(storage: ReturnToStorage) {
  storage.removeItem(AUTH_RETURN_TO_STORAGE_KEY);
}

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
