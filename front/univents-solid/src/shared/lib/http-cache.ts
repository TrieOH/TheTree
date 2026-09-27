export const CACHE_PUBLIC_STATIC =
  "public, max-age=300, stale-while-revalidate=86400";
export const CACHE_PRIVATE_NO_STORE = "private, no-store";

export function privateJsonResponse(data: unknown, init?: ResponseInit) {
  const headers = new Headers(init?.headers);
  headers.set("Cache-Control", CACHE_PRIVATE_NO_STORE);
  return Response.json(data, { ...init, headers });
}
