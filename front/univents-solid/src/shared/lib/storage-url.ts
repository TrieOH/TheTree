/**
 * Resolves a stored image path (or external/legacy absolute URL) to a publicly accessible URL.
 *
 * - Returns empty string for falsy/empty input.
 * - Leaves absolute URLs (http://, https://), protocol-relative (//), data URIs, and blob URLs untouched.
 * - Prepends the public storage base URL (from VITE_STORAGE_URL) to relative paths.
 */
export function resolveStorageUrl(pathOrUrl: string | null | undefined): string {
  if (!pathOrUrl) return "";
  const trimmed = pathOrUrl.trim();
  if (!trimmed) return "";

  // If already absolute or local blob/data URL, keep as is
  if (/^(https?:|data:|blob:|\/\/)/i.test(trimmed)) {
    return trimmed;
  }

  const rawBase =
    typeof import.meta !== "undefined" && import.meta.env
      ? (import.meta.env.VITE_STORAGE_URL || "")
      : "";

  const base = String(rawBase).trim().replace(/\/+$/, "");
  const cleanPath = trimmed.replace(/^\/+/, "");

  return base ? `${base}/${cleanPath}` : `/${cleanPath}`;
}
