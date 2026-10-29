const CACHE_PRIVATE_NO_STORE = "private, no-store";

export function privateJsonResponse(data: unknown, init?: ResponseInit) {
  const headers = new Headers(init?.headers);
  headers.set("Cache-Control", CACHE_PRIVATE_NO_STORE);
  return Response.json(data, { ...init, headers });
}
