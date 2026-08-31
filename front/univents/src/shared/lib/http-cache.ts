import { setResponseHeader } from "@tanstack/react-start/server";

export const CACHE_PUBLIC_STATIC =
  "public, max-age=300, stale-while-revalidate=86400";
export const CACHE_PRIVATE_NO_STORE = "private, no-store";

export function preventResponseCaching() {
  setResponseHeader("Cache-Control", CACHE_PRIVATE_NO_STORE);
}
